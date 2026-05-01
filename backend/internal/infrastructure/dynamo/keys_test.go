package dynamo_test

import (
	"testing"
	"time"

	"github.com/josimar/ranking/backend/internal/infrastructure/dynamo"
	"github.com/stretchr/testify/require"
)

func TestRankingPK(t *testing.T) {
	t.Parallel()
	require.Equal(t, "RANKING#abc-123", dynamo.RankingPK("abc-123"))
}

func TestMetadataSK(t *testing.T) {
	t.Parallel()
	require.Equal(t, "METADATA", dynamo.MetadataSK())
}

func TestItemSK(t *testing.T) {
	t.Parallel()
	require.Equal(t, "ITEM#item-1", dynamo.ItemSK("item-1"))
}

func TestTagSK(t *testing.T) {
	t.Parallel()
	require.Equal(t, "TAG#games", dynamo.TagSK("games"))
}

func TestRatingPK(t *testing.T) {
	t.Parallel()
	require.Equal(t, "RATING#r1#i1", dynamo.RatingPK("r1", "i1"))
}

func TestUserSK(t *testing.T) {
	t.Parallel()
	require.Equal(t, "USER#u1", dynamo.UserSK("u1"))
}

func TestOwnerGSI1PK(t *testing.T) {
	t.Parallel()
	require.Equal(t, "OWNER#u1", dynamo.OwnerGSI1PK("u1"))
}

func TestRankingGSI1SK(t *testing.T) {
	t.Parallel()
	ts := time.Date(2025, 5, 1, 12, 0, 0, 123456789, time.UTC)
	require.Equal(t, "RANKING#2025-05-01T12:00:00.123456789Z", dynamo.RankingGSI1SK(ts))
}

func TestPublicGSI2PK(t *testing.T) {
	t.Parallel()
	require.Equal(t, "PUBLIC", dynamo.PublicGSI2PK())
}

func TestTagGSI3PK(t *testing.T) {
	t.Parallel()
	require.Equal(t, "TAG#movies", dynamo.TagGSI3PK("movies"))
}

func TestTagGSI3SK(t *testing.T) {
	t.Parallel()
	ts := time.Date(2025, 5, 1, 12, 0, 0, 123456789, time.UTC)
	require.Equal(t, "2025-05-01T12:00:00.123456789Z#r1", dynamo.TagGSI3SK(ts, "r1"))
}
