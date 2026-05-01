package s3_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	infraS3 "github.com/josimar/ranking/backend/internal/infrastructure/s3"
)

func TestNewS3Client_ReturnsBucketName(t *testing.T) {
	t.Setenv("S3_BUCKET_IMAGES", "test-bucket")
	t.Setenv("AWS_REGION", "us-east-1")
	t.Setenv("AWS_ACCESS_KEY_ID", "test")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "test")

	client, bucket, err := infraS3.NewS3Client(context.Background())

	require.NoError(t, err)
	require.NotNil(t, client)
	require.Equal(t, "test-bucket", bucket)
}

func TestNewS3Client_MissingBucket_ReturnsError(t *testing.T) {
	t.Setenv("S3_BUCKET_IMAGES", "")
	t.Setenv("AWS_REGION", "us-east-1")
	t.Setenv("AWS_ACCESS_KEY_ID", "test")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "test")

	client, bucket, err := infraS3.NewS3Client(context.Background())

	require.Error(t, err)
	require.Nil(t, client)
	require.Empty(t, bucket)
	require.Contains(t, err.Error(), "S3_BUCKET_IMAGES")
}

func TestNewS3Client_WithEndpointURL(t *testing.T) {
	t.Setenv("S3_BUCKET_IMAGES", "test-bucket")
	t.Setenv("AWS_ENDPOINT_URL", "http://localhost:4566")
	t.Setenv("AWS_REGION", "us-east-1")
	t.Setenv("AWS_ACCESS_KEY_ID", "test")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "test")

	client, bucket, err := infraS3.NewS3Client(context.Background())

	require.NoError(t, err)
	require.NotNil(t, client)
	require.Equal(t, "test-bucket", bucket)
}
