package ranking

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	domain "github.com/josimar/ranking/backend/internal/domain/ranking"
	"github.com/josimar/ranking/backend/pkg/apperror"
	"github.com/stretchr/testify/require"
)

func newTestRankingWithName(id, name string) *domain.Ranking {
	return &domain.Ranking{ID: id, Name: name}
}

// --- Task 7.5: ListPublicRankingsUseCase ---

func TestListPublicRankings_DefaultSort(t *testing.T) {
	t.Parallel()
	rankings := []*domain.Ranking{newTestRankingWithName("1", "Alpha")}
	repo := &mockRankingRepo{
		listPublicFn: func(_ context.Context, limit int, cursor, sortField, sortDir string) ([]*domain.Ranking, string, error) {
			require.Equal(t, 20, limit)
			require.Empty(t, cursor)
			require.Empty(t, sortField)
			require.Empty(t, sortDir)
			return rankings, "", nil
		},
	}
	uc := NewListPublicRankingsUseCase(repo)
	out, err := uc.Execute(context.Background(), 0, "", "")
	require.NoError(t, err)
	require.Len(t, out.Rankings, 1)
	require.Empty(t, out.NextCursor)
}

func TestListPublicRankings_CustomSort(t *testing.T) {
	t.Parallel()
	repo := &mockRankingRepo{
		listPublicFn: func(_ context.Context, _ int, _, sortField, sortDir string) ([]*domain.Ranking, string, error) {
			require.Equal(t, "createdAt", sortField)
			require.Equal(t, "desc", sortDir)
			return nil, "", nil
		},
	}
	uc := NewListPublicRankingsUseCase(repo)
	out, err := uc.Execute(context.Background(), 5, "", "-createdAt")
	require.NoError(t, err)
	require.Empty(t, out.Rankings)
}

func TestListPublicRankings_InvalidSort(t *testing.T) {
	t.Parallel()
	repo := &mockRankingRepo{}
	uc := NewListPublicRankingsUseCase(repo)
	_, err := uc.Execute(context.Background(), 10, "", "invalid")
	require.Error(t, err)
	var appErr *apperror.AppError
	require.ErrorAs(t, err, &appErr)
	require.Equal(t, apperror.InvalidSortField, appErr.Code)
}

func TestListPublicRankings_Pagination(t *testing.T) {
	t.Parallel()
	repo := &mockRankingRepo{
		listPublicFn: func(_ context.Context, limit int, cursor, _, _ string) ([]*domain.Ranking, string, error) {
			require.Equal(t, 20, limit)
			require.Equal(t, "cursor123", cursor)
			return []*domain.Ranking{newTestRankingWithName("2", "Beta")}, "nextCur", nil
		},
	}
	uc := NewListPublicRankingsUseCase(repo)
	out, err := uc.Execute(context.Background(), 0, "cursor123", "")
	require.NoError(t, err)
	require.Len(t, out.Rankings, 1)
	require.Equal(t, "nextCur", out.NextCursor)
}

func TestListPublicRankings_ClampLimit(t *testing.T) {
	t.Parallel()
	repo := &mockRankingRepo{
		listPublicFn: func(_ context.Context, limit int, _, _, _ string) ([]*domain.Ranking, string, error) {
			require.Equal(t, 50, limit)
			return nil, "", nil
		},
	}
	uc := NewListPublicRankingsUseCase(repo)
	_, err := uc.Execute(context.Background(), 999, "", "")
	require.NoError(t, err)
}

// --- Task 7.6: GetRecentRankingsUseCase ---

func TestGetRecentRankings_ReturnsUpTo10(t *testing.T) {
	t.Parallel()
	rankings := make([]*domain.Ranking, 10)
	for i := range rankings {
		rankings[i] = newTestRankingWithName(strconv.Itoa(i), fmt.Sprintf("R%d", i))
	}
	repo := &mockRankingRepo{
		findRecentFn: func(_ context.Context) ([]*domain.Ranking, error) {
			return rankings, nil
		},
	}
	uc := NewGetRecentRankingsUseCase(repo)
	out, err := uc.Execute(context.Background())
	require.NoError(t, err)
	require.Len(t, out, 10)
}

func TestGetRecentRankings_EmptyResult(t *testing.T) {
	t.Parallel()
	repo := &mockRankingRepo{
		findRecentFn: func(_ context.Context) ([]*domain.Ranking, error) {
			return nil, nil
		},
	}
	uc := NewGetRecentRankingsUseCase(repo)
	out, err := uc.Execute(context.Background())
	require.NoError(t, err)
	require.Empty(t, out)
}

// --- Task 7.7: SearchRankingsUseCase ---

func TestSearchRankings_ByName(t *testing.T) {
	t.Parallel()
	repo := &mockRankingRepo{
		searchByNameFn: func(_ context.Context, term string, limit int, _ string) ([]*domain.Ranking, string, error) {
			require.Equal(t, "video games", term)
			require.Equal(t, 20, limit)
			return []*domain.Ranking{newTestRankingWithName("1", "Video Games")}, "", nil
		},
	}
	uc := NewSearchRankingsUseCase(repo)
	out, err := uc.Execute(context.Background(), "Vídeo Games", "", 0, "")
	require.NoError(t, err)
	require.Len(t, out.Rankings, 1)
}

func TestSearchRankings_ByTag(t *testing.T) {
	t.Parallel()
	repo := &mockRankingRepo{
		searchByTagFn: func(_ context.Context, tag string, _ int, _ string) ([]*domain.Ranking, string, error) {
			require.Equal(t, "games", tag)
			return []*domain.Ranking{newTestRankingWithName("1", "Games")}, "next", nil
		},
	}
	uc := NewSearchRankingsUseCase(repo)
	out, err := uc.Execute(context.Background(), "", "games", 0, "")
	require.NoError(t, err)
	require.Len(t, out.Rankings, 1)
	require.Equal(t, "next", out.NextCursor)
}

func TestSearchRankings_MissingParams(t *testing.T) {
	t.Parallel()
	repo := &mockRankingRepo{}
	uc := NewSearchRankingsUseCase(repo)
	_, err := uc.Execute(context.Background(), "", "", 0, "")
	require.Error(t, err)
	var appErr *apperror.AppError
	require.ErrorAs(t, err, &appErr)
	require.Equal(t, apperror.ValidationError, appErr.Code)
}

func TestSearchRankings_BothQAndTag_UsesSearchByName(t *testing.T) {
	t.Parallel()
	repo := &mockRankingRepo{
		searchByNameFn: func(_ context.Context, term string, _ int, _ string) ([]*domain.Ranking, string, error) {
			require.Equal(t, "games", term)
			return []*domain.Ranking{newTestRankingWithName("1", "Games")}, "", nil
		},
	}
	uc := NewSearchRankingsUseCase(repo)
	out, err := uc.Execute(context.Background(), "Games", "tag", 0, "")
	require.NoError(t, err)
	require.Len(t, out.Rankings, 1)
}

// --- Task 7.8: ListMyRankingsUseCase ---

func TestListMyRankings_WithResults(t *testing.T) {
	t.Parallel()
	repo := &mockRankingRepo{
		listByOwnerFn: func(_ context.Context, userID string, limit int, _ string) ([]*domain.Ranking, string, error) {
			require.Equal(t, "user-1", userID)
			require.Equal(t, 20, limit)
			return []*domain.Ranking{newTestRankingWithName("1", "My Ranking")}, "", nil
		},
	}
	uc := NewListMyRankingsUseCase(repo)
	out, err := uc.Execute(context.Background(), "user-1", 0, "")
	require.NoError(t, err)
	require.Len(t, out.Rankings, 1)
}

func TestListMyRankings_Empty(t *testing.T) {
	t.Parallel()
	repo := &mockRankingRepo{
		listByOwnerFn: func(_ context.Context, _ string, _ int, _ string) ([]*domain.Ranking, string, error) {
			return nil, "", nil
		},
	}
	uc := NewListMyRankingsUseCase(repo)
	out, err := uc.Execute(context.Background(), "user-2", 0, "")
	require.NoError(t, err)
	require.Empty(t, out.Rankings)
}

func TestListMyRankings_Pagination(t *testing.T) {
	t.Parallel()
	repo := &mockRankingRepo{
		listByOwnerFn: func(_ context.Context, _ string, limit int, cursor string) ([]*domain.Ranking, string, error) {
			require.Equal(t, 10, limit)
			require.Equal(t, "cur1", cursor)
			return []*domain.Ranking{newTestRankingWithName("3", "Third")}, "cur2", nil
		},
	}
	uc := NewListMyRankingsUseCase(repo)
	out, err := uc.Execute(context.Background(), "user-3", 10, "cur1")
	require.NoError(t, err)
	require.Len(t, out.Rankings, 1)
	require.Equal(t, "cur2", out.NextCursor)
}
