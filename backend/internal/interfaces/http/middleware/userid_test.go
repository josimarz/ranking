package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/josimar/ranking/backend/internal/interfaces/http/middleware"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestUserID_ValidUUID(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	var captured string
	r.Use(middleware.UserID())
	r.GET("/test", func(c *gin.Context) {
		captured = middleware.GetUserID(c)
		c.Status(http.StatusOK)
	})

	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
	c.Request.Header.Set("X-User-Id", "019bc9c7-b857-72c5-a491-7f9686c7989a")
	r.ServeHTTP(w, c.Request)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "019bc9c7-b857-72c5-a491-7f9686c7989a", captured)
}

func TestUserID_MissingHeader(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	var captured string
	r.Use(middleware.UserID())
	r.GET("/test", func(c *gin.Context) {
		captured = middleware.GetUserID(c)
		c.Status(http.StatusOK)
	})

	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, c.Request)

	require.Equal(t, http.StatusOK, w.Code)
	require.Empty(t, captured)
}

func TestUserID_InvalidUUID(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	var captured string
	r.Use(middleware.UserID())
	r.GET("/test", func(c *gin.Context) {
		captured = middleware.GetUserID(c)
		c.Status(http.StatusOK)
	})

	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
	c.Request.Header.Set("X-User-Id", "not-a-uuid")
	r.ServeHTTP(w, c.Request)

	require.Equal(t, http.StatusOK, w.Code)
	require.Empty(t, captured)
}
