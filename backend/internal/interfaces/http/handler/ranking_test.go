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
	ranking "github.com/josimar/ranking/backend/internal/application/ranking"
	domain "github.com/josimar/ranking/backend/internal/domain/ranking"
	"github.com/josimar/ranking/backend/internal/interfaces/http/handler"
	"github.com/josimar/ranking/backend/internal/interfaces/http/middleware"
	"github.com/stretchr/testify/require"
)

const (
	testUserID    = "019bc9c7-b857-72c5-a491-7f9686c7989a"
	testRankingID = "019bc9c7-b857-72c5-a491-7f9686c7989b"
)

// --- Mocks ---

type mockCreateUC struct {
	output *ranking.Output
	err    error
	input  ranking.CreateRankingInput
}

func (m *mockCreateUC) Execute(ctx context.Context, input ranking.CreateRankingInput) (*ranking.Output, error) {
	m.input = input
	return m.output, m.err
}

type mockGetUC struct {
	output    *ranking.Output
	err       error
	rankingID string
	userID    string
}

func (m *mockGetUC) Execute(ctx context.Context, rankingID, userID string) (*ranking.Output, error) {
	m.rankingID = rankingID
	m.userID = userID
	return m.output, m.err
}

type mockUpdateUC struct {
	output *ranking.Output
	err    error
	input  ranking.UpdateRankingInput
}

func (m *mockUpdateUC) Execute(ctx context.Context, input ranking.UpdateRankingInput) (*ranking.Output, error) {
	m.input = input
	return m.output, m.err
}

type mockDeleteUC struct {
	err       error
	rankingID string
	userID    string
}

func (m *mockDeleteUC) Execute(ctx context.Context, rankingID, userID string) error {
	m.rankingID = rankingID
	m.userID = userID
	return m.err
}

type mockListPublicUC struct {
	output *ranking.PaginatedOutput
	err    error
	limit  int
	cursor string
	sort   string
}

func (m *mockListPublicUC) Execute(ctx context.Context, limit int, cursor, sort string) (*ranking.PaginatedOutput, error) {
	m.limit = limit
	m.cursor = cursor
	m.sort = sort
	return m.output, m.err
}

type mockRecentUC struct {
	output []*domain.Ranking
	err    error
}

func (m *mockRecentUC) Execute(ctx context.Context) ([]*domain.Ranking, error) {
	return m.output, m.err
}

type mockSearchUC struct {
	output *ranking.PaginatedOutput
	err    error
	q      string
	tag    string
	limit  int
	cursor string
}

func (m *mockSearchUC) Execute(ctx context.Context, q, tag string, limit int, cursor string) (*ranking.PaginatedOutput, error) {
	m.q = q
	m.tag = tag
	m.limit = limit
	m.cursor = cursor
	return m.output, m.err
}

type mockMyRankingsUC struct {
	output *ranking.PaginatedOutput
	err    error
	userID string
	limit  int
	cursor string
}

func (m *mockMyRankingsUC) Execute(ctx context.Context, userID string, limit int, cursor string) (*ranking.PaginatedOutput, error) {
	m.userID = userID
	m.limit = limit
	m.cursor = cursor
	return m.output, m.err
}

// --- Helpers ---

func sampleOutput() *ranking.Output {
	return &ranking.Output{
		ID:          testRankingID,
		Name:        "Video Games",
		Description: "Best consoles",
		Visibility:  "public",
		Tags:        []string{"games"},
		Attributes: []ranking.AttributeOutput{
			{ID: "attr-1", Name: "Graphics", Description: "Visual quality", Active: true},
		},
		OwnerUserID: testUserID,
		CreatedAt:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		IsOwner:     false,
	}
}

func setupRouter(h *handler.RankingHandler) *gin.Engine {
	r := gin.New()
	r.Use(middleware.UserID())
	r.Use(middleware.ErrorHandler())
	r.POST("/api/v1/rankings", middleware.RequireUserID(), h.Create)
	r.GET("/api/v1/rankings", h.List)
	r.GET("/api/v1/rankings/recent", h.Recent)
	r.GET("/api/v1/rankings/search", h.Search)
	r.GET("/api/v1/rankings/mine", middleware.RequireUserID(), h.Mine)
	r.GET("/api/v1/rankings/:id", h.Get)
	r.PUT("/api/v1/rankings/:id", middleware.RequireUserID(), h.Update)
	r.DELETE("/api/v1/rankings/:id", middleware.RequireUserID(), h.Delete)
	return r
}

func newHandler(
	create *mockCreateUC,
	get *mockGetUC,
	update *mockUpdateUC,
	del *mockDeleteUC,
	listPublic *mockListPublicUC,
	recent *mockRecentUC,
	search *mockSearchUC,
	myRankings *mockMyRankingsUC,
) *handler.RankingHandler {
	return handler.NewRankingHandler(create, get, update, del, listPublic, recent, search, myRankings)
}

// --- Tests ---

func TestCreate_Success(t *testing.T) {
	t.Parallel()
	out := sampleOutput()
	out.IsOwner = true
	mock := &mockCreateUC{output: out}
	h := newHandler(mock, nil, nil, nil, nil, nil, nil, nil)
	r := setupRouter(h)

	body := `{"name":"Video Games","description":"Best consoles","visibility":"public","tags":["games"],"attributes":[{"name":"Graphics","description":"Visual quality"}]}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/rankings", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Id", testUserID)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
	require.Equal(t, testUserID, mock.input.UserID)
	require.Equal(t, "Video Games", mock.input.Name)

	var resp ranking.Output
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.Equal(t, testRankingID, resp.ID)
}

func TestCreate_InvalidJSON(t *testing.T) {
	t.Parallel()
	mock := &mockCreateUC{}
	h := newHandler(mock, nil, nil, nil, nil, nil, nil, nil)
	r := setupRouter(h)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/rankings", bytes.NewBufferString(`{invalid`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Id", testUserID)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreate_UseCaseError(t *testing.T) {
	t.Parallel()
	mock := &mockCreateUC{err: errors.New("some error")}
	h := newHandler(mock, nil, nil, nil, nil, nil, nil, nil)
	r := setupRouter(h)

	body := `{"name":"Test","visibility":"public","attributes":[{"name":"A"}]}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/rankings", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Id", testUserID)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGet_Success(t *testing.T) {
	t.Parallel()
	out := sampleOutput()
	mock := &mockGetUC{output: out}
	h := newHandler(nil, mock, nil, nil, nil, nil, nil, nil)
	r := setupRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/rankings/"+testRankingID, nil)
	req.Header.Set("X-User-Id", testUserID)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, testRankingID, mock.rankingID)
	require.Equal(t, testUserID, mock.userID)
}

func TestGet_WithoutUserID(t *testing.T) {
	t.Parallel()
	out := sampleOutput()
	mock := &mockGetUC{output: out}
	h := newHandler(nil, mock, nil, nil, nil, nil, nil, nil)
	r := setupRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/rankings/"+testRankingID, nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "", mock.userID)
}

func TestGet_Error(t *testing.T) {
	t.Parallel()
	mock := &mockGetUC{err: errors.New("not found")}
	h := newHandler(nil, mock, nil, nil, nil, nil, nil, nil)
	r := setupRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/rankings/"+testRankingID, nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUpdate_Success(t *testing.T) {
	t.Parallel()
	out := sampleOutput()
	out.IsOwner = true
	mock := &mockUpdateUC{output: out}
	h := newHandler(nil, nil, mock, nil, nil, nil, nil, nil)
	r := setupRouter(h)

	body := `{"name":"Updated","description":"New desc","visibility":"private","tags":["new"],"attributes":[{"name":"Sound"}]}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/rankings/"+testRankingID, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Id", testUserID)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, testRankingID, mock.input.RankingID)
	require.Equal(t, testUserID, mock.input.UserID)
	require.Equal(t, "Updated", mock.input.Name)
}

func TestUpdate_InvalidJSON(t *testing.T) {
	t.Parallel()
	mock := &mockUpdateUC{}
	h := newHandler(nil, nil, mock, nil, nil, nil, nil, nil)
	r := setupRouter(h)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/rankings/"+testRankingID, bytes.NewBufferString(`{bad`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Id", testUserID)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDelete_Success(t *testing.T) {
	t.Parallel()
	mock := &mockDeleteUC{}
	h := newHandler(nil, nil, nil, mock, nil, nil, nil, nil)
	r := setupRouter(h)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/rankings/"+testRankingID, nil)
	req.Header.Set("X-User-Id", testUserID)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
	require.Equal(t, testRankingID, mock.rankingID)
	require.Equal(t, testUserID, mock.userID)
}

func TestDelete_Error(t *testing.T) {
	t.Parallel()
	mock := &mockDeleteUC{err: errors.New("fail")}
	h := newHandler(nil, nil, nil, mock, nil, nil, nil, nil)
	r := setupRouter(h)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/rankings/"+testRankingID, nil)
	req.Header.Set("X-User-Id", testUserID)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestList_Success(t *testing.T) {
	t.Parallel()
	mock := &mockListPublicUC{
		output: &ranking.PaginatedOutput{
			Rankings:   []*domain.Ranking{{ID: "r1", Name: "Test"}},
			NextCursor: "abc",
		},
	}
	h := newHandler(nil, nil, nil, nil, mock, nil, nil, nil)
	r := setupRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/rankings?limit=5&cursor=xyz&sort=-name", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, 5, mock.limit)
	require.Equal(t, "xyz", mock.cursor)
	require.Equal(t, "-name", mock.sort)

	var body map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	require.Contains(t, body, "data")
	require.Contains(t, body, "pagination")
}

func TestList_Error(t *testing.T) {
	t.Parallel()
	mock := &mockListPublicUC{err: errors.New("fail")}
	h := newHandler(nil, nil, nil, nil, mock, nil, nil, nil)
	r := setupRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/rankings", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestRecent_Success(t *testing.T) {
	t.Parallel()
	mock := &mockRecentUC{
		output: []*domain.Ranking{{ID: "r1", Name: "Recent"}},
	}
	h := newHandler(nil, nil, nil, nil, nil, mock, nil, nil)
	r := setupRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/rankings/recent", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}

func TestRecent_Error(t *testing.T) {
	t.Parallel()
	mock := &mockRecentUC{err: errors.New("fail")}
	h := newHandler(nil, nil, nil, nil, nil, mock, nil, nil)
	r := setupRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/rankings/recent", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestSearch_Success(t *testing.T) {
	t.Parallel()
	mock := &mockSearchUC{
		output: &ranking.PaginatedOutput{
			Rankings: []*domain.Ranking{{ID: "r1"}},
		},
	}
	h := newHandler(nil, nil, nil, nil, nil, nil, mock, nil)
	r := setupRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/rankings/search?q=games&limit=10&cursor=c1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "games", mock.q)
	require.Equal(t, 10, mock.limit)
	require.Equal(t, "c1", mock.cursor)
}

func TestSearch_WithTag(t *testing.T) {
	t.Parallel()
	mock := &mockSearchUC{
		output: &ranking.PaginatedOutput{
			Rankings: []*domain.Ranking{{ID: "r1"}},
		},
	}
	h := newHandler(nil, nil, nil, nil, nil, nil, mock, nil)
	r := setupRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/rankings/search?tag=movies", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "movies", mock.tag)
}

func TestSearch_Error(t *testing.T) {
	t.Parallel()
	mock := &mockSearchUC{err: errors.New("fail")}
	h := newHandler(nil, nil, nil, nil, nil, nil, mock, nil)
	r := setupRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/rankings/search?q=test", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestMine_Success(t *testing.T) {
	t.Parallel()
	mock := &mockMyRankingsUC{
		output: &ranking.PaginatedOutput{
			Rankings:   []*domain.Ranking{{ID: "r1"}},
			NextCursor: "next",
		},
	}
	h := newHandler(nil, nil, nil, nil, nil, nil, nil, mock)
	r := setupRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/rankings/mine?limit=15&cursor=c2", nil)
	req.Header.Set("X-User-Id", testUserID)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, testUserID, mock.userID)
	require.Equal(t, 15, mock.limit)
	require.Equal(t, "c2", mock.cursor)
}

func TestMine_Error(t *testing.T) {
	t.Parallel()
	mock := &mockMyRankingsUC{err: errors.New("fail")}
	h := newHandler(nil, nil, nil, nil, nil, nil, nil, mock)
	r := setupRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/rankings/mine", nil)
	req.Header.Set("X-User-Id", testUserID)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
}
