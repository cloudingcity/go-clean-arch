package v1

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPingRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	NewPingRoutes(r.Group("v1"))

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/ping", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}
