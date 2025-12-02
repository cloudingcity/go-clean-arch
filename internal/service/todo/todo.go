package todo

import (
	"errors"

	"github.com/cloudingcity/todo/internal/entity"
	"github.com/cloudingcity/todo/internal/repo"
	"github.com/cloudingcity/todo/internal/service"
)

type Service struct {
	repo repo.Todo
}

func NewService(repo repo.Todo) service.Todo {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(title, description string) (*entity.Todo, error) {
	return s.repo.Create(title, description)
}

func (s *Service) List() ([]entity.Todo, error) {
	return s.repo.List()
}

func (s *Service) Get(id int) (*entity.Todo, error) {
	todo, err := s.repo.Get(id)
	if errors.Is(err, repo.ErrNotFound) {
		return nil, service.ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return todo, nil
}

func (s *Service) Update(id int, input entity.UpdateTodoInput) error {
	if err := s.repo.Update(id, input); errors.Is(err, repo.ErrNotFound) {
		return service.ErrNotFound
	} else if err != nil {
		return err
	}
	return nil
}

func (s *Service) Delete(id int) error {
	if err := s.repo.Delete(id); errors.Is(err, repo.ErrNotFound) {
		return service.ErrNotFound
	} else if err != nil {
		return err
	}
	return nil
}
