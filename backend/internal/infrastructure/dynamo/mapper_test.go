package dynamo_test

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/josimar/ranking/backend/internal/domain/item"
	"github.com/josimar/ranking/backend/internal/domain/ranking"
	"github.com/josimar/ranking/backend/internal/domain/rating"
	"github.com/josimar/ranking/backend/internal/infrastructure/dynamo"
	"github.com/stretchr/testify/require"
)

func newTestRanking() *ranking.Ranking {
	now := time.Date(2025, 5, 1, 12, 0, 0, 0, time.UTC)
	return &ranking.Ranking{
		ID:          "r1",
		Name:        "Video Games",
		NameLower:   "video games",
		Description: "A ranking of consoles",
		Visibility:  ranking.Public,
		Tags:        []string{"games", "consoles"},
		Attributes: []ranking.Attribute{
			{ID: "a1", Name: "Graphics", Description: "Visual quality", Active: true},
		},
		OwnerUserID: "u1",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func TestMarshalUnmarshalRanking(t *testing.T) {
	t.Parallel()

	t.Run("public ranking round-trip", func(t *testing.T) {
		t.Parallel()

		r := newTestRanking()
		m := dynamo.MarshalRanking(r)

		require.Equal(t, "RANKING#r1", m["PK"].(*types.AttributeValueMemberS).Value)
		require.Equal(t, "METADATA", m["SK"].(*types.AttributeValueMemberS).Value)
		require.Equal(t, "OWNER#u1", m["GSI1PK"].(*types.AttributeValueMemberS).Value)
		require.Equal(t, "PUBLIC", m["GSI2PK"].(*types.AttributeValueMemberS).Value)
		require.Equal(t, "Video Games", m["Name"].(*types.AttributeValueMemberS).Value)

		result, err := dynamo.UnmarshalRanking(m)
		require.NoError(t, err)
		require.Equal(t, r.ID, result.ID)
		require.Equal(t, r.Name, result.Name)
		require.Equal(t, r.NameLower, result.NameLower)
		require.Equal(t, r.Description, result.Description)
		require.Equal(t, r.Visibility, result.Visibility)
		require.Equal(t, r.Tags, result.Tags)
		require.Equal(t, r.OwnerUserID, result.OwnerUserID)
		require.Equal(t, r.CreatedAt.UTC(), result.CreatedAt.UTC())
		require.Equal(t, r.UpdatedAt.UTC(), result.UpdatedAt.UTC())
		require.Len(t, result.Attributes, 1)
		require.Equal(t, r.Attributes[0].ID, result.Attributes[0].ID)
		require.Equal(t, r.Attributes[0].Name, result.Attributes[0].Name)
		require.Equal(t, r.Attributes[0].Active, result.Attributes[0].Active)
	})

	t.Run("private ranking has no GSI2PK", func(t *testing.T) {
		t.Parallel()

		r := newTestRanking()
		r.Visibility = ranking.Private
		m := dynamo.MarshalRanking(r)

		_, hasGSI2 := m["GSI2PK"]
		require.False(t, hasGSI2)
	})
}

func TestMarshalUnmarshalItem(t *testing.T) {
	t.Parallel()

	now := time.Date(2025, 5, 1, 12, 0, 0, 0, time.UTC)
	i := &item.Item{
		ID:        "i1",
		Name:      "Sega Mega Drive",
		ImageKey:  "images/r1/i1.jpg",
		RankingID: "r1",
		CreatedBy: "u1",
		CreatedAt: now,
		UpdatedAt: now,
	}

	m := dynamo.MarshalItem(i)

	require.Equal(t, "RANKING#r1", m["PK"].(*types.AttributeValueMemberS).Value)
	require.Equal(t, "ITEM#i1", m["SK"].(*types.AttributeValueMemberS).Value)
	require.Equal(t, "Sega Mega Drive", m["Name"].(*types.AttributeValueMemberS).Value)

	result, err := dynamo.UnmarshalItem(m)
	require.NoError(t, err)
	require.Equal(t, i.ID, result.ID)
	require.Equal(t, i.Name, result.Name)
	require.Equal(t, i.ImageKey, result.ImageKey)
	require.Equal(t, i.RankingID, result.RankingID)
	require.Equal(t, i.CreatedBy, result.CreatedBy)
	require.Equal(t, i.CreatedAt.UTC(), result.CreatedAt.UTC())
}

func TestMarshalUnmarshalRating(t *testing.T) {
	t.Parallel()

	now := time.Date(2025, 5, 1, 12, 0, 0, 0, time.UTC)
	r := &rating.Rating{
		RankingID: "r1",
		ItemID:    "i1",
		UserID:    "u1",
		Scores:    map[string]int{"a1": 90, "a2": 80},
		CreatedAt: now,
		UpdatedAt: now,
	}

	m := dynamo.MarshalRating(r)

	require.Equal(t, "RATING#r1#i1", m["PK"].(*types.AttributeValueMemberS).Value)
	require.Equal(t, "USER#u1", m["SK"].(*types.AttributeValueMemberS).Value)

	result, err := dynamo.UnmarshalRating(m)
	require.NoError(t, err)
	require.Equal(t, r.RankingID, result.RankingID)
	require.Equal(t, r.ItemID, result.ItemID)
	require.Equal(t, r.UserID, result.UserID)
	require.Equal(t, r.Scores, result.Scores)
	require.Equal(t, r.CreatedAt.UTC(), result.CreatedAt.UTC())
}

func TestMarshalRankingTag(t *testing.T) {
	t.Parallel()

	ts := time.Date(2025, 5, 1, 12, 0, 0, 0, time.UTC)
	m := dynamo.MarshalRankingTag("r1", "games", ts, "Video Games")

	require.Equal(t, "RANKING#r1", m["PK"].(*types.AttributeValueMemberS).Value)
	require.Equal(t, "TAG#games", m["SK"].(*types.AttributeValueMemberS).Value)
	require.Equal(t, "TAG#games", m["GSI3PK"].(*types.AttributeValueMemberS).Value)
	require.Contains(t, m["GSI3SK"].(*types.AttributeValueMemberS).Value, "r1")
	require.Equal(t, "Video Games", m["RankingName"].(*types.AttributeValueMemberS).Value)
}
