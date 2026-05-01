// Package image provides image processing infrastructure.
package image

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg" // register JPEG decoder
	_ "image/png"  // register PNG decoder

	"github.com/chai2010/webp"
	"github.com/josimar/ranking/backend/pkg/apperror"
	"golang.org/x/image/draw"
	_ "golang.org/x/image/webp" // register WebP decoder
)

const (
	maxFileSize      = 10 * 1024 * 1024 // 10 MB
	maxOriginalDim   = 800
	thumbnailDim     = 150
	originalQuality  = 80
	thumbnailQuality = 75
)

// WebPProcessor implements the domain ImageProcessor interface.
type WebPProcessor struct{}

// NewWebPProcessor creates a new WebPProcessor.
func NewWebPProcessor() *WebPProcessor {
	return &WebPProcessor{}
}

// ValidateFormat checks that data has a valid image magic bytes and is within size limits.
func (p *WebPProcessor) ValidateFormat(data []byte) error {
	if len(data) > maxFileSize {
		return apperror.NewBadRequestError(apperror.ImageTooLarge, "image exceeds maximum size of 10 MB")
	}

	if !isJPEG(data) && !isPNG(data) && !isWebP(data) {
		return apperror.NewBadRequestError(apperror.InvalidImageFormat, "image must be JPEG, PNG, or WebP")
	}

	return nil
}

// Process decodes the image, resizes to original and thumbnail, and encodes both as WebP.
func (p *WebPProcessor) Process(data []byte) ([]byte, []byte, error) {
	src, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode image: %w", err)
	}

	original, err := encodeResized(src, maxOriginalDim, false, originalQuality)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to encode original: %w", err)
	}

	thumbnail, err := encodeResized(src, thumbnailDim, true, thumbnailQuality)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to encode thumbnail: %w", err)
	}

	return original, thumbnail, nil
}

func encodeResized(src image.Image, maxDim int, crop bool, quality int) ([]byte, error) {
	var dst *image.RGBA

	if crop {
		dst = centerCrop(src, maxDim)
	} else {
		dst = fitResize(src, maxDim)
	}

	var buf bytes.Buffer
	if err := webp.Encode(&buf, dst, &webp.Options{Quality: float32(quality)}); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func fitResize(src image.Image, maxDim int) *image.RGBA {
	bounds := src.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	if w <= maxDim && h <= maxDim {
		dst := image.NewRGBA(image.Rect(0, 0, w, h))
		draw.Copy(dst, image.Point{}, src, bounds, draw.Src, nil)

		return dst
	}

	newW, newH := scaleDimensions(w, h, maxDim)
	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, bounds, draw.Over, nil)

	return dst
}

func centerCrop(src image.Image, size int) *image.RGBA {
	bounds := src.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	// Determine crop region (square from center)
	cropSize := min(w, h)
	x0 := bounds.Min.X + (w-cropSize)/2
	y0 := bounds.Min.Y + (h-cropSize)/2
	cropRect := image.Rect(x0, y0, x0+cropSize, y0+cropSize)

	dst := image.NewRGBA(image.Rect(0, 0, size, size))
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, cropRect, draw.Over, nil)

	return dst
}

func scaleDimensions(w, h, maxDim int) (int, int) {
	if w >= h {
		return maxDim, maxDim * h / w
	}

	return maxDim * w / h, maxDim
}

func isJPEG(data []byte) bool {
	return len(data) >= 3 && data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF
}

func isPNG(data []byte) bool {
	return len(data) >= 4 && data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47
}

func isWebP(data []byte) bool {
	return len(data) >= 12 &&
		data[0] == 'R' && data[1] == 'I' && data[2] == 'F' && data[3] == 'F' &&
		data[8] == 'W' && data[9] == 'E' && data[10] == 'B' && data[11] == 'P'
}
