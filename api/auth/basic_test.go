package auth

import (
	"net/http"
	"testing"

	"github.com/mdlayher/deltaiota/data"
	"github.com/mdlayher/deltaiota/ditest"
)

// Test_basicAuthenticate verifies that basicAuthenticate properly authenticates
// users using the HTTP Basic Authorization header.
func Test_basicAuthenticate(t *testing.T) {
	// Establish temporary database for test
	err := ditest.WithTemporaryDB(func(db *data.DB) {
		// Build context
		ac := NewContext(db)

		// Create a mock user for authentication
		user := ditest.MockUser()

		// Set a known password for tests
		password := "test"
		if err := user.SetPassword(password); err != nil {
			t.Error(err)
			return
		}

		// Store user in temporary database
		if err := ac.db.InsertUser(user); err != nil {
			t.Error(err)
			return
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
				t.Error(err)
				return
			}

			// Set credentials for HTTP Basic
			req.SetBasicAuth(test.username, test.password)

			// Attempt authentication
			_, _, cErr, sErr := ac.basicAuthenticate(req)

			// Fail tests on any server error
			if sErr != nil {
				t.Error(err)
				return
			}

			// Check for expected client error
			if cErr != test.err {
				t.Errorf("unexpected err: %v != %v", cErr, test.err)
				return
			}
		}
	})

	if err != nil {
		t.Fatal("ditest.WithTemporaryDB:", err)
	}
}

// Test_basicCredentials verifies that basicCredentials produces a correct
// username and password pair for input HTTP Basic Authorization header.
func Test_basicCredentials(t *testing.T) {
	var tests = []struct {
		input    string
		username string
		password string
		err      *Error
	}{
		// Empty input
		{"", "", "", errNoAuthorizationHeader},
		// No Authorization type
		{"abcdef012346789", "", "", errNoAuthorizationType},
		// Not HTTP Basic
		{"Hello abcdef012346789", "", "", errNotBasicAuthorization},
		// Invalid base64
		{"Basic !@#", "", "", errInvalidBase64Authorization},
		// Valid pair
		{"Basic dGVzdDp0ZXN0", "test", "test", nil},
	}

	for _, test := range tests {
		// Split input header into credentials
		username, password, err := basicCredentials(test.input)
		if err != nil {
			// Check for expected error
			if err == test.err {
				continue
			}
		}

		// Verify username
		if username != test.username {
			t.Fatalf("unexpected username: %v != %v", username, test.username)
		}

		// Verify password
		if password != test.password {
			t.Fatalf("unexpected password: %v != %v", password, test.password)
		}
	}
}
