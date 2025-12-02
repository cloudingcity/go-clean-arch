package memory

import (
	"slices"
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

func (r *todoRepo) List() ([]entity.Todo, error) {
	return r.store, nil
}

func (r *todoRepo) Get(id int) (*entity.Todo, error) {
	idx := slices.IndexFunc(r.store, func(todo entity.Todo) bool {
		return todo.ID == id
	})

	if idx == -1 {
		return nil, repo.ErrNotFound
	}

	return &r.store[idx], nil
}

func (r *todoRepo) Update(id int, input entity.UpdateTodoInput) error {
	idx := slices.IndexFunc(r.store, func(todo entity.Todo) bool {
		return todo.ID == id
	})
	if idx == -1 {
		return repo.ErrNotFound
	}
	if input.Title != nil {
		r.store[idx].Title = *input.Title
	}
	if input.Description != nil {
		r.store[idx].Description = *input.Description
	}
	if input.IsCompleted != nil {
		r.store[idx].IsCompleted = *input.IsCompleted
	}
	r.store[idx].UpdatedAt = timeNow()
	return nil
}
