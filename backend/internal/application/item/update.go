package item

import (
	"context"
	"fmt"

	domainitem "github.com/josimar/ranking/backend/internal/domain/item"
	domainranking "github.com/josimar/ranking/backend/internal/domain/ranking"
)

// UpdateItemInput is the input DTO for updating an item.
type UpdateItemInput struct {
	RankingID string
	ItemID    string
	UserID    string
	Name      string
	ImageData []byte
}

// UpdateItemUseCase handles updating an item.
type UpdateItemUseCase struct {
	rankingRepo domainranking.Repository
	itemRepo    domainitem.Repository
	imageStore  domainitem.ImageStore
	imageProc   domainitem.ImageProcessor
}

// NewUpdateItemUseCase creates a new UpdateItemUseCase.
func NewUpdateItemUseCase(
	rankingRepo domainranking.Repository,
	itemRepo domainitem.Repository,
	imageStore domainitem.ImageStore,
	imageProc domainitem.ImageProcessor,
) *UpdateItemUseCase {
	return &UpdateItemUseCase{
		rankingRepo: rankingRepo,
		itemRepo:    itemRepo,
		imageStore:  imageStore,
		imageProc:   imageProc,
	}
}

// Execute updates an existing item.
func (uc *UpdateItemUseCase) Execute(ctx context.Context, input UpdateItemInput) (*Output, error) {
	r, err := uc.rankingRepo.FindByID(ctx, input.RankingID)
	if err != nil {
		return nil, err
	}

	existing, err := uc.itemRepo.FindByID(ctx, input.RankingID, input.ItemID)
	if err != nil {
		return nil, err
	}

	if !existing.CanBeModifiedBy(input.UserID, r.OwnerUserID) {
		return nil, domainitem.ErrNotAuthorized
	}

	if len(input.ImageData) > 0 {
		if err := uc.imageProc.ValidateFormat(input.ImageData); err != nil {
			return nil, err
		}

		original, thumbnail, procErr := uc.imageProc.Process(input.ImageData)
		if procErr != nil {
			return nil, procErr
		}

		if err := uc.imageStore.Store(ctx, input.RankingID, input.ItemID, original, thumbnail); err != nil {
			return nil, err
		}

		existing.ImageKey = fmt.Sprintf("items/%s/%s", input.RankingID, input.ItemID)
	}

	if err := existing.Update(input.Name); err != nil {
		return nil, err
	}

	if err := uc.itemRepo.Update(ctx, existing); err != nil {
		return nil, err
	}

	return toOutput(existing), nil
}
