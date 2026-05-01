package item

import (
	"context"
	"testing"

	domainitem "github.com/josimar/ranking/backend/internal/domain/item"
	domainranking "github.com/josimar/ranking/backend/internal/domain/ranking"
	"github.com/josimar/ranking/backend/pkg/apperror"
	"github.com/stretchr/testify/require"
)

func TestAddItemUseCase_SuccessWithoutImage(t *testing.T) {
	t.Parallel()

	var saved *domainitem.Item
	rankingRepo := &mockRankingRepo{
		findByIDFn: func(_ context.Context, _ string) (*domainranking.Ranking, error) {
			return newTestRanking(), nil
		},
	}
	itemRepo := &mockItemRepo{
		countByRankingFn: func(_ context.Context, _ string) (int, error) { return 5, nil },
		saveFn: func(_ context.Context, i *domainitem.Item) error {
			saved = i
			return nil
		},
	}

	uc := NewAddItemUseCase(rankingRepo, itemRepo, &mockImageStore{}, &mockImageProcessor{})
	input := AddItemInput{
		RankingID: testRankingID,
		UserID:    testCreatorID,
		Name:      "PlayStation 5",
	}

	out, err := uc.Execute(context.Background(), input)

	require.NoError(t, err)
	require.NotEmpty(t, out.ID)
	require.Equal(t, "PlayStation 5", out.Name)
	require.Empty(t, out.ImageKey)
	require.Equal(t, testRankingID, out.RankingID)
	require.Equal(t, testCreatorID, out.CreatedBy)
	require.NotNil(t, saved)
}

func TestAddItemUseCase_SuccessWithImageData(t *testing.T) {
	t.Parallel()

	imageStored := false
	rankingRepo := &mockRankingRepo{
		findByIDFn: func(_ context.Context, _ string) (*domainranking.Ranking, error) {
			return newTestRanking(), nil
		},
	}
	itemRepo := &mockItemRepo{
		countByRankingFn: func(_ context.Context, _ string) (int, error) { return 0, nil },
		saveFn:           func(_ context.Context, _ *domainitem.Item) error { return nil },
	}
	imgStore := &mockImageStore{
		storeFn: func(_ context.Context, _, _ string, _, _ []byte) error {
			imageStored = true
			return nil
		},
	}

	uc := NewAddItemUseCase(rankingRepo, itemRepo, imgStore, &mockImageProcessor{})
	input := AddItemInput{
		RankingID: testRankingID,
		UserID:    testCreatorID,
		Name:      "Xbox Series X",
		ImageData: []byte("fake-image-data"),
	}

	out, err := uc.Execute(context.Background(), input)

	require.NoError(t, err)
	require.NotEmpty(t, out.ImageKey)
	require.True(t, imageStored)
}

func TestAddItemUseCase_MaxItemsReached(t *testing.T) {
	t.Parallel()

	rankingRepo := &mockRankingRepo{
		findByIDFn: func(_ context.Context, _ string) (*domainranking.Ranking, error) {
			return newTestRanking(), nil
		},
	}
	itemRepo := &mockItemRepo{
		countByRankingFn: func(_ context.Context, _ string) (int, error) { return 100, nil },
	}

	uc := NewAddItemUseCase(rankingRepo, itemRepo, &mockImageStore{}, &mockImageProcessor{})
	input := AddItemInput{
		RankingID: testRankingID,
		UserID:    testCreatorID,
		Name:      "Too Many",
	}

	_, err := uc.Execute(context.Background(), input)

	require.Error(t, err)
	var appErr *apperror.AppError
	require.ErrorAs(t, err, &appErr)
	require.Equal(t, apperror.MaxItemsReached, appErr.Code)
}

func TestAddItemUseCase_RankingNotFound(t *testing.T) {
	t.Parallel()

	rankingRepo := &mockRankingRepo{
		findByIDFn: func(_ context.Context, _ string) (*domainranking.Ranking, error) {
			return nil, domainranking.ErrRankingNotFound
		},
	}

	uc := NewAddItemUseCase(rankingRepo, &mockItemRepo{}, &mockImageStore{}, &mockImageProcessor{})
	input := AddItemInput{
		RankingID: "nonexistent",
		UserID:    testCreatorID,
		Name:      "Test",
	}

	_, err := uc.Execute(context.Background(), input)

	require.Error(t, err)
	var appErr *apperror.AppError
	require.ErrorAs(t, err, &appErr)
	require.Equal(t, apperror.RankingNotFound, appErr.Code)
}
