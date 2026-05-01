// Package s3 provides S3-based implementations for domain image storage.
package s3

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// NewS3Client creates an S3 client and returns the configured bucket name.
// If AWS_ENDPOINT_URL is set, it uses a custom endpoint with path-style access.
func NewS3Client(ctx context.Context) (*s3.Client, string, error) {
	bucket := os.Getenv("S3_BUCKET_IMAGES")
	if bucket == "" {
		return nil, "", errors.New("S3_BUCKET_IMAGES environment variable is required")
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("loading AWS config: %w", err)
	}

	var opts []func(*s3.Options)

	if endpoint := os.Getenv("AWS_ENDPOINT_URL"); endpoint != "" {
		opts = append(opts, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(endpoint)
			o.UsePathStyle = true
		})
	}

	client := s3.NewFromConfig(cfg, opts...)

	return client, bucket, nil
}
