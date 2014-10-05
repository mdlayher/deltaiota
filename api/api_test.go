package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mdlayher/deltaiota/api/v0"
	"github.com/mdlayher/deltaiota/data"
	"github.com/mdlayher/deltaiota/ditest"
)

// TestNewServeMux verifies that NewServeMux properly sets up the root of the API.
func TestNewServeMux(t *testing.T) {
	// Set up temporary database for test
	err := ditest.WithTemporaryDB(func(db *data.DB) {
		// Set up HTTP test server
		srv := httptest.NewServer(NewServeMux(db))
		defer srv.Close()

		// Set up tests to perform against server
		tests := []struct {
			method string
			path   string
			code   int
		}{
			// Root path
			{"GET", "/", http.StatusNotFound},
			// API v0 root
			{"GET", v0.APIPrefix, http.StatusNotFound},
			// API v1 root
			{"GET", "/api/v1", http.StatusNotFound},
		}

		// Iterate and perform requests for all tests
		for i, test := range tests {
			// Point test at HTTP test server
			test.path = srv.URL + test.path

			// Set up logging prefix
			logPrefix := fmt.Sprintf("[%02d] [%s %s]", i, test.method, test.path)

			// Generate HTTP request
			req, err := http.NewRequest(test.method, test.path, nil)
			if err != nil {
				t.Error(err)
			}

			// Receive HTTP response
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Error(err)
			}

			// Check for expected status code
			if res.StatusCode != test.code {
				t.Errorf("%s: unexpected code: %v != %v", logPrefix, res.StatusCode, test.code)
			}
		}
	})

	// Fail on errors from database setup/cleanup
	if err != nil {
		t.Fatal("ditest.WithTemporaryDB:", err)
	}
}
