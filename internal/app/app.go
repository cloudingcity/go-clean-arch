package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Run() error {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	if err := r.Run(":8080"); err != nil {
		return err
	}

	return nil
}
