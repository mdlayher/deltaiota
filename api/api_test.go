package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mdlayher/deltaiota/api/v0"
	"github.com/mdlayher/deltaiota/data"
	"github.com/mdlayher/deltaiota/ditest"
)

// TestNewServeMuxGETRootNotFound verifies that nothing exists at the root
// of the HTTP API server.
func TestNewServeMuxGETRootNotFound(t *testing.T) {
	testNewServeMux(t, "GET", "/", http.StatusNotFound)
}

// TestNewServeMuxGETAPIv0RootNotFound verifies that nothing exists at the root
// of the v0 HTTP API server.
func TestNewServeMuxGETAPIv0RootNotFound(t *testing.T) {
	testNewServeMux(t, "GET", v0.APIPrefix, http.StatusNotFound)
}

// TestNewServeMuxGETAPIV1RootNotFound verifies that nothing exists at the root
// of the v1 HTTP API server.
func TestNewServeMuxGETAPIv1RootNotFound(t *testing.T) {
	testNewServeMux(t, "GET", "/api/v1", http.StatusNotFound)
}

// testNewServeMux is a helper which verifies that an HTTP request with the
// given path returns the expected HTTP status code.
func testNewServeMux(t *testing.T, method string, path string, code int) {
	// Set up temporary database for test
	ditest.WithTemporaryDBNew(t, func(t *testing.T, db *data.DB) {
		// Set up HTTP test server
		srv := httptest.NewServer(NewServeMux(db))
		defer srv.Close()

		// Generate HTTP request, point at test server
		path = srv.URL + path
		req, err := http.NewRequest(method, path, nil)
		if err != nil {
			t.Fatal(err)
		}

		// Receive HTTP response
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		// Check for expected status code
		if res.StatusCode != code {
			t.Errorf("unexpected code: %v != %v", res.StatusCode, code)
		}
	})
}
