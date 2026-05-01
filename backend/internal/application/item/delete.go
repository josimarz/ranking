package item

import (
	"context"

	domainitem "github.com/josimar/ranking/backend/internal/domain/item"
	domainranking "github.com/josimar/ranking/backend/internal/domain/ranking"
	domainrating "github.com/josimar/ranking/backend/internal/domain/rating"
)

// DeleteItemUseCase handles deleting an item and its associated data.
type DeleteItemUseCase struct {
	rankingRepo domainranking.Repository
	itemRepo    domainitem.Repository
	ratingRepo  domainrating.Repository
	imageStore  domainitem.ImageStore
}

// NewDeleteItemUseCase creates a new DeleteItemUseCase.
func NewDeleteItemUseCase(
	rankingRepo domainranking.Repository,
	itemRepo domainitem.Repository,
	ratingRepo domainrating.Repository,
	imageStore domainitem.ImageStore,
) *DeleteItemUseCase {
	return &DeleteItemUseCase{
		rankingRepo: rankingRepo,
		itemRepo:    itemRepo,
		ratingRepo:  ratingRepo,
		imageStore:  imageStore,
	}
}

// Execute deletes an item and all associated ratings and images.
func (uc *DeleteItemUseCase) Execute(ctx context.Context, rankingID, itemID, userID string) error {
	r, err := uc.rankingRepo.FindByID(ctx, rankingID)
	if err != nil {
		return err
	}

	existing, err := uc.itemRepo.FindByID(ctx, rankingID, itemID)
	if err != nil {
		return err
	}

	if !existing.CanBeModifiedBy(userID, r.OwnerUserID) {
		return domainitem.ErrNotAuthorized
	}

	if err := uc.ratingRepo.DeleteByItem(ctx, rankingID, itemID); err != nil {
		return err
	}

	if err := uc.imageStore.Delete(ctx, rankingID, itemID); err != nil {
		return err
	}

	return uc.itemRepo.Delete(ctx, rankingID, itemID)
}
