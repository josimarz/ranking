// Package item defines the Item domain entity and related interfaces.
package item

import (
	"context"
	"time"

	"github.com/josimar/ranking/backend/pkg/apperror"
	"github.com/josimar/ranking/backend/pkg/uuid"
)

const maxNameLen = 100

// Item represents an element to be rated within a ranking.
type Item struct {
	ID        string
	Name      string
	ImageKey  string
	RankingID string
	CreatedBy string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewItem creates a new Item with a generated ID and validated name.
func NewItem(name, rankingID, createdBy string) (*Item, error) {
	if err := validateName(name); err != nil {
		return nil, err
	}

	now := time.Now()

	return &Item{
		ID:        uuid.New(),
		Name:      name,
		RankingID: rankingID,
		CreatedBy: createdBy,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// Update changes the item name after validation.
func (i *Item) Update(name string) error {
	if err := validateName(name); err != nil {
		return err
	}

	i.Name = name
	i.UpdatedAt = time.Now()

	return nil
}

// CanBeModifiedBy returns true if the user is the item creator or the ranking owner.
func (i *Item) CanBeModifiedBy(userID, rankingOwnerUserID string) bool {
	return userID == i.CreatedBy || userID == rankingOwnerUserID
}

func validateName(name string) error {
	if len(name) == 0 || len(name) > maxNameLen {
		return ErrItemNameRequired
	}

	return nil
}

// Domain errors.
var (
	ErrItemNameRequired = apperror.NewValidationError("item name is required and must be between 1 and 100 characters")
	ErrMaxItemsReached  = apperror.NewBadRequestError(apperror.MaxItemsReached, "maximum number of items reached for this ranking")
	ErrItemNotFound     = apperror.NewNotFoundError(apperror.ItemNotFound, "item not found")
	ErrNotAuthorized    = apperror.NewForbiddenError(apperror.NotAuthorized, "not authorized to modify this item")
)

// Repository defines persistence operations for items.
type Repository interface {
	Save(ctx context.Context, item *Item) error
	FindByID(ctx context.Context, rankingID, itemID string) (*Item, error)
	FindByRanking(ctx context.Context, rankingID string) ([]*Item, error)
	Update(ctx context.Context, item *Item) error
	Delete(ctx context.Context, rankingID, itemID string) error
	DeleteByRanking(ctx context.Context, rankingID string) error
	CountByRanking(ctx context.Context, rankingID string) (int, error)
}

// ImageStore defines operations for storing and retrieving item images.
type ImageStore interface {
	Store(ctx context.Context, rankingID, itemID string, original, thumbnail []byte) error
	GenerateURLs(ctx context.Context, rankingID, itemID string) (originalURL, thumbnailURL string, err error)
	Delete(ctx context.Context, rankingID, itemID string) error
	DeleteByRanking(ctx context.Context, rankingID string) error
}

// ImageProcessor defines operations for processing item images.
type ImageProcessor interface {
	Process(data []byte) (original, thumbnail []byte, err error)
	ValidateFormat(data []byte) error
}
