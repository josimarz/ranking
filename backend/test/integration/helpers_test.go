//go:build integration

package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func doRequest(t *testing.T, env *testEnv, method, path string, body any, userID string) *http.Response {
	t.Helper()

	url := fmt.Sprintf("%s%s", env.server.URL, path)

	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		require.NoError(t, err)
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(env.ctx, method, url, bodyReader)
	require.NoError(t, err)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if userID != "" {
		req.Header.Set("X-User-Id", userID)
	}

	resp, err := env.client.Do(req)
	require.NoError(t, err)

	return resp
}

func readBody(t *testing.T, resp *http.Response) []byte {
	t.Helper()
	defer func() { _ = resp.Body.Close() }()

	data, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return data
}

func decodeJSON(t *testing.T, resp *http.Response, target any) {
	t.Helper()

	data := readBody(t, resp)
	require.NoError(t, json.Unmarshal(data, target))
}
