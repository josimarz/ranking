package rating

import (
	"context"
	"testing"
	"time"

	domainrating "github.com/josimar/ranking/backend/internal/domain/rating"
	"github.com/stretchr/testify/require"
)

func TestGetMyRatingUseCase_Rated(t *testing.T) {
	t.Parallel()

	existing := &domainrating.Rating{
		RankingID: testRankingID,
		ItemID:    testItemID,
		UserID:    testUserID,
		Scores:    map[string]int{testAttrID1: 80, testAttrID2: 90},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	uc := NewGetMyRatingUseCase(
		&mockRatingRepo{findByUserAndItemFn: func(_ context.Context, rankingID, itemID, userID string) (*domainrating.Rating, error) {
			require.Equal(t, testRankingID, rankingID)
			require.Equal(t, testItemID, itemID)
			require.Equal(t, testUserID, userID)
			return existing, nil
		}},
	)

	result, err := uc.Execute(context.Background(), testRankingID, testItemID, testUserID)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, existing.Scores, result.Scores)
}

func TestGetMyRatingUseCase_NotRated(t *testing.T) {
	t.Parallel()

	uc := NewGetMyRatingUseCase(
		&mockRatingRepo{findByUserAndItemFn: func(_ context.Context, _, _, _ string) (*domainrating.Rating, error) {
			return nil, nil
		}},
	)

	result, err := uc.Execute(context.Background(), testRankingID, testItemID, testUserID)

	require.NoError(t, err)
	require.Nil(t, result)
}
