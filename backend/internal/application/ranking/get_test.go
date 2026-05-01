package ranking

import (
	"context"
	"testing"

	domain "github.com/josimar/ranking/backend/internal/domain/ranking"
	"github.com/stretchr/testify/require"
)

func TestGetRankingUseCase_Found(t *testing.T) {
	t.Parallel()

	r := newFullTestRanking()
	repo := &mockRankingRepo{
		findByIDFn: func(_ context.Context, id string) (*domain.Ranking, error) {
			if id == testRankingID {
				return r, nil
			}
			return nil, domain.ErrRankingNotFound
		},
	}

	uc := NewGetRankingUseCase(repo)
	out, err := uc.Execute(context.Background(), testRankingID, "")

	require.NoError(t, err)
	require.Equal(t, testRankingID, out.ID)
	require.Equal(t, "Video Games", out.Name)
	require.False(t, out.IsOwner)
}

func TestGetRankingUseCase_NotFound(t *testing.T) {
	t.Parallel()

	repo := &mockRankingRepo{
		findByIDFn: func(_ context.Context, _ string) (*domain.Ranking, error) {
			return nil, domain.ErrRankingNotFound
		},
	}

	uc := NewGetRankingUseCase(repo)
	_, err := uc.Execute(context.Background(), "nonexistent", "")

	require.Error(t, err)
	require.ErrorIs(t, err, domain.ErrRankingNotFound)
}

func TestGetRankingUseCase_IsOwnerTrue(t *testing.T) {
	t.Parallel()

	r := newFullTestRanking()
	repo := &mockRankingRepo{
		findByIDFn: func(_ context.Context, _ string) (*domain.Ranking, error) {
			return r, nil
		},
	}

	uc := NewGetRankingUseCase(repo)
	out, err := uc.Execute(context.Background(), testRankingID, testOwnerID)

	require.NoError(t, err)
	require.True(t, out.IsOwner)
}

func TestGetRankingUseCase_IsOwnerFalse(t *testing.T) {
	t.Parallel()

	r := newFullTestRanking()
	repo := &mockRankingRepo{
		findByIDFn: func(_ context.Context, _ string) (*domain.Ranking, error) {
			return r, nil
		},
	}

	uc := NewGetRankingUseCase(repo)
	out, err := uc.Execute(context.Background(), testRankingID, "other-user")

	require.NoError(t, err)
	require.False(t, out.IsOwner)
}
