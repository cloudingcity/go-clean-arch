package http

import (
	"github.com/cloudingcity/todo/internal/handler/http/v1"
	"github.com/gin-gonic/gin"
)

func NewRouter(r *gin.Engine) {
	v1Group := r.Group("/v1")
	{
		v1.NewPingRoutes(v1Group)
	}
}
