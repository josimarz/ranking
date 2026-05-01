package handler

import (
	"context"

	"github.com/gin-gonic/gin"
	domainrating "github.com/josimar/ranking/backend/internal/domain/rating"
	"github.com/josimar/ranking/backend/internal/interfaces/http/middleware"
	"github.com/josimar/ranking/backend/internal/interfaces/http/response"
	"github.com/josimar/ranking/backend/pkg/apperror"
)

// SubmitRatingExecutor submits a user's rating for an item.
type SubmitRatingExecutor interface {
	Execute(ctx context.Context, rankingID, itemID, userID string, scores map[string]int) (*domainrating.Rating, error)
}

// GetMyRatingExecutor retrieves the current user's rating for an item.
type GetMyRatingExecutor interface {
	Execute(ctx context.Context, rankingID, itemID, userID string) (*domainrating.Rating, error)
}

// SubmitRatingRequest is the request body for submitting a rating.
type SubmitRatingRequest struct {
	Scores map[string]int `json:"scores" binding:"required"`
}

// RatingHandler handles rating HTTP requests.
type RatingHandler struct {
	submit  SubmitRatingExecutor
	getMine GetMyRatingExecutor
}

// NewRatingHandler creates a new RatingHandler.
func NewRatingHandler(submit SubmitRatingExecutor, getMine GetMyRatingExecutor) *RatingHandler {
	return &RatingHandler{submit: submit, getMine: getMine}
}

// Submit handles POST requests to submit a rating.
//
//	@Summary		Submit a rating
//	@Description	Submit or update scores for an item in a ranking
//	@Tags			ratings
//	@Accept			json
//	@Produce		json
//	@Param			X-User-Id	header		string				true	"User ID"
//	@Param			rankingId	path		string				true	"Ranking ID"
//	@Param			itemId		path		string				true	"Item ID"
//	@Param			body		body		SubmitRatingRequest	true	"Rating scores"
//	@Success		200			{object}	RatingResponse
//	@Failure		400			{object}	ErrorResponse
//	@Failure		404			{object}	ErrorResponse
//	@Router			/rankings/{rankingId}/items/{itemId}/rating [put]
func (h *RatingHandler) Submit(c *gin.Context) {
	rankingID := c.Param("rankingId")
	itemID := c.Param("itemId")
	userID := middleware.GetUserID(c)

	var req SubmitRatingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperror.NewValidationError(err.Error()))
		c.Abort()
		return
	}

	result, err := h.submit.Execute(c.Request.Context(), rankingID, itemID, userID, req.Scores)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.RespondOK(c, result)
}

// GetMine handles GET requests to retrieve the current user's rating.
//
//	@Summary		Get my rating
//	@Description	Retrieve the current user's rating for an item
//	@Tags			ratings
//	@Produce		json
//	@Param			X-User-Id	header		string	true	"User ID"
//	@Param			rankingId	path		string	true	"Ranking ID"
//	@Param			itemId		path		string	true	"Item ID"
//	@Success		200			{object}	RatingResponse
//	@Failure		404			{object}	ErrorResponse
//	@Router			/rankings/{rankingId}/items/{itemId}/rating/mine [get]
func (h *RatingHandler) GetMine(c *gin.Context) {
	rankingID := c.Param("rankingId")
	itemID := c.Param("itemId")
	userID := middleware.GetUserID(c)

	result, err := h.getMine.Execute(c.Request.Context(), rankingID, itemID, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.RespondOK(c, result)
}
