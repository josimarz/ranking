package sorting

import (
	"errors"
	"fmt"
	"testing"

	"github.com/josimar/ranking/backend/pkg/apperror"
	"github.com/stretchr/testify/require"
)

const fieldName = "name"
const fieldOverall = "overall"

var allowedFields = []string{fieldName, fieldOverall, "createdAt"}

func TestParseSort_EmptyString(t *testing.T) {
	t.Parallel()
	result, err := ParseSort("", allowedFields)
	require.NoError(t, err)
	require.Nil(t, result)
}

func TestParseSort_SingleFieldAscending(t *testing.T) {
	t.Parallel()
	result, err := ParseSort(fieldName, allowedFields)
	require.NoError(t, err)
	require.Equal(t, []SortField{{Field: fieldName, Direction: Ascending}}, result)
}

func TestParseSort_SingleFieldDescending(t *testing.T) {
	t.Parallel()
	result, err := ParseSort(fmt.Sprintf("-%s", fieldOverall), allowedFields)
	require.NoError(t, err)
	require.Equal(t, []SortField{{Field: fieldOverall, Direction: Descending}}, result)
}

func TestParseSort_ExplicitAscendingPrefix(t *testing.T) {
	t.Parallel()
	result, err := ParseSort("+createdAt", allowedFields)
	require.NoError(t, err)
	require.Equal(t, []SortField{{Field: "createdAt", Direction: Ascending}}, result)
}

func TestParseSort_MultipleFields(t *testing.T) {
	t.Parallel()
	result, err := ParseSort(fmt.Sprintf("-%s,%s", fieldOverall, fieldName), allowedFields)
	require.NoError(t, err)
	require.Equal(t, []SortField{
		{Field: fieldOverall, Direction: Descending},
		{Field: fieldName, Direction: Ascending},
	}, result)
}

func TestParseSort_InvalidField(t *testing.T) {
	t.Parallel()
	result, err := ParseSort("invalid", allowedFields)
	require.Nil(t, result)
	require.Error(t, err)

	var appErr *apperror.AppError
	require.True(t, errors.As(err, &appErr))
	require.Equal(t, apperror.InvalidSortField, appErr.Code)
}
