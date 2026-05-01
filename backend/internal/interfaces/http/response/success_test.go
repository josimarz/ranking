package response_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/josimar/ranking/backend/internal/interfaces/http/response"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestRespondOK(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	data := map[string]string{"name": "test"}
	response.RespondOK(c, data)

	require.Equal(t, http.StatusOK, w.Code)

	var body map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)
	require.Equal(t, "test", body["name"])
}

func TestRespondCreated(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	data := map[string]string{"id": "123"}
	response.RespondCreated(c, data)

	require.Equal(t, http.StatusCreated, w.Code)

	var body map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)
	require.Equal(t, "123", body["id"])
}

func TestRespondNoContent(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	response.RespondNoContent(c)

	require.Equal(t, http.StatusNoContent, c.Writer.Status())
	require.Empty(t, w.Body.Bytes())
}
