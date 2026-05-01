package ranking

import (
	"context"

	domainitem "github.com/josimar/ranking/backend/internal/domain/item"
	domain "github.com/josimar/ranking/backend/internal/domain/ranking"
	domainrating "github.com/josimar/ranking/backend/internal/domain/rating"
)

// DeleteRankingUseCase handles ranking deletion with cascading cleanup.
type DeleteRankingUseCase struct {
	rankingRepo domain.Repository
	itemRepo    domainitem.Repository
	ratingRepo  domainrating.Repository
	imageStore  domainitem.ImageStore
}

// NewDeleteRankingUseCase creates a new DeleteRankingUseCase.
func NewDeleteRankingUseCase(
	rankingRepo domain.Repository,
	itemRepo domainitem.Repository,
	ratingRepo domainrating.Repository,
	imageStore domainitem.ImageStore,
) *DeleteRankingUseCase {
	return &DeleteRankingUseCase{
		rankingRepo: rankingRepo,
		itemRepo:    itemRepo,
		ratingRepo:  ratingRepo,
		imageStore:  imageStore,
	}
}

// Execute deletes a ranking and all associated data after verifying ownership.
func (uc *DeleteRankingUseCase) Execute(ctx context.Context, rankingID, userID string) error {
	r, err := uc.rankingRepo.FindByID(ctx, rankingID)
	if err != nil {
		return err
	}

	if !r.IsOwner(userID) {
		return domain.ErrNotOwner
	}

	items, err := uc.itemRepo.FindByRanking(ctx, rankingID)
	if err != nil {
		return err
	}

	itemIDs := make([]string, len(items))
	for i, it := range items {
		itemIDs[i] = it.ID
	}

	if len(itemIDs) > 0 {
		if err := uc.ratingRepo.DeleteByRanking(ctx, rankingID, itemIDs); err != nil {
			return err
		}
	}

	if err := uc.imageStore.DeleteByRanking(ctx, rankingID); err != nil {
		return err
	}

	if err := uc.itemRepo.DeleteByRanking(ctx, rankingID); err != nil {
		return err
	}

	return uc.rankingRepo.Delete(ctx, rankingID)
}
