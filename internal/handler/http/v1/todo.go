package v1

import (
	"net/http"
	"time"

	"github.com/cloudingcity/todo/internal/service"
	"github.com/gin-gonic/gin"
)

type todoHandler struct {
	srv service.Todo
}

func NewTodoRoutes(rg *gin.RouterGroup, srv service.Todo) {
	h := &todoHandler{
		srv: srv,
	}
	rg.POST("/todos", h.create)
}

type createTodoReq struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type createTodoResp struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	IsCompleted bool      `json:"isCompleted"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (h *todoHandler) create(c *gin.Context) {
	var req createTodoReq

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todo, err := h.srv.Create(req.Title, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createTodoResp{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
		IsCompleted: todo.IsCompleted,
		CreatedAt:   todo.CreatedAt,
		UpdatedAt:   todo.UpdatedAt,
	})
}
