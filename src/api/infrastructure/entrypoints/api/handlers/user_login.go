package handlers

import (
	"net/http"

	"iycds2025_api/src/api/core/entities"
	"iycds2025_api/src/api/core/errors"
	"iycds2025_api/src/api/core/usecases/login"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserLogin struct {
	UseCase login.UserLogin
}

func (handler *UserLogin) Handle(c *gin.Context) {
	validate := validator.New()

	var loginRequest entities.Login
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	if err := validate.Struct(loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data", "details": err.Error()})
		return
	}

	response, err := handler.UseCase.Execute(c.Request.Context(), &loginRequest)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.Code, gin.H{"error": apiErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": response})
}
