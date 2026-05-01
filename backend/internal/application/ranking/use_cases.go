package ranking

import (
	"context"

	domain "github.com/josimar/ranking/backend/internal/domain/ranking"
	"github.com/josimar/ranking/backend/pkg/apperror"
	"github.com/josimar/ranking/backend/pkg/pagination"
	"github.com/josimar/ranking/backend/pkg/sorting"
)

// PaginatedOutput holds a paginated list of rankings.
type PaginatedOutput struct {
	Rankings   []*domain.Ranking
	NextCursor string
}

// --- Task 7.5: ListPublicRankingsUseCase ---

var listPublicAllowedSortFields = []string{"name", "createdAt", "updatedAt"}

// ListPublicRankingsUseCase lists public rankings with pagination and sorting.
type ListPublicRankingsUseCase struct {
	repo domain.Repository
}

// NewListPublicRankingsUseCase creates a new ListPublicRankingsUseCase.
func NewListPublicRankingsUseCase(repo domain.Repository) *ListPublicRankingsUseCase {
	return &ListPublicRankingsUseCase{repo: repo}
}

// Execute runs the list public rankings use case.
func (uc *ListPublicRankingsUseCase) Execute(ctx context.Context, limit int, cursor, sort string) (*PaginatedOutput, error) {
	limit = pagination.ClampLimit(limit, 20, 50)

	var sortField, sortDir string
	if sort != "" {
		fields, err := sorting.ParseSort(sort, listPublicAllowedSortFields)
		if err != nil {
			return nil, err
		}
		if len(fields) > 0 {
			sortField = fields[0].Field
			sortDir = string(fields[0].Direction)
		}
	}

	rankings, nextCursor, err := uc.repo.ListPublic(ctx, limit, cursor, sortField, sortDir)
	if err != nil {
		return nil, err
	}

	return &PaginatedOutput{Rankings: rankings, NextCursor: nextCursor}, nil
}

// --- Task 7.6: GetRecentRankingsUseCase ---

// GetRecentRankingsUseCase retrieves the most recent public rankings.
type GetRecentRankingsUseCase struct {
	repo domain.Repository
}

// NewGetRecentRankingsUseCase creates a new GetRecentRankingsUseCase.
func NewGetRecentRankingsUseCase(repo domain.Repository) *GetRecentRankingsUseCase {
	return &GetRecentRankingsUseCase{repo: repo}
}

// Execute runs the get recent rankings use case.
func (uc *GetRecentRankingsUseCase) Execute(ctx context.Context) ([]*domain.Ranking, error) {
	return uc.repo.FindRecent(ctx)
}

// --- Task 7.7: SearchRankingsUseCase ---

// SearchRankingsUseCase searches rankings by name or tag.
type SearchRankingsUseCase struct {
	repo domain.Repository
}

// NewSearchRankingsUseCase creates a new SearchRankingsUseCase.
func NewSearchRankingsUseCase(repo domain.Repository) *SearchRankingsUseCase {
	return &SearchRankingsUseCase{repo: repo}
}

// Execute runs the search rankings use case.
func (uc *SearchRankingsUseCase) Execute(ctx context.Context, q, tag string, limit int, cursor string) (*PaginatedOutput, error) {
	if q == "" && tag == "" {
		return nil, apperror.NewValidationError("at least one of q or tag is required")
	}

	limit = pagination.ClampLimit(limit, 20, 50)

	var rankings []*domain.Ranking
	var nextCursor string
	var err error

	if q != "" {
		normalized := domain.NormalizeName(q)
		rankings, nextCursor, err = uc.repo.SearchByName(ctx, normalized, limit, cursor)
	} else {
		rankings, nextCursor, err = uc.repo.SearchByTag(ctx, tag, limit, cursor)
	}

	if err != nil {
		return nil, err
	}

	return &PaginatedOutput{Rankings: rankings, NextCursor: nextCursor}, nil
}

// --- Task 7.8: ListMyRankingsUseCase ---

// ListMyRankingsUseCase lists rankings owned by a specific user.
type ListMyRankingsUseCase struct {
	repo domain.Repository
}

// NewListMyRankingsUseCase creates a new ListMyRankingsUseCase.
func NewListMyRankingsUseCase(repo domain.Repository) *ListMyRankingsUseCase {
	return &ListMyRankingsUseCase{repo: repo}
}

// Execute runs the list my rankings use case.
func (uc *ListMyRankingsUseCase) Execute(ctx context.Context, userID string, limit int, cursor string) (*PaginatedOutput, error) {
	limit = pagination.ClampLimit(limit, 20, 50)

	rankings, nextCursor, err := uc.repo.ListByOwner(ctx, userID, limit, cursor)
	if err != nil {
		return nil, err
	}

	return &PaginatedOutput{Rankings: rankings, NextCursor: nextCursor}, nil
}
