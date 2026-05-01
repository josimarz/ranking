package ranking

import (
	"context"
	"time"

	domain "github.com/josimar/ranking/backend/internal/domain/ranking"
)

// AttributeInput represents input for creating/updating an attribute.
type AttributeInput struct {
	Name        string
	Description string
}

// AttributeOutput represents an attribute in use case output.
type AttributeOutput struct {
	ID          string
	Name        string
	Description string
	Active      bool
}

// Output is the output DTO for ranking use cases.
type Output struct {
	ID          string
	Name        string
	Description string
	Visibility  string
	Tags        []string
	Attributes  []AttributeOutput
	OwnerUserID string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	IsOwner     bool
}

// CreateRankingInput is the input DTO for creating a ranking.
type CreateRankingInput struct {
	Name        string
	Description string
	Visibility  string
	Tags        []string
	Attributes  []AttributeInput
	UserID      string
}

// CreateRankingUseCase handles ranking creation.
type CreateRankingUseCase struct {
	repo domain.Repository
}

// NewCreateRankingUseCase creates a new CreateRankingUseCase.
func NewCreateRankingUseCase(repo domain.Repository) *CreateRankingUseCase {
	return &CreateRankingUseCase{repo: repo}
}

// Execute creates a new ranking.
func (uc *CreateRankingUseCase) Execute(ctx context.Context, input CreateRankingInput) (*Output, error) {
	attrs := make([]domain.AttributeInput, len(input.Attributes))
	for i, a := range input.Attributes {
		attrs[i] = domain.AttributeInput{Name: a.Name, Description: a.Description}
	}

	r, err := domain.NewRanking(input.Name, input.Description, domain.Visibility(input.Visibility), input.Tags, attrs, input.UserID)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Save(ctx, r); err != nil {
		return nil, err
	}

	return toOutput(r, false), nil
}

func toOutput(r *domain.Ranking, isOwner bool) *Output {
	attrs := make([]AttributeOutput, len(r.Attributes))
	for i, a := range r.Attributes {
		attrs[i] = AttributeOutput{ID: a.ID, Name: a.Name, Description: a.Description, Active: a.Active}
	}

	return &Output{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Visibility:  string(r.Visibility),
		Tags:        r.Tags,
		Attributes:  attrs,
		OwnerUserID: r.OwnerUserID,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
		IsOwner:     isOwner,
	}
}
