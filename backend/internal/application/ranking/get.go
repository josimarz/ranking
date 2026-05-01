package ranking

import (
	"context"

	domain "github.com/josimar/ranking/backend/internal/domain/ranking"
)

// GetRankingUseCase handles fetching a ranking by ID.
type GetRankingUseCase struct {
	repo domain.Repository
}

// NewGetRankingUseCase creates a new GetRankingUseCase.
func NewGetRankingUseCase(repo domain.Repository) *GetRankingUseCase {
	return &GetRankingUseCase{repo: repo}
}

// Execute fetches a ranking and determines ownership.
func (uc *GetRankingUseCase) Execute(ctx context.Context, rankingID, userID string) (*Output, error) {
	r, err := uc.repo.FindByID(ctx, rankingID)
	if err != nil {
		return nil, err
	}

	isOwner := userID != "" && r.IsOwner(userID)

	return toOutput(r, isOwner), nil
}
