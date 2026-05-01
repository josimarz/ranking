package ranking

import (
	"context"
	"time"

	domainitem "github.com/josimar/ranking/backend/internal/domain/item"
	domain "github.com/josimar/ranking/backend/internal/domain/ranking"
	domainrating "github.com/josimar/ranking/backend/internal/domain/rating"
)

const (
	testRankingID = "ranking-1"
	testOwnerID   = "owner-1"
	testUserID    = "user-123"
	testPublic    = "public"
	testGamesTag  = "games"
)

type mockRankingRepo struct {
	saveFn         func(ctx context.Context, r *domain.Ranking) error
	findByIDFn     func(ctx context.Context, id string) (*domain.Ranking, error)
	updateFn       func(ctx context.Context, r *domain.Ranking) error
	deleteFn       func(ctx context.Context, id string) error
	listPublicFn   func(ctx context.Context, limit int, cursor, sortField, sortDir string) ([]*domain.Ranking, string, error)
	listByOwnerFn  func(ctx context.Context, userID string, limit int, cursor string) ([]*domain.Ranking, string, error)
	searchByNameFn func(ctx context.Context, term string, limit int, cursor string) ([]*domain.Ranking, string, error)
	searchByTagFn  func(ctx context.Context, tag string, limit int, cursor string) ([]*domain.Ranking, string, error)
	findRecentFn   func(ctx context.Context) ([]*domain.Ranking, error)
}

func (m *mockRankingRepo) Save(ctx context.Context, r *domain.Ranking) error {
	if m.saveFn != nil {
		return m.saveFn(ctx, r)
	}
	return nil
}

func (m *mockRankingRepo) FindByID(ctx context.Context, id string) (*domain.Ranking, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *mockRankingRepo) Update(ctx context.Context, r *domain.Ranking) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, r)
	}
	return nil
}

func (m *mockRankingRepo) Delete(ctx context.Context, id string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return nil
}

func (m *mockRankingRepo) ListPublic(ctx context.Context, limit int, cursor, sortField, sortDir string) ([]*domain.Ranking, string, error) {
	if m.listPublicFn != nil {
		return m.listPublicFn(ctx, limit, cursor, sortField, sortDir)
	}
	return nil, "", nil
}

func (m *mockRankingRepo) ListByOwner(ctx context.Context, userID string, limit int, cursor string) ([]*domain.Ranking, string, error) {
	if m.listByOwnerFn != nil {
		return m.listByOwnerFn(ctx, userID, limit, cursor)
	}
	return nil, "", nil
}

func (m *mockRankingRepo) SearchByName(ctx context.Context, term string, limit int, cursor string) ([]*domain.Ranking, string, error) {
	if m.searchByNameFn != nil {
		return m.searchByNameFn(ctx, term, limit, cursor)
	}
	return nil, "", nil
}

func (m *mockRankingRepo) SearchByTag(ctx context.Context, tag string, limit int, cursor string) ([]*domain.Ranking, string, error) {
	if m.searchByTagFn != nil {
		return m.searchByTagFn(ctx, tag, limit, cursor)
	}
	return nil, "", nil
}

func (m *mockRankingRepo) FindRecent(ctx context.Context) ([]*domain.Ranking, error) {
	if m.findRecentFn != nil {
		return m.findRecentFn(ctx)
	}
	return nil, nil
}

func newFullTestRanking() *domain.Ranking {
	return &domain.Ranking{
		ID:          testRankingID,
		Name:        "Video Games",
		Description: "Best consoles",
		Visibility:  domain.Public,
		Tags:        []string{testGamesTag},
		Attributes:  []domain.Attribute{{ID: "a1", Name: "Graphics", Active: true}},
		OwnerUserID: testOwnerID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

type mockItemRepo struct {
	findByRankingFn   func(ctx context.Context, rankingID string) ([]*domainitem.Item, error)
	deleteByRankingFn func(ctx context.Context, rankingID string) error
}

func (m *mockItemRepo) Save(_ context.Context, _ *domainitem.Item) error { return nil }

func (m *mockItemRepo) FindByID(_ context.Context, _ string, _ string) (*domainitem.Item, error) {
	return nil, nil
}

func (m *mockItemRepo) FindByRanking(ctx context.Context, rankingID string) ([]*domainitem.Item, error) {
	if m.findByRankingFn != nil {
		return m.findByRankingFn(ctx, rankingID)
	}
	return nil, nil
}

func (m *mockItemRepo) Update(_ context.Context, _ *domainitem.Item) error { return nil }
func (m *mockItemRepo) Delete(_ context.Context, _ string, _ string) error { return nil }

func (m *mockItemRepo) DeleteByRanking(ctx context.Context, rankingID string) error {
	if m.deleteByRankingFn != nil {
		return m.deleteByRankingFn(ctx, rankingID)
	}
	return nil
}

func (m *mockItemRepo) CountByRanking(_ context.Context, _ string) (int, error) { return 0, nil }

type mockRatingRepo struct {
	deleteByRankingFn func(ctx context.Context, rankingID string, itemIDs []string) error
}

func (m *mockRatingRepo) Save(_ context.Context, _ *domainrating.Rating) error { return nil }

func (m *mockRatingRepo) FindByUserAndItem(_ context.Context, _ string, _ string, _ string) (*domainrating.Rating, error) {
	return nil, nil
}

func (m *mockRatingRepo) FindAllByItem(_ context.Context, _ string, _ string) ([]*domainrating.Rating, error) {
	return nil, nil
}

func (m *mockRatingRepo) DeleteByItem(_ context.Context, _ string, _ string) error { return nil }

func (m *mockRatingRepo) DeleteByRanking(ctx context.Context, rankingID string, itemIDs []string) error {
	if m.deleteByRankingFn != nil {
		return m.deleteByRankingFn(ctx, rankingID, itemIDs)
	}
	return nil
}

type mockImageStore struct {
	deleteByRankingFn func(ctx context.Context, rankingID string) error
}

func (m *mockImageStore) Store(_ context.Context, _ string, _ string, _ []byte, _ []byte) error {
	return nil
}

func (m *mockImageStore) GenerateURLs(_ context.Context, _ string, _ string) (string, string, error) {
	return "", "", nil
}

func (m *mockImageStore) Delete(_ context.Context, _ string, _ string) error { return nil }

func (m *mockImageStore) DeleteByRanking(ctx context.Context, rankingID string) error {
	if m.deleteByRankingFn != nil {
		return m.deleteByRankingFn(ctx, rankingID)
	}
	return nil
}
