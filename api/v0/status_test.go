package v0

import (
	"encoding/json"
	"net/http"
	"runtime"
	"testing"

	"github.com/mdlayher/deltaiota/api/util"
	"github.com/mdlayher/deltaiota/data"
	"github.com/mdlayher/deltaiota/ditest"
)

// TestStatusAPI verifies that StatusAPI correctly routes requests to
// other Status API handlers, using the input HTTP request.
func TestStatusAPI(t *testing.T) {
	// Invoke tests with temporary database
	err := ditest.WithTemporaryDB(func(db *data.DB) {
		// Build context
		c := &Context{
			db: db,
		}

		var tests = []struct {
			method string
			code   int
		}{
			// GetStatus
			{"GET", http.StatusOK},
			{"HEAD", http.StatusOK},
			// Method not allowed
			{"POST", http.StatusMethodNotAllowed},
			{"PUT", http.StatusMethodNotAllowed},
			{"PATCH", http.StatusMethodNotAllowed},
			{"DELETE", http.StatusMethodNotAllowed},
			{"CAT", http.StatusMethodNotAllowed},
		}

		for _, test := range tests {
			// Generate HTTP request
			r, err := http.NewRequest(test.method, "/", nil)
			if err != nil {
				t.Error(err)
				return
			}

			// Delegate to appropriate handler
			code, _, err := c.StatusAPI(r, util.Vars{})
			if err != nil {
				t.Error(err)
				return
			}

			// Ensure proper HTTP status code
			if code != test.code {
				t.Errorf("unexpected code: %v != %v", code, test.code)
				return
			}
		}
	})

	// Check for errors from database setup/cleanup
	if err != nil {
		t.Fatal("ditest.WithTemporaryDB:", err)
	}
}

// TestGetStatus verifies that GetStatus returns accurate server status values.
func TestGetStatus(t *testing.T) {
	// Invoke tests with temporary database
	err := ditest.WithTemporaryDB(func(db *data.DB) {
		// Build context
		c := &Context{
			db: db,
		}

		// Fetch server status
		code, body, err := c.GetStatus(nil, util.Vars{})
		if err != nil {
			t.Error(err)
			return
		}

		// Ensure proper HTTP status code
		if code != http.StatusOK {
			t.Errorf("unexpected code: %v != %v", code, http.StatusOK)
			return
		}

		// Unmarshal response body
		var res StatusResponse
		if err := json.Unmarshal(body, &res); err != nil {
			t.Error(err)
			return
		}

		// Verify status returned
		if res.Status == nil {
			t.Error("empty Status object in response")
			return
		}

		// Verify that some static values are the same
		if res.Status.Architecture != runtime.GOARCH {
			t.Errorf("unexpected Architecture: %v != %v", res.Status.Architecture, runtime.GOARCH)
			return
		}
		if res.Status.NumCPU != runtime.NumCPU() {
			t.Errorf("unexpected NumCPU: %v != %v", res.Status.NumCPU, runtime.NumCPU())
			return
		}
		if res.Status.Platform != runtime.GOOS {
			t.Errorf("unexpected Platform: %v != %v", res.Status.Platform, runtime.GOOS)
			return
		}
	})

	// Check for errors from database setup/cleanup
	if err != nil {
		t.Fatal("ditest.WithTemporaryDB:", err)
	}
}
