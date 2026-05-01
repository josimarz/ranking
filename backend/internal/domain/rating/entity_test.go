package rating_test

import (
	"errors"
	"testing"
	"time"

	"github.com/josimar/ranking/backend/internal/domain/rating"
	"github.com/josimar/ranking/backend/pkg/apperror"
	"github.com/stretchr/testify/require"
)

const (
	testAttr1 = "attr1"
	testAttr2 = "attr2"
	testAttr3 = "attr3"
)

func TestNewRating_Success(t *testing.T) {
	t.Parallel()

	scores := map[string]int{testAttr1: 80, testAttr2: 90}
	activeIDs := []string{testAttr1, testAttr2}
	before := time.Now()

	r, err := rating.NewRating("ranking-1", "item-1", "user-1", scores, activeIDs)

	require.NoError(t, err)
	require.Equal(t, "ranking-1", r.RankingID)
	require.Equal(t, "item-1", r.ItemID)
	require.Equal(t, "user-1", r.UserID)
	require.Equal(t, scores, r.Scores)
	require.False(t, r.CreatedAt.Before(before))
	require.False(t, r.UpdatedAt.Before(before))
}

func TestNewRating_ErrIncompleteRatings(t *testing.T) {
	t.Parallel()

	scores := map[string]int{testAttr1: 80}
	activeIDs := []string{testAttr1, testAttr2}

	r, err := rating.NewRating("ranking-1", "item-1", "user-1", scores, activeIDs)

	require.Nil(t, r)
	require.Error(t, err)
	var appErr *apperror.AppError
	require.True(t, errors.As(err, &appErr))
	require.Equal(t, apperror.IncompleteRatings, appErr.Code)
}

func TestNewRating_ErrInvalidScore_TooHigh(t *testing.T) {
	t.Parallel()

	scores := map[string]int{testAttr1: 101}
	activeIDs := []string{testAttr1}

	r, err := rating.NewRating("ranking-1", "item-1", "user-1", scores, activeIDs)

	require.Nil(t, r)
	require.Error(t, err)
	var appErr *apperror.AppError
	require.True(t, errors.As(err, &appErr))
	require.Equal(t, apperror.InvalidScore, appErr.Code)
}

func TestNewRating_ErrInvalidScore_Negative(t *testing.T) {
	t.Parallel()

	scores := map[string]int{testAttr1: -1}
	activeIDs := []string{testAttr1}

	r, err := rating.NewRating("ranking-1", "item-1", "user-1", scores, activeIDs)

	require.Nil(t, r)
	require.Error(t, err)
	var appErr *apperror.AppError
	require.True(t, errors.As(err, &appErr))
	require.Equal(t, apperror.InvalidScore, appErr.Code)
}

func TestNewRating_BoundaryScores(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		score int
	}{
		{"score 0 is valid", 0},
		{"score 100 is valid", 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			scores := map[string]int{testAttr1: tt.score}
			activeIDs := []string{testAttr1}

			r, err := rating.NewRating("ranking-1", "item-1", "user-1", scores, activeIDs)

			require.NoError(t, err)
			require.NotNil(t, r)
		})
	}
}

func TestCalculateOverall_WithActiveAttributes(t *testing.T) {
	t.Parallel()

	scores := map[string]int{testAttr1: 80, testAttr2: 90, testAttr3: 70}
	activeIDs := []string{testAttr1, testAttr2, testAttr3}

	r, err := rating.NewRating("ranking-1", "item-1", "user-1", scores, activeIDs)
	require.NoError(t, err)

	overall := r.CalculateOverall(activeIDs)
	require.InDelta(t, 80.0, overall, 0.001)
}

func TestCalculateOverall_SubsetOfAttributes(t *testing.T) {
	t.Parallel()

	scores := map[string]int{testAttr1: 80, testAttr2: 90, testAttr3: 70}
	activeIDs := []string{testAttr1, testAttr2, testAttr3}

	r, err := rating.NewRating("ranking-1", "item-1", "user-1", scores, activeIDs)
	require.NoError(t, err)

	overall := r.CalculateOverall([]string{testAttr1, testAttr3})
	require.InDelta(t, 75.0, overall, 0.001)
}

func TestCalculateOverall_NoActiveAttributes(t *testing.T) {
	t.Parallel()

	scores := map[string]int{testAttr1: 80}
	activeIDs := []string{testAttr1}

	r, err := rating.NewRating("ranking-1", "item-1", "user-1", scores, activeIDs)
	require.NoError(t, err)

	overall := r.CalculateOverall([]string{})
	require.Equal(t, 0.0, overall)
}

func TestCalculateOverall_NilActiveAttributes(t *testing.T) {
	t.Parallel()

	scores := map[string]int{testAttr1: 80}
	activeIDs := []string{testAttr1}

	r, err := rating.NewRating("ranking-1", "item-1", "user-1", scores, activeIDs)
	require.NoError(t, err)

	overall := r.CalculateOverall(nil)
	require.Equal(t, 0.0, overall)
}
