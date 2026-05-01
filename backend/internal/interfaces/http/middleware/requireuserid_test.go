package middleware_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/josimar/ranking/backend/internal/interfaces/http/middleware"
	"github.com/stretchr/testify/require"
)

func TestRequireUserID_Present(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	r.Use(middleware.UserID())
	r.Use(middleware.RequireUserID())
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-User-Id", "019bc9c7-b857-72c5-a491-7f9686c7989a")
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}

func TestRequireUserID_Missing(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	r.Use(middleware.ErrorHandler())
	r.Use(middleware.UserID())
	r.Use(middleware.RequireUserID())
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)

	var body map[string]map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)
	require.Equal(t, "MISSING_USER_ID", body["error"]["code"])
}
