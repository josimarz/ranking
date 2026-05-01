package rating

import "context"

// Repository defines persistence operations for ratings.
type Repository interface {
	Save(ctx context.Context, r *Rating) error
	FindByUserAndItem(ctx context.Context, rankingID, itemID, userID string) (*Rating, error)
	FindAllByItem(ctx context.Context, rankingID, itemID string) ([]*Rating, error)
	DeleteByItem(ctx context.Context, rankingID, itemID string) error
	DeleteByRanking(ctx context.Context, rankingID string, itemIDs []string) error
}
