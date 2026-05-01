package ranking

import (
	"context"
	"testing"

	domain "github.com/josimar/ranking/backend/internal/domain/ranking"
	"github.com/stretchr/testify/require"
)

func TestUpdateRankingUseCase_Success(t *testing.T) {
	t.Parallel()

	r := newFullTestRanking()
	var updated *domain.Ranking
	repo := &mockRankingRepo{
		findByIDFn: func(_ context.Context, _ string) (*domain.Ranking, error) {
			return r, nil
		},
		updateFn: func(_ context.Context, ranking *domain.Ranking) error {
			updated = ranking
			return nil
		},
	}

	uc := NewUpdateRankingUseCase(repo)
	input := UpdateRankingInput{
		RankingID:   testRankingID,
		UserID:      testOwnerID,
		Name:        "Updated Name",
		Description: "Updated desc",
		Visibility:  "private",
		Tags:        []string{"updated"},
		Attributes:  []AttributeInput{{Name: "Sound", Description: "Audio"}},
	}

	out, err := uc.Execute(context.Background(), input)

	require.NoError(t, err)
	require.Equal(t, "Updated Name", out.Name)
	require.Equal(t, "Updated desc", out.Description)
	require.Equal(t, "private", out.Visibility)
	require.NotNil(t, updated)
}

func TestUpdateRankingUseCase_NotOwner(t *testing.T) {
	t.Parallel()

	r := newFullTestRanking()
	repo := &mockRankingRepo{
		findByIDFn: func(_ context.Context, _ string) (*domain.Ranking, error) {
			return r, nil
		},
	}

	uc := NewUpdateRankingUseCase(repo)
	input := UpdateRankingInput{
		RankingID:  testRankingID,
		UserID:     "other-user",
		Name:       "X",
		Visibility: testPublic,
		Attributes: []AttributeInput{{Name: "A"}},
	}

	_, err := uc.Execute(context.Background(), input)
	require.Error(t, err)
	require.ErrorIs(t, err, domain.ErrNotOwner)
}

func TestUpdateRankingUseCase_NotFound(t *testing.T) {
	t.Parallel()

	repo := &mockRankingRepo{
		findByIDFn: func(_ context.Context, _ string) (*domain.Ranking, error) {
			return nil, domain.ErrRankingNotFound
		},
	}

	uc := NewUpdateRankingUseCase(repo)
	input := UpdateRankingInput{
		RankingID:  "nonexistent",
		UserID:     "user-1",
		Name:       "X",
		Visibility: testPublic,
		Attributes: []AttributeInput{{Name: "A"}},
	}

	_, err := uc.Execute(context.Background(), input)
	require.Error(t, err)
	require.ErrorIs(t, err, domain.ErrRankingNotFound)
}

func TestUpdateRankingUseCase_ValidationError(t *testing.T) {
	t.Parallel()

	r := newFullTestRanking()
	repo := &mockRankingRepo{
		findByIDFn: func(_ context.Context, _ string) (*domain.Ranking, error) {
			return r, nil
		},
	}

	uc := NewUpdateRankingUseCase(repo)
	input := UpdateRankingInput{
		RankingID:  testRankingID,
		UserID:     testOwnerID,
		Name:       "",
		Visibility: testPublic,
		Attributes: []AttributeInput{{Name: "A"}},
	}

	_, err := uc.Execute(context.Background(), input)
	require.Error(t, err)
	require.Contains(t, err.Error(), "name is required")
}
