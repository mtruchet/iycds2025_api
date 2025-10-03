package handlers

import (
	"net/http"
	"strconv"

	"iycds2025_api/src/api/core/usecases/appointment"

	"github.com/gin-gonic/gin"
)

type ServiceAvailabilityHandler struct {
	GetServiceAvailability appointment.GetServiceAvailability
}

func (h *ServiceAvailabilityHandler) Handle(c *gin.Context) {
	// Obtener ID del servicio desde URL
	idStr := c.Param("id")
	serviceID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid service ID",
		})
		return
	}

	// Obtener fecha desde query parameter
	date := c.Query("date")
	if date == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Date parameter is required (format: YYYY-MM-DD)",
		})
		return
	}

	// Ejecutar use case
	response, err := h.GetServiceAvailability.Execute(c.Request.Context(), serviceID, date)
	if err != nil {
		// Manejar diferentes tipos de errores
		switch {
		case err.Error() == "Service not found":
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
		case err.Error() == "Service is not active":
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

	c.JSON(http.StatusOK, gin.H{
		"message": "Service availability retrieved successfully",
		"data":    response,
	})
}
