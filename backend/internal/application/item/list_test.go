package item

import (
	"context"
	"fmt"
	"testing"
	"time"

	domainitem "github.com/josimar/ranking/backend/internal/domain/item"
	domainranking "github.com/josimar/ranking/backend/internal/domain/ranking"
	domainrating "github.com/josimar/ranking/backend/internal/domain/rating"
	"github.com/josimar/ranking/backend/pkg/apperror"
	"github.com/stretchr/testify/require"
)

const (
	testItemID2   = "item-2"
	testAttrID1   = "a1"
	testAttrID2   = "a2"
	testAttrGfx   = "Graphics"
	testAttrSound = "Sound"
	testNamePS5   = "PlayStation 5"
	testNameXbox  = "Xbox Series X"
)

func newListRanking() *domainranking.Ranking {
	return &domainranking.Ranking{
		ID:         testRankingID,
		Name:       "Video Games",
		Visibility: domainranking.Public,
		Attributes: []domainranking.Attribute{
			{ID: testAttrID1, Name: testAttrGfx, Active: true},
			{ID: testAttrID2, Name: testAttrSound, Active: true},
			{ID: "a3", Name: "Old", Active: false},
		},
		OwnerUserID: testOwnerID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func newListItems() []*domainitem.Item {
	now := time.Now()
	return []*domainitem.Item{
		{ID: testItemID, Name: testNamePS5, ImageKey: "items/ranking-1/item-1", RankingID: testRankingID, CreatedBy: testCreatorID, CreatedAt: now, UpdatedAt: now},
		{ID: testItemID2, Name: testNameXbox, RankingID: testRankingID, CreatedBy: testCreatorID, CreatedAt: now, UpdatedAt: now},
	}
}

func TestListItemsUseCase_AvgModeWithRatings(t *testing.T) {
	t.Parallel()

	rankingRepo := &mockRankingRepo{
		findByIDFn: func(_ context.Context, _ string) (*domainranking.Ranking, error) {
			return newListRanking(), nil
		},
	}
	itemRepo := &mockItemRepo{
		findByRankingFn: func(_ context.Context, _ string) ([]*domainitem.Item, error) {
			return newListItems(), nil
		},
	}
	ratingRepo := &mockRatingRepo{
		findAllByItemFn: func(_ context.Context, _, itemID string) ([]*domainrating.Rating, error) {
			if itemID == testItemID {
				return []*domainrating.Rating{
					{RankingID: testRankingID, ItemID: testItemID, UserID: "u1", Scores: map[string]int{testAttrID1: 80, testAttrID2: 60}},
					{RankingID: testRankingID, ItemID: testItemID, UserID: "u2", Scores: map[string]int{testAttrID1: 100, testAttrID2: 80}},
				}, nil
			}
			return []*domainrating.Rating{
				{RankingID: testRankingID, ItemID: testItemID2, UserID: "u1", Scores: map[string]int{testAttrID1: 50, testAttrID2: 50}},
			}, nil
		},
	}
	imgStore := &mockImageStore{
		generateURLsFn: func(_ context.Context, _, itemID string) (string, string, error) {
			return fmt.Sprintf("https://img/%s", itemID), fmt.Sprintf("https://thumb/%s", itemID), nil
		},
	}

	uc := NewListItemsUseCase(rankingRepo, itemRepo, ratingRepo, imgStore)
	result, err := uc.Execute(context.Background(), ListItemsInput{
		RankingID: testRankingID,
		Mode:      modeAvg,
	})

	require.NoError(t, err)
	require.Len(t, result, 2)

	// item-1: Graphics avg=(80+100)/2=90, Sound avg=(60+80)/2=70, overall=(90+70)/2=80
	ps5 := findListOutput(result, testItemID)
	require.NotNil(t, ps5)
	require.Equal(t, 90.0, ps5.Scores[testAttrGfx])
	require.Equal(t, 70.0, ps5.Scores[testAttrSound])
	require.Equal(t, 80.0, ps5.Overall)
	require.Equal(t, "https://img/item-1", ps5.ImageURL)
	require.Equal(t, "https://thumb/item-1", ps5.ThumbnailURL)

	// item-2: Graphics=50, Sound=50, overall=50
	xbox := findListOutput(result, testItemID2)
	require.NotNil(t, xbox)
	require.Equal(t, 50.0, xbox.Scores[testAttrGfx])
	require.Equal(t, 50.0, xbox.Scores[testAttrSound])
	require.Equal(t, 50.0, xbox.Overall)
	require.Empty(t, xbox.ImageURL)
}

func TestListItemsUseCase_UserModeWithRatings(t *testing.T) {
	t.Parallel()

	rankingRepo := &mockRankingRepo{
		findByIDFn: func(_ context.Context, _ string) (*domainranking.Ranking, error) {
			return newListRanking(), nil
		},
	}
	itemRepo := &mockItemRepo{
		findByRankingFn: func(_ context.Context, _ string) ([]*domainitem.Item, error) {
			return newListItems(), nil
		},
	}
	ratingRepo := &mockRatingRepo{
		findByUserAndItemFn: func(_ context.Context, _, itemID, _ string) (*domainrating.Rating, error) {
			if itemID == testItemID {
				return &domainrating.Rating{
					RankingID: testRankingID, ItemID: testItemID, UserID: "u1",
					Scores: map[string]int{testAttrID1: 90, testAttrID2: 70},
				}, nil
			}
			return &domainrating.Rating{
				RankingID: testRankingID, ItemID: testItemID2, UserID: "u1",
				Scores: map[string]int{testAttrID1: 40, testAttrID2: 60},
			}, nil
		},
	}

	uc := NewListItemsUseCase(rankingRepo, itemRepo, ratingRepo, &mockImageStore{})
	result, err := uc.Execute(context.Background(), ListItemsInput{
		RankingID: testRankingID,
		Mode:      modeUser,
		UserID:    "u1",
	})

	require.NoError(t, err)
	require.Len(t, result, 2)

	ps5 := findListOutput(result, testItemID)
	require.NotNil(t, ps5)
	require.Equal(t, 90.0, ps5.Scores[testAttrGfx])
	require.Equal(t, 70.0, ps5.Scores[testAttrSound])
	require.Equal(t, 80.0, ps5.Overall)

	xbox := findListOutput(result, testItemID2)
	require.NotNil(t, xbox)
	require.Equal(t, 40.0, xbox.Scores[testAttrGfx])
	require.Equal(t, 60.0, xbox.Scores[testAttrSound])
	require.Equal(t, 50.0, xbox.Overall)
}

func TestListItemsUseCase_UserModeNoRatings(t *testing.T) {
	t.Parallel()

	rankingRepo := &mockRankingRepo{
		findByIDFn: func(_ context.Context, _ string) (*domainranking.Ranking, error) {
			return newListRanking(), nil
		},
	}
	itemRepo := &mockItemRepo{
		findByRankingFn: func(_ context.Context, _ string) ([]*domainitem.Item, error) {
			return newListItems(), nil
		},
	}
	ratingRepo := &mockRatingRepo{
		findByUserAndItemFn: func(_ context.Context, _, _, _ string) (*domainrating.Rating, error) {
			return nil, nil
		},
	}

	uc := NewListItemsUseCase(rankingRepo, itemRepo, ratingRepo, &mockImageStore{})
	result, err := uc.Execute(context.Background(), ListItemsInput{
		RankingID: testRankingID,
		Mode:      modeUser,
		UserID:    "u1",
	})

	require.NoError(t, err)
	require.Len(t, result, 2)

	for _, out := range result {
		require.Equal(t, 0.0, out.Overall)
		require.Empty(t, out.Scores)
	}
}

func TestListItemsUseCase_SortByOverallDesc(t *testing.T) {
	t.Parallel()

	rankingRepo := &mockRankingRepo{
		findByIDFn: func(_ context.Context, _ string) (*domainranking.Ranking, error) {
			return newListRanking(), nil
		},
	}
	itemRepo := &mockItemRepo{
		findByRankingFn: func(_ context.Context, _ string) ([]*domainitem.Item, error) {
			return newListItems(), nil
		},
	}
	ratingRepo := &mockRatingRepo{
		findAllByItemFn: func(_ context.Context, _, itemID string) ([]*domainrating.Rating, error) {
			if itemID == testItemID {
				return []*domainrating.Rating{
					{Scores: map[string]int{testAttrID1: 40, testAttrID2: 40}},
				}, nil
			}
			return []*domainrating.Rating{
				{Scores: map[string]int{testAttrID1: 90, testAttrID2: 90}},
			}, nil
		},
	}

	uc := NewListItemsUseCase(rankingRepo, itemRepo, ratingRepo, &mockImageStore{})
	result, err := uc.Execute(context.Background(), ListItemsInput{
		RankingID: testRankingID,
		Mode:      modeAvg,
		Sort:      "-overall",
	})

	require.NoError(t, err)
	require.Len(t, result, 2)
	// Xbox (90) should come first, then PS5 (40)
	require.Equal(t, testItemID2, result[0].ID)
	require.Equal(t, testItemID, result[1].ID)
}

func TestListItemsUseCase_SortByAttribute(t *testing.T) {
	t.Parallel()

	rankingRepo := &mockRankingRepo{
		findByIDFn: func(_ context.Context, _ string) (*domainranking.Ranking, error) {
			return newListRanking(), nil
		},
	}
	itemRepo := &mockItemRepo{
		findByRankingFn: func(_ context.Context, _ string) ([]*domainitem.Item, error) {
			return newListItems(), nil
		},
	}
	ratingRepo := &mockRatingRepo{
		findAllByItemFn: func(_ context.Context, _, itemID string) ([]*domainrating.Rating, error) {
			if itemID == testItemID {
				return []*domainrating.Rating{
					{Scores: map[string]int{testAttrID1: 30, testAttrID2: 90}},
				}, nil
			}
			return []*domainrating.Rating{
				{Scores: map[string]int{testAttrID1: 80, testAttrID2: 10}},
			}, nil
		},
	}

	uc := NewListItemsUseCase(rankingRepo, itemRepo, ratingRepo, &mockImageStore{})
	result, err := uc.Execute(context.Background(), ListItemsInput{
		RankingID: testRankingID,
		Mode:      modeAvg,
		Sort:      testAttrGfx,
	})

	require.NoError(t, err)
	require.Len(t, result, 2)
	// PS5 Graphics=30, Xbox Graphics=80 → ascending: PS5 first
	require.Equal(t, testItemID, result[0].ID)
	require.Equal(t, testItemID2, result[1].ID)
}

func TestListItemsUseCase_InvalidSortField(t *testing.T) {
	t.Parallel()

	rankingRepo := &mockRankingRepo{
		findByIDFn: func(_ context.Context, _ string) (*domainranking.Ranking, error) {
			return newListRanking(), nil
		},
	}
	itemRepo := &mockItemRepo{
		findByRankingFn: func(_ context.Context, _ string) ([]*domainitem.Item, error) {
			return newListItems(), nil
		},
	}

	uc := NewListItemsUseCase(rankingRepo, itemRepo, &mockRatingRepo{}, &mockImageStore{})
	_, err := uc.Execute(context.Background(), ListItemsInput{
		RankingID: testRankingID,
		Mode:      modeAvg,
		Sort:      "-invalid",
	})

	require.Error(t, err)
	var appErr *apperror.AppError
	require.ErrorAs(t, err, &appErr)
	require.Equal(t, apperror.InvalidSortField, appErr.Code)
}

func findListOutput(items []ListOutput, id string) *ListOutput {
	for i := range items {
		if items[i].ID == id {
			return &items[i]
		}
	}
	return nil
}
