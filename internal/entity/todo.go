package entity

import "time"

type Todo struct {
	ID          int
	Title       string
	Description string
	IsCompleted bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type UpdateTodoInput struct {
	Title       *string
	Description *string
	IsCompleted *bool
}
