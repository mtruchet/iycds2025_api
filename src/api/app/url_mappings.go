package app

import (
	"iycds2025_api/src/api/infrastructure/dependencies"
	"iycds2025_api/src/api/middleware"

	"github.com/gin-gonic/gin"
)

func configureURLMappings(router *gin.Engine, handlers *dependencies.HandlerContainer) {
	// Endpoint simple ping/pong
	router.GET("/ping", handlers.Ping.Handle)

	// Grupo de API
	group := router.Group("/api")

	// Endpoints públicos para autenticación con rate limiting
	group.POST("/user/login", middleware.StrictRateLimit(), handlers.UserLogin.Handle)
	group.POST("/user/register", middleware.StandardRateLimit(), handlers.UserRegister.Handle)
	
	group.POST("/user/forgot-password", middleware.StrictRateLimit(), handlers.PasswordForgot.Handle)
	group.POST("/user/reset-password", middleware.StandardRateLimit(), handlers.PasswordReset.Handle)
}
