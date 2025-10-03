package handlers

import (
	"net/http"

	"iycds2025_api/src/api/core/usecases/appointment"

	"github.com/gin-gonic/gin"
)

type AppointmentListHandler struct {
	ListMyAppointments appointment.ListMyAppointments
}

func (h *AppointmentListHandler) Handle(c *gin.Context) {
	// Obtener userID del contexto
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	// Ejecutar use case
	response, err := h.ListMyAppointments.Execute(c.Request.Context(), userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Appointments retrieved successfully",
		"data":    response,
	})
}
