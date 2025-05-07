package handlers

import (
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/edkadigital/startmeup/config"
	"github.com/edkadigital/startmeup/pkg/services"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

var (
	srv *httptest.Server
	c   *services.Container
)

func TestMain(m *testing.M) {
	// Set the environment to test
	config.SwitchEnvironment(config.EnvTest)

	// Start a new container
	c = services.NewContainer()

	// Start a test HTTP server
	if err := BuildRouter(c); err != nil {
		panic(err)
	}
	srv = httptest.NewServer(c.Web)

	// Run tests
	exitVal := m.Run()

	// Shutdown the container and test server
	if err := c.Shutdown(); err != nil {
		panic(err)
	}
	srv.Close()

	os.Exit(exitVal)
}

type httpRequest struct {
	route  string
	client http.Client
	body   url.Values
	t      *testing.T
}

func request(t *testing.T) *httpRequest {
	jar, err := cookiejar.New(nil)
	require.NoError(t, err)
	r := httpRequest{
		t:    t,
		body: url.Values{},
		client: http.Client{
			Jar: jar,
		},
	}
	return &r
}

func (h *httpRequest) setRoute(route string, params ...any) *httpRequest {
	h.route = srv.URL + c.Web.Reverse(route, params)
	return h
}

func (h *httpRequest) get() *httpResponse {
	resp, err := h.client.Get(h.route)
	require.NoError(h.t, err)
	r := httpResponse{
		t:        h.t,
		Response: resp,
	}
	return &r
}

type httpResponse struct {
	*http.Response
	t *testing.T
}

func (h *httpResponse) assertStatusCode(code int) *httpResponse {
	assert.Equal(h.t, code, h.StatusCode)
	return h
}

func (h *httpResponse) toDoc() *goquery.Document {
	doc, err := goquery.NewDocumentFromReader(h.Body)
	require.NoError(h.t, err)
	err = h.Body.Close()
	assert.NoError(h.t, err)
	return doc
}
