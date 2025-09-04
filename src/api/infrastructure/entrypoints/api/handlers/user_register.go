package handlers

import (
	"net/http"

	"iycds2025_api/src/api/core/entities"
	"iycds2025_api/src/api/core/errors"
	"iycds2025_api/src/api/core/usecases/register"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserRegister struct {
	UseCase register.UserRegister
}

func (handler *UserRegister) Handle(c *gin.Context) {
	validate := validator.New()

	var registerRequest entities.UserRegister
	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	if err := validate.Struct(registerRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data", "details": err.Error()})
		return
	}

	user, err := handler.UseCase.Execute(c.Request.Context(), &registerRequest)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.Code, gin.H{"error": apiErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// No devolvemos la contraseña en la respuesta
	responseUser := gin.H{
		"id":         user.ID,
		"name":       user.Name,
		"email":      user.Email,
		"locality":   user.Locality,
		"province":   user.Province,
		"phone":      user.Phone,
		"created_at": user.CreatedAt,
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user":    responseUser,
	})
}
