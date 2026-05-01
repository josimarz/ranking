package middleware_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/josimar/ranking/backend/internal/interfaces/http/middleware"
	"github.com/josimar/ranking/backend/pkg/apperror"
	"github.com/stretchr/testify/require"
)

func TestErrorHandler_AppError(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	r.Use(middleware.ErrorHandler())
	r.GET("/test", func(c *gin.Context) {
		_ = c.Error(apperror.NewValidationError("bad input"))
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var body map[string]map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)
	require.Equal(t, "VALIDATION_ERROR", body["error"]["code"])
	require.Equal(t, "bad input", body["error"]["message"])
	require.Equal(t, float64(http.StatusBadRequest), body["error"]["status"])
}

func TestErrorHandler_UnknownError(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	r.Use(middleware.ErrorHandler())
	r.GET("/test", func(c *gin.Context) {
		_ = c.Error(errors.New("something unexpected"))
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)

	var body map[string]map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)
	require.Equal(t, "INTERNAL_ERROR", body["error"]["code"])
}

func TestErrorHandler_NoError(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	r.Use(middleware.ErrorHandler())
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}
