package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewPingRoutes(rg *gin.RouterGroup) {
	rg.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
}
