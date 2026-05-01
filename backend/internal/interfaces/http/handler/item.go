package handler

import (
	"context"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
	item "github.com/josimar/ranking/backend/internal/application/item"
	"github.com/josimar/ranking/backend/internal/interfaces/http/middleware"
	"github.com/josimar/ranking/backend/internal/interfaces/http/response"
	"github.com/josimar/ranking/backend/pkg/apperror"
)

// AddItemExecutor adds an item to a ranking.
type AddItemExecutor interface {
	Execute(ctx context.Context, input item.AddItemInput) (*item.Output, error)
}

// ListItemsExecutor lists items for a ranking.
type ListItemsExecutor interface {
	Execute(ctx context.Context, input item.ListItemsInput) ([]item.ListOutput, error)
}

// UpdateItemExecutor updates an item.
type UpdateItemExecutor interface {
	Execute(ctx context.Context, input item.UpdateItemInput) (*item.Output, error)
}

// DeleteItemExecutor deletes an item.
type DeleteItemExecutor interface {
	Execute(ctx context.Context, rankingID, itemID, userID string) error
}

// ItemHandler handles item HTTP requests.
type ItemHandler struct {
	add    AddItemExecutor
	list   ListItemsExecutor
	update UpdateItemExecutor
	del    DeleteItemExecutor
}

// NewItemHandler creates a new ItemHandler.
func NewItemHandler(
	add AddItemExecutor,
	list ListItemsExecutor,
	update UpdateItemExecutor,
	del DeleteItemExecutor,
) *ItemHandler {
	return &ItemHandler{add: add, list: list, update: update, del: del}
}

type createItemRequest struct {
	Name     string `json:"name" binding:"required"`
	ImageURL string `json:"imageUrl"`
}

type updateItemRequest struct {
	Name string `json:"name" binding:"required"`
}

// Create handles POST /api/v1/rankings/:rankingId/items.
//
//	@Summary		Add an item to a ranking
//	@Description	Add a new item with optional image (JSON or multipart)
//	@Tags			items
//	@Accept			json,mpfd
//	@Produce		json
//	@Param			X-User-Id	header		string				true	"User ID"
//	@Param			rankingId	path		string				true	"Ranking ID"
//	@Param			body		body		createItemRequest	false	"Item data (JSON)"
//	@Success		201			{object}	ItemResponse
//	@Failure		400			{object}	ErrorResponse
//	@Failure		404			{object}	ErrorResponse
//	@Router			/rankings/{rankingId}/items [post]
func (h *ItemHandler) Create(c *gin.Context) {
	rankingID := c.Param("rankingId")
	userID := middleware.GetUserID(c)

	input := item.AddItemInput{
		RankingID: rankingID,
		UserID:    userID,
	}

	if isMultipart(c) {
		name := c.PostForm("name")
		if name == "" {
			_ = c.Error(apperror.NewValidationError("name is required"))
			c.Abort()
			return
		}
		input.Name = name

		file, _, err := c.Request.FormFile("image")
		if err == nil {
			defer func() { _ = file.Close() }()
			data, readErr := io.ReadAll(file)
			if readErr != nil {
				_ = c.Error(readErr)
				return
			}
			input.ImageData = data
		}
	} else {
		var req createItemRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			_ = c.Error(apperror.NewValidationError(err.Error()))
			c.Abort()
			return
		}
		input.Name = req.Name
		input.ImageURL = req.ImageURL
	}

	result, err := h.add.Execute(c.Request.Context(), input)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.RespondCreated(c, result)
}

// List handles GET /api/v1/rankings/:rankingId/items.
//
//	@Summary		List items in a ranking
//	@Description	List items with computed scores (avg or user mode)
//	@Tags			items
//	@Produce		json
//	@Param			X-User-Id	header		string	false	"User ID"
//	@Param			rankingId	path		string	true	"Ranking ID"
//	@Param			mode		query		string	false	"Score mode (avg or user)"	default(avg)
//	@Param			sort		query		string	false	"Sort field (e.g. -overall)"
//	@Success		200			{array}		ItemListResponse
//	@Failure		400			{object}	ErrorResponse
//	@Failure		404			{object}	ErrorResponse
//	@Router			/rankings/{rankingId}/items [get]
func (h *ItemHandler) List(c *gin.Context) {
	result, err := h.list.Execute(c.Request.Context(), item.ListItemsInput{
		RankingID: c.Param("rankingId"),
		Mode:      c.DefaultQuery("mode", "avg"),
		Sort:      c.Query("sort"),
		UserID:    middleware.GetUserID(c),
	})
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.RespondOK(c, result)
}

// Update handles PUT /api/v1/rankings/:rankingId/items/:itemId.
//
//	@Summary		Update an item
//	@Description	Update an existing item with optional image (JSON or multipart)
//	@Tags			items
//	@Accept			json,mpfd
//	@Produce		json
//	@Param			X-User-Id	header		string				true	"User ID"
//	@Param			rankingId	path		string				true	"Ranking ID"
//	@Param			itemId		path		string				true	"Item ID"
//	@Param			body		body		updateItemRequest	false	"Item data (JSON)"
//	@Success		200			{object}	ItemResponse
//	@Failure		400			{object}	ErrorResponse
//	@Failure		403			{object}	ErrorResponse
//	@Failure		404			{object}	ErrorResponse
//	@Router			/rankings/{rankingId}/items/{itemId} [put]
func (h *ItemHandler) Update(c *gin.Context) {
	rankingID := c.Param("rankingId")
	itemID := c.Param("itemId")
	userID := middleware.GetUserID(c)

	input := item.UpdateItemInput{
		RankingID: rankingID,
		ItemID:    itemID,
		UserID:    userID,
	}

	if isMultipart(c) {
		name := c.PostForm("name")
		if name == "" {
			_ = c.Error(apperror.NewValidationError("name is required"))
			c.Abort()
			return
		}
		input.Name = name

		file, _, err := c.Request.FormFile("image")
		if err == nil {
			defer func() { _ = file.Close() }()
			data, readErr := io.ReadAll(file)
			if readErr != nil {
				_ = c.Error(readErr)
				return
			}
			input.ImageData = data
		}
	} else {
		var req updateItemRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			_ = c.Error(apperror.NewValidationError(err.Error()))
			c.Abort()
			return
		}
		input.Name = req.Name
	}

	result, err := h.update.Execute(c.Request.Context(), input)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.RespondOK(c, result)
}

// Delete handles DELETE /api/v1/rankings/:rankingId/items/:itemId.
//
//	@Summary		Delete an item
//	@Description	Delete an item and its associated ratings and images
//	@Tags			items
//	@Produce		json
//	@Param			X-User-Id	header	string	true	"User ID"
//	@Param			rankingId	path	string	true	"Ranking ID"
//	@Param			itemId		path	string	true	"Item ID"
//	@Success		204			"No Content"
//	@Failure		403			{object}	ErrorResponse
//	@Failure		404			{object}	ErrorResponse
//	@Router			/rankings/{rankingId}/items/{itemId} [delete]
func (h *ItemHandler) Delete(c *gin.Context) {
	if err := h.del.Execute(
		c.Request.Context(),
		c.Param("rankingId"),
		c.Param("itemId"),
		middleware.GetUserID(c),
	); err != nil {
		_ = c.Error(err)
		return
	}

	response.RespondNoContent(c)
}

func isMultipart(c *gin.Context) bool {
	ct := c.ContentType()
	return strings.HasPrefix(ct, "multipart/form-data")
}
