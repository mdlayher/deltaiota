package v0

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mdlayher/deltaiota/data"
	"github.com/mdlayher/deltaiota/ditest"
)

// TestNewServeMux verifies that NewServeMux properly sets up API v0.
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
			body   []byte
		}{
			// Root path
			{"GET", "/", http.StatusNotFound, nil},

			// Users API - tests MUST be run in order

			// All users
			{"GET", "/users", http.StatusOK, nil},

			// One user, bad user ID
			{"GET", "/users/-1", http.StatusBadRequest, nil},
			{"GET", "/users/matt", http.StatusBadRequest, nil},
			// One user, which does not exist yet
			{"GET", "/users/1", http.StatusNotFound, nil},

			// Create user, malformed JSON
			{"POST", "/users", http.StatusBadRequest, []byte(`{`)},
			// Create user, no password field
			{"POST", "/users", http.StatusBadRequest, []byte(`{}`)},
			// Create user, empty password field
			{"POST", "/users", http.StatusBadRequest, []byte(`{"password":""}`)},
			// Create user, missing username
			{"POST", "/users", http.StatusBadRequest, []byte(`{"password":"test","firstName":"test","lastName":"test"}`)},
			// Create user, missing first name
			{"POST", "/users", http.StatusBadRequest, []byte(`{"password":"test","lastName":"test","username":"test"}`)},
			// Create user, missing last name
			{"POST", "/users", http.StatusBadRequest, []byte(`{"password":"test","firstName":"test","username":"test"}`)},
			// Create user, valid request
			{"POST", "/users", http.StatusCreated, []byte(`{"password":"test","firstName":"test","lastName":"test","username":"test"}`)},
			// Create user, duplicate username
			{"POST", "/users", http.StatusConflict, []byte(`{"password":"test","firstName":"test2","lastName":"test2","username":"test"}`)},
			// Create user, duplicate first name
			{"POST", "/users", http.StatusConflict, []byte(`{"password":"test","firstName":"test","lastName":"test2","username":"test2"}`)},
			// Create user, duplicate last name
			{"POST", "/users", http.StatusConflict, []byte(`{"password":"test","firstName":"test2","lastName":"test","username":"test2"}`)},

			// Get user, which now exists
			{"GET", "/users/1", http.StatusOK, nil},
		}

		// Iterate and perform requests for all tests
		for i, test := range tests {
			// Point test at HTTP test server
			test.path = srv.URL + APIPrefix + test.path

			// Set up logging prefix
			logPrefix := fmt.Sprintf("[%02d] [%s %s]", i, test.method, test.path)

			// Generate HTTP request
			req, err := http.NewRequest(test.method, test.path, bytes.NewReader(test.body))
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
