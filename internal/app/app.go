package app

import (
	"github.com/cloudingcity/todo/internal/handler/http"
	"github.com/gin-gonic/gin"
)

func Run() error {
	r := gin.Default()

	http.NewRouter(r)

	if err := r.Run(":8080"); err != nil {
		return err
	}

	return nil
}
