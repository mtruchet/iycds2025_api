package handlers

import (
	"net/http"
	"strconv"

	"iycds2025_api/src/api/core/entities"
	"iycds2025_api/src/api/core/usecases/appointment"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AppointmentUpdateStatusHandler struct {
	UpdateAppointmentStatus appointment.UpdateAppointmentStatus
}

func (h *AppointmentUpdateStatusHandler) Handle(c *gin.Context) {
	// Obtener userID del contexto
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	// Obtener ID del appointment desde URL
	idStr := c.Param("id")
	appointmentID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid appointment ID",
		})
		return
	}

	// Parsear el body de la request
	var req entities.AppointmentUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body: " + err.Error(),
		})
		return
	}

	// Validar datos
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input data: " + err.Error(),
		})
		return
	}

	// Ejecutar use case
	err = h.UpdateAppointmentStatus.Execute(c.Request.Context(), appointmentID, req.Status, userID.(int64))
	if err != nil {
		// Manejar diferentes tipos de errores
		switch {
		case err.Error() == "Appointment not found":
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
		case err.Error() == "You don't have permission to cancel this appointment" ||
			 err.Error() == "Only the service provider can accept or reject appointments" ||
			 err.Error() == "Only the service provider can mark appointments as completed":
			c.JSON(http.StatusForbidden, gin.H{
				"error": err.Error(),
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		return
	}

	// Mensaje dinámico basado en el status
	message := "Appointment " + req.Status + " successfully"
	c.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}
