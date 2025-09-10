package handlers

import (
	"net/http"
	"strconv"

	"iycds2025_api/src/api/core/usecases/service"

	"github.com/gin-gonic/gin"
)

type ServiceGetByIDHandler struct {
	GetServiceByID service.GetServiceByID
}

func (h *ServiceGetByIDHandler) Handle(c *gin.Context) {
	// Obtener ID del servicio desde URL
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid service ID",
		})
		return
	}

	// Ejecutar use case
	response, err := h.GetServiceByID.Execute(c.Request.Context(), id)
	if err != nil {
		// Manejar diferentes tipos de errores
		switch err.Error() {
		case "Service not found":
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Service not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Service retrieved successfully",
		"data":    response,
	})
}
