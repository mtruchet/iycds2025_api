package handlers

import (
	"net/http"

	"iycds2025_api/src/api/core/entities"
	"iycds2025_api/src/api/core/usecases/user"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserUpdateHandler struct {
	UpdateUser user.UpdateUser
}

func (h *UserUpdateHandler) Handle(c *gin.Context) {
	// Obtener userID del contexto (establecido por AuthMiddleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	// Parsear request body
	var userUpdate entities.UserUpdate
	if err := c.ShouldBindJSON(&userUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body: " + err.Error(),
		})
		return
	}

	// Validar datos
	validate := validator.New()
	if err := validate.Struct(userUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input data: " + err.Error(),
		})
		return
	}

	// Ejecutar use case
	response, err := h.UpdateUser.Execute(c.Request.Context(), userID.(int64), &userUpdate)
	if err != nil {
		// Manejar diferentes tipos de errores
		switch {
		case err.Error() == "User not found":
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
		case err.Error() == "Email is already in use":
			c.JSON(http.StatusConflict, gin.H{
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
		"message": "User updated successfully",
		"data":    response,
	})
}
