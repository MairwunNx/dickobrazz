package server

import (
	"net/http"

	"dickobrazz/server/config"
	"dickobrazz/server/handlers"
	"dickobrazz/server/middleware"
	"github.com/gin-gonic/gin"
)

func newRouter(cfg config.Config, handler *handlers.Handler) *gin.Engine {
	router := gin.New()
	router.Use(middleware.RequestID(), middleware.CORS(cfg), gin.Logger(), gin.Recovery())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := router.Group("/api/v1")
	api.POST("/auth/login", handler.Login)

	protected := api.Group("/")
	protected.Use(middleware.RequireAuth(cfg))
	protected.GET("/me", handler.Me)
	protected.POST("/cock/size", handler.CockSize)
	protected.GET("/cock/ruler", handler.CockRuler)
	protected.GET("/cock/race", handler.CockRace)
	protected.GET("/cock/dynamic", handler.CockDynamic)
	protected.GET("/cock/seasons", handler.CockSeasons)
	protected.GET("/cock/achievements", handler.CockAchievements)
	protected.GET("/cock/ladder", handler.CockLadder)

	return router
}
