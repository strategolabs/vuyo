package router

import (
	"context"
	"errors"
	"net/http"

	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	"github.com/pterodactyl/wings/router/middleware"
	"github.com/pterodactyl/wings/server"
	"github.com/pterodactyl/wings/server/installer"
)

// Returns all the servers that are registered and configured correctly on
// this wings instance.
func getAllServers(c *gin.Context) {
	servers := middleware.ExtractManager(c).All()
	out := make([]server.APIResponse, len(servers), len(servers))
	for i, v := range servers {
		out[i] = v.ToAPIResponse()
	}
	c.JSON(http.StatusOK, out)
}

// Creates a new server on the wings daemon and begins the installation process
// for it.
func postCreateServer(c *gin.Context) {
	manager := middleware.ExtractManager(c)

	details := installer.ServerDetails{}
	if err := c.BindJSON(&details); err != nil {
		return
	}

	install, err := installer.New(c.Request.Context(), manager, details)
	if err != nil {
		if installer.IsValidationError(err) {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"error": "The data provided in the request could not be validated.",
			})
			return
		}

		middleware.CaptureAndAbort(c, err)
		return
	}

	// Plop that server instance onto the request so that it can be referenced in
	// requests from here-on out.
	manager.Add(install.Server())

	// Begin the installation process in the background to not block the request
	// cycle. If there are any errors they will be logged and communicated back
	// to the Panel where a reinstall may take place.
	go func(i *installer.Installer) {
		if err := i.Server().CreateEnvironment(); err != nil {
			i.Server().Log().WithField("error", err).Error("failed to create server environment during install process")
			return
		}

		if err := i.Server().Install(); err != nil {
			log.WithFields(log.Fields{"server": i.Server().ID(), "error": err}).Error("failed to run install process for server")
			return
		}

		if i.StartOnCompletion {
			log.WithField("server_id", i.Server().ID()).Debug("starting server after successful installation")
			if err := i.Server().HandlePowerAction(server.PowerActionStart, 30); err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					log.WithFields(log.Fields{"server_id": i.Server().ID(), "action": "start"}).Warn("could not acquire a lock while attempting to perform a power action")
				} else {
					log.WithFields(log.Fields{"server_id": i.Server().ID(), "action": "start", "error": err}).Error("encountered error processing a server power action in the background")
				}
			}
		} else {
			log.WithField("server_id", i.Server().ID()).Debug("skipping automatic start after successful server installation")
		}
	}(install)

	c.Status(http.StatusAccepted)
}
