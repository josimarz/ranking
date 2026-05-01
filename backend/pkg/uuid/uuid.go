// Package uuid provides UUID v7 generation and validation.
package uuid

import (
	"regexp"

	googleuuid "github.com/google/uuid"
)

var uuidRegex = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

// New generates a new UUID v7 string.
func New() string {
	return googleuuid.Must(googleuuid.NewV7()).String()
}

// IsValid checks whether the given string is a valid UUID in standard format.
func IsValid(id string) bool {
	return uuidRegex.MatchString(id)
}
