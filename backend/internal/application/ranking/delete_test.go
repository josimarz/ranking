package ranking

import (
	"context"
	"testing"

	domainitem "github.com/josimar/ranking/backend/internal/domain/item"
	domain "github.com/josimar/ranking/backend/internal/domain/ranking"
	"github.com/stretchr/testify/require"
)

func TestDeleteRankingUseCase_Success(t *testing.T) {
	t.Parallel()

	r := newFullTestRanking()
	items := []*domainitem.Item{
		{ID: "item-1", RankingID: testRankingID},
		{ID: "item-2", RankingID: testRankingID},
	}

	var deletedRatingItemIDs []string
	var deletedImages, deletedItems, deletedRanking bool

	rankingRepo := &mockRankingRepo{
		findByIDFn: func(_ context.Context, _ string) (*domain.Ranking, error) {
			return r, nil
		},
		deleteFn: func(_ context.Context, _ string) error {
			deletedRanking = true
			return nil
		},
	}

	itemRepo := &mockItemRepo{
		findByRankingFn: func(_ context.Context, _ string) ([]*domainitem.Item, error) {
			return items, nil
		},
		deleteByRankingFn: func(_ context.Context, _ string) error {
			deletedItems = true
			return nil
		},
	}

	ratingRepo := &mockRatingRepo{
		deleteByRankingFn: func(_ context.Context, _ string, itemIDs []string) error {
			deletedRatingItemIDs = itemIDs
			return nil
		},
	}

	imageStore := &mockImageStore{
		deleteByRankingFn: func(_ context.Context, _ string) error {
			deletedImages = true
			return nil
		},
	}

	uc := NewDeleteRankingUseCase(rankingRepo, itemRepo, ratingRepo, imageStore)
	err := uc.Execute(context.Background(), testRankingID, testOwnerID)

	require.NoError(t, err)
	require.True(t, deletedRanking)
	require.True(t, deletedItems)
	require.True(t, deletedImages)
	require.Equal(t, []string{"item-1", "item-2"}, deletedRatingItemIDs)
}

func TestDeleteRankingUseCase_NotOwner(t *testing.T) {
	t.Parallel()

	r := newFullTestRanking()
	rankingRepo := &mockRankingRepo{
		findByIDFn: func(_ context.Context, _ string) (*domain.Ranking, error) {
			return r, nil
		},
	}

	uc := NewDeleteRankingUseCase(rankingRepo, &mockItemRepo{}, &mockRatingRepo{}, &mockImageStore{})
	err := uc.Execute(context.Background(), testRankingID, "other-user")

	require.Error(t, err)
	require.ErrorIs(t, err, domain.ErrNotOwner)
}

func TestDeleteRankingUseCase_NotFound(t *testing.T) {
	t.Parallel()

	rankingRepo := &mockRankingRepo{
		findByIDFn: func(_ context.Context, _ string) (*domain.Ranking, error) {
			return nil, domain.ErrRankingNotFound
		},
	}

	uc := NewDeleteRankingUseCase(rankingRepo, &mockItemRepo{}, &mockRatingRepo{}, &mockImageStore{})
	err := uc.Execute(context.Background(), "nonexistent", "user-1")

	require.Error(t, err)
	require.ErrorIs(t, err, domain.ErrRankingNotFound)
}
