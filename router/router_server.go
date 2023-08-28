package router

import (
	"net/http"
	"os"
	"strconv"

	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	"github.com/pterodactyl/wings/router/downloader"
	"github.com/pterodactyl/wings/router/middleware"
	"github.com/pterodactyl/wings/server"
	"github.com/pterodactyl/wings/server/transfer"
)

// Returns a single server from the collection of servers.
func getServer(c *gin.Context) {
	c.JSON(http.StatusOK, ExtractServer(c).ToAPIResponse())
}

// Returns the logs for a given server instance.
func getServerLogs(c *gin.Context) {
	s := ExtractServer(c)

	l, _ := strconv.Atoi(c.DefaultQuery("size", "100"))
	if l <= 0 {
		l = 100
	} else if l > 100 {
		l = 100
	}

	out, err := s.ReadLogfile(l)
	if err != nil {
		middleware.CaptureAndAbort(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": out})
}

// Deletes a server from the wings daemon and dissociate its objects.
func deleteServer(c *gin.Context) {
	s := middleware.ExtractServer(c)

	// Immediately suspend the server to prevent a user from attempting
	// to start it while this process is running.
	s.Config().SetSuspended(true)

	// Notify all websocket clients that the server is being deleted.
	// This is useful for two reasons, one to tell clients not to bother
	// retrying to connect to the websocket.  And two, for transfers when
	// the server has been successfully transferred to another node, and
	// the client needs to switch to the new node.
	if s.IsTransferring() {
		s.Events().Publish(server.TransferStatusEvent, transfer.StatusCompleted)
	}
	s.Events().Publish(server.DeletedEvent, nil)

	s.CleanupForDestroy()

	// Remove any pending remote file downloads for the server.
	for _, dl := range downloader.ByServer(s.ID()) {
		dl.Cancel()
	}

	// Destroy the environment; in Docker this will handle a running container and
	// forcibly terminate it before removing the container, so we do not need to handle
	// that here.
	if err := s.Environment.Destroy(); err != nil {
		middleware.CaptureAndAbort(c, err)
		return
	}

	// Once the environment is terminated, remove the server files from the system. This is
	// done in a separate process since failure is not the end of the world and can be
	// manually cleaned up after the fact.
	//
	// In addition, servers with large amounts of files can take some time to finish deleting,
	// so we don't want to block the HTTP call while waiting on this.
	go func(p string) {
		if err := os.RemoveAll(p); err != nil {
			log.WithFields(log.Fields{"path": p, "error": err}).Warn("failed to remove server files during deletion process")
		}
	}(s.Filesystem().Path())

	middleware.ExtractManager(c).Remove(func(server *server.Server) bool {
		return server.ID() == s.ID()
	})

	// Deallocate the reference to this server.
	s = nil

	c.Status(http.StatusNoContent)
}
