package v1

import (
	"errors"
	"net/http"
	"time"

	"github.com/cloudingcity/todo/internal/entity"
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
	rg.GET("/todos", h.list)
	rg.GET("/todos/:id", h.get)
	rg.PATCH("/todos/:id", h.update)
	rg.DELETE("/todos/:id", h.remove)
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

type listTodoResp struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	IsCompleted bool      `json:"isCompleted"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (h *todoHandler) list(c *gin.Context) {
	todos, err := h.srv.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := make([]listTodoResp, len(todos))
	for i, todo := range todos {
		resp[i] = listTodoResp{
			ID:          todo.ID,
			Title:       todo.Title,
			Description: todo.Description,
			IsCompleted: todo.IsCompleted,
			CreatedAt:   todo.CreatedAt,
			UpdatedAt:   todo.UpdatedAt,
		}
	}
	c.JSON(http.StatusOK, resp)
}

type getTodoReq struct {
	ID int `uri:"id" binding:"required"`
}

type getTodoResp struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	IsCompleted bool      `json:"isCompleted"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (h *todoHandler) get(c *gin.Context) {
	var req getTodoReq
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todo, err := h.srv.Get(req.ID)
	if errors.Is(err, service.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, getTodoResp{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
		IsCompleted: todo.IsCompleted,
		CreatedAt:   todo.CreatedAt,
		UpdatedAt:   todo.UpdatedAt,
	})
}

type updateTodoReq struct {
	ID          int     `uri:"id" binding:"required"`
	Title       *string `json:"title"`
	Description *string `json:"description"`
	IsCompleted *bool   `json:"isCompleted"`
}

func (h *todoHandler) update(c *gin.Context) {
	var req updateTodoReq
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := entity.UpdateTodoInput{
		Title:       req.Title,
		Description: req.Description,
		IsCompleted: req.IsCompleted,
	}
	if err := h.srv.Update(req.ID, input); errors.Is(err, service.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

type removeTodoReq struct {
	ID int `uri:"id" binding:"required"`
}

func (h *todoHandler) remove(c *gin.Context) {
	var req removeTodoReq
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.srv.Delete(req.ID); errors.Is(err, service.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
