package uuid_test

import (
	"regexp"
	"testing"

	"github.com/josimar/ranking/backend/pkg/uuid"
	"github.com/stretchr/testify/require"
)

var uuidRegex = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

func TestNew_ReturnsValidUUID(t *testing.T) {
	t.Parallel()

	id := uuid.New()

	require.Regexp(t, uuidRegex, id)
}

func TestNew_ReturnsUUIDv7(t *testing.T) {
	t.Parallel()

	id := uuid.New()

	require.Equal(t, "7", string(id[14]), "expected UUID version 7")
}

func TestNew_Uniqueness(t *testing.T) {
	t.Parallel()

	seen := make(map[string]struct{}, 1000)
	for range 1000 {
		id := uuid.New()
		_, exists := seen[id]
		require.False(t, exists, "duplicate UUID found: %s", id)
		seen[id] = struct{}{}
	}
}

func TestIsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "valid UUID v7",
			input: uuid.New(),
			want:  true,
		},
		{
			name:  "valid UUID v4",
			input: "550e8400-e29b-41d4-a716-446655440000",
			want:  true,
		},
		{
			name:  "empty string",
			input: "",
			want:  false,
		},
		{
			name:  "random string",
			input: "not-a-uuid",
			want:  false,
		},
		{
			name:  "too short",
			input: "550e8400-e29b-41d4-a716",
			want:  false,
		},
		{
			name:  "invalid characters",
			input: "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
			want:  false,
		},
		{
			name:  "missing dashes",
			input: "550e8400e29b41d4a716446655440000",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := uuid.IsValid(tt.input)
			require.Equal(t, tt.want, got)
		})
	}
}
