package http

import (
	pingV1 "github.com/cloudingcity/todo/internal/handler/http/v1/ping"
	"github.com/gin-gonic/gin"
)

func NewRouter(r *gin.Engine) {
	v1Group := r.Group("/v1")
	{
		pingV1.NewRoutes(v1Group)
	}
}
