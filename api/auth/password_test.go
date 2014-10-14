package auth

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/mdlayher/deltaiota/data"
	"github.com/mdlayher/deltaiota/ditest"
)

// Test_passwordAuthenticate verifies that passwordAuthenticate properly authenticates
// users using the HTTP Basic Authorization header with a username and password pair.
func Test_passwordAuthenticate(t *testing.T) {
	// Establish temporary database for test
	err := ditest.WithTemporaryDB(func(db *data.DB) error {
		// Build context
		ac := NewContext(db)

		// Create a mock user for authentication
		user := ditest.MockUser()

		// Set a known password for tests
		password := "test"
		if err := user.SetPassword(password); err != nil {
			return err
		}

		// Store user in temporary database
		if err := ac.db.InsertUser(user); err != nil {
			return err
		}

		var tests = []struct {
			username string
			password string
			err      error
		}{
			// Empty username and password
			{"", "", errNoUsername},
			// Username and empty password
			{user.Username, "", errNoPassword},
			// Invalid username
			{"test2", password, errInvalidUsername},
			// Invalid password
			{user.Username, "test2", errInvalidPassword},
			// Valid credentials
			{user.Username, password, nil},
		}

		for _, test := range tests {
			// Create mock HTTP request
			req, err := http.NewRequest("POST", "/", nil)
			if err != nil {
				return err
			}

			// Set credentials for HTTP Basic
			req.SetBasicAuth(test.username, test.password)

			// Attempt authentication
			_, _, cErr, sErr := ac.passwordAuthenticate(req)

			// Fail tests on any server error
			if sErr != nil {
				return sErr
			}

			// Check for expected client error
			if cErr != test.err {
				return fmt.Errorf("unexpected err: %v != %v", cErr, test.err)
			}
		}

		return nil
	})

	if err != nil {
		t.Fatal("ditest.WithTemporaryDB:", err)
	}
}
