package router

import (
	"log"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/strategodev/vuyo/router/middleware"
)

// Configure configures the routing infrastructure for this daemon instance.
func Configure() *gin.Engine {
	// gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())

	config := cors.DefaultConfig()
	config.AllowCredentials = true
	config.AllowOrigins = []string{"http://localhost:3000"}

	router.Use(cors.New(config))

	authMiddleware := middleware.AuthMiddleware()

	router.POST("/login", authMiddleware.LoginHandler)

	router.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	auth := router.Group("/auth")
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/hello", middleware.HelloHandler)
	}

	// This route is special it sits above all the other requests because we are
	// using a JWT to authorize access to it, therefore it needs to be publicly
	// accessible.
	// router.GET("/api/servers/:server/ws", middleware.ServerExists(), getServerWebsocket)

	// This request is called by another daemon when a server is going to be transferred out.
	// This request does not need the AuthorizationMiddleware as the panel should never call it
	// and requests are authenticated through a JWT the panel issues to the other daemon.
	// router.POST("/api/transfers", postTransfers)

	// All the routes beyond this mount will use an authorization middleware
	// and will not be accessible without the correct Authorization header provided.
	// protected := router.Use(middleware.RequireAuthorization())
	// protected.POST("/api/update", postUpdateConfiguration)
	// protected.GET("/api/system", getSystemInformation)
	// protected.GET("/api/servers", getAllServers)
	// protected.POST("/api/servers", postCreateServer)
	// protected.DELETE("/api/transfers/:server", deleteTransfer)

	// These are server specific routes, and require that the request be authorized, and
	// that the server exist on the Daemon.
	server := router.Group("/api/servers/:server")
	// server.Use(middleware.RequireAuthorization(), middleware.ServerExists())
	{
		server.GET("", getServer)
		server.DELETE("", deleteServer)

		server.GET("/logs", getServerLogs)
		// server.POST("/power", postServerPower)
		// server.POST("/commands", postServerCommands)
		// server.POST("/install", postServerInstall)
		// server.POST("/reinstall", postServerReinstall)
		// server.POST("/sync", postServerSync)
		// server.POST("/ws/deny", postServerDenyWSTokens)

		// This archive request causes the archive to start being created
		// this should only be triggered by the panel.
		// server.POST("/transfer", postServerTransfer)
		// server.DELETE("/transfer", deleteServerTransfer)

		// files := server.Group("/files")
		// {
		// 	files.GET("/contents", getServerFileContents)
		// 	files.GET("/list-directory", getServerListDirectory)
		// 	files.PUT("/rename", putServerRenameFiles)
		// 	files.POST("/copy", postServerCopyFile)
		// 	files.POST("/write", postServerWriteFile)
		// 	files.POST("/create-directory", postServerCreateDirectory)
		// 	files.POST("/delete", postServerDeleteFiles)
		// 	files.POST("/compress", postServerCompressFiles)
		// 	files.POST("/decompress", postServerDecompressFiles)
		// 	files.POST("/chmod", postServerChmodFile)

		// 	files.GET("/pull", middleware.RemoteDownloadEnabled(), getServerPullingFiles)
		// 	files.POST("/pull", middleware.RemoteDownloadEnabled(), postServerPullRemoteFile)
		// 	files.DELETE("/pull/:download", middleware.RemoteDownloadEnabled(), deleteServerPullRemoteFile)
		// }

		// backup := server.Group("/backup")
		// {
		// 	backup.POST("", postServerBackup)
		// 	backup.POST("/:backup/restore", postServerRestoreBackup)
		// 	backup.DELETE("/:backup", deleteServerBackup)
		// }
	}

	return router
}
