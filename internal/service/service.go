package service

import (
	"errors"

	"github.com/cloudingcity/todo/internal/entity"
)

var (
	ErrNotFound = errors.New("not found")
)

//go:generate mockgen -source=service.go -destination mocks/service.go -package mocks
type Todo interface {
	Create(title, description string) (*entity.Todo, error)
	List() ([]entity.Todo, error)
	Get(id int) (*entity.Todo, error)
	Update(id int, input entity.UpdateTodoInput) error
	Delete(id int) error
}
