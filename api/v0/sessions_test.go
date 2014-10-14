package v0

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/mdlayher/deltaiota/api/auth"
	"github.com/mdlayher/deltaiota/api/util"
	"github.com/mdlayher/deltaiota/data"
	"github.com/mdlayher/deltaiota/data/models"
	"github.com/mdlayher/deltaiota/ditest"
)

// TestSessionsAPI verifies that SessionsAPI correctly routes requests to
// other Sessions API handlers, using the input HTTP request.
func TestSessionsAPI(t *testing.T) {
	withDBContextUser(t, func(db *data.DB, c *Context, user *models.User) error {
		// Generate and store mock session
		session := &models.Session{
			UserID: user.ID,
			Key:    ditest.RandomString(32),
			Expire: uint64(time.Now().Unix()),
		}
		if err := c.db.InsertSession(session); err != nil {
			return err
		}

		var tests = []struct {
			method string
			code   int
		}{
			// DeleteSession
			{"DELETE", http.StatusNoContent},
			// Method not allowed
			{"GET", http.StatusMethodNotAllowed},
			{"HEAD", http.StatusMethodNotAllowed},
			{"PUT", http.StatusMethodNotAllowed},
			{"PATCH", http.StatusMethodNotAllowed},
			{"CAT", http.StatusMethodNotAllowed},
			// Method not allowed from the UsersAPI call as a safety mechanism,
			// even though it is valid with password authentication.
			{"POST", http.StatusMethodNotAllowed},
		}

		for _, test := range tests {
			// Generate HTTP request
			r, err := http.NewRequest(test.method, "/", nil)
			if err != nil {
				return err
			}

			// Store mock-authenticated session
			auth.SetSession(r, session)

			// Delegate to appropriate handler
			code, _, err := c.SessionsAPI(r, util.Vars{})
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

// TestPostSession verifies that PostSession returns the appropriate HTTP status
// code, body, and any errors which occur.
func TestPostSession(t *testing.T) {
	withDBContextUser(t, func(db *data.DB, c *Context, user *models.User) error {
		// Generate HTTP request
		r, err := http.NewRequest("POST", "/", nil)
		if err != nil {
			return err
		}

		// Store mock-authenticated user
		auth.SetUser(r, user)

		// Invoke PostSession with HTTP request
		code, body, err := c.PostSession(r, util.Vars{})
		if err != nil {
			return err
		}

		// Ensure proper HTTP status code
		if code != http.StatusOK {
			return fmt.Errorf("expected HTTP OK, got code: %v", code)
		}

		// Unmarshal response body
		var res SessionsResponse
		if err := json.Unmarshal(body, &res); err != nil {
			return err
		}

		// Verify session belongs to this user
		if res.Session.UserID != user.ID {
			return fmt.Errorf("unexpected user ID: %v != %v", res.Session.UserID, user.ID)
		}

		return nil
	})
}

// TestDeleteSession verifies that DeleteSession returns the appropriate HTTP status
// code, body, and any errors which occur.
func TestDeleteSession(t *testing.T) {
	withDBContextUser(t, func(db *data.DB, c *Context, user *models.User) error {
		// Generate and store mock session
		session := &models.Session{
			UserID: user.ID,
			Key:    ditest.RandomString(32),
			Expire: uint64(time.Now().Unix()),
		}
		if err := c.db.InsertSession(session); err != nil {
			return err
		}

		// Generate HTTP request
		r, err := http.NewRequest("DELETE", "/", nil)
		if err != nil {
			return err
		}

		// Store mock-authenticated session
		auth.SetSession(r, session)

		// Invoke DeleteSession with HTTP request
		code, _, err := c.DeleteSession(r, util.Vars{})
		if err != nil {
			return err
		}

		// Ensure proper HTTP status code
		if code != http.StatusNoContent {
			return fmt.Errorf("expected HTTP OK, got code: %v", code)
		}

		// Ensure session was deleted
		if _, err := c.db.SelectSessionByKey(session.Key); err != sql.ErrNoRows {
			return fmt.Errorf("called DeleteSession, but session still exists: %v", session)
		}

		return nil
	})
}
