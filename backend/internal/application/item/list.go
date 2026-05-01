package item

import (
	"context"
	"slices"
	"strings"

	domainitem "github.com/josimar/ranking/backend/internal/domain/item"
	domainranking "github.com/josimar/ranking/backend/internal/domain/ranking"
	domainrating "github.com/josimar/ranking/backend/internal/domain/rating"
	"github.com/josimar/ranking/backend/pkg/sorting"
)

// ListOutput is the output DTO for listing items with scores.
type ListOutput struct {
	ID           string
	Name         string
	ImageURL     string
	ThumbnailURL string
	Scores       map[string]float64
	Overall      float64
}

const (
	modeAvg  = "avg"
	modeUser = "user"
)

// ListItemsInput is the input DTO for listing items.
type ListItemsInput struct {
	RankingID string
	Mode      string
	UserID    string
	Sort      string
}

// ListItemsUseCase handles listing items with computed scores.
type ListItemsUseCase struct {
	rankingRepo domainranking.Repository
	itemRepo    domainitem.Repository
	ratingRepo  domainrating.Repository
	imageStore  domainitem.ImageStore
}

// NewListItemsUseCase creates a new ListItemsUseCase.
func NewListItemsUseCase(
	rankingRepo domainranking.Repository,
	itemRepo domainitem.Repository,
	ratingRepo domainrating.Repository,
	imageStore domainitem.ImageStore,
) *ListItemsUseCase {
	return &ListItemsUseCase{
		rankingRepo: rankingRepo,
		itemRepo:    itemRepo,
		ratingRepo:  ratingRepo,
		imageStore:  imageStore,
	}
}

// Execute lists items for a ranking with computed scores and sorting.
func (uc *ListItemsUseCase) Execute(ctx context.Context, input ListItemsInput) ([]ListOutput, error) {
	r, err := uc.rankingRepo.FindByID(ctx, input.RankingID)
	if err != nil {
		return nil, err
	}

	activeAttrs := r.ActiveAttributes()
	attrIDs := make([]string, len(activeAttrs))
	attrNames := make([]string, len(activeAttrs))
	idToName := make(map[string]string, len(activeAttrs))
	for i, a := range activeAttrs {
		attrIDs[i] = a.ID
		attrNames[i] = a.Name
		idToName[a.ID] = a.Name
	}

	allowedSortFields := append([]string{"name", "overall"}, attrNames...)
	sortFields, err := sorting.ParseSort(input.Sort, allowedSortFields)
	if err != nil {
		return nil, err
	}

	items, err := uc.itemRepo.FindByRanking(ctx, input.RankingID)
	if err != nil {
		return nil, err
	}

	result := make([]ListOutput, len(items))
	for i, it := range items {
		scores, overall := uc.computeScores(ctx, input, it, attrIDs, idToName)

		var imageURL, thumbURL string
		if it.ImageKey != "" {
			imageURL, thumbURL, _ = uc.imageStore.GenerateURLs(ctx, it.RankingID, it.ID)
		}

		result[i] = ListOutput{
			ID:           it.ID,
			Name:         it.Name,
			ImageURL:     imageURL,
			ThumbnailURL: thumbURL,
			Scores:       scores,
			Overall:      overall,
		}
	}

	if len(sortFields) > 0 {
		slices.SortFunc(result, func(a, b ListOutput) int {
			for _, sf := range sortFields {
				cmp := compareBySortField(a, b, sf)
				if cmp != 0 {
					return cmp
				}
			}
			return 0
		})
	}

	return result, nil
}

func (uc *ListItemsUseCase) computeScores(
	ctx context.Context,
	input ListItemsInput,
	it *domainitem.Item,
	attrIDs []string,
	idToName map[string]string,
) (map[string]float64, float64) {
	if input.Mode == modeUser {
		return uc.computeUserScores(ctx, input, it, attrIDs, idToName)
	}
	return uc.computeAvgScores(ctx, it, attrIDs, idToName)
}

func (uc *ListItemsUseCase) computeAvgScores(
	ctx context.Context,
	it *domainitem.Item,
	attrIDs []string,
	idToName map[string]string,
) (map[string]float64, float64) {
	ratings, err := uc.ratingRepo.FindAllByItem(ctx, it.RankingID, it.ID)
	if err != nil || len(ratings) == 0 {
		return nil, 0
	}

	scores := make(map[string]float64, len(attrIDs))
	for _, id := range attrIDs {
		sum := 0
		count := 0
		for _, rt := range ratings {
			if v, ok := rt.Scores[id]; ok {
				sum += v
				count++
			}
		}
		if count > 0 {
			scores[idToName[id]] = float64(sum) / float64(count)
		}
	}

	if len(scores) == 0 {
		return nil, 0
	}

	var total float64
	for _, v := range scores {
		total += v
	}

	return scores, total / float64(len(scores))
}

func (uc *ListItemsUseCase) computeUserScores(
	ctx context.Context,
	input ListItemsInput,
	it *domainitem.Item,
	attrIDs []string,
	idToName map[string]string,
) (map[string]float64, float64) {
	rt, err := uc.ratingRepo.FindByUserAndItem(ctx, it.RankingID, it.ID, input.UserID)
	if err != nil || rt == nil {
		return nil, 0
	}

	scores := make(map[string]float64, len(attrIDs))
	for _, id := range attrIDs {
		if v, ok := rt.Scores[id]; ok {
			scores[idToName[id]] = float64(v)
		}
	}

	if len(scores) == 0 {
		return nil, 0
	}

	var total float64
	for _, v := range scores {
		total += v
	}

	return scores, total / float64(len(scores))
}

func compareBySortField(a, b ListOutput, sf sorting.SortField) int {
	var cmp int
	if strings.EqualFold(sf.Field, "name") {
		cmp = strings.Compare(a.Name, b.Name)
	} else if strings.EqualFold(sf.Field, "overall") {
		cmp = compareFloat(a.Overall, b.Overall)
	} else {
		cmp = compareFloat(a.Scores[sf.Field], b.Scores[sf.Field])
	}

	if sf.Direction == sorting.Descending {
		cmp = -cmp
	}

	return cmp
}

func compareFloat(a, b float64) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}
