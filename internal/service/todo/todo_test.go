package todo

import (
	"errors"
	"testing"
	"time"

	"github.com/cloudingcity/todo/internal/entity"
	"github.com/cloudingcity/todo/internal/repo"
	"github.com/cloudingcity/todo/internal/repo/mocks"
	"github.com/cloudingcity/todo/internal/service"
	"github.com/samber/lo"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

var (
	mockErr = errors.New("something wrong")
)

type todoSuite struct {
	suite.Suite
	srv      service.Todo
	mockRepo *mocks.MockTodo
}

func (s *todoSuite) SetupSubTest() {
	ctrl := gomock.NewController(s.T())
	s.mockRepo = mocks.NewMockTodo(ctrl)
	s.srv = NewService(s.mockRepo)
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
		want        *entity.Todo
		wantErr     error
	}{
		{
			desc:        "success",
			title:       "title-1",
			description: "desc-1",
			setup: func() {
				s.mockRepo.EXPECT().Create("title-1", "desc-1").Return(&entity.Todo{
					ID:          1,
					Title:       "title-1",
					Description: "desc-1",
					IsCompleted: false,
					CreatedAt:   time.Unix(123456789, 0),
					UpdatedAt:   time.Unix(123456789, 0),
				}, nil).Times(1)
			},
			want: &entity.Todo{
				ID:          1,
				Title:       "title-1",
				Description: "desc-1",
				IsCompleted: false,
				CreatedAt:   time.Unix(123456789, 0),
				UpdatedAt:   time.Unix(123456789, 0),
			},
			wantErr: nil,
		},
		{
			desc:        "create failed",
			title:       "title-1",
			description: "desc-1",
			setup: func() {
				s.mockRepo.EXPECT().Create("title-1", "desc-1").Return(nil, mockErr).Times(1)
			},
			want:    nil,
			wantErr: mockErr,
		},
	}
	for _, tt := range tests {
		s.Run(tt.desc, func() {
			tt.setup()
			got, err := s.srv.Create(tt.title, tt.description)
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
				s.mockRepo.EXPECT().List().Return([]entity.Todo{
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
				}, nil).Times(1)
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
		{
			desc: "list failed",
			setup: func() {
				s.mockRepo.EXPECT().List().Return(nil, mockErr).Times(1)
			},
			want:    nil,
			wantErr: mockErr,
		},
	}
	for _, tt := range tests {
		s.Run(tt.desc, func() {
			tt.setup()
			got, err := s.srv.List()
			s.Equal(tt.want, got)
			s.ErrorIs(tt.wantErr, err)
		})
	}
}

func (s *todoSuite) TestGet() {
	tests := []struct {
		desc    string
		id      int
		setup   func()
		want    *entity.Todo
		wantErr error
	}{
		{
			desc: "success",
			id:   1,
			setup: func() {
				s.mockRepo.EXPECT().Get(1).Return(&entity.Todo{
					ID:          1,
					Title:       "title-1",
					Description: "desc-1",
					IsCompleted: false,
					CreatedAt:   time.Unix(123456789, 0),
					UpdatedAt:   time.Unix(123456789, 0),
				}, nil).Times(1)
			},
			want: &entity.Todo{
				ID:          1,
				Title:       "title-1",
				Description: "desc-1",
				IsCompleted: false,
				CreatedAt:   time.Unix(123456789, 0),
				UpdatedAt:   time.Unix(123456789, 0),
			},
			wantErr: nil,
		},
		{
			desc: "not found",
			id:   1,
			setup: func() {
				s.mockRepo.EXPECT().Get(1).Return(nil, repo.ErrNotFound).Times(1)
			},
			want:    nil,
			wantErr: service.ErrNotFound,
		},
	}
	for _, tt := range tests {
		s.Run(tt.desc, func() {
			tt.setup()
			got, err := s.srv.Get(tt.id)
			s.Equal(tt.want, got)
			s.ErrorIs(tt.wantErr, err)
		})
	}
}

func (s *todoSuite) TestUpdate() {
	tests := []struct {
		desc    string
		id      int
		input   entity.UpdateTodoInput
		setup   func()
		want    *entity.Todo
		wantErr error
	}{
		{
			desc: "success",
			id:   1,
			input: entity.UpdateTodoInput{
				Title:       lo.ToPtr("title-update"),
				Description: lo.ToPtr("desc-update"),
				IsCompleted: lo.ToPtr(true),
			},
			setup: func() {
				s.mockRepo.EXPECT().Update(1, entity.UpdateTodoInput{
					Title:       lo.ToPtr("title-update"),
					Description: lo.ToPtr("desc-update"),
					IsCompleted: lo.ToPtr(true),
				}).Return(nil).Times(1)
			},
			wantErr: nil,
		},
		{
			desc: "not found",
			id:   1,
			input: entity.UpdateTodoInput{
				Title:       lo.ToPtr("title-update"),
				Description: lo.ToPtr("desc-update"),
				IsCompleted: lo.ToPtr(true),
			},
			setup: func() {
				s.mockRepo.EXPECT().Update(1, entity.UpdateTodoInput{
					Title:       lo.ToPtr("title-update"),
					Description: lo.ToPtr("desc-update"),
					IsCompleted: lo.ToPtr(true),
				}).Return(repo.ErrNotFound).Times(1)
			},
			wantErr: service.ErrNotFound,
		},
	}
	for _, tt := range tests {
		s.Run(tt.desc, func() {
			tt.setup()
			err := s.srv.Update(tt.id, tt.input)
			s.ErrorIs(tt.wantErr, err)
		})
	}
}

func (s *todoSuite) TestDelete() {
	tests := []struct {
		desc    string
		id      int
		setup   func()
		wantErr error
	}{
		{
			desc: "success",
			id:   1,
			setup: func() {
				s.mockRepo.EXPECT().Delete(1).Return(nil).Times(1)
			},
			wantErr: nil,
		},
		{
			desc: "not found",
			id:   1,
			setup: func() {
				s.mockRepo.EXPECT().Delete(1).Return(repo.ErrNotFound).Times(1)
			},
			wantErr: service.ErrNotFound,
		},
	}
	for _, tt := range tests {
		s.Run(tt.desc, func() {
			tt.setup()
			err := s.srv.Delete(tt.id)
			s.ErrorIs(tt.wantErr, err)
		})
	}
}
