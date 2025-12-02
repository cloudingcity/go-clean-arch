package memory

import (
	"testing"
	"time"

	"github.com/cloudingcity/todo/internal/entity"
	"github.com/cloudingcity/todo/internal/repo"
	"github.com/stretchr/testify/suite"
)

type todoSuite struct {
	suite.Suite
	repo repo.Todo
}

func (s *todoSuite) SetupSubTest() {
	s.repo = NewTodoRepo()
}

func (s *todoSuite) TearDownSubTest() {
}

func TestTodoSuite(t *testing.T) {
	suite.Run(t, new(todoSuite))
}

func (s *todoSuite) TestCreate() {
	tests := []struct {
		desc        string
		title       string
		description string
		setup       func()
		want        entity.Todo
		wantErr     error
	}{
		{
			desc:        "success",
			title:       "title-1",
			description: "desc-1",
			setup: func() {
				timeNow = func() time.Time {
					return time.Unix(123456789, 0)
				}
			},
			want: entity.Todo{
				ID:          1,
				Title:       "title-1",
				Description: "desc-1",
				IsCompleted: false,
				CreatedAt:   time.Unix(123456789, 0),
				UpdatedAt:   time.Unix(123456789, 0),
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		s.Run(tt.desc, func() {
			tt.setup()
			got, err := s.repo.Create(tt.title, tt.description)
			s.Equal(tt.want, got)
			s.ErrorIs(tt.wantErr, err)
		})
	}
}

func (s *todoSuite) TestList() {
	tests := []struct {
		desc    string
		setup   func()
		want    []entity.Todo
		wantErr error
	}{
		{
			desc: "success",
			setup: func() {
				timeNow = func() time.Time {
					return time.Unix(123456789, 0)
				}
				_, _ = s.repo.Create("title-1", "desc-1")
				_, _ = s.repo.Create("title-2", "desc-2")
				_, _ = s.repo.Create("title-3", "desc-3")
			},
			want: []entity.Todo{
				{
					ID:          1,
					Title:       "title-1",
					Description: "desc-1",
					IsCompleted: false,
					CreatedAt:   time.Unix(123456789, 0),
					UpdatedAt:   time.Unix(123456789, 0),
				},
				{
					ID:          2,
					Title:       "title-2",
					Description: "desc-2",
					IsCompleted: false,
					CreatedAt:   time.Unix(123456789, 0),
					UpdatedAt:   time.Unix(123456789, 0),
				},
				{
					ID:          3,
					Title:       "title-3",
					Description: "desc-3",
					IsCompleted: false,
					CreatedAt:   time.Unix(123456789, 0),
					UpdatedAt:   time.Unix(123456789, 0),
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		s.Run(tt.desc, func() {
			tt.setup()
			got, err := s.repo.List()
			s.Equal(tt.want, got)
			s.ErrorIs(tt.wantErr, err)
		})
	}
}
