package middleware_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/josimar/ranking/backend/internal/interfaces/http/middleware"
	"github.com/stretchr/testify/require"
)

func TestLogging_JSONOutput(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	r.Use(middleware.Logging(&buf, "prod", "ranking"))
	r.Use(middleware.UserID())
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-User-Id", "019bc9c7-b857-72c5-a491-7f9686c7989a")
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var logEntry map[string]any
	err := json.Unmarshal(buf.Bytes(), &logEntry)
	require.NoError(t, err)
	require.Equal(t, "GET", logEntry["method"])
	require.Equal(t, "/test", logEntry["path"])
	require.Equal(t, float64(http.StatusOK), logEntry["statusCode"])
	require.Equal(t, "019bc9c7-b857-72c5-a491-7f9686c7989a", logEntry["userId"])
	require.Equal(t, "prod", logEntry["env"])
	require.Equal(t, "ranking", logEntry["service"])
	require.NotEmpty(t, logEntry["requestId"])
	require.NotEmpty(t, logEntry["latency"])
}

func TestLogging_TextOutput(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	r.Use(middleware.Logging(&buf, "local", "ranking"))
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, buf.String(), "GET")
	require.Contains(t, buf.String(), "/test")
}
