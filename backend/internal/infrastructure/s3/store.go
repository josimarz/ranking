package s3

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const presignExpiry = 1 * time.Hour

// API defines the S3 operations used by the image store.
type API interface {
	PutObject(ctx context.Context, input *s3.PutObjectInput, opts ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	DeleteObject(ctx context.Context, input *s3.DeleteObjectInput, opts ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
	DeleteObjects(ctx context.Context, input *s3.DeleteObjectsInput, opts ...func(*s3.Options)) (*s3.DeleteObjectsOutput, error)
	ListObjectsV2(ctx context.Context, input *s3.ListObjectsV2Input, opts ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
}

// Presigner defines the presigning operations used by the image store.
type Presigner interface {
	PresignGetObject(ctx context.Context, input *s3.GetObjectInput, opts ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error)
}

// ImageStore implements the domain ImageStore interface using S3.
type ImageStore struct {
	client    API
	presigner Presigner
	bucket    string
}

// NewS3ImageStore creates a new ImageStore.
func NewS3ImageStore(client API, presigner Presigner, bucket string) *ImageStore {
	return &ImageStore{
		client:    client,
		presigner: presigner,
		bucket:    bucket,
	}
}

// Store uploads original and thumbnail images to S3.
func (s *ImageStore) Store(ctx context.Context, rankingID, itemID string, original, thumbnail []byte) error {
	originalKey := objectKey(rankingID, itemID, "original.webp")
	thumbnailKey := objectKey(rankingID, itemID, "thumbnail.webp")

	if err := s.putObject(ctx, originalKey, original); err != nil {
		return fmt.Errorf("storing original image: %w", err)
	}

	if err := s.putObject(ctx, thumbnailKey, thumbnail); err != nil {
		return fmt.Errorf("storing thumbnail image: %w", err)
	}

	return nil
}

// GenerateURLs creates pre-signed GetObject URLs with 1-hour expiry.
func (s *ImageStore) GenerateURLs(ctx context.Context, rankingID, itemID string) (string, string, error) {
	originalKey := objectKey(rankingID, itemID, "original.webp")
	thumbnailKey := objectKey(rankingID, itemID, "thumbnail.webp")

	origReq, err := s.presigner.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(originalKey),
	}, func(o *s3.PresignOptions) {
		o.Expires = presignExpiry
	})
	if err != nil {
		return "", "", fmt.Errorf("presigning original URL: %w", err)
	}

	thumbReq, err := s.presigner.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(thumbnailKey),
	}, func(o *s3.PresignOptions) {
		o.Expires = presignExpiry
	})
	if err != nil {
		return "", "", fmt.Errorf("presigning thumbnail URL: %w", err)
	}

	return origReq.URL, thumbReq.URL, nil
}

// Delete removes both original and thumbnail images from S3.
func (s *ImageStore) Delete(ctx context.Context, rankingID, itemID string) error {
	originalKey := objectKey(rankingID, itemID, "original.webp")
	thumbnailKey := objectKey(rankingID, itemID, "thumbnail.webp")

	if _, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(originalKey),
	}); err != nil {
		return fmt.Errorf("deleting original image: %w", err)
	}

	if _, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(thumbnailKey),
	}); err != nil {
		return fmt.Errorf("deleting thumbnail image: %w", err)
	}

	return nil
}

// DeleteByRanking removes all images for a ranking using ListObjectsV2 + DeleteObjects.
func (s *ImageStore) DeleteByRanking(ctx context.Context, rankingID string) error {
	prefix := fmt.Sprintf("items/%s/", rankingID)

	listOutput, err := s.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		return fmt.Errorf("listing objects for ranking: %w", err)
	}

	if len(listOutput.Contents) == 0 {
		return nil
	}

	objects := make([]types.ObjectIdentifier, len(listOutput.Contents))
	for i := range listOutput.Contents {
		objects[i] = types.ObjectIdentifier{Key: listOutput.Contents[i].Key}
	}

	if _, err := s.client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
		Bucket: aws.String(s.bucket),
		Delete: &types.Delete{
			Objects: objects,
			Quiet:   aws.Bool(true),
		},
	}); err != nil {
		return fmt.Errorf("batch deleting objects: %w", err)
	}

	return nil
}

func (s *ImageStore) putObject(ctx context.Context, key string, data []byte) error {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String("image/webp"),
	})

	return err
}

func objectKey(rankingID, itemID, filename string) string {
	return fmt.Sprintf("items/%s/%s/%s", rankingID, itemID, filename)
}
