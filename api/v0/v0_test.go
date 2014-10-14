package v0

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mdlayher/deltaiota/data"
	"github.com/mdlayher/deltaiota/data/models"
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

// withContext sets up a test context with an API context wrapping a
// temporary database.
func withContext(t *testing.T, fn func(c *Context) error) {
	// Invoke tests with temporary database
	err := ditest.WithTemporaryDB(func(db *data.DB) {
		// Build context
		c := &Context{
			db: db,
		}

		// Invoke test
		if err := fn(c); err != nil {
			t.Error(err)
			return
		}
	})

	// Check for errors from database setup/cleanup
	if err != nil {
		t.Fatal("ditest.WithTemporaryDB:", err)
	}
}

// withContextUser builds upon withContext, adding a mock user.
func withContextUser(t *testing.T, fn func(c *Context, user *models.User) error) {
	withContext(t, func(c *Context) error {
		// Generate mock user
		user := ditest.MockUser()
		if err := c.db.InsertUser(user); err != nil {
			return err
		}

		// Invoke test
		return fn(c, user)
	})
}
