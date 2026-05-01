package rating

import (
	"context"

	domainitem "github.com/josimar/ranking/backend/internal/domain/item"
	domainranking "github.com/josimar/ranking/backend/internal/domain/ranking"
	domainrating "github.com/josimar/ranking/backend/internal/domain/rating"
)

// RankingFinder finds a ranking by ID.
type RankingFinder interface {
	FindByID(ctx context.Context, id string) (*domainranking.Ranking, error)
}

// ItemFinder finds an item by ranking and item ID.
type ItemFinder interface {
	FindByID(ctx context.Context, rankingID, itemID string) (*domainitem.Item, error)
}

// Saver persists and retrieves ratings.
type Saver interface {
	Save(ctx context.Context, r *domainrating.Rating) error
	FindByUserAndItem(ctx context.Context, rankingID, itemID, userID string) (*domainrating.Rating, error)
}

// SubmitRatingUseCase handles submitting a user's rating for an item.
type SubmitRatingUseCase struct {
	rankings RankingFinder
	items    ItemFinder
	ratings  Saver
}

// NewSubmitRatingUseCase creates a new SubmitRatingUseCase.
func NewSubmitRatingUseCase(rankings RankingFinder, items ItemFinder, ratings Saver) *SubmitRatingUseCase {
	return &SubmitRatingUseCase{rankings: rankings, items: items, ratings: ratings}
}

// Execute validates and saves a user's rating for an item in a ranking.
func (uc *SubmitRatingUseCase) Execute(ctx context.Context, rankingID, itemID, userID string, scores map[string]int) (*domainrating.Rating, error) {
	rnk, err := uc.rankings.FindByID(ctx, rankingID)
	if err != nil {
		return nil, err
	}

	if _, err := uc.items.FindByID(ctx, rankingID, itemID); err != nil {
		return nil, err
	}

	activeAttrs := rnk.ActiveAttributes()
	attrIDs := make([]string, len(activeAttrs))
	for i, a := range activeAttrs {
		attrIDs[i] = a.ID
	}

	r, err := domainrating.NewRating(rankingID, itemID, userID, scores, attrIDs)
	if err != nil {
		return nil, err
	}

	if err := uc.ratings.Save(ctx, r); err != nil {
		return nil, err
	}

	return r, nil
}
