package v1

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/cloudingcity/todo/internal/entity"
	"github.com/cloudingcity/todo/internal/service/mocks"
	"github.com/gin-gonic/gin"
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
