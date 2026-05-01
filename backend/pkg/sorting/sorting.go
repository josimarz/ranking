package sorting

import (
	"fmt"
	"slices"
	"strings"

	"github.com/josimar/ranking/backend/pkg/apperror"
)

// Direction represents a sort direction.
type Direction string

const (
	// Ascending sort order.
	Ascending Direction = "asc"
	// Descending sort order.
	Descending Direction = "desc"
)

// SortField represents a single parsed sort field.
type SortField struct {
	Field     string
	Direction Direction
}

// ParseSort parses a JSON:API sort parameter into sort fields.
// Returns nil when param is empty. Returns an INVALID_SORT_FIELD error
// for fields not in allowedFields.
func ParseSort(param string, allowedFields []string) ([]SortField, error) {
	if param == "" {
		return nil, nil
	}

	parts := strings.Split(param, ",")
	fields := make([]SortField, 0, len(parts))

	for _, p := range parts {
		dir := Ascending
		field := p

		if after, found := strings.CutPrefix(field, "-"); found {
			dir = Descending
			field = after
		} else if after, found := strings.CutPrefix(field, "+"); found {
			field = after
		}

		if !slices.Contains(allowedFields, field) {
			return nil, apperror.NewBadRequestError(
				apperror.InvalidSortField,
				fmt.Sprintf("invalid sort field: %s", field),
			)
		}

		fields = append(fields, SortField{Field: field, Direction: dir})
	}

	return fields, nil
}
