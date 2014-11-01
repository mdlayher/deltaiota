package auth

import (
	"database/sql"
	"net/http"
	"testing"
	"time"

	"github.com/mdlayher/deltaiota/data"
	"github.com/mdlayher/deltaiota/data/models"
	"github.com/mdlayher/deltaiota/ditest"
)

// Test_keyAuthenticateOK verifies that keyAuthenticate works properly with a
// valid username and key pair.
func Test_keyAuthenticateOK(t *testing.T) {
	test_keyAuthenticate(t, nil, nil)
}

// Test_keyAuthenticateNoUsername verifies that keyAuthenticate returns a client
// error when no username is set.
func Test_keyAuthenticateNoUsername(t *testing.T) {
	test_keyAuthenticate(t, errNoUsername, func(t *testing.T, ac *Context, user *models.User, session *models.Session) {
		// Empty username
		user.Username = ""
	})
}

// Test_keyAuthenticateNoKey verifies that keyAuthenticate returns a client
// error when no key is set.
func Test_keyAuthenticateNoKey(t *testing.T) {
	test_keyAuthenticate(t, errNoKey, func(t *testing.T, ac *Context, user *models.User, session *models.Session) {
		// Empty key
		session.Key = ""
	})
}

// Test_keyAuthenticateInvalidUsername verifies that keyAuthenticate returns a client
// error when an invalid username is set.
func Test_keyAuthenticateInvalidUsername(t *testing.T) {
	test_keyAuthenticate(t, errInvalidUsername, func(t *testing.T, ac *Context, user *models.User, session *models.Session) {
		// Invalid username
		user.Username = ditest.RandomString(8)
	})
}

// Test_keyAuthenticateInvalidKey verifies that keyAuthenticate returns a client
// error when an invalid key is set.
func Test_keyAuthenticateInvalidKey(t *testing.T) {
	test_keyAuthenticate(t, errInvalidKey, func(t *testing.T, ac *Context, user *models.User, session *models.Session) {
		// Invalid key
		session.Key = ditest.RandomString(8)
	})
}

// Test_keyAuthenticateWrongKeyForUser verifies that keyAuthenticate returns a client
// error when a valid user attempts to use another user's key.
func Test_keyAuthenticateWrongKeyForUser(t *testing.T) {
	test_keyAuthenticate(t, errInvalidKey, func(t *testing.T, ac *Context, user *models.User, session *models.Session) {
		// Generate another mock user
		user2 := ditest.MockUser()
		if err := ac.db.InsertUser(user2); err != nil {
			t.Fatal(err)
		}

		// Set original user's username to new user's username,
		// so that the key is valid, but it does not belong to
		// the new user
		user.Username = user2.Username
	})
}

// Test_keyAuthenticateExpiredSession verifies that keyAuthenticate returns a client
// error when a valid user attempts to use an expired session.
func Test_keyAuthenticateExpiredSession(t *testing.T) {
	test_keyAuthenticate(t, errExpiredKey, func(t *testing.T, ac *Context, user *models.User, session *models.Session) {
		// Expire session immediately
		session.SetExpire(time.Now().Add(-1 * time.Hour))
		if err := ac.db.UpdateSession(session); err != nil {
			t.Fatal(err)
		}
	})
}

// test_keyAuthenticate is a test helper which aids in testing the keyAuthenticate
// handler.  It establishes test context, performs a setup function which can be used
// to manipulate test data, and finally expects a certain error to occur on authentication.
func test_keyAuthenticate(t *testing.T, expErr error, fn func(t *testing.T, ac *Context, user *models.User, session *models.Session)) {
	ditest.WithTemporaryDBNew(t, func(t *testing.T, db *data.DB) {
		// Build context
		ac := NewContext(db)

		// Create and store mock user in temporary database
		user := ditest.MockUser()
		if err := ac.db.InsertUser(user); err != nil {
			t.Fatal(err)
		}

		// Generate a session for first user
		session, err := user.NewSession(time.Now().Add(1 * time.Minute))
		if err != nil {
			t.Fatal(err)
		}

		// Store session in temporary database
		if err := ac.db.InsertSession(session); err != nil {
			t.Fatal(err)
		}

		// Create mock HTTP request
		r, err := http.NewRequest("POST", "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		// If set, perform test setup closure, to manipulate data and test
		// for certain failure conditions
		if fn != nil {
			fn(t, ac, user, session)
		}

		// Set credentials for HTTP Basic
		r.SetBasicAuth(user.Username, session.Key)

		// Attempt authentication
		_, _, cErr, sErr := ac.keyAuthenticate(r)

		// Fail tests on any server error
		if sErr != nil {
			t.Fatal(sErr)
		}

		// Check for expected client error
		if cErr != expErr {
			t.Fatalf("unexpected client err: %v != %v", cErr, expErr)
		}

		// Ensure any expired sessions were deleted
		if cErr == errExpiredKey {
			if _, err := ac.db.SelectSessionByKey(session.Key); err != sql.ErrNoRows {
				t.Fatalf("session expired, but still in database")
			}
		}
	})
}
