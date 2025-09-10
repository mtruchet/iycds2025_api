package handlers

import (
	"net/http"

	"iycds2025_api/src/api/utils"

	"github.com/gin-gonic/gin"
)

type CategoriesHandler struct{}

func (h *CategoriesHandler) Handle(c *gin.Context) {
	categories := utils.GetValidCategories()
	
	c.JSON(http.StatusOK, gin.H{
		"message":    "Categories retrieved successfully",
		"categories": categories,
	})
}
