package app

import (
	"iycds2025_api/src/api/infrastructure/dependencies"

	"github.com/gin-gonic/gin"
)

func configureURLMappings(router *gin.Engine, handlers *dependencies.HandlerContainer) {
	// Endpoint simple ping/pong
	router.GET("/ping", handlers.Ping.Handle)

	// Grupo de API
	group := router.Group("/api")

	// Endpoints públicos para autenticación
	group.POST("/user/login", handlers.UserLogin.Handle)
	group.POST("/user/register", handlers.UserRegister.Handle)
	group.POST("/user/forgot-password", handlers.PasswordForgot.Handle)
	group.POST("/user/reset-password", handlers.PasswordReset.Handle)
}
