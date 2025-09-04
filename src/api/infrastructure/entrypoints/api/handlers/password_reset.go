package handlers

import (
	"net/http"

	"iycds2025_api/src/api/core/entities"
	"iycds2025_api/src/api/core/errors"
	"iycds2025_api/src/api/core/usecases/password"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type PasswordReset struct {
	UseCase password.ResetPassword
}

func (handler *PasswordReset) Handle(c *gin.Context) {
	validate := validator.New()

	var resetRequest entities.ResetPassword
	if err := c.ShouldBindJSON(&resetRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	if err := validate.Struct(resetRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data", "details": err.Error()})
		return
	}

	// Ejecutar el caso de uso
	err := handler.UseCase.Execute(c.Request.Context(), &resetRequest)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.Code, gin.H{"error": apiErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}
