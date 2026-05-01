package apperror

import (
	"fmt"
	"net/http"
)

// Error code constants for all application error types.
const (
	ValidationError     = "VALIDATION_ERROR"
	MissingUserID       = "MISSING_USER_ID"
	RankingNotFound     = "RANKING_NOT_FOUND"
	ItemNotFound        = "ITEM_NOT_FOUND"
	NotOwner            = "NOT_OWNER"
	NotAuthorized       = "NOT_AUTHORIZED"
	InvalidSortField    = "INVALID_SORT_FIELD"
	InvalidCursor       = "INVALID_CURSOR"
	IncompleteRatings   = "INCOMPLETE_RATINGS"
	InvalidScore        = "INVALID_SCORE"
	InvalidImageFormat  = "INVALID_IMAGE_FORMAT"
	ImageTooLarge       = "IMAGE_TOO_LARGE"
	ImageDownloadFailed = "IMAGE_DOWNLOAD_FAILED"
	MaxItemsReached     = "MAX_ITEMS_REACHED"
	InternalError       = "INTERNAL_ERROR"
)

// AppError represents a structured application error.
type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
}

// Error satisfies the error interface.
func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// NewValidationError creates a 400 error with the VALIDATION_ERROR code.
func NewValidationError(msg string) *AppError {
	return &AppError{Code: ValidationError, Message: msg, Status: http.StatusBadRequest}
}

// NewNotFoundError creates a 404 error with the given code.
func NewNotFoundError(code, msg string) *AppError {
	return &AppError{Code: code, Message: msg, Status: http.StatusNotFound}
}

// NewForbiddenError creates a 403 error with the given code.
func NewForbiddenError(code, msg string) *AppError {
	return &AppError{Code: code, Message: msg, Status: http.StatusForbidden}
}

// NewBadRequestError creates a 400 error with the given code.
func NewBadRequestError(code, msg string) *AppError {
	return &AppError{Code: code, Message: msg, Status: http.StatusBadRequest}
}

// NewInternalError creates a 500 internal server error.
func NewInternalError() *AppError {
	return &AppError{Code: InternalError, Message: "internal server error", Status: http.StatusInternalServerError}
}
