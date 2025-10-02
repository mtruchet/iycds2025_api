package handlers

import (
	"net/http"
	"strconv"

	"iycds2025_api/src/api/core/usecases/service"

	"github.com/gin-gonic/gin"
)

type ServiceDeleteHandler struct {
	DeleteService service.DeleteService
}

func (h *ServiceDeleteHandler) Handle(c *gin.Context) {
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
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid service ID",
		})
		return
	}

	// Ejecutar use case
	err = h.DeleteService.Execute(c.Request.Context(), id, userID.(int64))
	if err != nil {
		// Manejar diferentes tipos de errores
		switch {
		case err.Error() == "Service not found" || err.Error() == "Service not found or already deleted":
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
		case err.Error() == "You don't have permission to delete this service":
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
		"message": "Service deleted successfully",
	})
}
