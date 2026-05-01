package rating

import (
	"fmt"
	"time"

	"github.com/josimar/ranking/backend/pkg/apperror"
)

// Rating represents a user's scores for an item in a ranking.
type Rating struct {
	RankingID string
	ItemID    string
	UserID    string
	Scores    map[string]int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewRating creates a Rating after validating that all activeAttributeIDs have scores in range 0–100.
func NewRating(rankingID, itemID, userID string, scores map[string]int, activeAttributeIDs []string) (*Rating, error) {
	for _, id := range activeAttributeIDs {
		if _, ok := scores[id]; !ok {
			return nil, apperror.NewBadRequestError(apperror.IncompleteRatings, fmt.Sprintf("missing score for attribute %s", id))
		}
	}

	for _, id := range activeAttributeIDs {
		s := scores[id]
		if s < 0 || s > 100 {
			return nil, apperror.NewBadRequestError(apperror.InvalidScore, fmt.Sprintf("score for attribute %s must be 0–100, got %d", id, s))
		}
	}

	now := time.Now()
	return &Rating{
		RankingID: rankingID,
		ItemID:    itemID,
		UserID:    userID,
		Scores:    scores,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// CalculateOverall returns the arithmetic mean of scores for the given active attribute IDs.
func (r *Rating) CalculateOverall(activeAttributeIDs []string) float64 {
	if len(activeAttributeIDs) == 0 {
		return 0
	}

	sum := 0
	for _, id := range activeAttributeIDs {
		sum += r.Scores[id]
	}

	return float64(sum) / float64(len(activeAttributeIDs))
}
