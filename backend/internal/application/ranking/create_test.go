package ranking

import (
	"context"
	"errors"
	"testing"

	domain "github.com/josimar/ranking/backend/internal/domain/ranking"
	"github.com/stretchr/testify/require"
)

func TestCreateRankingUseCase_Success(t *testing.T) {
	t.Parallel()

	var saved *domain.Ranking
	repo := &mockRankingRepo{
		saveFn: func(_ context.Context, r *domain.Ranking) error {
			saved = r
			return nil
		},
	}

	uc := NewCreateRankingUseCase(repo)
	input := CreateRankingInput{
		Name:        "Video Games",
		Description: "Best consoles",
		Visibility:  testPublic,
		Tags:        []string{testGamesTag},
		Attributes:  []AttributeInput{{Name: "Graphics", Description: "Visual quality"}},
		UserID:      testUserID,
	}

	out, err := uc.Execute(context.Background(), input)

	require.NoError(t, err)
	require.NotEmpty(t, out.ID)
	require.Equal(t, "Video Games", out.Name)
	require.Equal(t, "Best consoles", out.Description)
	require.Equal(t, testPublic, out.Visibility)
	require.Equal(t, []string{testGamesTag}, out.Tags)
	require.Len(t, out.Attributes, 1)
	require.Equal(t, "Graphics", out.Attributes[0].Name)
	require.Equal(t, testUserID, out.OwnerUserID)
	require.NotNil(t, saved)
}

func TestCreateRankingUseCase_ValidationError(t *testing.T) {
	t.Parallel()

	repo := &mockRankingRepo{}
	uc := NewCreateRankingUseCase(repo)

	input := CreateRankingInput{
		Name:       "",
		Visibility: testPublic,
		Attributes: []AttributeInput{{Name: "A"}},
		UserID:     testUserID,
	}

	_, err := uc.Execute(context.Background(), input)
	require.Error(t, err)
	require.Contains(t, err.Error(), "name is required")
}

func TestCreateRankingUseCase_RepoError(t *testing.T) {
	t.Parallel()

	repo := &mockRankingRepo{
		saveFn: func(_ context.Context, _ *domain.Ranking) error {
			return errors.New("db error")
		},
	}

	uc := NewCreateRankingUseCase(repo)
	input := CreateRankingInput{
		Name:       "Test",
		Visibility: testPublic,
		Attributes: []AttributeInput{{Name: "A"}},
		UserID:     testUserID,
	}

	_, err := uc.Execute(context.Background(), input)
	require.Error(t, err)
	require.Contains(t, err.Error(), "db error")
}
