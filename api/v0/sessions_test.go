package v0

import (
	"database/sql"
	"encoding/json"
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
	// Invoke tests with temporary database
	err := ditest.WithTemporaryDB(func(db *data.DB) {
		// Build context
		c := &Context{
			db: db,
		}

		// Generate and store mock user
		user := ditest.MockUser()
		if err := c.db.InsertUser(user); err != nil {
			t.Error(err)
			return
		}

		// Generate and store mock session
		session := &models.Session{
			UserID: user.ID,
			Key:    ditest.RandomString(32),
			Expire: uint64(time.Now().Unix()),
		}
		if err := c.db.InsertSession(session); err != nil {
			t.Error(err)
			return
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
				t.Error(err)
				return
			}

			// Store mock-authenticated session
			auth.SetSession(r, session)

			// Delegate to appropriate handler
			code, _, err := c.SessionsAPI(r, util.Vars{})
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

// TestPostSession verifies that PostSession returns the appropriate HTTP status
// code, body, and any errors which occur.
func TestPostSession(t *testing.T) {
	// Invoke tests with temporary database
	err := ditest.WithTemporaryDB(func(db *data.DB) {
		// Build context
		c := &Context{
			db: db,
		}

		// Generate and store mock user
		user := ditest.MockUser()
		if err := c.db.InsertUser(user); err != nil {
			t.Error(err)
			return
		}

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

// TestDeleteSession verifies that DeleteSession returns the appropriate HTTP status
// code, body, and any errors which occur.
func TestDeleteSession(t *testing.T) {
	// Invoke tests with temporary database
	err := ditest.WithTemporaryDB(func(db *data.DB) {
		// Build context
		c := &Context{
			db: db,
		}

		// Generate and store mock user
		user := ditest.MockUser()
		if err := c.db.InsertUser(user); err != nil {
			t.Error(err)
			return
		}

		// Generate and store mock session
		session := &models.Session{
			UserID: user.ID,
			Key:    ditest.RandomString(32),
			Expire: uint64(time.Now().Unix()),
		}
		if err := c.db.InsertSession(session); err != nil {
			t.Error(err)
			return
		}

		// Generate HTTP request
		r, err := http.NewRequest("DELETE", "/", nil)
		if err != nil {
			t.Error(err)
			return
		}

		// Store mock-authenticated session
		auth.SetSession(r, session)

		// Invoke DeleteSession with HTTP request
		code, _, err := c.DeleteSession(r, util.Vars{})
		if err != nil {
			t.Error(err)
			return
		}

		// Ensure proper HTTP status code
		if code != http.StatusNoContent {
			t.Error("expected HTTP OK, got code:", code)
			return
		}

		// Ensure session was deleted
		if _, err := c.db.SelectSessionByKey(session.Key); err != sql.ErrNoRows {
			t.Errorf("called DeleteSession, but session still exists: %v", session)
			return
		}
	})

	// Check for errors from database setup/cleanup
	if err != nil {
		t.Fatal("ditest.WithTemporaryDB:", err)
	}
}
