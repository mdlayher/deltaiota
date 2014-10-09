package v0

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/mdlayher/deltaiota/api/auth"
	"github.com/mdlayher/deltaiota/api/util"
	"github.com/mdlayher/deltaiota/data"
	"github.com/mdlayher/deltaiota/ditest"
)

// TestPostSession verifies that PostSession returns the appropriate HTTP status
// code, body, and any errors which occur.
func TestPostSession(t *testing.T) {
	// Invoke tests with temporary database
	err := ditest.WithTemporaryDB(func(db *data.DB) {
		// Build context
		c := &Context{
			db: db,
		}

		// Generate mock user
		user := ditest.MockUser()
		user.ID = 1000

		// Generate HTTP request
		r, err := http.NewRequest("POST", "/", nil)
		if err != nil {
			t.Error(err)
			return
		}

		// Store mock-authenticated user
		auth.SetUser(r, user)

		// Invoke PostSession with HTTP request
		code, body, err := c.PostSession(r, util.Vars{})
		if err != nil {
			t.Error(err)
			return
		}

		// Ensure proper HTTP status code
		if code != http.StatusOK {
			t.Error("expected HTTP OK, got code:", code)
			return
		}

		// Unmarshal response body
		var res SessionsResponse
		if err := json.Unmarshal(body, &res); err != nil {
			t.Error(err)
			return
		}

		// Verify session belongs to this user
		if res.Session.UserID != user.ID {
			t.Errorf("unexpected user ID: %v != %v", res.Session.UserID, user.ID)
			return
		}
	})

	// Check for errors from database setup/cleanup
	if err != nil {
		t.Fatal("ditest.WithTemporaryDB:", err)
	}
}
