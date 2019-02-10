package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vinothsparrow/scanner/config"
	"github.com/vinothsparrow/scanner/server"
)

func init() {
	config.Init("test")
}

var apiKey = "test"

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestHealth(t *testing.T) {
	// Build our expected body
	body := "Ok"
	// Grab our router
	router := server.NewRouter()
	// Perform a GET request with that handler.
	w := performRequest(router, "GET", "/health")
	// Assert we encoded correctly,
	// the request gives a 200
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, body, w.Body.String())
}

func TestDeny(t *testing.T) {
	// Grab our router
	router := server.NewRouter()
	// Perform a GET request with that handler.
	w := performRequest(router, "GET", "/v1/scan/git")
	// Assert we encoded correctly,
	// the request gives a 200
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestApiKey(t *testing.T) {
	// Grab our router
	router := server.NewRouter()
	// Perform a GET request with that handler.
	w := performRequest(router, "GET", "/v1/scan/git?api_key="+apiKey)
	// Assert we encoded correctly,
	// the request gives a 200
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
