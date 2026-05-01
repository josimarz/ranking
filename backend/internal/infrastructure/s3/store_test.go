package s3_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/stretchr/testify/require"

	infraS3 "github.com/josimar/ranking/backend/internal/infrastructure/s3"
)

// mockS3API implements infraS3.S3API for testing.
type mockS3API struct {
	putObjectFn     func(ctx context.Context, input *s3.PutObjectInput, opts ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	deleteObjectFn  func(ctx context.Context, input *s3.DeleteObjectInput, opts ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
	deleteObjectsFn func(ctx context.Context, input *s3.DeleteObjectsInput, opts ...func(*s3.Options)) (*s3.DeleteObjectsOutput, error)
	listObjectsV2Fn func(ctx context.Context, input *s3.ListObjectsV2Input, opts ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
}

func (m *mockS3API) PutObject(ctx context.Context, input *s3.PutObjectInput, opts ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	return m.putObjectFn(ctx, input, opts...)
}

func (m *mockS3API) DeleteObject(ctx context.Context, input *s3.DeleteObjectInput, opts ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
	return m.deleteObjectFn(ctx, input, opts...)
}

func (m *mockS3API) DeleteObjects(ctx context.Context, input *s3.DeleteObjectsInput, opts ...func(*s3.Options)) (*s3.DeleteObjectsOutput, error) {
	return m.deleteObjectsFn(ctx, input, opts...)
}

func (m *mockS3API) ListObjectsV2(ctx context.Context, input *s3.ListObjectsV2Input, opts ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
	return m.listObjectsV2Fn(ctx, input, opts...)
}

// mockPresigner implements infraS3.Presigner for testing.
type mockPresigner struct {
	presignGetObjectFn func(ctx context.Context, input *s3.GetObjectInput, opts ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error)
}

func (m *mockPresigner) PresignGetObject(ctx context.Context, input *s3.GetObjectInput, opts ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error) {
	return m.presignGetObjectFn(ctx, input, opts...)
}

const (
	testBucket    = "test-bucket"
	testRankingID = "ranking-123"
	testItemID    = "item-456"
)

func TestS3ImageStore_Store(t *testing.T) {
	t.Parallel()

	t.Run("stores original and thumbnail", func(t *testing.T) {
		t.Parallel()

		var putCalls []string
		mock := &mockS3API{
			putObjectFn: func(_ context.Context, input *s3.PutObjectInput, _ ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
				putCalls = append(putCalls, *input.Key)
				require.Equal(t, testBucket, *input.Bucket)
				require.Equal(t, "image/webp", *input.ContentType)
				return &s3.PutObjectOutput{}, nil
			},
		}

		store := infraS3.NewS3ImageStore(mock, nil, testBucket)
		err := store.Store(context.Background(), testRankingID, testItemID, []byte("orig"), []byte("thumb"))

		require.NoError(t, err)
		require.Len(t, putCalls, 2)
		require.Equal(t, fmt.Sprintf("items/%s/%s/original.webp", testRankingID, testItemID), putCalls[0])
		require.Equal(t, fmt.Sprintf("items/%s/%s/thumbnail.webp", testRankingID, testItemID), putCalls[1])
	})

	t.Run("returns error on PutObject failure", func(t *testing.T) {
		t.Parallel()

		mock := &mockS3API{
			putObjectFn: func(_ context.Context, _ *s3.PutObjectInput, _ ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
				return nil, errors.New("s3 error")
			},
		}

		store := infraS3.NewS3ImageStore(mock, nil, testBucket)
		err := store.Store(context.Background(), testRankingID, testItemID, []byte("orig"), []byte("thumb"))

		require.Error(t, err)
		require.Contains(t, err.Error(), "s3 error")
	})
}

func TestS3ImageStore_GenerateURLs(t *testing.T) {
	t.Parallel()

	t.Run("generates presigned URLs", func(t *testing.T) {
		t.Parallel()

		var keysRequested []string
		presigner := &mockPresigner{
			presignGetObjectFn: func(_ context.Context, input *s3.GetObjectInput, _ ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error) {
				keysRequested = append(keysRequested, *input.Key)
				require.Equal(t, testBucket, *input.Bucket)
				return &v4.PresignedHTTPRequest{URL: fmt.Sprintf("https://s3/%s", *input.Key)}, nil
			},
		}

		store := infraS3.NewS3ImageStore(nil, presigner, testBucket)
		origURL, thumbURL, err := store.GenerateURLs(context.Background(), testRankingID, testItemID)

		require.NoError(t, err)
		require.Contains(t, origURL, "original.webp")
		require.Contains(t, thumbURL, "thumbnail.webp")
		require.Len(t, keysRequested, 2)
	})

	t.Run("returns error on presign failure", func(t *testing.T) {
		t.Parallel()

		presigner := &mockPresigner{
			presignGetObjectFn: func(_ context.Context, _ *s3.GetObjectInput, _ ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error) {
				return nil, errors.New("presign error")
			},
		}

		store := infraS3.NewS3ImageStore(nil, presigner, testBucket)
		origURL, thumbURL, err := store.GenerateURLs(context.Background(), testRankingID, testItemID)

		require.Error(t, err)
		require.Empty(t, origURL)
		require.Empty(t, thumbURL)
	})
}

func TestS3ImageStore_Delete(t *testing.T) {
	t.Parallel()

	t.Run("deletes both objects", func(t *testing.T) {
		t.Parallel()

		var deletedKeys []string
		mock := &mockS3API{
			deleteObjectFn: func(_ context.Context, input *s3.DeleteObjectInput, _ ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
				deletedKeys = append(deletedKeys, *input.Key)
				require.Equal(t, testBucket, *input.Bucket)
				return &s3.DeleteObjectOutput{}, nil
			},
		}

		store := infraS3.NewS3ImageStore(mock, nil, testBucket)
		err := store.Delete(context.Background(), testRankingID, testItemID)

		require.NoError(t, err)
		require.Len(t, deletedKeys, 2)
		require.Equal(t, fmt.Sprintf("items/%s/%s/original.webp", testRankingID, testItemID), deletedKeys[0])
		require.Equal(t, fmt.Sprintf("items/%s/%s/thumbnail.webp", testRankingID, testItemID), deletedKeys[1])
	})

	t.Run("returns error on delete failure", func(t *testing.T) {
		t.Parallel()

		mock := &mockS3API{
			deleteObjectFn: func(_ context.Context, _ *s3.DeleteObjectInput, _ ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
				return nil, errors.New("delete error")
			},
		}

		store := infraS3.NewS3ImageStore(mock, nil, testBucket)
		err := store.Delete(context.Background(), testRankingID, testItemID)

		require.Error(t, err)
	})
}

func TestS3ImageStore_DeleteByRanking(t *testing.T) {
	t.Parallel()

	t.Run("lists and deletes all objects for ranking", func(t *testing.T) {
		t.Parallel()

		mock := &mockS3API{
			listObjectsV2Fn: func(_ context.Context, input *s3.ListObjectsV2Input, _ ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
				require.Equal(t, testBucket, *input.Bucket)
				require.Equal(t, fmt.Sprintf("items/%s/", testRankingID), *input.Prefix)
				return &s3.ListObjectsV2Output{
					Contents: []s3types.Object{
						{Key: aws.String(fmt.Sprintf("items/%s/item1/original.webp", testRankingID))},
						{Key: aws.String(fmt.Sprintf("items/%s/item1/thumbnail.webp", testRankingID))},
					},
				}, nil
			},
			deleteObjectsFn: func(_ context.Context, input *s3.DeleteObjectsInput, _ ...func(*s3.Options)) (*s3.DeleteObjectsOutput, error) {
				require.Equal(t, testBucket, *input.Bucket)
				require.Len(t, input.Delete.Objects, 2)
				return &s3.DeleteObjectsOutput{}, nil
			},
		}

		store := infraS3.NewS3ImageStore(mock, nil, testBucket)
		err := store.DeleteByRanking(context.Background(), testRankingID)

		require.NoError(t, err)
	})

	t.Run("no-op when no objects found", func(t *testing.T) {
		t.Parallel()

		mock := &mockS3API{
			listObjectsV2Fn: func(_ context.Context, _ *s3.ListObjectsV2Input, _ ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
				return &s3.ListObjectsV2Output{Contents: []s3types.Object{}}, nil
			},
		}

		store := infraS3.NewS3ImageStore(mock, nil, testBucket)
		err := store.DeleteByRanking(context.Background(), testRankingID)

		require.NoError(t, err)
	})

	t.Run("returns error on list failure", func(t *testing.T) {
		t.Parallel()

		mock := &mockS3API{
			listObjectsV2Fn: func(_ context.Context, _ *s3.ListObjectsV2Input, _ ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
				return nil, errors.New("list error")
			},
		}

		store := infraS3.NewS3ImageStore(mock, nil, testBucket)
		err := store.DeleteByRanking(context.Background(), testRankingID)

		require.Error(t, err)
		require.Contains(t, err.Error(), "list error")
	})

	t.Run("returns error on batch delete failure", func(t *testing.T) {
		t.Parallel()

		mock := &mockS3API{
			listObjectsV2Fn: func(_ context.Context, _ *s3.ListObjectsV2Input, _ ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
				return &s3.ListObjectsV2Output{
					Contents: []s3types.Object{
						{Key: aws.String("items/r/i/original.webp")},
					},
				}, nil
			},
			deleteObjectsFn: func(_ context.Context, _ *s3.DeleteObjectsInput, _ ...func(*s3.Options)) (*s3.DeleteObjectsOutput, error) {
				return nil, errors.New("batch delete error")
			},
		}

		store := infraS3.NewS3ImageStore(mock, nil, testBucket)
		err := store.DeleteByRanking(context.Background(), testRankingID)

		require.Error(t, err)
		require.Contains(t, err.Error(), "batch delete error")
	})
}
