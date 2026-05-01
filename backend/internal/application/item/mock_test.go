package item

import (
	"context"
	"time"

	domainitem "github.com/josimar/ranking/backend/internal/domain/item"
	domainranking "github.com/josimar/ranking/backend/internal/domain/ranking"
	domainrating "github.com/josimar/ranking/backend/internal/domain/rating"
)

const (
	testRankingID = "ranking-1"
	testItemID    = "item-1"
	testOwnerID   = "owner-1"
	testCreatorID = "creator-1"
	testOtherUser = "other-user"
)

type mockRankingRepo struct {
	findByIDFn func(ctx context.Context, id string) (*domainranking.Ranking, error)
}

func (m *mockRankingRepo) Save(_ context.Context, _ *domainranking.Ranking) error { return nil }

func (m *mockRankingRepo) FindByID(ctx context.Context, id string) (*domainranking.Ranking, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *mockRankingRepo) Update(_ context.Context, _ *domainranking.Ranking) error { return nil }
func (m *mockRankingRepo) Delete(_ context.Context, _ string) error                 { return nil }

func (m *mockRankingRepo) ListPublic(_ context.Context, _ int, _, _, _ string) ([]*domainranking.Ranking, string, error) {
	return nil, "", nil
}

func (m *mockRankingRepo) ListByOwner(_ context.Context, _ string, _ int, _ string) ([]*domainranking.Ranking, string, error) {
	return nil, "", nil
}

func (m *mockRankingRepo) SearchByName(_ context.Context, _ string, _ int, _ string) ([]*domainranking.Ranking, string, error) {
	return nil, "", nil
}

func (m *mockRankingRepo) SearchByTag(_ context.Context, _ string, _ int, _ string) ([]*domainranking.Ranking, string, error) {
	return nil, "", nil
}

func (m *mockRankingRepo) FindRecent(_ context.Context) ([]*domainranking.Ranking, error) {
	return nil, nil
}

type mockItemRepo struct {
	saveFn            func(ctx context.Context, item *domainitem.Item) error
	findByIDFn        func(ctx context.Context, rankingID, itemID string) (*domainitem.Item, error)
	updateFn          func(ctx context.Context, item *domainitem.Item) error
	deleteFn          func(ctx context.Context, rankingID, itemID string) error
	countByRankingFn  func(ctx context.Context, rankingID string) (int, error)
	findByRankingFn   func(ctx context.Context, rankingID string) ([]*domainitem.Item, error)
	deleteByRankingFn func(ctx context.Context, rankingID string) error
}

func (m *mockItemRepo) Save(ctx context.Context, item *domainitem.Item) error {
	if m.saveFn != nil {
		return m.saveFn(ctx, item)
	}
	return nil
}

func (m *mockItemRepo) FindByID(ctx context.Context, rankingID, itemID string) (*domainitem.Item, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, rankingID, itemID)
	}
	return nil, nil
}

func (m *mockItemRepo) FindByRanking(ctx context.Context, rankingID string) ([]*domainitem.Item, error) {
	if m.findByRankingFn != nil {
		return m.findByRankingFn(ctx, rankingID)
	}
	return nil, nil
}

func (m *mockItemRepo) Update(ctx context.Context, item *domainitem.Item) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, item)
	}
	return nil
}

func (m *mockItemRepo) Delete(ctx context.Context, rankingID, itemID string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, rankingID, itemID)
	}
	return nil
}

func (m *mockItemRepo) DeleteByRanking(ctx context.Context, rankingID string) error {
	if m.deleteByRankingFn != nil {
		return m.deleteByRankingFn(ctx, rankingID)
	}
	return nil
}

func (m *mockItemRepo) CountByRanking(ctx context.Context, rankingID string) (int, error) {
	if m.countByRankingFn != nil {
		return m.countByRankingFn(ctx, rankingID)
	}
	return 0, nil
}

type mockImageStore struct {
	storeFn        func(ctx context.Context, rankingID, itemID string, original, thumbnail []byte) error
	deleteFn       func(ctx context.Context, rankingID, itemID string) error
	generateURLsFn func(ctx context.Context, rankingID, itemID string) (string, string, error)
}

func (m *mockImageStore) Store(ctx context.Context, rankingID, itemID string, original, thumbnail []byte) error {
	if m.storeFn != nil {
		return m.storeFn(ctx, rankingID, itemID, original, thumbnail)
	}
	return nil
}

func (m *mockImageStore) GenerateURLs(ctx context.Context, rankingID, itemID string) (string, string, error) {
	if m.generateURLsFn != nil {
		return m.generateURLsFn(ctx, rankingID, itemID)
	}
	return "", "", nil
}

func (m *mockImageStore) Delete(ctx context.Context, rankingID, itemID string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, rankingID, itemID)
	}
	return nil
}

func (m *mockImageStore) DeleteByRanking(_ context.Context, _ string) error { return nil }

type mockImageProcessor struct {
	processFn        func(data []byte) ([]byte, []byte, error)
	validateFormatFn func(data []byte) error
}

func (m *mockImageProcessor) Process(data []byte) ([]byte, []byte, error) {
	if m.processFn != nil {
		return m.processFn(data)
	}
	return []byte("original"), []byte("thumb"), nil
}

func (m *mockImageProcessor) ValidateFormat(data []byte) error {
	if m.validateFormatFn != nil {
		return m.validateFormatFn(data)
	}
	return nil
}

type mockRatingRepo struct {
	deleteByItemFn      func(ctx context.Context, rankingID, itemID string) error
	findAllByItemFn     func(ctx context.Context, rankingID, itemID string) ([]*domainrating.Rating, error)
	findByUserAndItemFn func(ctx context.Context, rankingID, itemID, userID string) (*domainrating.Rating, error)
}

func (m *mockRatingRepo) Save(_ context.Context, _ *domainrating.Rating) error { return nil }

func (m *mockRatingRepo) FindByUserAndItem(ctx context.Context, rankingID, itemID, userID string) (*domainrating.Rating, error) {
	if m.findByUserAndItemFn != nil {
		return m.findByUserAndItemFn(ctx, rankingID, itemID, userID)
	}
	return nil, nil
}

func (m *mockRatingRepo) FindAllByItem(ctx context.Context, rankingID, itemID string) ([]*domainrating.Rating, error) {
	if m.findAllByItemFn != nil {
		return m.findAllByItemFn(ctx, rankingID, itemID)
	}
	return nil, nil
}

func (m *mockRatingRepo) DeleteByItem(ctx context.Context, rankingID, itemID string) error {
	if m.deleteByItemFn != nil {
		return m.deleteByItemFn(ctx, rankingID, itemID)
	}
	return nil
}

func (m *mockRatingRepo) DeleteByRanking(_ context.Context, _ string, _ []string) error {
	return nil
}

func newTestRanking() *domainranking.Ranking {
	return &domainranking.Ranking{
		ID:          testRankingID,
		Name:        "Video Games",
		Visibility:  domainranking.Public,
		Attributes:  []domainranking.Attribute{{ID: "a1", Name: "Graphics", Active: true}},
		OwnerUserID: testOwnerID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func newTestItem() *domainitem.Item {
	now := time.Now()
	return &domainitem.Item{
		ID:        testItemID,
		Name:      "PlayStation 5",
		ImageKey:  "",
		RankingID: testRankingID,
		CreatedBy: testCreatorID,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
