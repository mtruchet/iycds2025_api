package app

import (
	"iycds2025_api/src/api/infrastructure/dependencies"

	"github.com/gin-gonic/gin"
)

func configureURLMappings(router *gin.Engine, handlers *dependencies.HandlerContainer) {
	// Endpoint simple ping/pong
	router.GET("/ping", handlers.Ping.Handle)
}
