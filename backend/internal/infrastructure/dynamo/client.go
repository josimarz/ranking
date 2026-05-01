// Package dynamo provides DynamoDB infrastructure utilities.
package dynamo

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// NewDynamoClient creates a DynamoDB client. If AWS_ENDPOINT_URL is set,
// the client uses that endpoint (for LocalStack).
func NewDynamoClient(ctx context.Context) (*dynamodb.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	var opts []func(*dynamodb.Options)
	if endpoint := os.Getenv("AWS_ENDPOINT_URL"); endpoint != "" {
		opts = append(opts, func(o *dynamodb.Options) {
			o.BaseEndpoint = aws.String(endpoint)
		})
	}

	return dynamodb.NewFromConfig(cfg, opts...), nil
}
