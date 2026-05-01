package ranking

import (
	"context"

	domain "github.com/josimar/ranking/backend/internal/domain/ranking"
)

// UpdateRankingInput is the input DTO for updating a ranking.
type UpdateRankingInput struct {
	RankingID   string
	UserID      string
	Name        string
	Description string
	Visibility  string
	Tags        []string
	Attributes  []AttributeInput
}

// UpdateRankingUseCase handles ranking updates.
type UpdateRankingUseCase struct {
	repo domain.Repository
}

// NewUpdateRankingUseCase creates a new UpdateRankingUseCase.
func NewUpdateRankingUseCase(repo domain.Repository) *UpdateRankingUseCase {
	return &UpdateRankingUseCase{repo: repo}
}

// Execute updates an existing ranking after verifying ownership.
func (uc *UpdateRankingUseCase) Execute(ctx context.Context, input UpdateRankingInput) (*Output, error) {
	r, err := uc.repo.FindByID(ctx, input.RankingID)
	if err != nil {
		return nil, err
	}

	if !r.IsOwner(input.UserID) {
		return nil, domain.ErrNotOwner
	}

	attrs := make([]domain.AttributeInput, len(input.Attributes))
	for i, a := range input.Attributes {
		attrs[i] = domain.AttributeInput{Name: a.Name, Description: a.Description}
	}

	if err := r.Update(input.Name, input.Description, domain.Visibility(input.Visibility), input.Tags, attrs); err != nil {
		return nil, err
	}

	if err := uc.repo.Update(ctx, r); err != nil {
		return nil, err
	}

	return toOutput(r, true), nil
}
