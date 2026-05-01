package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/josimar/ranking/backend/internal/interfaces/http/middleware"
	"github.com/stretchr/testify/require"
)

func TestCORS_PreflightOptions(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	r.Use(middleware.CORS("http://localhost:3000"))
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodOptions, "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
	require.Equal(t, "http://localhost:3000", w.Header().Get("Access-Control-Allow-Origin"))
	require.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "GET")
	require.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "POST")
	require.Contains(t, w.Header().Get("Access-Control-Allow-Headers"), "X-User-Id")
	require.Contains(t, w.Header().Get("Access-Control-Allow-Headers"), "X-Api-Key")
}

func TestCORS_RegularRequest(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	r.Use(middleware.CORS("http://localhost:3000"))
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "http://localhost:3000", w.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORS_WildcardOrigin(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	r.Use(middleware.CORS("*"))
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Origin", "http://any-origin.com")
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
}
