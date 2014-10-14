package v0

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"testing"

	"github.com/mdlayher/deltaiota/api/util"
)

// TestStatusAPI verifies that StatusAPI correctly routes requests to
// other Status API handlers, using the input HTTP request.
func TestStatusAPI(t *testing.T) {
	withContext(t, func(c *Context) error {
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
				return err
			}

			// Delegate to appropriate handler
			code, _, err := c.StatusAPI(r, util.Vars{})
			if err != nil {
				return err
			}

			// Ensure proper HTTP status code
			if code != test.code {
				return fmt.Errorf("unexpected code: %v != %v", code, test.code)
			}
		}

		return nil
	})
}

// TestGetStatus verifies that GetStatus returns accurate server status values.
func TestGetStatus(t *testing.T) {
	withContext(t, func(c *Context) error {
		// Fetch server status
		code, body, err := c.GetStatus(nil, util.Vars{})
		if err != nil {
			return err
		}

		// Ensure proper HTTP status code
		if code != http.StatusOK {
			return fmt.Errorf("unexpected code: %v != %v", code, http.StatusOK)
		}

		// Unmarshal response body
		var res StatusResponse
		if err := json.Unmarshal(body, &res); err != nil {
			return err
		}

		// Verify status returned
		if res.Status == nil {
			return fmt.Errorf("empty Status object in response")
		}

		// Verify that some static values are the same
		if res.Status.Architecture != runtime.GOARCH {
			return fmt.Errorf("unexpected Architecture: %v != %v", res.Status.Architecture, runtime.GOARCH)
		}
		if res.Status.NumCPU != runtime.NumCPU() {
			return fmt.Errorf("unexpected NumCPU: %v != %v", res.Status.NumCPU, runtime.NumCPU())
		}
		if res.Status.Platform != runtime.GOOS {
			return fmt.Errorf("unexpected Platform: %v != %v", res.Status.Platform, runtime.GOOS)
		}

		return nil
	})
}
