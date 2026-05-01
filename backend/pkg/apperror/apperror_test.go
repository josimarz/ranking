package apperror_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/josimar/ranking/backend/pkg/apperror"
)

func TestAppError_SatisfiesErrorInterface(t *testing.T) {
	t.Parallel()

	var _ error = apperror.NewValidationError("test")
}

func TestAppError_ErrorMethod(t *testing.T) {
	t.Parallel()

	err := apperror.NewValidationError("invalid field")
	expected := fmt.Sprintf("[%s] invalid field", apperror.ValidationError)
	if err.Error() != expected {
		t.Errorf("got %q, want %q", err.Error(), expected)
	}
}

func TestErrorCodeConstants(t *testing.T) {
	t.Parallel()

	codes := []string{
		apperror.ValidationError,
		apperror.MissingUserID,
		apperror.RankingNotFound,
		apperror.ItemNotFound,
		apperror.NotOwner,
		apperror.NotAuthorized,
		apperror.InvalidSortField,
		apperror.InvalidCursor,
		apperror.IncompleteRatings,
		apperror.InvalidScore,
		apperror.InvalidImageFormat,
		apperror.ImageTooLarge,
		apperror.ImageDownloadFailed,
		apperror.MaxItemsReached,
		apperror.InternalError,
	}

	seen := make(map[string]bool)
	for _, code := range codes {
		if code == "" {
			t.Error("error code must not be empty")
		}
		if seen[code] {
			t.Errorf("duplicate error code: %s", code)
		}
		seen[code] = true
	}
}

func TestNewValidationError(t *testing.T) {
	t.Parallel()

	err := apperror.NewValidationError("name is required")

	if err.Code != apperror.ValidationError {
		t.Errorf("got code %q, want %q", err.Code, apperror.ValidationError)
	}
	if err.Message != "name is required" {
		t.Errorf("got message %q, want %q", err.Message, "name is required")
	}
	if err.Status != http.StatusBadRequest {
		t.Errorf("got status %d, want %d", err.Status, http.StatusBadRequest)
	}
}

func TestNewNotFoundError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		code string
		msg  string
	}{
		{"ranking not found", apperror.RankingNotFound, "ranking not found"},
		{"item not found", apperror.ItemNotFound, "item not found"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := apperror.NewNotFoundError(tt.code, tt.msg)

			if err.Code != tt.code {
				t.Errorf("got code %q, want %q", err.Code, tt.code)
			}
			if err.Message != tt.msg {
				t.Errorf("got message %q, want %q", err.Message, tt.msg)
			}
			if err.Status != http.StatusNotFound {
				t.Errorf("got status %d, want %d", err.Status, http.StatusNotFound)
			}
		})
	}
}

func TestNewForbiddenError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		code string
		msg  string
	}{
		{"not owner", apperror.NotOwner, "not the owner"},
		{"not authorized", apperror.NotAuthorized, "not authorized"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := apperror.NewForbiddenError(tt.code, tt.msg)

			if err.Code != tt.code {
				t.Errorf("got code %q, want %q", err.Code, tt.code)
			}
			if err.Message != tt.msg {
				t.Errorf("got message %q, want %q", err.Message, tt.msg)
			}
			if err.Status != http.StatusForbidden {
				t.Errorf("got status %d, want %d", err.Status, http.StatusForbidden)
			}
		})
	}
}

func TestNewBadRequestError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		code string
		msg  string
	}{
		{"invalid sort", apperror.InvalidSortField, "invalid sort field"},
		{"invalid cursor", apperror.InvalidCursor, "bad cursor"},
		{"incomplete ratings", apperror.IncompleteRatings, "missing ratings"},
		{"invalid image format", apperror.InvalidImageFormat, "bad format"},
		{"image too large", apperror.ImageTooLarge, "too large"},
		{"image download failed", apperror.ImageDownloadFailed, "download failed"},
		{"max items reached", apperror.MaxItemsReached, "limit reached"},
		{"missing user id", apperror.MissingUserID, "user id required"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := apperror.NewBadRequestError(tt.code, tt.msg)

			if err.Code != tt.code {
				t.Errorf("got code %q, want %q", err.Code, tt.code)
			}
			if err.Message != tt.msg {
				t.Errorf("got message %q, want %q", err.Message, tt.msg)
			}
			if err.Status != http.StatusBadRequest {
				t.Errorf("got status %d, want %d", err.Status, http.StatusBadRequest)
			}
		})
	}
}

func TestNewInternalError(t *testing.T) {
	t.Parallel()

	err := apperror.NewInternalError()

	if err.Code != apperror.InternalError {
		t.Errorf("got code %q, want %q", err.Code, apperror.InternalError)
	}
	if err.Message != "internal server error" {
		t.Errorf("got message %q, want %q", err.Message, "internal server error")
	}
	if err.Status != http.StatusInternalServerError {
		t.Errorf("got status %d, want %d", err.Status, http.StatusInternalServerError)
	}
}
