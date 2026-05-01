package item

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	domainitem "github.com/josimar/ranking/backend/internal/domain/item"
	domainranking "github.com/josimar/ranking/backend/internal/domain/ranking"
	"github.com/josimar/ranking/backend/pkg/apperror"
)

const maxItems = 100

// Output is the output DTO for item use cases.
type Output struct {
	ID        string
	Name      string
	ImageKey  string
	RankingID string
	CreatedBy string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// AddItemInput is the input DTO for adding an item.
type AddItemInput struct {
	RankingID string
	UserID    string
	Name      string
	ImageData []byte
	ImageURL  string
}

// AddItemUseCase handles adding an item to a ranking.
type AddItemUseCase struct {
	rankingRepo domainranking.Repository
	itemRepo    domainitem.Repository
	imageStore  domainitem.ImageStore
	imageProc   domainitem.ImageProcessor
}

// NewAddItemUseCase creates a new AddItemUseCase.
func NewAddItemUseCase(
	rankingRepo domainranking.Repository,
	itemRepo domainitem.Repository,
	imageStore domainitem.ImageStore,
	imageProc domainitem.ImageProcessor,
) *AddItemUseCase {
	return &AddItemUseCase{
		rankingRepo: rankingRepo,
		itemRepo:    itemRepo,
		imageStore:  imageStore,
		imageProc:   imageProc,
	}
}

// Execute adds a new item to a ranking.
func (uc *AddItemUseCase) Execute(ctx context.Context, input AddItemInput) (*Output, error) {
	if _, err := uc.rankingRepo.FindByID(ctx, input.RankingID); err != nil {
		return nil, err
	}

	count, err := uc.itemRepo.CountByRanking(ctx, input.RankingID)
	if err != nil {
		return nil, err
	}

	if count >= maxItems {
		return nil, domainitem.ErrMaxItemsReached
	}

	imageData := input.ImageData
	if len(imageData) == 0 && input.ImageURL != "" {
		downloaded, dlErr := downloadImage(input.ImageURL)
		if dlErr != nil {
			return nil, dlErr
		}
		imageData = downloaded
	}

	newItem, err := domainitem.NewItem(input.Name, input.RankingID, input.UserID)
	if err != nil {
		return nil, err
	}

	if len(imageData) > 0 {
		if err := uc.imageProc.ValidateFormat(imageData); err != nil {
			return nil, err
		}

		original, thumbnail, procErr := uc.imageProc.Process(imageData)
		if procErr != nil {
			return nil, procErr
		}

		if err := uc.imageStore.Store(ctx, input.RankingID, newItem.ID, original, thumbnail); err != nil {
			return nil, err
		}

		newItem.ImageKey = fmt.Sprintf("items/%s/%s", input.RankingID, newItem.ID)
	}

	if err := uc.itemRepo.Save(ctx, newItem); err != nil {
		return nil, err
	}

	return toOutput(newItem), nil
}

func downloadImage(url string) ([]byte, error) {
	resp, err := http.Get(url) //nolint:gosec,noctx // URL is user-provided input for image download
	if err != nil {
		return nil, apperror.NewBadRequestError(apperror.ImageDownloadFailed, fmt.Sprintf("failed to download image: %s", err.Error()))
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, apperror.NewBadRequestError(apperror.ImageDownloadFailed, fmt.Sprintf("image download returned status %d", resp.StatusCode))
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, apperror.NewBadRequestError(apperror.ImageDownloadFailed, fmt.Sprintf("failed to read image data: %s", err.Error()))
	}

	return data, nil
}

func toOutput(i *domainitem.Item) *Output {
	return &Output{
		ID:        i.ID,
		Name:      i.Name,
		ImageKey:  i.ImageKey,
		RankingID: i.RankingID,
		CreatedBy: i.CreatedBy,
		CreatedAt: i.CreatedAt,
		UpdatedAt: i.UpdatedAt,
	}
}
