package repo

import "github.com/cloudingcity/todo/internal/entity"

type Todo interface {
	Create(title, description string) (entity.Todo, error)
}
