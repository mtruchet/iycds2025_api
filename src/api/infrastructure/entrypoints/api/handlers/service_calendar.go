package handlers

import (
	"net/http"
	"strconv"

	"iycds2025_api/src/api/core/usecases/service"

	"github.com/gin-gonic/gin"
)

type ServiceCalendarHandler struct {
	GetServiceCalendar *service.GetServiceCalendarUseCase
}

func (h *ServiceCalendarHandler) Handle(c *gin.Context) {
	// Obtener ID del servicio desde la URL
	serviceIDParam := c.Param("id")
	serviceID, err := strconv.Atoi(serviceIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid service ID",
		})
		return
	}

	// Ejecutar use case
	result, apiErr := h.GetServiceCalendar.Execute(serviceID)

	if apiErr != nil {
		c.JSON(apiErr.Code, gin.H{
			"error": apiErr.Message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Service calendar retrieved successfully",
		"data":    result,
	})
}
