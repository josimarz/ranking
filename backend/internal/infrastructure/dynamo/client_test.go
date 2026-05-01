package dynamo_test

import (
	"context"
	"testing"

	"github.com/josimar/ranking/backend/internal/infrastructure/dynamo"
	"github.com/stretchr/testify/require"
)

func TestNewDynamoClient(t *testing.T) {
	t.Run("returns a non-nil client", func(t *testing.T) {
		t.Parallel()

		client, err := dynamo.NewDynamoClient(context.Background())

		require.NoError(t, err)
		require.NotNil(t, client)
	})

	t.Run("returns a client when AWS_ENDPOINT_URL is set", func(t *testing.T) {
		t.Setenv("AWS_ENDPOINT_URL", "http://localhost:4566")

		client, err := dynamo.NewDynamoClient(context.Background())

		require.NoError(t, err)
		require.NotNil(t, client)
	})
}
