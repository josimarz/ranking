// Package handler provides HTTP handlers for the API.
package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"
)

// HealthChecker abstracts DynamoDB connectivity checks.
type HealthChecker interface {
	DescribeTable(ctx context.Context, params *dynamodb.DescribeTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DescribeTableOutput, error)
}

// HealthHandler handles health check requests.
type HealthHandler struct {
	tableName string
	checker   HealthChecker
}

// NewHealthHandler creates a new HealthHandler.
func NewHealthHandler(tableName string, checker HealthChecker) *HealthHandler {
	return &HealthHandler{
		tableName: tableName,
		checker:   checker,
	}
}

type healthResponse struct {
	Status      string `json:"status"`
	Version     string `json:"version"`
	Environment string `json:"environment"`
	DynamoDB    string `json:"dynamodb"`
}

// Check verifies DynamoDB connectivity and returns health status.
//
//	@Summary		Health check
//	@Description	Verifies DynamoDB connectivity and returns service health status
//	@Tags			health
//	@Produce		json
//	@Success		200	{object}	healthResponse
//	@Router			/health [get]
func (h *HealthHandler) Check(c *gin.Context) {
	resp := healthResponse{
		Status:      "healthy",
		Version:     "1.0.0",
		Environment: os.Getenv("ENV"),
		DynamoDB:    "connected",
	}

	_, err := h.checker.DescribeTable(c.Request.Context(), &dynamodb.DescribeTableInput{
		TableName: aws.String(h.tableName),
	})
	if err != nil {
		resp.Status = "degraded"
		resp.DynamoDB = fmt.Sprintf("error: %s", err.Error())
	}

	c.JSON(http.StatusOK, resp)
}
