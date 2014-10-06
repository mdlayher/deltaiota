package v0

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/mdlayher/deltaiota/data"
	"github.com/mdlayher/deltaiota/data/models"
	"github.com/mdlayher/deltaiota/ditest"
)

// TestListUsers verifies that ListUsers returns the appropriate HTTP status
// code, body, and any errors which occur.
func TestListUsers(t *testing.T) {
	// Invoke tests with temporary database
	err := ditest.WithTemporaryDB(func(db *data.DB) {
		// Build context
		c := &context{
			db: db,
		}

		// Loop twice to check for empty users, add a user, then verify
		// that the same user is returned
		var lastUser *models.User
		for i := 0; i < 1; i++ {
			// Fetch list of current users
			code, body, err := c.ListUsers(nil)
			if err != nil {
				t.Fatal(err)
			}

			// Ensure proper HTTP status code
			if code != http.StatusOK {
				t.Fatalf("unexpected code: %v != %v", code, http.StatusOK)
			}

			// Unmarshal response body
			var res UsersResponse
			if err := json.Unmarshal(body, &res); err != nil {
				t.Fatal(err)
			}

			// Verify length of users slice
			if len(res.Users) != i {
				t.Fatalf("unexpected number of users returned: %v", res.Users)
			}

			// Check if generated user returned and break loop on second run
			if i == 1 {
				if res.Users[0] != lastUser {
					t.Fatalf("unexpected User: %v != %v", res.Users[0], lastUser)
				}

				break
			}

			// Generate and save a mock user in the database
			lastUser = ditest.MockUser()
			if err := c.db.InsertUser(lastUser); err != nil {
				t.Fatal(err)
			}
		}
	})

	// Check for errors from database setup/cleanup
	if err != nil {
		t.Fatalf("ditest.WithTemporaryDB:", err)
	}
}