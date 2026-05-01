package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	domainrating "github.com/josimar/ranking/backend/internal/domain/rating"
	"github.com/josimar/ranking/backend/internal/interfaces/http/handler"
	"github.com/stretchr/testify/require"
)

type mockSubmitRating struct {
	result *domainrating.Rating
	err    error
}

func (m *mockSubmitRating) Execute(_ context.Context, _, _, _ string, _ map[string]int) (*domainrating.Rating, error) {
	return m.result, m.err
}

type mockGetMyRating struct {
	result *domainrating.Rating
	err    error
}

func (m *mockGetMyRating) Execute(_ context.Context, _, _, _ string) (*domainrating.Rating, error) {
	return m.result, m.err
}

func setupRatingContext(method, path string, body []byte, rankingID, itemID, userID string) (*httptest.ResponseRecorder, *gin.Context) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	var req *http.Request
	if body != nil {
		req = httptest.NewRequest(method, path, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}

	c.Request = req
	c.Params = gin.Params{
		{Key: "rankingId", Value: rankingID},
		{Key: "itemId", Value: itemID},
	}
	if userID != "" {
		c.Set("userId", userID)
	}

	return w, c
}

func TestRatingHandler_Submit_Success(t *testing.T) {
	t.Parallel()

	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	scores := map[string]int{"attr1": 80, "attr2": 90}

	mock := &mockSubmitRating{
		result: &domainrating.Rating{
			RankingID: "ranking-1",
			ItemID:    "item-1",
			UserID:    "user-1",
			Scores:    scores,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	h := handler.NewRatingHandler(mock, nil)

	body, _ := json.Marshal(handler.SubmitRatingRequest{Scores: scores})
	w, c := setupRatingContext(http.MethodPost, "/api/v1/rankings/ranking-1/items/item-1/ratings", body, "ranking-1", "item-1", "user-1")

	h.Submit(c)

	require.Equal(t, http.StatusOK, w.Code)

	var resp domainrating.Rating
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	require.Equal(t, "ranking-1", resp.RankingID)
	require.Equal(t, scores, resp.Scores)
}

func TestRatingHandler_Submit_UseCaseError(t *testing.T) {
	t.Parallel()

	mock := &mockSubmitRating{err: errors.New("use case failed")}
	h := handler.NewRatingHandler(mock, nil)

	body, _ := json.Marshal(handler.SubmitRatingRequest{Scores: map[string]int{"a": 1}})
	w, c := setupRatingContext(http.MethodPost, "/api/v1/rankings/r1/items/i1/ratings", body, "r1", "i1", "user-1")

	h.Submit(c)

	require.Len(t, c.Errors, 1)
	require.Equal(t, "use case failed", c.Errors[0].Err.Error())
	require.Equal(t, http.StatusOK, w.Code) // error middleware would handle status
}

func TestRatingHandler_Submit_InvalidJSON(t *testing.T) {
	t.Parallel()

	h := handler.NewRatingHandler(&mockSubmitRating{}, nil)

	_, c := setupRatingContext(http.MethodPost, "/api/v1/rankings/r1/items/i1/ratings", []byte(`{invalid`), "r1", "i1", "user-1")

	h.Submit(c)

	require.Len(t, c.Errors, 1)
	require.True(t, c.IsAborted())
}

func TestRatingHandler_GetMine_Success(t *testing.T) {
	t.Parallel()

	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	scores := map[string]int{"attr1": 75}

	mock := &mockGetMyRating{
		result: &domainrating.Rating{
			RankingID: "ranking-1",
			ItemID:    "item-1",
			UserID:    "user-1",
			Scores:    scores,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	h := handler.NewRatingHandler(nil, mock)

	w, c := setupRatingContext(http.MethodGet, "/api/v1/rankings/ranking-1/items/item-1/ratings/mine", nil, "ranking-1", "item-1", "user-1")

	h.GetMine(c)

	require.Equal(t, http.StatusOK, w.Code)

	var resp domainrating.Rating
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	require.Equal(t, "ranking-1", resp.RankingID)
	require.Equal(t, scores, resp.Scores)
}

func TestRatingHandler_GetMine_NotRated(t *testing.T) {
	t.Parallel()

	mock := &mockGetMyRating{result: nil}
	h := handler.NewRatingHandler(nil, mock)

	w, c := setupRatingContext(http.MethodGet, "/api/v1/rankings/r1/items/i1/ratings/mine", nil, "r1", "i1", "user-1")

	h.GetMine(c)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "null", w.Body.String())
}

func TestRatingHandler_GetMine_UseCaseError(t *testing.T) {
	t.Parallel()

	mock := &mockGetMyRating{err: errors.New("db error")}
	h := handler.NewRatingHandler(nil, mock)

	_, c := setupRatingContext(http.MethodGet, "/api/v1/rankings/r1/items/i1/ratings/mine", nil, "r1", "i1", "user-1")

	h.GetMine(c)

	require.Len(t, c.Errors, 1)
	require.Equal(t, "db error", c.Errors[0].Err.Error())
}
