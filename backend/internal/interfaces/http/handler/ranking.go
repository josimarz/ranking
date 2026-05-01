package handler

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
	ranking "github.com/josimar/ranking/backend/internal/application/ranking"
	domain "github.com/josimar/ranking/backend/internal/domain/ranking"
	"github.com/josimar/ranking/backend/internal/interfaces/http/middleware"
	"github.com/josimar/ranking/backend/internal/interfaces/http/response"
	"github.com/josimar/ranking/backend/pkg/apperror"
)

// CreateRankingExecutor creates a ranking.
type CreateRankingExecutor interface {
	Execute(ctx context.Context, input ranking.CreateRankingInput) (*ranking.Output, error)
}

// GetRankingExecutor retrieves a ranking.
type GetRankingExecutor interface {
	Execute(ctx context.Context, rankingID, userID string) (*ranking.Output, error)
}

// UpdateRankingExecutor updates a ranking.
type UpdateRankingExecutor interface {
	Execute(ctx context.Context, input ranking.UpdateRankingInput) (*ranking.Output, error)
}

// DeleteRankingExecutor deletes a ranking.
type DeleteRankingExecutor interface {
	Execute(ctx context.Context, rankingID, userID string) error
}

// ListPublicRankingsExecutor lists public rankings.
type ListPublicRankingsExecutor interface {
	Execute(ctx context.Context, limit int, cursor, sort string) (*ranking.PaginatedOutput, error)
}

// GetRecentRankingsExecutor retrieves recent rankings.
type GetRecentRankingsExecutor interface {
	Execute(ctx context.Context) ([]*domain.Ranking, error)
}

// SearchRankingsExecutor searches rankings.
type SearchRankingsExecutor interface {
	Execute(ctx context.Context, q, tag string, limit int, cursor string) (*ranking.PaginatedOutput, error)
}

// ListMyRankingsExecutor lists user's rankings.
type ListMyRankingsExecutor interface {
	Execute(ctx context.Context, userID string, limit int, cursor string) (*ranking.PaginatedOutput, error)
}

// RankingHandler handles ranking HTTP requests.
type RankingHandler struct {
	create     CreateRankingExecutor
	get        GetRankingExecutor
	update     UpdateRankingExecutor
	del        DeleteRankingExecutor
	listPublic ListPublicRankingsExecutor
	recent     GetRecentRankingsExecutor
	search     SearchRankingsExecutor
	mine       ListMyRankingsExecutor
}

// NewRankingHandler creates a new RankingHandler.
func NewRankingHandler(
	create CreateRankingExecutor,
	get GetRankingExecutor,
	update UpdateRankingExecutor,
	del DeleteRankingExecutor,
	listPublic ListPublicRankingsExecutor,
	recent GetRecentRankingsExecutor,
	search SearchRankingsExecutor,
	mine ListMyRankingsExecutor,
) *RankingHandler {
	return &RankingHandler{
		create: create, get: get, update: update, del: del,
		listPublic: listPublic, recent: recent, search: search, mine: mine,
	}
}

type createRankingRequest struct {
	Name        string                    `json:"name" binding:"required"`
	Description string                    `json:"description"`
	Visibility  string                    `json:"visibility" binding:"required"`
	Tags        []string                  `json:"tags"`
	Attributes  []rankingAttributeRequest `json:"attributes" binding:"required"`
}

type rankingAttributeRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type updateRankingRequest struct {
	Name        string                    `json:"name" binding:"required"`
	Description string                    `json:"description"`
	Visibility  string                    `json:"visibility" binding:"required"`
	Tags        []string                  `json:"tags"`
	Attributes  []rankingAttributeRequest `json:"attributes" binding:"required"`
}

// Create handles POST /api/v1/rankings.
//
//	@Summary		Create a new ranking
//	@Description	Create a new ranking with attributes
//	@Tags			rankings
//	@Accept			json
//	@Produce		json
//	@Param			X-User-Id	header		string				true	"User ID"
//	@Param			body		body		createRankingRequest	true	"Ranking data"
//	@Success		201			{object}	RankingResponse
//	@Failure		400			{object}	ErrorResponse
//	@Failure		401			{object}	ErrorResponse
//	@Router			/rankings [post]
func (h *RankingHandler) Create(c *gin.Context) {
	var req createRankingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperror.NewValidationError(err.Error()))
		c.Abort()
		return
	}

	attrs := make([]ranking.AttributeInput, len(req.Attributes))
	for i, a := range req.Attributes {
		attrs[i] = ranking.AttributeInput{Name: a.Name, Description: a.Description}
	}

	result, err := h.create.Execute(c.Request.Context(), ranking.CreateRankingInput{
		Name: req.Name, Description: req.Description, Visibility: req.Visibility,
		Tags: req.Tags, Attributes: attrs, UserID: middleware.GetUserID(c),
	})
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.RespondCreated(c, result)
}

// Get handles GET /api/v1/rankings/:id.
//
//	@Summary		Get a ranking
//	@Description	Retrieve a ranking by ID
//	@Tags			rankings
//	@Produce		json
//	@Param			X-User-Id	header		string	false	"User ID"
//	@Param			id			path		string	true	"Ranking ID"
//	@Success		200			{object}	RankingResponse
//	@Failure		404			{object}	ErrorResponse
//	@Router			/rankings/{id} [get]
func (h *RankingHandler) Get(c *gin.Context) {
	result, err := h.get.Execute(c.Request.Context(), c.Param("id"), middleware.GetUserID(c))
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.RespondOK(c, result)
}

// Update handles PUT /api/v1/rankings/:id.
//
//	@Summary		Update a ranking
//	@Description	Update an existing ranking
//	@Tags			rankings
//	@Accept			json
//	@Produce		json
//	@Param			X-User-Id	header		string				true	"User ID"
//	@Param			id			path		string				true	"Ranking ID"
//	@Param			body		body		updateRankingRequest	true	"Ranking data"
//	@Success		200			{object}	RankingResponse
//	@Failure		400			{object}	ErrorResponse
//	@Failure		403			{object}	ErrorResponse
//	@Failure		404			{object}	ErrorResponse
//	@Router			/rankings/{id} [put]
func (h *RankingHandler) Update(c *gin.Context) {
	var req updateRankingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperror.NewValidationError(err.Error()))
		c.Abort()
		return
	}

	attrs := make([]ranking.AttributeInput, len(req.Attributes))
	for i, a := range req.Attributes {
		attrs[i] = ranking.AttributeInput{Name: a.Name, Description: a.Description}
	}

	result, err := h.update.Execute(c.Request.Context(), ranking.UpdateRankingInput{
		RankingID: c.Param("id"), UserID: middleware.GetUserID(c),
		Name: req.Name, Description: req.Description, Visibility: req.Visibility,
		Tags: req.Tags, Attributes: attrs,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.RespondOK(c, result)
}

// Delete handles DELETE /api/v1/rankings/:id.
//
//	@Summary		Delete a ranking
//	@Description	Delete a ranking by ID
//	@Tags			rankings
//	@Produce		json
//	@Param			X-User-Id	header	string	true	"User ID"
//	@Param			id			path	string	true	"Ranking ID"
//	@Success		204			"No Content"
//	@Failure		403			{object}	ErrorResponse
//	@Failure		404			{object}	ErrorResponse
//	@Router			/rankings/{id} [delete]
func (h *RankingHandler) Delete(c *gin.Context) {
	if err := h.del.Execute(c.Request.Context(), c.Param("id"), middleware.GetUserID(c)); err != nil {
		_ = c.Error(err)
		return
	}

	response.RespondNoContent(c)
}

// List handles GET /api/v1/rankings.
//
//	@Summary		List public rankings
//	@Description	List public rankings with pagination and sorting
//	@Tags			rankings
//	@Produce		json
//	@Param			limit	query		int		false	"Page size"		default(20)
//	@Param			cursor	query		string	false	"Pagination cursor"
//	@Param			sort	query		string	false	"Sort field (e.g. -createdAt)"
//	@Success		200		{object}	PaginatedRankingResponse
//	@Failure		400		{object}	ErrorResponse
//	@Router			/rankings [get]
func (h *RankingHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	result, err := h.listPublic.Execute(c.Request.Context(), limit, c.Query("cursor"), c.Query("sort"))
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.RespondOK(c, response.NewPaginatedResponse(result.Rankings, result.NextCursor))
}

// Recent handles GET /api/v1/rankings/recent.
//
//	@Summary		Get recent rankings
//	@Description	Retrieve the 10 most recent public rankings
//	@Tags			rankings
//	@Produce		json
//	@Success		200	{array}		RankingResponse
//	@Router			/rankings/recent [get]
func (h *RankingHandler) Recent(c *gin.Context) {
	result, err := h.recent.Execute(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.RespondOK(c, result)
}

// Search handles GET /api/v1/rankings/search.
//
//	@Summary		Search rankings
//	@Description	Search rankings by name or tag
//	@Tags			rankings
//	@Produce		json
//	@Param			q		query		string	false	"Search query"
//	@Param			tag		query		string	false	"Tag filter"
//	@Param			limit	query		int		false	"Page size"		default(20)
//	@Param			cursor	query		string	false	"Pagination cursor"
//	@Success		200		{object}	PaginatedRankingResponse
//	@Failure		400		{object}	ErrorResponse
//	@Router			/rankings/search [get]
func (h *RankingHandler) Search(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	result, err := h.search.Execute(c.Request.Context(), c.Query("q"), c.Query("tag"), limit, c.Query("cursor"))
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.RespondOK(c, response.NewPaginatedResponse(result.Rankings, result.NextCursor))
}

// Mine handles GET /api/v1/rankings/mine.
//
//	@Summary		List my rankings
//	@Description	List rankings owned by the current user
//	@Tags			rankings
//	@Produce		json
//	@Param			X-User-Id	header		string	true	"User ID"
//	@Param			limit		query		int		false	"Page size"		default(20)
//	@Param			cursor		query		string	false	"Pagination cursor"
//	@Success		200			{object}	PaginatedRankingResponse
//	@Failure		401			{object}	ErrorResponse
//	@Router			/rankings/mine [get]
func (h *RankingHandler) Mine(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	result, err := h.mine.Execute(c.Request.Context(), middleware.GetUserID(c), limit, c.Query("cursor"))
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.RespondOK(c, response.NewPaginatedResponse(result.Rankings, result.NextCursor))
}
