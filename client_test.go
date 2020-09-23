package expensify

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	partnerUserID     = "testUserID"
	partnerUserSecret = "testUserSecret-123-abc"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient(partnerUserID, partnerUserSecret)
	require.NoError(t, err)
	require.NotNil(t, client)

	// Are endpoints/resources present?
	assert.Equal(t, &expenseService{client: client}, client.Expense)

	// Is default configuration present?
	if assert.NotNil(t, client.baseURL) {
		assert.Equal(t, baseURL, client.baseURL.String())
	}
	assert.NotEmpty(t, client.userAgent)

	assert.Equal(t, partnerUserID, client.partnerUserID)
	assert.Equal(t, partnerUserSecret, client.partnerUserSecret)

	assert.Equal(t, defaultHTTPClient, client.httpClient)
}

func TestDo(t *testing.T) {
	hf := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		_, _ = fmt.Fprint(w, `{"A":"a"}`)
	}

	client, teardown := setup(t, hf)
	defer teardown()

	type foo struct {
		A string
	}

	req, err := client.newRequest(context.Background(), "", "", nil)
	require.NoError(t, err)

	var body foo
	err = client.do(req, &body)
	require.NoError(t, err)

	assert.Equal(t, foo{"a"}, body)
}

func TestDo_HTTPError(t *testing.T) {
	hf := func(w http.ResponseWriter, r *http.Request) {
		httpErr := Error{
			StatusCode: http.StatusBadRequest,
			Message:    "This was a bad request!",
		}
		err := json.NewEncoder(w).Encode(httpErr)
		require.NoError(t, err)
	}

	client, teardown := setup(t, hf)
	defer teardown()

	req, err := client.newRequest(context.Background(), "", "", nil)
	require.NoError(t, err)

	err = client.do(req, nil)
	require.EqualError(t, err, "Bad Request: This was a bad request!")
}

func TestDo_RedirectLoop(t *testing.T) {
	hf := func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	client, teardown := setup(t, hf)
	defer teardown()

	req, err := client.newRequest(context.Background(), "", "", nil)
	require.NoError(t, err)

	err = client.do(req, nil)
	require.Error(t, err)

	assert.IsType(t, err, &url.Error{})
}

// setup sets up a test HTTP server along with a client that is configured to
// talk to that test server. Tests should pass a handler function which provides
// the response for the API method being tested.
func setup(t *testing.T, hf http.HandlerFunc) (*Client, func()) {
	t.Helper()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hf.ServeHTTP(w, r)
	}))

	client, err := NewClient(partnerUserID, partnerUserSecret, SetClient(srv.Client()))
	require.NoError(t, err)

	srvURL, err := url.ParseRequestURI(srv.URL)
	require.NoError(t, err)

	client.baseURL = srvURL

	return client, func() { srv.Close() }
}
