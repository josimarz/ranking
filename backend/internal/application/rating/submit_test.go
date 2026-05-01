package rating

import (
	"context"
	"testing"

	domainitem "github.com/josimar/ranking/backend/internal/domain/item"
	domainranking "github.com/josimar/ranking/backend/internal/domain/ranking"
	domainrating "github.com/josimar/ranking/backend/internal/domain/rating"
	"github.com/josimar/ranking/backend/pkg/apperror"
	"github.com/stretchr/testify/require"
)

func TestSubmitRatingUseCase_Success(t *testing.T) {
	t.Parallel()

	rnk := &domainranking.Ranking{
		ID: testRankingID,
		Attributes: []domainranking.Attribute{
			{ID: testAttrID1, Name: testGraphics, Active: true},
			{ID: testAttrID2, Name: testSound, Active: true},
		},
	}

	itm := &domainitem.Item{ID: testItemID, RankingID: testRankingID}

	var savedRating *domainrating.Rating

	uc := NewSubmitRatingUseCase(
		&mockRankingRepo{findByIDFn: func(_ context.Context, id string) (*domainranking.Ranking, error) {
			require.Equal(t, testRankingID, id)
			return rnk, nil
		}},
		&mockItemRepo{findByIDFn: func(_ context.Context, rankingID, itemID string) (*domainitem.Item, error) {
			require.Equal(t, testRankingID, rankingID)
			require.Equal(t, testItemID, itemID)
			return itm, nil
		}},
		&mockRatingRepo{saveFn: func(_ context.Context, r *domainrating.Rating) error {
			savedRating = r
			return nil
		}},
	)

	scores := map[string]int{testAttrID1: 80, testAttrID2: 90}
	result, err := uc.Execute(context.Background(), testRankingID, testItemID, testUserID, scores)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, testRankingID, result.RankingID)
	require.Equal(t, testItemID, result.ItemID)
	require.Equal(t, testUserID, result.UserID)
	require.Equal(t, scores, result.Scores)
	require.NotNil(t, savedRating)
}

func TestSubmitRatingUseCase_IncompleteScores(t *testing.T) {
	t.Parallel()

	rnk := &domainranking.Ranking{
		ID: testRankingID,
		Attributes: []domainranking.Attribute{
			{ID: testAttrID1, Name: testGraphics, Active: true},
			{ID: testAttrID2, Name: testSound, Active: true},
		},
	}

	itm := &domainitem.Item{ID: testItemID, RankingID: testRankingID}

	uc := NewSubmitRatingUseCase(
		&mockRankingRepo{findByIDFn: func(_ context.Context, _ string) (*domainranking.Ranking, error) {
			return rnk, nil
		}},
		&mockItemRepo{findByIDFn: func(_ context.Context, _, _ string) (*domainitem.Item, error) {
			return itm, nil
		}},
		&mockRatingRepo{},
	)

	scores := map[string]int{testAttrID1: 80} // missing a2
	result, err := uc.Execute(context.Background(), testRankingID, testItemID, testUserID, scores)

	require.Nil(t, result)
	require.Error(t, err)

	var appErr *apperror.AppError
	require.ErrorAs(t, err, &appErr)
	require.Equal(t, apperror.IncompleteRatings, appErr.Code)
}

func TestSubmitRatingUseCase_InvalidScoreRange(t *testing.T) {
	t.Parallel()

	rnk := &domainranking.Ranking{
		ID: testRankingID,
		Attributes: []domainranking.Attribute{
			{ID: testAttrID1, Name: testGraphics, Active: true},
		},
	}

	itm := &domainitem.Item{ID: testItemID, RankingID: testRankingID}

	uc := NewSubmitRatingUseCase(
		&mockRankingRepo{findByIDFn: func(_ context.Context, _ string) (*domainranking.Ranking, error) {
			return rnk, nil
		}},
		&mockItemRepo{findByIDFn: func(_ context.Context, _, _ string) (*domainitem.Item, error) {
			return itm, nil
		}},
		&mockRatingRepo{},
	)

	scores := map[string]int{testAttrID1: 150}
	result, err := uc.Execute(context.Background(), testRankingID, testItemID, testUserID, scores)

	require.Nil(t, result)
	require.Error(t, err)

	var appErr *apperror.AppError
	require.ErrorAs(t, err, &appErr)
	require.Equal(t, apperror.InvalidScore, appErr.Code)
}

func TestSubmitRatingUseCase_RankingNotFound(t *testing.T) {
	t.Parallel()

	uc := NewSubmitRatingUseCase(
		&mockRankingRepo{findByIDFn: func(_ context.Context, _ string) (*domainranking.Ranking, error) {
			return nil, domainranking.ErrRankingNotFound
		}},
		&mockItemRepo{},
		&mockRatingRepo{},
	)

	result, err := uc.Execute(context.Background(), testRankingID, testItemID, testUserID, map[string]int{testAttrID1: 80})

	require.Nil(t, result)
	require.Error(t, err)

	var appErr *apperror.AppError
	require.ErrorAs(t, err, &appErr)
	require.Equal(t, apperror.RankingNotFound, appErr.Code)
}

func TestSubmitRatingUseCase_ItemNotFound(t *testing.T) {
	t.Parallel()

	rnk := &domainranking.Ranking{
		ID: testRankingID,
		Attributes: []domainranking.Attribute{
			{ID: testAttrID1, Name: testGraphics, Active: true},
		},
	}

	uc := NewSubmitRatingUseCase(
		&mockRankingRepo{findByIDFn: func(_ context.Context, _ string) (*domainranking.Ranking, error) {
			return rnk, nil
		}},
		&mockItemRepo{findByIDFn: func(_ context.Context, _, _ string) (*domainitem.Item, error) {
			return nil, apperror.NewNotFoundError(apperror.ItemNotFound, "item not found")
		}},
		&mockRatingRepo{},
	)

	result, err := uc.Execute(context.Background(), testRankingID, testItemID, testUserID, map[string]int{testAttrID1: 80})

	require.Nil(t, result)
	require.Error(t, err)

	var appErr *apperror.AppError
	require.ErrorAs(t, err, &appErr)
	require.Equal(t, apperror.ItemNotFound, appErr.Code)
}
