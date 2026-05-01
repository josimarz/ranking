package rating

import (
	"context"

	domainitem "github.com/josimar/ranking/backend/internal/domain/item"
	domainranking "github.com/josimar/ranking/backend/internal/domain/ranking"
	domainrating "github.com/josimar/ranking/backend/internal/domain/rating"
)

const (
	testRankingID = "ranking-1"
	testItemID    = "item-1"
	testUserID    = "user-1"
	testAttrID1   = "a1"
	testAttrID2   = "a2"
	testGraphics  = "Graphics"
	testSound     = "Sound"
)

type mockRankingRepo struct {
	findByIDFn func(ctx context.Context, id string) (*domainranking.Ranking, error)
}

func (m *mockRankingRepo) FindByID(ctx context.Context, id string) (*domainranking.Ranking, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, id)
	}
	return nil, nil
}

type mockItemRepo struct {
	findByIDFn func(ctx context.Context, rankingID, itemID string) (*domainitem.Item, error)
}

func (m *mockItemRepo) FindByID(ctx context.Context, rankingID, itemID string) (*domainitem.Item, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, rankingID, itemID)
	}
	return nil, nil
}

type mockRatingRepo struct {
	saveFn              func(ctx context.Context, r *domainrating.Rating) error
	findByUserAndItemFn func(ctx context.Context, rankingID, itemID, userID string) (*domainrating.Rating, error)
}

func (m *mockRatingRepo) Save(ctx context.Context, r *domainrating.Rating) error {
	if m.saveFn != nil {
		return m.saveFn(ctx, r)
	}
	return nil
}

func (m *mockRatingRepo) FindByUserAndItem(ctx context.Context, rankingID, itemID, userID string) (*domainrating.Rating, error) {
	if m.findByUserAndItemFn != nil {
		return m.findByUserAndItemFn(ctx, rankingID, itemID, userID)
	}
	return nil, nil
}
