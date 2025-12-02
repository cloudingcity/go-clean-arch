package app

import (
	"github.com/cloudingcity/todo/internal/handler/http"
	"github.com/cloudingcity/todo/internal/repo/memory"
	"github.com/cloudingcity/todo/internal/service/todo"
	"github.com/gin-gonic/gin"
)

func Run() error {
	r := gin.Default()

	todoRepo := memory.NewTodoRepo()
	todoSrv := todo.NewService(todoRepo)
	http.NewRouter(r, todoSrv)

	if err := r.Run(":8080"); err != nil {
		return err
	}

	return nil
}
