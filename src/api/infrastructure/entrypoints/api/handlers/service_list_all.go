package handlers

import (
	"net/http"

	"iycds2025_api/src/api/core/usecases/service"

	"github.com/gin-gonic/gin"
)

type ServiceListAllHandler struct {
	ListAllServices service.ListAllServices
}

func (h *ServiceListAllHandler) Handle(c *gin.Context) {
	// Ejecutar use case
	response, err := h.ListAllServices.Execute(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Services retrieved successfully",
		"data":    response,
	})
}
