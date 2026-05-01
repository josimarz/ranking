// Package ranking defines the Ranking aggregate root and related value objects.
package ranking

import (
	"context"
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/josimar/ranking/backend/pkg/apperror"
	"github.com/josimar/ranking/backend/pkg/uuid"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// Visibility represents the visibility type of a ranking.
type Visibility string

const (
	// Public rankings are discoverable by anyone.
	Public Visibility = "public"
	// Private rankings are accessible only via direct link.
	Private Visibility = "private"
)

// ValidateVisibility checks whether a Visibility value is valid.
func ValidateVisibility(v Visibility) error {
	if v != Public && v != Private {
		return apperror.NewValidationError(fmt.Sprintf("visibility must be %q or %q", Public, Private))
	}
	return nil
}

// Attribute is a value object representing a rating criterion.
type Attribute struct {
	ID          string
	Name        string
	Description string
	Active      bool
}

// AttributeInput is used when creating or updating attributes.
type AttributeInput struct {
	Name        string
	Description string
}

// Ranking is the aggregate root for the ranking domain.
type Ranking struct {
	ID          string
	Name        string
	NameLower   string
	Description string
	Visibility  Visibility
	Tags        []string
	Attributes  []Attribute
	OwnerUserID string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Domain errors.
var (
	ErrRankingNotFound = apperror.NewNotFoundError(apperror.RankingNotFound, "ranking not found")
	ErrNotOwner        = apperror.NewForbiddenError(apperror.NotOwner, "not the ranking owner")
)

// NormalizeName converts a name to lowercase and strips diacritical marks.
func NormalizeName(name string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, name)
	return strings.ToLower(result)
}

// NewRanking creates a new Ranking with validated inputs.
func NewRanking(name, description string, visibility Visibility, tags []string, attributes []AttributeInput, ownerUserID string) (*Ranking, error) {
	if err := validateRankingFields(name, description, visibility, tags, attributes, ownerUserID); err != nil {
		return nil, err
	}

	now := time.Now()
	attrs := make([]Attribute, len(attributes))
	for i, a := range attributes {
		attrs[i] = Attribute{ID: uuid.New(), Name: a.Name, Description: a.Description, Active: true}
	}

	return &Ranking{
		ID:          uuid.New(),
		Name:        name,
		NameLower:   NormalizeName(name),
		Description: description,
		Visibility:  visibility,
		Tags:        tags,
		Attributes:  attrs,
		OwnerUserID: ownerUserID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// Update modifies the ranking fields and performs soft-delete on removed attributes.
func (r *Ranking) Update(name, description string, visibility Visibility, tags []string, newAttributes []AttributeInput) error {
	if err := validateRankingFields(name, description, visibility, tags, newAttributes, r.OwnerUserID); err != nil {
		return err
	}

	r.Name = name
	r.NameLower = NormalizeName(name)
	r.Description = description
	r.Visibility = visibility
	r.Tags = tags
	r.UpdatedAt = time.Now()

	// Build lookup of new attribute names
	newNames := make(map[string]AttributeInput, len(newAttributes))
	for _, a := range newAttributes {
		newNames[a.Name] = a
	}

	// Mark existing attributes: keep if in newNames, soft-delete otherwise
	kept := make(map[string]bool, len(r.Attributes))
	for i := range r.Attributes {
		if _, exists := newNames[r.Attributes[i].Name]; exists {
			r.Attributes[i].Active = true
			kept[r.Attributes[i].Name] = true
		} else {
			r.Attributes[i].Active = false
		}
	}

	// Add truly new attributes
	for _, a := range newAttributes {
		if !kept[a.Name] {
			r.Attributes = append(r.Attributes, Attribute{
				ID: uuid.New(), Name: a.Name, Description: a.Description, Active: true,
			})
		}
	}

	return nil
}

// IsOwner checks whether the given user ID matches the ranking owner.
func (r *Ranking) IsOwner(userID string) bool {
	return r.OwnerUserID == userID
}

// ActiveAttributes returns only attributes that are active.
func (r *Ranking) ActiveAttributes() []Attribute {
	var result []Attribute
	for _, a := range r.Attributes {
		if a.Active {
			result = append(result, a)
		}
	}
	return result
}

// Repository defines persistence operations for rankings.
type Repository interface {
	Save(ctx context.Context, ranking *Ranking) error
	FindByID(ctx context.Context, id string) (*Ranking, error)
	Update(ctx context.Context, ranking *Ranking) error
	Delete(ctx context.Context, id string) error
	ListPublic(ctx context.Context, limit int, cursor string, sortField string, sortDir string) ([]*Ranking, string, error)
	ListByOwner(ctx context.Context, userID string, limit int, cursor string) ([]*Ranking, string, error)
	SearchByName(ctx context.Context, term string, limit int, cursor string) ([]*Ranking, string, error)
	SearchByTag(ctx context.Context, tag string, limit int, cursor string) ([]*Ranking, string, error)
	FindRecent(ctx context.Context) ([]*Ranking, error)
}

func validateRankingFields(name, description string, visibility Visibility, tags []string, attributes []AttributeInput, ownerUserID string) error {
	if len(strings.TrimSpace(name)) == 0 {
		return apperror.NewValidationError("name is required")
	}
	if len(name) > 100 {
		return apperror.NewValidationError("name must be at most 100 characters")
	}
	if len(description) > 500 {
		return apperror.NewValidationError("description must be at most 500 characters")
	}
	if err := ValidateVisibility(visibility); err != nil {
		return err
	}
	if len(tags) > 10 {
		return apperror.NewValidationError("tags must be at most 10")
	}
	if len(attributes) == 0 {
		return apperror.NewValidationError("at least 1 attribute is required")
	}
	if len(attributes) > 6 {
		return apperror.NewValidationError("at most 6 attributes are allowed")
	}
	if len(strings.TrimSpace(ownerUserID)) == 0 {
		return apperror.NewValidationError("owner user ID is required")
	}
	for _, a := range attributes {
		if len(strings.TrimSpace(a.Name)) == 0 {
			return apperror.NewValidationError("attribute name is required")
		}
		if len(a.Name) > 50 {
			return apperror.NewValidationError("attribute name must be at most 50 characters")
		}
		if len(a.Description) > 200 {
			return apperror.NewValidationError("attribute description must be at most 200 characters")
		}
	}
	return nil
}
