package pagination_test

import (
	"encoding/base64"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/josimar/ranking/backend/pkg/apperror"
	"github.com/josimar/ranking/backend/pkg/pagination"
	"github.com/stretchr/testify/require"
)

func TestEncodeCursor_RoundTrip(t *testing.T) {
	t.Parallel()

	key := map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{Value: "RANKING#123"},
		"SK": &types.AttributeValueMemberS{Value: "ITEM#456"},
	}

	cursor := pagination.EncodeCursor(key)
	require.NotEmpty(t, cursor)

	decoded, err := pagination.DecodeCursor(cursor)
	require.NoError(t, err)
	require.Equal(t, key, decoded)
}

func TestEncodeCursor_NilKey(t *testing.T) {
	t.Parallel()

	cursor := pagination.EncodeCursor(nil)
	require.Empty(t, cursor)
}

func TestDecodeCursor_EmptyCursor(t *testing.T) {
	t.Parallel()

	decoded, err := pagination.DecodeCursor("")
	require.NoError(t, err)
	require.Nil(t, decoded)
}

func TestDecodeCursor_InvalidBase64(t *testing.T) {
	t.Parallel()

	_, err := pagination.DecodeCursor("not-valid-base64!!!")
	require.Error(t, err)

	var appErr *apperror.AppError
	require.ErrorAs(t, err, &appErr)
	require.Equal(t, apperror.InvalidCursor, appErr.Code)
}

func TestDecodeCursor_InvalidJSON(t *testing.T) {
	t.Parallel()

	encoded := base64.URLEncoding.EncodeToString([]byte("not-json"))
	_, err := pagination.DecodeCursor(encoded)
	require.Error(t, err)

	var appErr *apperror.AppError
	require.ErrorAs(t, err, &appErr)
	require.Equal(t, apperror.InvalidCursor, appErr.Code)
}

func TestClampLimit(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		limit        int
		defaultLimit int
		maxLimit     int
		want         int
	}{
		{"zero returns default", 0, 10, 50, 10},
		{"negative returns default", -5, 10, 50, 10},
		{"exceeds max returns max", 100, 10, 50, 50},
		{"within range returns limit", 25, 10, 50, 25},
		{"equal to max returns max", 50, 10, 50, 50},
		{"equal to one returns one", 1, 10, 50, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := pagination.ClampLimit(tt.limit, tt.defaultLimit, tt.maxLimit)
			require.Equal(t, tt.want, got)
		})
	}
}
