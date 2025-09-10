package handlers

import (
	"net/http"

	"iycds2025_api/src/api/core/entities"
	"iycds2025_api/src/api/core/usecases/service"

	"github.com/gin-gonic/gin"
)

type ServiceCreateHandler struct {
	CreateService service.CreateService
}

func (h *ServiceCreateHandler) Handle(c *gin.Context) {
	// Obtener userID del contexto (establecido por AuthMiddleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	// Parsear request body
	var serviceReq entities.ServiceCreate
	if err := c.ShouldBindJSON(&serviceReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body: " + err.Error(),
		})
		return
	}

	// Ejecutar use case
	response, err := h.CreateService.Execute(c.Request.Context(), &serviceReq, userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Service created successfully",
		"data":    response,
	})
}
