package handlers

import (
	"net/http"
	"strconv"

	"iycds2025_api/src/api/core/usecases/appointment"

	"github.com/gin-gonic/gin"
)

type ServiceAppointmentsHandler struct {
	ListServiceAppointments appointment.ListServiceAppointments
}

func (h *ServiceAppointmentsHandler) Handle(c *gin.Context) {
	// Obtener userID del contexto
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	// Obtener ID del servicio desde URL
	idStr := c.Param("id")
	serviceID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid service ID",
		})
		return
	}

	// Ejecutar use case
	response, err := h.ListServiceAppointments.Execute(c.Request.Context(), serviceID, userID.(int64))
	if err != nil {
		// Manejar diferentes tipos de errores
		switch {
		case err.Error() == "Service not found":
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
		case err.Error() == "You don't have permission to view appointments for this service":
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

	c.JSON(http.StatusOK, gin.H{
		"message": "Service appointments retrieved successfully",
		"data":    response,
	})
}
