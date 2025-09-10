package handlers

import (
	"net/http"
	"strconv"

	"iycds2025_api/src/api/core/entities"
	"iycds2025_api/src/api/core/usecases/service"

	"github.com/gin-gonic/gin"
)

type ServiceUpdateHandler struct {
	UpdateService service.UpdateService
}

func (h *ServiceUpdateHandler) Handle(c *gin.Context) {
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

	// Parsear request body
	var serviceReq entities.ServiceUpdate
	if err := c.ShouldBindJSON(&serviceReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body: " + err.Error(),
		})
		return
	}

	// Ejecutar use case
	response, err := h.UpdateService.Execute(c.Request.Context(), id, &serviceReq, userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Service updated successfully",
		"data":    response,
	})
}
