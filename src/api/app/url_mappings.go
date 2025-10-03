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

	// Endpoint público para obtener categorías
	group.GET("/categories", handlers.Categories.Handle)

	// Endpoint público para obtener todos los servicios
	group.GET("/services", handlers.ServiceListAll.Handle)

	// Endpoint público para obtener un servicio por ID
	group.GET("/services/:id", handlers.ServiceGetByID.Handle)

	// Endpoint público para obtener disponibilidad de un servicio
	group.GET("/services/:id/availability", handlers.ServiceAvailability.Handle)

	// Endpoint público para obtener calendario de un servicio (30 días)
	group.GET("/services/:id/calendar", handlers.ServiceCalendar.Handle)

	// Endpoints protegidos que requieren autenticación
	protected := group.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		// Actualización de perfil de usuario
		protected.PUT("/user/profile", middleware.StandardRateLimit(), handlers.UserUpdate.Handle)
		
		// CRUD de servicios
		protected.POST("/services", middleware.StandardRateLimit(), handlers.ServiceCreate.Handle)
		protected.PUT("/services/:id", middleware.StandardRateLimit(), handlers.ServiceUpdate.Handle)
		protected.DELETE("/services/:id", middleware.StandardRateLimit(), handlers.ServiceDelete.Handle)
		protected.PATCH("/services/:id/status", middleware.StandardRateLimit(), handlers.ServiceUpdateStatus.Handle)
		protected.GET("/my-services", middleware.StandardRateLimit(), handlers.ServiceList.Handle)
		
		// Appointments de servicios (para proveedores)
		protected.GET("/services/:id/appointments", middleware.StandardRateLimit(), handlers.ServiceAppointments.Handle)
		
		// CRUD de appointments/citas
		protected.POST("/appointments", middleware.StandardRateLimit(), handlers.AppointmentCreate.Handle)
		protected.GET("/my-appointments", middleware.StandardRateLimit(), handlers.AppointmentList.Handle)
		protected.PUT("/appointments/:id/status", middleware.StandardRateLimit(), handlers.AppointmentUpdateStatus.Handle)
	}
}
