package http

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	methodGET    = "GET"
	methodPOST   = "POST"
	methodPUT    = "PUT"
	methodDELETE = "DELETE"

	pathRankingByID = "/api/v1/rankings/:id"
)

func TestNewRouter_RegistersAllRoutes(t *testing.T) {
	t.Parallel()

	router := NewRouter(nil, nil, nil, nil, "*", "local")

	// Gin uses :id for all routes at the rankings/:id level.
	// Item/rating handlers receive rankingId via paramAlias middleware.
	expectedRoutes := []struct {
		method string
		path   string
	}{
		{methodGET, "/api/v1/health"},
		{methodPOST, "/api/v1/rankings"},
		{methodGET, "/api/v1/rankings"},
		{methodGET, "/api/v1/rankings/recent"},
		{methodGET, "/api/v1/rankings/search"},
		{methodGET, "/api/v1/rankings/mine"},
		{methodGET, pathRankingByID},
		{methodPUT, pathRankingByID},
		{methodDELETE, pathRankingByID},
		{methodPOST, "/api/v1/rankings/:id/items"},
		{methodGET, "/api/v1/rankings/:id/items"},
		{methodPUT, "/api/v1/rankings/:id/items/:itemId"},
		{methodDELETE, "/api/v1/rankings/:id/items/:itemId"},
		{methodPUT, "/api/v1/rankings/:id/items/:itemId/rating"},
		{methodGET, "/api/v1/rankings/:id/items/:itemId/rating/mine"},
	}

	routes := router.Routes()
	registered := make(map[string]bool)
	for _, r := range routes {
		key := fmt.Sprintf("%s %s", r.Method, r.Path)
		registered[key] = true
	}

	for _, exp := range expectedRoutes {
		key := fmt.Sprintf("%s %s", exp.method, exp.path)
		assert.True(t, registered[key], fmt.Sprintf("route not registered: %s", key))
	}

	require.Len(t, expectedRoutes, 15)
}
