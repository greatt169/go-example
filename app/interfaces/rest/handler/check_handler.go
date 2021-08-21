package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// health - handler for check service status
func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
