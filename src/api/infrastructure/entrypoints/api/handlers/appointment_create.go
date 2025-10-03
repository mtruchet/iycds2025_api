package handlers

import (
	"net/http"

	"iycds2025_api/src/api/core/entities"
	"iycds2025_api/src/api/core/usecases/appointment"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AppointmentCreateHandler struct {
	CreateAppointment appointment.CreateAppointment
}

func (h *AppointmentCreateHandler) Handle(c *gin.Context) {
	// Obtener userID del contexto (cliente que crea la cita)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	// Parsear request body
	var appointmentReq entities.AppointmentCreate
	if err := c.ShouldBindJSON(&appointmentReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body: " + err.Error(),
		})
		return
	}

	// Validar datos
	validate := validator.New()
	if err := validate.Struct(appointmentReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input data: " + err.Error(),
		})
		return
	}

	// Ejecutar use case
	response, err := h.CreateAppointment.Execute(c.Request.Context(), &appointmentReq, userID.(int64))
	if err != nil {
		// Manejar diferentes tipos de errores
		switch {
		case err.Error() == "Service not found":
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
		case err.Error() == "Time slot is already occupied or service not found":
			c.JSON(http.StatusConflict, gin.H{
				"error": "Time slot is already occupied",
			})
		case err.Error() == "Cannot create appointment for your own service":
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Appointment created successfully",
		"data":    response,
	})
}
