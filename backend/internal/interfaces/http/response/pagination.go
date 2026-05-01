package response

// Pagination holds cursor-based pagination metadata.
type Pagination struct {
	NextCursor string `json:"next_cursor"`
	HasMore    bool   `json:"has_more"`
}

// PaginatedResponse wraps data with pagination metadata.
type PaginatedResponse struct {
	Data       any        `json:"data"`
	Pagination Pagination `json:"pagination"`
}

// NewPaginatedResponse creates a PaginatedResponse, setting HasMore based on cursor presence.
func NewPaginatedResponse(data any, nextCursor string) PaginatedResponse {
	return PaginatedResponse{
		Data: data,
		Pagination: Pagination{
			NextCursor: nextCursor,
			HasMore:    nextCursor != "",
		},
	}
}
