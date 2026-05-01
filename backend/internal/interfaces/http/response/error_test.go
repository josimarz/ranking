package response_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/josimar/ranking/backend/internal/interfaces/http/response"
	"github.com/josimar/ranking/backend/pkg/apperror"
	"github.com/stretchr/testify/require"
)

func TestRespondError_WithAppError(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	appErr := apperror.NewValidationError("invalid name")
	response.RespondError(c, appErr)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var body map[string]map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)
	require.Equal(t, apperror.ValidationError, body["error"]["code"])
	require.Equal(t, "invalid name", body["error"]["message"])
	require.InDelta(t, float64(http.StatusBadRequest), body["error"]["status"], 0)
}

func TestRespondError_WithNotFoundError(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	appErr := apperror.NewNotFoundError(apperror.RankingNotFound, "ranking not found")
	response.RespondError(c, appErr)

	require.Equal(t, http.StatusNotFound, w.Code)

	var body map[string]map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)
	require.Equal(t, apperror.RankingNotFound, body["error"]["code"])
	require.Equal(t, "ranking not found", body["error"]["message"])
	require.InDelta(t, float64(http.StatusNotFound), body["error"]["status"], 0)
}

func TestRespondError_WithGenericError(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	genericErr := errors.New("something broke")
	response.RespondError(c, genericErr)

	// Should not write a response — delegates to error middleware via c.Error
	require.Equal(t, http.StatusOK, w.Code) // default status, not written by RespondError
	require.Len(t, c.Errors, 1)
	require.Equal(t, "something broke", c.Errors[0].Err.Error())
}
