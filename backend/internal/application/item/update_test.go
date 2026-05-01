package item

import (
	"context"
	"testing"

	domainitem "github.com/josimar/ranking/backend/internal/domain/item"
	domainranking "github.com/josimar/ranking/backend/internal/domain/ranking"
	"github.com/josimar/ranking/backend/pkg/apperror"
	"github.com/stretchr/testify/require"
)

func TestUpdateItemUseCase_SuccessByOwner(t *testing.T) {
	t.Parallel()

	var updated *domainitem.Item
	rankingRepo := &mockRankingRepo{
		findByIDFn: func(_ context.Context, _ string) (*domainranking.Ranking, error) {
			return newTestRanking(), nil
		},
	}
	itemRepo := &mockItemRepo{
		findByIDFn: func(_ context.Context, _, _ string) (*domainitem.Item, error) {
			return newTestItem(), nil
		},
		updateFn: func(_ context.Context, i *domainitem.Item) error {
			updated = i
			return nil
		},
	}

	uc := NewUpdateItemUseCase(rankingRepo, itemRepo, &mockImageStore{}, &mockImageProcessor{})
	input := UpdateItemInput{
		RankingID: testRankingID,
		ItemID:    testItemID,
		UserID:    testOwnerID,
		Name:      "Updated Name",
	}

	out, err := uc.Execute(context.Background(), input)

	require.NoError(t, err)
	require.Equal(t, "Updated Name", out.Name)
	require.NotNil(t, updated)
	require.Equal(t, "Updated Name", updated.Name)
}

func TestUpdateItemUseCase_SuccessByItemCreator(t *testing.T) {
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
		updateFn: func(_ context.Context, _ *domainitem.Item) error { return nil },
	}

	uc := NewUpdateItemUseCase(rankingRepo, itemRepo, &mockImageStore{}, &mockImageProcessor{})
	input := UpdateItemInput{
		RankingID: testRankingID,
		ItemID:    testItemID,
		UserID:    testCreatorID,
		Name:      "Creator Update",
	}

	out, err := uc.Execute(context.Background(), input)

	require.NoError(t, err)
	require.Equal(t, "Creator Update", out.Name)
}

func TestUpdateItemUseCase_NotAuthorized(t *testing.T) {
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

	uc := NewUpdateItemUseCase(rankingRepo, itemRepo, &mockImageStore{}, &mockImageProcessor{})
	input := UpdateItemInput{
		RankingID: testRankingID,
		ItemID:    testItemID,
		UserID:    testOtherUser,
		Name:      "Unauthorized",
	}

	_, err := uc.Execute(context.Background(), input)

	require.Error(t, err)
	var appErr *apperror.AppError
	require.ErrorAs(t, err, &appErr)
	require.Equal(t, apperror.NotAuthorized, appErr.Code)
}

func TestUpdateItemUseCase_ItemNotFound(t *testing.T) {
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

	uc := NewUpdateItemUseCase(rankingRepo, itemRepo, &mockImageStore{}, &mockImageProcessor{})
	input := UpdateItemInput{
		RankingID: testRankingID,
		ItemID:    "nonexistent",
		UserID:    testOwnerID,
		Name:      "Test",
	}

	_, err := uc.Execute(context.Background(), input)

	require.Error(t, err)
	var appErr *apperror.AppError
	require.ErrorAs(t, err, &appErr)
	require.Equal(t, apperror.ItemNotFound, appErr.Code)
}

func TestUpdateItemUseCase_WithNewImage(t *testing.T) {
	t.Parallel()

	imageStored := false
	rankingRepo := &mockRankingRepo{
		findByIDFn: func(_ context.Context, _ string) (*domainranking.Ranking, error) {
			return newTestRanking(), nil
		},
	}
	itemRepo := &mockItemRepo{
		findByIDFn: func(_ context.Context, _, _ string) (*domainitem.Item, error) {
			return newTestItem(), nil
		},
		updateFn: func(_ context.Context, _ *domainitem.Item) error { return nil },
	}
	imgStore := &mockImageStore{
		storeFn: func(_ context.Context, _, _ string, _, _ []byte) error {
			imageStored = true
			return nil
		},
	}

	uc := NewUpdateItemUseCase(rankingRepo, itemRepo, imgStore, &mockImageProcessor{})
	input := UpdateItemInput{
		RankingID: testRankingID,
		ItemID:    testItemID,
		UserID:    testOwnerID,
		Name:      "With Image",
		ImageData: []byte("new-image"),
	}

	out, err := uc.Execute(context.Background(), input)

	require.NoError(t, err)
	require.NotEmpty(t, out.ImageKey)
	require.True(t, imageStored)
}
