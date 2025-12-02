package v1

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/cloudingcity/todo/internal/entity"
	"github.com/cloudingcity/todo/internal/service"
	"github.com/cloudingcity/todo/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type todoSuite struct {
	suite.Suite
	router  *gin.Engine
	mockSrv *mocks.MockTodo
}

func (s *todoSuite) SetupSubTest() {
	ctrl := gomock.NewController(s.T())
	s.mockSrv = mocks.NewMockTodo(ctrl)

	gin.SetMode(gin.TestMode)
	s.router = gin.Default()
	NewTodoRoutes(s.router.Group("v1"), s.mockSrv)
}

func TestTodoSuite(t *testing.T) {
	suite.Run(t, new(todoSuite))
}

func (s *todoSuite) TestCreate() {
	tests := []struct {
		desc     string
		body     string
		mock     func()
		wantCode int
		wantResp string
	}{
		{
			desc: "success",
			body: `{"title": "title-1", "description": "desc-1"}`,
			mock: func() {
				s.mockSrv.EXPECT().Create("title-1", "desc-1").Return(&entity.Todo{
					ID:          999,
					Title:       "title-1",
					Description: "desc-1",
					IsCompleted: false,
					CreatedAt:   time.Unix(123456789, 0),
					UpdatedAt:   time.Unix(123456789, 0),
				}, nil).Times(1)
			},
			wantCode: http.StatusCreated,
			wantResp: `{
			  "id": 999,
			  "title": "title-1",
			  "description": "desc-1",
			  "isCompleted": false,
			  "createdAt": "1973-11-30T05:33:09+08:00",
			  "updatedAt": "1973-11-30T05:33:09+08:00"
			}`,
		},
		{
			desc: "invalid request",
			body: `{"title": 999, "description": 0}`,
			mock: func() {
			},
			wantCode: http.StatusBadRequest,
			wantResp: `{
				"error": "json: cannot unmarshal number into Go struct field createTodoReq.title of type string"
			}`,
		},
		{
			desc: "service create failed",
			body: `{"title": "title-1", "description": "desc-1"}`,
			mock: func() {
				s.mockSrv.EXPECT().Create("title-1", "desc-1").Return(nil, errors.New("something wrong")).Times(1)
			},
			wantCode: http.StatusInternalServerError,
			wantResp: `{
				"error": "something wrong"
			}`,
		},
	}

	for _, tt := range tests {
		s.Run(tt.desc, func() {
			if tt.mock != nil {
				tt.mock()
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/v1/todos", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			s.router.ServeHTTP(w, req)

			s.Equal(tt.wantCode, w.Code)
			s.JSONEq(tt.wantResp, w.Body.String())
		})
	}
}

func (s *todoSuite) TestList() {
	tests := []struct {
		desc     string
		mock     func()
		wantCode int
		wantResp string
	}{
		{
			desc: "success",
			mock: func() {
				s.mockSrv.EXPECT().List().Return([]entity.Todo{
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
						IsCompleted: true,
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
			wantCode: http.StatusOK,
			wantResp: `[
				{
				  "id": 1,
				  "title": "title-1",
				  "description": "desc-1",
				  "isCompleted": false,
				  "createdAt": "1973-11-30T05:33:09+08:00",
				  "updatedAt": "1973-11-30T05:33:09+08:00"
				},
				{
				  "id": 2,
				  "title": "title-2",
				  "description": "desc-2",
				  "isCompleted": true,
				  "createdAt": "1973-11-30T05:33:09+08:00",
				  "updatedAt": "1973-11-30T05:33:09+08:00"
				},
				{
				  "id": 3,
				  "title": "title-3",
				  "description": "desc-3",
				  "isCompleted": false,
				  "createdAt": "1973-11-30T05:33:09+08:00",
				  "updatedAt": "1973-11-30T05:33:09+08:00"
				}
            ]`,
		},
		{
			desc: "service list failed",
			mock: func() {
				s.mockSrv.EXPECT().List().Return(nil, errors.New("something wrong")).Times(1)
			},
			wantCode: http.StatusInternalServerError,
			wantResp: `{
				"error": "something wrong"
			}`,
		},
	}

	for _, tt := range tests {
		s.Run(tt.desc, func() {
			if tt.mock != nil {
				tt.mock()
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/v1/todos", nil)
			s.router.ServeHTTP(w, req)

			s.Equal(tt.wantCode, w.Code)
			s.JSONEq(tt.wantResp, w.Body.String())
		})
	}
}

func (s *todoSuite) TestGet() {
	tests := []struct {
		desc     string
		id       string
		mock     func()
		wantCode int
		wantResp string
	}{
		{
			desc: "success",
			id:   "1",
			mock: func() {
				s.mockSrv.EXPECT().Get(1).Return(&entity.Todo{
					ID:          1,
					Title:       "title-1",
					Description: "desc-1",
					IsCompleted: false,
					CreatedAt:   time.Unix(123456789, 0),
					UpdatedAt:   time.Unix(123456789, 0),
				}, nil).Times(1)
			},
			wantCode: http.StatusOK,
			wantResp: `{
			  "id": 1,
			  "title": "title-1",
			  "description": "desc-1",
			  "isCompleted": false,
			  "createdAt": "1973-11-30T05:33:09+08:00",
			  "updatedAt": "1973-11-30T05:33:09+08:00"
            }`,
		},
		{
			desc:     "wrong id",
			id:       "wrong-id",
			mock:     func() {},
			wantCode: http.StatusBadRequest,
			wantResp: `{
				"error": "strconv.ParseInt: parsing \"wrong-id\": invalid syntax"
			}`,
		},
		{
			desc: "not found",
			id:   "1",
			mock: func() {
				s.mockSrv.EXPECT().Get(1).Return(nil, service.ErrNotFound).Times(1)
			},
			wantCode: http.StatusNotFound,
			wantResp: `{
				"error": "not found"
			}`,
		},
		{
			desc: "service get failed",
			id:   "1",
			mock: func() {
				s.mockSrv.EXPECT().Get(1).Return(nil, errors.New("something wrong")).Times(1)
			},
			wantCode: http.StatusInternalServerError,
			wantResp: `{
				"error": "something wrong"
			}`,
		},
	}

	for _, tt := range tests {
		s.Run(tt.desc, func() {
			if tt.mock != nil {
				tt.mock()
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/todos/%s", tt.id), nil)
			s.router.ServeHTTP(w, req)

			s.Equal(tt.wantCode, w.Code)
			s.JSONEq(tt.wantResp, w.Body.String())
		})
	}
}

func (s *todoSuite) TestUpdate() {
	tests := []struct {
		desc     string
		id       string
		body     string
		mock     func()
		wantCode int
		wantResp string
	}{
		{
			desc: "success",
			id:   "1",
			body: `{"title": "title-1", "description": "desc-1", "isCompleted": true}`,
			mock: func() {
				s.mockSrv.EXPECT().Update(1, entity.UpdateTodoInput{
					Title:       lo.ToPtr("title-1"),
					Description: lo.ToPtr("desc-1"),
					IsCompleted: lo.ToPtr(true),
				}).Return(nil).Times(1)
			},
			wantCode: http.StatusNoContent,
			wantResp: "",
		},
		{
			desc: "wrong id",
			id:   "wrong-id",
			body: `{"title": "title-1", "description": "desc-1", "isCompleted": true}`,
			mock: func() {
			},
			wantCode: http.StatusBadRequest,
			wantResp: `{
				"error": "strconv.ParseInt: parsing \"wrong-id\": invalid syntax"
			}`,
		},
		{
			desc: "invalid body",
			id:   "1",
			body: `{"title": 999, "description": "desc-1", "isCompleted": true}`,
			mock: func() {
			},
			wantCode: http.StatusBadRequest,
			wantResp: `{
				"error": "json: cannot unmarshal number into Go struct field updateTodoReq.title of type string"
			}`,
		},
		{
			desc: "not found",
			id:   "1",
			body: `{"title": "title-1", "description": "desc-1", "isCompleted": true}`,
			mock: func() {
				s.mockSrv.EXPECT().Update(1, entity.UpdateTodoInput{
					Title:       lo.ToPtr("title-1"),
					Description: lo.ToPtr("desc-1"),
					IsCompleted: lo.ToPtr(true),
				}).Return(service.ErrNotFound).Times(1)
			},
			wantCode: http.StatusNotFound,
			wantResp: `{
				"error": "not found"
			}`,
		},
		{
			desc: "service update failed",
			id:   "1",
			body: `{"title": "title-1", "description": "desc-1", "isCompleted": true}`,
			mock: func() {
				s.mockSrv.EXPECT().Update(1, entity.UpdateTodoInput{
					Title:       lo.ToPtr("title-1"),
					Description: lo.ToPtr("desc-1"),
					IsCompleted: lo.ToPtr(true),
				}).Return(errors.New("something wrong")).Times(1)
			},
			wantCode: http.StatusInternalServerError,
			wantResp: `{
				"error": "something wrong"
			}`,
		},
	}

	for _, tt := range tests {
		s.Run(tt.desc, func() {
			if tt.mock != nil {
				tt.mock()
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/v1/todos/%s", tt.id), strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			s.router.ServeHTTP(w, req)

			s.Equal(tt.wantCode, w.Code)
			if tt.wantResp != "" {
				s.JSONEq(tt.wantResp, w.Body.String())
			}
		})
	}
}

func (s *todoSuite) TestDelete() {
	tests := []struct {
		desc     string
		id       string
		mock     func()
		wantCode int
		wantResp string
	}{
		{
			desc: "success",
			id:   "1",
			mock: func() {
				s.mockSrv.EXPECT().Delete(1).Return(nil).Times(1)
			},
			wantCode: http.StatusNoContent,
			wantResp: "",
		},
		{
			desc:     "wrong id",
			id:       "wrong-id",
			mock:     func() {},
			wantCode: http.StatusBadRequest,
			wantResp: `{
				"error": "strconv.ParseInt: parsing \"wrong-id\": invalid syntax"
			}`,
		},
		{
			desc: "not found",
			id:   "1",
			mock: func() {
				s.mockSrv.EXPECT().Delete(1).Return(service.ErrNotFound).Times(1)
			},
			wantCode: http.StatusNotFound,
			wantResp: `{
				"error": "not found"
			}`,
		},
		{
			desc: "service delete failed",
			id:   "1",
			mock: func() {
				s.mockSrv.EXPECT().Delete(1).Return(errors.New("something wrong")).Times(1)
			},
			wantCode: http.StatusInternalServerError,
			wantResp: `{
				"error": "something wrong"
			}`,
		},
	}

	for _, tt := range tests {
		s.Run(tt.desc, func() {
			if tt.mock != nil {
				tt.mock()
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/todos/%s", tt.id), nil)
			s.router.ServeHTTP(w, req)

			s.Equal(tt.wantCode, w.Code)
		})
	}
}
