package response_test

import (
	"testing"

	"github.com/josimar/ranking/backend/internal/interfaces/http/response"
	"github.com/stretchr/testify/require"
)

func TestNewPaginatedResponse_WithCursor(t *testing.T) {
	t.Parallel()

	data := []string{"a", "b"}
	resp := response.NewPaginatedResponse(data, "cursor123")

	require.Equal(t, data, resp.Data)
	require.Equal(t, "cursor123", resp.Pagination.NextCursor)
	require.True(t, resp.Pagination.HasMore)
}

func TestNewPaginatedResponse_WithoutCursor(t *testing.T) {
	t.Parallel()

	data := []string{"a"}
	resp := response.NewPaginatedResponse(data, "")

	require.Equal(t, data, resp.Data)
	require.Empty(t, resp.Pagination.NextCursor)
	require.False(t, resp.Pagination.HasMore)
}
