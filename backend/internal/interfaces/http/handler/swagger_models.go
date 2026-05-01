package handler

import "time"

// --- Swagger response models ---

// ErrorResponse represents a structured API error response.
type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

// ErrorBody contains error details.
type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// RankingResponse represents a ranking in API responses.
type RankingResponse struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Visibility  string              `json:"visibility"`
	Tags        []string            `json:"tags"`
	Attributes  []AttributeResponse `json:"attributes"`
	OwnerUserID string              `json:"ownerUserId"`
	CreatedAt   time.Time           `json:"createdAt"`
	UpdatedAt   time.Time           `json:"updatedAt"`
	IsOwner     bool                `json:"isOwner"`
}

// AttributeResponse represents an attribute in API responses.
type AttributeResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
}

// ItemResponse represents an item in API responses.
type ItemResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	ImageKey  string    `json:"imageKey"`
	RankingID string    `json:"rankingId"`
	CreatedBy string    `json:"createdBy"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ItemListResponse represents an item with scores in list responses.
type ItemListResponse struct {
	ID           string             `json:"id"`
	Name         string             `json:"name"`
	ImageURL     string             `json:"imageUrl"`
	ThumbnailURL string             `json:"thumbnailUrl"`
	Scores       map[string]float64 `json:"scores"`
	Overall      float64            `json:"overall"`
}

// RatingResponse represents a rating in API responses.
type RatingResponse struct {
	RankingID string         `json:"rankingId"`
	ItemID    string         `json:"itemId"`
	UserID    string         `json:"userId"`
	Scores    map[string]int `json:"scores"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
}

// PaginatedRankingResponse wraps paginated ranking data.
type PaginatedRankingResponse struct {
	Data       []RankingResponse  `json:"data"`
	Pagination PaginationResponse `json:"pagination"`
}

// PaginationResponse holds cursor-based pagination metadata.
type PaginationResponse struct {
	NextCursor string `json:"next_cursor"`
	HasMore    bool   `json:"has_more"`
}
