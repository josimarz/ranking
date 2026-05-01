package rating

import (
	"context"

	domainrating "github.com/josimar/ranking/backend/internal/domain/rating"
)

// Finder retrieves a user's rating.
type Finder interface {
	FindByUserAndItem(ctx context.Context, rankingID, itemID, userID string) (*domainrating.Rating, error)
}

// GetMyRatingUseCase retrieves the current user's rating for an item.
type GetMyRatingUseCase struct {
	ratings Finder
}

// NewGetMyRatingUseCase creates a new GetMyRatingUseCase.
func NewGetMyRatingUseCase(ratings Finder) *GetMyRatingUseCase {
	return &GetMyRatingUseCase{ratings: ratings}
}

// Execute retrieves the user's rating, returning nil if not rated.
func (uc *GetMyRatingUseCase) Execute(ctx context.Context, rankingID, itemID, userID string) (*domainrating.Rating, error) {
	return uc.ratings.FindByUserAndItem(ctx, rankingID, itemID, userID)
}
