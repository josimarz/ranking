package handler_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	item "github.com/josimar/ranking/backend/internal/application/item"
	"github.com/josimar/ranking/backend/internal/interfaces/http/handler"
	"github.com/josimar/ranking/backend/internal/interfaces/http/middleware"
	"github.com/stretchr/testify/require"
)

// --- Item Mocks ---

type mockAddItem struct {
	output *item.Output
	err    error
	input  item.AddItemInput
}

func (m *mockAddItem) Execute(_ context.Context, input item.AddItemInput) (*item.Output, error) {
	m.input = input
	return m.output, m.err
}

type mockListItems struct {
	output []item.ListOutput
	err    error
	input  item.ListItemsInput
}

func (m *mockListItems) Execute(_ context.Context, input item.ListItemsInput) ([]item.ListOutput, error) {
	m.input = input
	return m.output, m.err
}

type mockUpdateItem struct {
	output *item.Output
	err    error
	input  item.UpdateItemInput
}

func (m *mockUpdateItem) Execute(_ context.Context, input item.UpdateItemInput) (*item.Output, error) {
	m.input = input
	return m.output, m.err
}

type mockDeleteItem struct {
	err       error
	rankingID string
	itemID    string
	userID    string
}

func (m *mockDeleteItem) Execute(_ context.Context, rankingID, itemID, userID string) error {
	m.rankingID = rankingID
	m.itemID = itemID
	m.userID = userID
	return m.err
}

// --- Helpers ---

const testItemID = "019bc9c7-b857-72c5-a491-7f9686c7989c"

func sampleItemOutput() *item.Output {
	return &item.Output{
		ID:        testItemID,
		Name:      "Sega Mega Drive",
		RankingID: testRankingID,
		CreatedBy: testUserID,
		CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
	}
}

func setupItemRouter(h *handler.ItemHandler) *gin.Engine {
	r := gin.New()
	r.Use(middleware.UserID())
	r.Use(middleware.ErrorHandler())
	g := r.Group("/api/v1/rankings/:rankingId/items")
	{
		g.POST("", middleware.RequireUserID(), h.Create)
		g.GET("", h.List)
		g.PUT("/:itemId", middleware.RequireUserID(), h.Update)
		g.DELETE("/:itemId", middleware.RequireUserID(), h.Delete)
	}
	return r
}

// --- Create Tests ---

func TestItemCreate_JSONSuccess(t *testing.T) {
	t.Parallel()
	out := sampleItemOutput()
	mock := &mockAddItem{output: out}
	h := handler.NewItemHandler(mock, nil, nil, nil)
	r := setupItemRouter(h)

	body := `{"name":"Sega Mega Drive","imageUrl":"https://example.com/img.png"}`
	url := fmt.Sprintf("/api/v1/rankings/%s/items", testRankingID)
	req := httptest.NewRequest(http.MethodPost, url, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Id", testUserID)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
	require.Equal(t, testRankingID, mock.input.RankingID)
	require.Equal(t, testUserID, mock.input.UserID)
	require.Equal(t, "Sega Mega Drive", mock.input.Name)
	require.Equal(t, "https://example.com/img.png", mock.input.ImageURL)
}

func TestItemCreate_MultipartSuccess(t *testing.T) {
	t.Parallel()
	out := sampleItemOutput()
	mock := &mockAddItem{output: out}
	h := handler.NewItemHandler(mock, nil, nil, nil)
	r := setupItemRouter(h)

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	_ = writer.WriteField("name", "Sega Mega Drive")
	part, _ := writer.CreateFormFile("image", "test.png")
	_, _ = part.Write([]byte("fake-image-data"))
	require.NoError(t, writer.Close())

	url := fmt.Sprintf("/api/v1/rankings/%s/items", testRankingID)
	req := httptest.NewRequest(http.MethodPost, url, &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-User-Id", testUserID)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
	require.Equal(t, "Sega Mega Drive", mock.input.Name)
	require.Equal(t, []byte("fake-image-data"), mock.input.ImageData)
}

func TestItemCreate_InvalidJSON(t *testing.T) {
	t.Parallel()
	h := handler.NewItemHandler(&mockAddItem{}, nil, nil, nil)
	r := setupItemRouter(h)

	url := fmt.Sprintf("/api/v1/rankings/%s/items", testRankingID)
	req := httptest.NewRequest(http.MethodPost, url, bytes.NewBufferString(`{invalid`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Id", testUserID)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestItemCreate_UseCaseError(t *testing.T) {
	t.Parallel()
	mock := &mockAddItem{err: errors.New("fail")}
	h := handler.NewItemHandler(mock, nil, nil, nil)
	r := setupItemRouter(h)

	body := `{"name":"Test"}`
	url := fmt.Sprintf("/api/v1/rankings/%s/items", testRankingID)
	req := httptest.NewRequest(http.MethodPost, url, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Id", testUserID)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestItemCreate_MissingName(t *testing.T) {
	t.Parallel()
	h := handler.NewItemHandler(&mockAddItem{}, nil, nil, nil)
	r := setupItemRouter(h)

	body := `{"imageUrl":"https://example.com/img.png"}`
	url := fmt.Sprintf("/api/v1/rankings/%s/items", testRankingID)
	req := httptest.NewRequest(http.MethodPost, url, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Id", testUserID)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

// --- List Tests ---

func TestItemList_Success(t *testing.T) {
	t.Parallel()
	mock := &mockListItems{
		output: []item.ListOutput{
			{ID: testItemID, Name: "Sega", Overall: 85.0, Scores: map[string]float64{"Graphics": 90}},
		},
	}
	h := handler.NewItemHandler(nil, mock, nil, nil)
	r := setupItemRouter(h)

	url := fmt.Sprintf("/api/v1/rankings/%s/items?mode=avg&sort=-overall", testRankingID)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("X-User-Id", testUserID)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, testRankingID, mock.input.RankingID)
	require.Equal(t, "avg", mock.input.Mode)
	require.Equal(t, "-overall", mock.input.Sort)
	require.Equal(t, testUserID, mock.input.UserID)
}

func TestItemList_DefaultMode(t *testing.T) {
	t.Parallel()
	mock := &mockListItems{output: []item.ListOutput{}}
	h := handler.NewItemHandler(nil, mock, nil, nil)
	r := setupItemRouter(h)

	url := fmt.Sprintf("/api/v1/rankings/%s/items", testRankingID)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "avg", mock.input.Mode)
}

func TestItemList_UserMode(t *testing.T) {
	t.Parallel()
	mock := &mockListItems{output: []item.ListOutput{}}
	h := handler.NewItemHandler(nil, mock, nil, nil)
	r := setupItemRouter(h)

	url := fmt.Sprintf("/api/v1/rankings/%s/items?mode=user", testRankingID)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("X-User-Id", testUserID)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "user", mock.input.Mode)
	require.Equal(t, testUserID, mock.input.UserID)
}

func TestItemList_Error(t *testing.T) {
	t.Parallel()
	mock := &mockListItems{err: errors.New("fail")}
	h := handler.NewItemHandler(nil, mock, nil, nil)
	r := setupItemRouter(h)

	url := fmt.Sprintf("/api/v1/rankings/%s/items", testRankingID)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
}

// --- Update Tests ---

func TestItemUpdate_JSONSuccess(t *testing.T) {
	t.Parallel()
	out := sampleItemOutput()
	out.Name = "Updated Name"
	mock := &mockUpdateItem{output: out}
	h := handler.NewItemHandler(nil, nil, mock, nil)
	r := setupItemRouter(h)

	body := `{"name":"Updated Name"}`
	url := fmt.Sprintf("/api/v1/rankings/%s/items/%s", testRankingID, testItemID)
	req := httptest.NewRequest(http.MethodPut, url, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Id", testUserID)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, testRankingID, mock.input.RankingID)
	require.Equal(t, testItemID, mock.input.ItemID)
	require.Equal(t, testUserID, mock.input.UserID)
	require.Equal(t, "Updated Name", mock.input.Name)
}

func TestItemUpdate_MultipartSuccess(t *testing.T) {
	t.Parallel()
	out := sampleItemOutput()
	mock := &mockUpdateItem{output: out}
	h := handler.NewItemHandler(nil, nil, mock, nil)
	r := setupItemRouter(h)

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	_ = writer.WriteField("name", "Updated")
	part, _ := writer.CreateFormFile("image", "new.png")
	_, _ = part.Write([]byte("new-image"))
	require.NoError(t, writer.Close())

	url := fmt.Sprintf("/api/v1/rankings/%s/items/%s", testRankingID, testItemID)
	req := httptest.NewRequest(http.MethodPut, url, &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-User-Id", testUserID)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "Updated", mock.input.Name)
	require.Equal(t, []byte("new-image"), mock.input.ImageData)
}

func TestItemUpdate_InvalidJSON(t *testing.T) {
	t.Parallel()
	h := handler.NewItemHandler(nil, nil, &mockUpdateItem{}, nil)
	r := setupItemRouter(h)

	url := fmt.Sprintf("/api/v1/rankings/%s/items/%s", testRankingID, testItemID)
	req := httptest.NewRequest(http.MethodPut, url, bytes.NewBufferString(`{bad`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Id", testUserID)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestItemUpdate_UseCaseError(t *testing.T) {
	t.Parallel()
	mock := &mockUpdateItem{err: errors.New("fail")}
	h := handler.NewItemHandler(nil, nil, mock, nil)
	r := setupItemRouter(h)

	body := `{"name":"Test"}`
	url := fmt.Sprintf("/api/v1/rankings/%s/items/%s", testRankingID, testItemID)
	req := httptest.NewRequest(http.MethodPut, url, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Id", testUserID)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
}

// --- Delete Tests ---

func TestItemDelete_Success(t *testing.T) {
	t.Parallel()
	mock := &mockDeleteItem{}
	h := handler.NewItemHandler(nil, nil, nil, mock)
	r := setupItemRouter(h)

	url := fmt.Sprintf("/api/v1/rankings/%s/items/%s", testRankingID, testItemID)
	req := httptest.NewRequest(http.MethodDelete, url, nil)
	req.Header.Set("X-User-Id", testUserID)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
	require.Equal(t, testRankingID, mock.rankingID)
	require.Equal(t, testItemID, mock.itemID)
	require.Equal(t, testUserID, mock.userID)
}

func TestItemDelete_Error(t *testing.T) {
	t.Parallel()
	mock := &mockDeleteItem{err: errors.New("fail")}
	h := handler.NewItemHandler(nil, nil, nil, mock)
	r := setupItemRouter(h)

	url := fmt.Sprintf("/api/v1/rankings/%s/items/%s", testRankingID, testItemID)
	req := httptest.NewRequest(http.MethodDelete, url, nil)
	req.Header.Set("X-User-Id", testUserID)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestItemDelete_MissingUserID(t *testing.T) {
	t.Parallel()
	mock := &mockDeleteItem{}
	h := handler.NewItemHandler(nil, nil, nil, mock)
	r := setupItemRouter(h)

	url := fmt.Sprintf("/api/v1/rankings/%s/items/%s", testRankingID, testItemID)
	req := httptest.NewRequest(http.MethodDelete, url, nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)
}
