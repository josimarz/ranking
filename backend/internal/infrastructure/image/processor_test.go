package image

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"testing"

	"github.com/josimar/ranking/backend/pkg/apperror"
	"github.com/stretchr/testify/require"
)

func createTestImage(width, height int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := range height {
		for x := range width {
			img.Set(x, y, color.RGBA{R: uint8(x % 256), G: uint8(y % 256), B: 128, A: 255})
		}
	}

	return img
}

func encodeJPEG(t *testing.T, img image.Image) []byte {
	t.Helper()

	var buf bytes.Buffer
	require.NoError(t, jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90}))

	return buf.Bytes()
}

func encodePNG(t *testing.T, img image.Image) []byte {
	t.Helper()

	var buf bytes.Buffer
	require.NoError(t, png.Encode(&buf, img))

	return buf.Bytes()
}

func TestValidateFormat_JPEG(t *testing.T) {
	t.Parallel()

	p := NewWebPProcessor()
	data := encodeJPEG(t, createTestImage(10, 10))

	err := p.ValidateFormat(data)

	require.NoError(t, err)
}

func TestValidateFormat_PNG(t *testing.T) {
	t.Parallel()

	p := NewWebPProcessor()
	data := encodePNG(t, createTestImage(10, 10))

	err := p.ValidateFormat(data)

	require.NoError(t, err)
}

func TestValidateFormat_WebP(t *testing.T) {
	t.Parallel()

	p := NewWebPProcessor()
	// Minimal WebP header: RIFF....WEBP
	data := []byte("RIFF\x00\x00\x00\x00WEBP")

	err := p.ValidateFormat(data)

	require.NoError(t, err)
}

func TestValidateFormat_InvalidFormat(t *testing.T) {
	t.Parallel()

	p := NewWebPProcessor()
	data := []byte("not an image at all")

	err := p.ValidateFormat(data)

	require.Error(t, err)

	var appErr *apperror.AppError
	require.ErrorAs(t, err, &appErr)
	require.Equal(t, apperror.InvalidImageFormat, appErr.Code)
}

func TestValidateFormat_EmptyData(t *testing.T) {
	t.Parallel()

	p := NewWebPProcessor()

	err := p.ValidateFormat([]byte{})

	require.Error(t, err)

	var appErr *apperror.AppError
	require.ErrorAs(t, err, &appErr)
	require.Equal(t, apperror.InvalidImageFormat, appErr.Code)
}

func TestValidateFormat_TooLarge(t *testing.T) {
	t.Parallel()

	p := NewWebPProcessor()
	// Create data just over 10 MB with valid JPEG header
	data := make([]byte, 10*1024*1024+1)
	data[0] = 0xFF
	data[1] = 0xD8
	data[2] = 0xFF

	err := p.ValidateFormat(data)

	require.Error(t, err)

	var appErr *apperror.AppError
	require.ErrorAs(t, err, &appErr)
	require.Equal(t, apperror.ImageTooLarge, appErr.Code)
}

func TestValidateFormat_ExactlyMaxSize(t *testing.T) {
	t.Parallel()

	p := NewWebPProcessor()
	// Exactly 10 MB with valid JPEG header
	data := make([]byte, 10*1024*1024)
	data[0] = 0xFF
	data[1] = 0xD8
	data[2] = 0xFF

	err := p.ValidateFormat(data)

	require.NoError(t, err)
}

func TestProcess_JPEG(t *testing.T) {
	t.Parallel()

	p := NewWebPProcessor()
	data := encodeJPEG(t, createTestImage(100, 80))

	original, thumbnail, err := p.Process(data)

	require.NoError(t, err)
	require.NotEmpty(t, original)
	require.NotEmpty(t, thumbnail)
}

func TestProcess_PNG(t *testing.T) {
	t.Parallel()

	p := NewWebPProcessor()
	data := encodePNG(t, createTestImage(100, 80))

	original, thumbnail, err := p.Process(data)

	require.NoError(t, err)
	require.NotEmpty(t, original)
	require.NotEmpty(t, thumbnail)
}

func TestProcess_ResizesLargeImage(t *testing.T) {
	t.Parallel()

	p := NewWebPProcessor()
	data := encodeJPEG(t, createTestImage(1600, 1200))

	original, thumbnail, err := p.Process(data)

	require.NoError(t, err)
	require.NotEmpty(t, original)
	require.NotEmpty(t, thumbnail)
	// Original should be smaller than input since it was resized
	require.Less(t, len(original), len(data))
}

func TestProcess_SmallImageNotUpscaled(t *testing.T) {
	t.Parallel()

	p := NewWebPProcessor()
	data := encodeJPEG(t, createTestImage(50, 40))

	original, thumbnail, err := p.Process(data)

	require.NoError(t, err)
	require.NotEmpty(t, original)
	require.NotEmpty(t, thumbnail)
}

func TestProcess_InvalidData(t *testing.T) {
	t.Parallel()

	p := NewWebPProcessor()

	_, _, err := p.Process([]byte("not an image"))

	require.Error(t, err)
}
