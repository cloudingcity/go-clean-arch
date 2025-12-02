package memory

import (
	"time"

	"github.com/cloudingcity/todo/internal/entity"
	"github.com/cloudingcity/todo/internal/repo"
)

var timeNow = time.Now

type todoRepo struct {
	idCounter int
	store     []entity.Todo
}

func NewTodoRepo() repo.Todo {
	return &todoRepo{
		idCounter: 1,
	}
}

func (r *todoRepo) Create(title, description string) (entity.Todo, error) {
	now := timeNow()
	todo := entity.Todo{
		ID:          r.idCounter,
		Title:       title,
		Description: description,
		IsCompleted: false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	r.store = append(r.store, todo)
	r.idCounter++
	return todo, nil
}
