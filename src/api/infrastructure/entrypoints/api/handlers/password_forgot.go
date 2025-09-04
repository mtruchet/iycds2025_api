package handlers

import (
	"net/http"

	"iycds2025_api/src/api/core/entities"
	"iycds2025_api/src/api/core/errors"
	"iycds2025_api/src/api/core/usecases/password"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type PasswordForgot struct {
	UseCase password.ForgotPassword
}

func (handler *PasswordForgot) Handle(c *gin.Context) {
	validate := validator.New()

	var forgotRequest entities.ForgotPassword
	if err := c.ShouldBindJSON(&forgotRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	if err := validate.Struct(forgotRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data", "details": err.Error()})
		return
	}

	// Ejecutar el caso de uso
	_, err := handler.UseCase.Execute(c.Request.Context(), &forgotRequest)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.Code, gin.H{"error": apiErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Por seguridad, siempre respondemos con éxito, independientemente de si el email existe o no
	c.JSON(http.StatusOK, gin.H{"message": "If your email is registered in our system, you will receive instructions to reset your password"})
}
