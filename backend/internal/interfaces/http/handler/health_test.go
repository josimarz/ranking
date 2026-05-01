package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"
	"github.com/josimar/ranking/backend/internal/interfaces/http/handler"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
}

type mockHealthChecker struct {
	err error
}

func (m *mockHealthChecker) DescribeTable(_ context.Context, _ *dynamodb.DescribeTableInput, _ ...func(*dynamodb.Options)) (*dynamodb.DescribeTableOutput, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &dynamodb.DescribeTableOutput{}, nil
}

type healthResponse struct {
	Status      string `json:"status"`
	Version     string `json:"version"`
	Environment string `json:"environment"`
	DynamoDB    string `json:"dynamodb"`
}

func TestHealthHandler_Check_Healthy(t *testing.T) {
	t.Setenv("ENV", "dev")

	h := handler.NewHealthHandler("test-table", &mockHealthChecker{})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)

	h.Check(c)

	require.Equal(t, http.StatusOK, w.Code)

	var body healthResponse
	err := json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)
	require.Equal(t, "healthy", body.Status)
	require.Equal(t, "1.0.0", body.Version)
	require.Equal(t, "dev", body.Environment)
	require.Equal(t, "connected", body.DynamoDB)
}

func TestHealthHandler_Check_Degraded(t *testing.T) {
	t.Setenv("ENV", "prod")

	h := handler.NewHealthHandler("test-table", &mockHealthChecker{
		err: errors.New("connection refused"),
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)

	h.Check(c)

	require.Equal(t, http.StatusOK, w.Code)

	var body healthResponse
	err := json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)
	require.Equal(t, "degraded", body.Status)
	require.Equal(t, "1.0.0", body.Version)
	require.Equal(t, "prod", body.Environment)
	require.Equal(t, "error: connection refused", body.DynamoDB)
}
