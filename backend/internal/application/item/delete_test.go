package item

import (
	"context"
	"testing"

	domainitem "github.com/josimar/ranking/backend/internal/domain/item"
	domainranking "github.com/josimar/ranking/backend/internal/domain/ranking"
	"github.com/josimar/ranking/backend/pkg/apperror"
	"github.com/stretchr/testify/require"
)

func TestDeleteItemUseCase_Success(t *testing.T) {
	t.Parallel()

	ratingsDeleted := false
	imageDeleted := false
	itemDeleted := false

	rankingRepo := &mockRankingRepo{
		findByIDFn: func(_ context.Context, _ string) (*domainranking.Ranking, error) {
			return newTestRanking(), nil
		},
	}
	itemRepo := &mockItemRepo{
		findByIDFn: func(_ context.Context, _, _ string) (*domainitem.Item, error) {
			return newTestItem(), nil
		},
		deleteFn: func(_ context.Context, _, _ string) error {
			itemDeleted = true
			return nil
		},
	}
	ratingRepo := &mockRatingRepo{
		deleteByItemFn: func(_ context.Context, _, _ string) error {
			ratingsDeleted = true
			return nil
		},
	}
	imgStore := &mockImageStore{
		deleteFn: func(_ context.Context, _, _ string) error {
			imageDeleted = true
			return nil
		},
	}

	uc := NewDeleteItemUseCase(rankingRepo, itemRepo, ratingRepo, imgStore)

	err := uc.Execute(context.Background(), testRankingID, testItemID, testOwnerID)

	require.NoError(t, err)
	require.True(t, ratingsDeleted)
	require.True(t, imageDeleted)
	require.True(t, itemDeleted)
}

func TestDeleteItemUseCase_NotAuthorized(t *testing.T) {
	t.Parallel()

	rankingRepo := &mockRankingRepo{
		findByIDFn: func(_ context.Context, _ string) (*domainranking.Ranking, error) {
			return newTestRanking(), nil
		},
	}
	itemRepo := &mockItemRepo{
		findByIDFn: func(_ context.Context, _, _ string) (*domainitem.Item, error) {
			return newTestItem(), nil
		},
	}

	uc := NewDeleteItemUseCase(rankingRepo, itemRepo, &mockRatingRepo{}, &mockImageStore{})

	err := uc.Execute(context.Background(), testRankingID, testItemID, testOtherUser)

	require.Error(t, err)
	var appErr *apperror.AppError
	require.ErrorAs(t, err, &appErr)
	require.Equal(t, apperror.NotAuthorized, appErr.Code)
}

func TestDeleteItemUseCase_ItemNotFound(t *testing.T) {
	t.Parallel()

	rankingRepo := &mockRankingRepo{
		findByIDFn: func(_ context.Context, _ string) (*domainranking.Ranking, error) {
			return newTestRanking(), nil
		},
	}
	itemRepo := &mockItemRepo{
		findByIDFn: func(_ context.Context, _, _ string) (*domainitem.Item, error) {
			return nil, domainitem.ErrItemNotFound
		},
	}

	uc := NewDeleteItemUseCase(rankingRepo, itemRepo, &mockRatingRepo{}, &mockImageStore{})

	err := uc.Execute(context.Background(), testRankingID, "nonexistent", testOwnerID)

	require.Error(t, err)
	var appErr *apperror.AppError
	require.ErrorAs(t, err, &appErr)
	require.Equal(t, apperror.ItemNotFound, appErr.Code)
}

func TestDeleteItemUseCase_RankingNotFound(t *testing.T) {
	t.Parallel()

	rankingRepo := &mockRankingRepo{
		findByIDFn: func(_ context.Context, _ string) (*domainranking.Ranking, error) {
			return nil, domainranking.ErrRankingNotFound
		},
	}

	uc := NewDeleteItemUseCase(rankingRepo, &mockItemRepo{}, &mockRatingRepo{}, &mockImageStore{})

	err := uc.Execute(context.Background(), "nonexistent", testItemID, testOwnerID)

	require.Error(t, err)
	var appErr *apperror.AppError
	require.ErrorAs(t, err, &appErr)
	require.Equal(t, apperror.RankingNotFound, appErr.Code)
}
