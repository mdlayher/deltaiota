package auth

import (
	"database/sql"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/mdlayher/deltaiota/data"
	"github.com/mdlayher/deltaiota/ditest"
)

// Test_keyAuthenticate verifies that keyAuthenticate properly authenticates
// a session using the HTTP Basic Authorization header with a username and
// session key pair.
func Test_keyAuthenticate(t *testing.T) {
	// Establish temporary database for test
	err := ditest.WithTemporaryDB(func(db *data.DB) error {
		// Build context
		ac := NewContext(db)

		// Create and store mock user in temporary database
		user := ditest.MockUser()
		if err := ac.db.InsertUser(user); err != nil {
			return err
		}

		// Create and store mock user without sessions in temporary database
		user2 := ditest.MockUser()
		if err := ac.db.InsertUser(user2); err != nil {
			return err
		}

		// Generate a session for first user
		session, err := user.NewSession(time.Now().Add(1 * time.Minute))
		if err != nil {
			return err
		}

		// Store session in temporary database
		if err := ac.db.InsertSession(session); err != nil {
			return err
		}

		// Generate an expired session for first user
		expSession, err := user.NewSession(time.Now().Add(-1 * time.Minute))
		if err != nil {
			return err
		}

		// Store session in temporary database
		if err := ac.db.InsertSession(expSession); err != nil {
			return err
		}

		var tests = []struct {
			username string
			key      string
			err      error
		}{
			// Empty username and key
			{"", "", errNoUsername},
			// Username and empty key
			{user.Username, "", errNoKey},
			// Invalid username
			{"test2", session.Key, errInvalidUsername},
			// Invalid key
			{user.Username, "test2", errInvalidKey},
			// Key does not belong to user
			{user2.Username, session.Key, errInvalidKey},
			// Expired session
			{user.Username, expSession.Key, errExpiredKey},
			// Valid credentials
			{user.Username, session.Key, nil},
		}

		for _, test := range tests {
			// Create mock HTTP request
			req, err := http.NewRequest("POST", "/", nil)
			if err != nil {
				return err
			}

			// Set credentials for HTTP Basic
			req.SetBasicAuth(test.username, test.key)

			// Attempt authentication
			_, _, cErr, sErr := ac.keyAuthenticate(req)

			// Fail tests on any server error
			if sErr != nil {
				return sErr
			}

			// Check for expected client error
			if cErr != test.err {
				return fmt.Errorf("unexpected err: %v != %v", cErr, test.err)
			}
		}

		// Ensure expired session was deleted
		if _, err := ac.db.SelectSessionByKey(expSession.Key); err != sql.ErrNoRows {
			return fmt.Errorf("session expired, but still in database")
		}

		return nil
	})

	if err != nil {
		t.Fatal("ditest.WithTemporaryDB:", err)
	}
}
