package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Ping struct{}

func (handler *Ping) Handle(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
