package http

import (
	"github.com/cloudingcity/todo/internal/handler/http/v1"
	"github.com/cloudingcity/todo/internal/service"
	"github.com/gin-gonic/gin"
)

func NewRouter(r *gin.Engine, todoSrv service.Todo) {
	v1Group := r.Group("/v1")
	{
		v1.NewPingRoutes(v1Group)
		v1.NewTodoRoutes(v1Group, todoSrv)
	}
}
