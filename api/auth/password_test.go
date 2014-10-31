package auth

import (
	"net/http"
	"testing"

	"github.com/mdlayher/deltaiota/data"
	"github.com/mdlayher/deltaiota/data/models"
	"github.com/mdlayher/deltaiota/ditest"
)

// Test_passwordAuthenticateOK verifies that passwordAuthenticate works properly
// with a valid username and password pair.
func Test_passwordAuthenticateOK(t *testing.T) {
	test_passwordAuthenticate(t, nil, nil)
}

// Test_passwordAuthenticateNoUsername verifies that passwordAuthenticate returns
// a client error when no username is set.
func Test_passwordAuthenticateNoUsername(t *testing.T) {
	test_passwordAuthenticate(t, errNoUsername, func(t *testing.T, ac *Context, user *models.User) {
		// Empty username
		user.Username = ""
	})
}

// Test_passwordAuthenticateNoPassword verifies that passwordAuthenticate returns
// a client error when no password is set.
func Test_passwordAuthenticateNoPassword(t *testing.T) {
	test_passwordAuthenticate(t, errNoPassword, func(t *testing.T, ac *Context, user *models.User) {
		// Empty password
		user.Password = ""
	})
}

// Test_passwordAuthenticateInvalidUsername verifies that passwordAuthenticate returns
// a client error when an invalid username is set.
func Test_passwordAuthenticateInvalidUsername(t *testing.T) {
	test_passwordAuthenticate(t, errInvalidUsername, func(t *testing.T, ac *Context, user *models.User) {
		// Invalid username
		user.Username = ditest.RandomString(8)
	})
}

// Test_passwordAuthenticateInvalidPassword verifies that passwordAuthenticate returns
// a client error when an invalid password is set.
func Test_passwordAuthenticateInvalidPassword(t *testing.T) {
	test_passwordAuthenticate(t, errInvalidPassword, func(t *testing.T, ac *Context, user *models.User) {
		// Invalid password
		user.Password = ditest.RandomString(8)
	})
}

// test_passwordAuthenticate is a test handler which aids in testing the passwordAuthenticate
// handler.  It establishes text context, performs a setup function which can be used
// to manipulate test data, and finally expects a certain error to occur on authentication.
func test_passwordAuthenticate(t *testing.T, expErr error, fn func(t *testing.T, ac *Context, user *models.User)) {
	ditest.WithTemporaryDBNew(t, func(t *testing.T, db *data.DB) {
		// Build context
		ac := NewContext(db)

		// Create a mock user for authentication
		user := ditest.MockUser()
		plainPass := user.Password
		if err := user.SetPassword(plainPass); err != nil {
			t.Fatal(err)
		}
		if err := ac.db.InsertUser(user); err != nil {
			t.Fatal(err)
		}

		// Restore plaintext user password for tests, since the hash is already
		// in database
		user.Password = plainPass

		// Create mock HTTP request
		r, err := http.NewRequest("POST", "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		// If set, perform test setup closure, to manipulate data and test
		// for certain failure conditions
		if fn != nil {
			fn(t, ac, user)
		}

		// Set credentials for HTTP Basic
		r.SetBasicAuth(user.Username, user.Password)

		// Attempt authentication
		_, _, cErr, sErr := ac.passwordAuthenticate(r)

		// Fail tests on any server error
		if sErr != nil {
			t.Fatal(sErr)
		}

		// Check for expected client error
		if cErr != expErr {
			t.Fatalf("unexpected client err: %v != %v", cErr, expErr)
		}
	})
}
