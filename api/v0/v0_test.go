package v0

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mdlayher/deltaiota/data"
	"github.com/mdlayher/deltaiota/data/models"
	"github.com/mdlayher/deltaiota/ditest"
)

// TestNewServeMuxGETRootNotFound verifies that nothing exists at the root
// of the HTTP API server.
func TestNewServeMuxGETRootNotFound(t *testing.T) {
	testNewServeMux(t, "GET", "/", http.StatusNotFound)
}

// TestNewServeMuxGETHEADNotificationsOK verifies that HTTP GET and HEAD
// methods return HTTP 200 on the Notifications API.
func TestNewServeMuxGETHEADNotificationsOK(t *testing.T) {
	for _, m := range []string{"GET", "HEAD"} {
		testNewServeMux(t, m, "/notifications", http.StatusOK)
	}
}

// TestNewServeMuxNotificationsMethodNotAllowed verifies that disallowed HTTP
// methods return HTTP 405 on the Notifications API.
func TestNewServeMuxNotificationsMethodNotAllowed(t *testing.T) {
	for _, m := range []string{"CAT", "DELETE", "PATCH", "POST", "PUT"} {
		testNewServeMux(t, m, "/notifications", http.StatusMethodNotAllowed)
	}
}

// TestNewServeMuxGETSessionsOK verifies that HTTP GET
// method returns HTTP 200 on the Sessions API.
func TestNewServeMuxGETSessionsOK(t *testing.T) {
	testNewServeMux(t, "GET", "/sessions", http.StatusOK)
}

// TestNewServeMuxPOSTSessionsOK verifies that HTTP POST
// method returns HTTP 200 on the Sessions API.
func TestNewServeMuxPOSTSessionsOK(t *testing.T) {
	testNewServeMux(t, "POST", "/sessions", http.StatusOK)
}

// TestNewServeMuxDELETESessionsNoContent verifies that HTTP DELETE
// method returns HTTP 204 on the Sessions API.
func TestNewServeMuxDELETESessionsNoContent(t *testing.T) {
	testNewServeMux(t, "DELETE", "/sessions", http.StatusNoContent)
}

// TestNewServeMuxSessionsMethodNotAllowed verifies that disallowed HTTP
// methods return HTTP 405 on the Sessions API.
func TestNewServeMuxSessionsMethodNotAllowed(t *testing.T) {
	for _, m := range []string{"PATCH", "PUT"} {
		testNewServeMux(t, m, "/sessions", http.StatusMethodNotAllowed)
	}
}

// TestNewServeMuxGETHEADStatusOK verifies that HTTP GET and HEAD
// methods return HTTP 200 on the Status API.
func TestNewServeMuxGETHEADStatusOK(t *testing.T) {
	for _, m := range []string{"GET", "HEAD"} {
		testNewServeMux(t, m, "/status", http.StatusOK)
	}
}

// TestNewServeMuxStatusMethodNotAllowed verifies that disallowed HTTP
// methods return HTTP 405 on the Status API.
func TestNewServeMuxStatusMethodNotAllowed(t *testing.T) {
	for _, m := range []string{"CAT", "DELETE", "PATCH", "POST", "PUT"} {
		testNewServeMux(t, m, "/status", http.StatusMethodNotAllowed)
	}
}

// TestNewServeMuxGETHEADUsersNoIDOK verifies that HTTP GET and HEAD
// methods return HTTP 200 on the Users API, with no ID.
func TestNewServeMuxGETHEADUsersNoIDOK(t *testing.T) {
	for _, m := range []string{"GET", "HEAD"} {
		testNewServeMux(t, m, "/users", http.StatusOK)
	}
}

// TestNewServeMuxGETHEADUsersWithIDOK verifies that HTTP GET and HEAD
// methods return HTTP 200 on the Users API, with ID.
func TestNewServeMuxGETHEADUsersWithIDOK(t *testing.T) {
	for _, m := range []string{"GET", "HEAD"} {
		testNewServeMux(t, m, "/users/1", http.StatusOK)
	}
}

// TestNewServeMuxPOSTUsersBadRequest verifies that the HTTP POST
// method returns HTTP 400 on the Users API with no request body.
func TestNewServeMuxPOSTUsersBadRequest(t *testing.T) {
	testNewServeMux(t, "POST", "/users", http.StatusBadRequest)
}

// TestNewServeMuxPUTUsersNoIDBadRequest verifies that the HTTP PUT
// method returns HTTP 400 on the Users API with no request body
// and no ID.
func TestNewServeMuxPUTUsersNoIDBadRequest(t *testing.T) {
	testNewServeMux(t, "PUT", "/users", http.StatusBadRequest)
}

// TestNewServeMuxPUTUsersWithIDBadRequest verifies that the HTTP PUT
// method returns HTTP 400 on the Users API with no request body
// and with an ID.
func TestNewServeMuxPUTUsersWithIDBadRequest(t *testing.T) {
	testNewServeMux(t, "PUT", "/users/1", http.StatusBadRequest)
}

// TestNewServeMuxDELETEUsersNoIDBadRequest verifies that the HTTP DELETE
// method returns HTTP 400 on the Users API with no ID.
func TestNewServeMuxDELETEUsersNoIDBadRequest(t *testing.T) {
	testNewServeMux(t, "DELETE", "/users", http.StatusBadRequest)
}

// TestNewServeMuxDELETEUsersWithIDNoContent verifies that the HTTP DELETE
// method returns HTTP 204 on the Users API with an ID.
func TestNewServeMuxDELETEUsersWithIDNoContent(t *testing.T) {
	testNewServeMux(t, "DELETE", "/users/1", http.StatusNoContent)
}

// testNewServeMux is a helper which verifies that an HTTP request with the
// given path returns the expected HTTP status code.
func testNewServeMux(t *testing.T, method string, path string, code int) {
	ditest.WithTemporaryDBNew(t, func(t *testing.T, db *data.DB) {
		// Set up HTTP test server
		srv := httptest.NewServer(NewServeMux(db))
		defer srv.Close()

		// Set up temporary user for authentication
		user := ditest.MockUser()
		if err := db.InsertUser(user); err != nil {
			t.Fatal(err)
		}

		// Special case, if POSTing to sessions, password authentication
		// is used
		usePassword := method == "POST" && path == "/sessions"
		var plainPass string
		if usePassword {
			// Hash password and update
			plainPass = user.Password
			if err := user.SetPassword(user.Password); err != nil {
				t.Fatal(err)
			}
			if err := db.UpdateUser(user); err != nil {
				t.Fatal(err)
			}
		}

		// Set up temporary session for authentication
		session, err := user.NewSession(time.Now().Add(1 * time.Minute))
		if err != nil {
			t.Fatal(err)
		}
		if err := db.InsertSession(session); err != nil {
			t.Fatal(err)
		}

		// Generate HTTP request, point at test server with
		// API v0 namespace
		path = srv.URL + APIPrefix + path
		req, err := http.NewRequest(method, path, nil)
		if err != nil {
			t.Fatal(err)
		}

		// Set up password or key authentication
		if usePassword {
			req.SetBasicAuth(user.Username, plainPass)
		} else {
			req.SetBasicAuth(user.Username, session.Key)
		}

		// Receive HTTP response
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		// Check for expected status code
		if res.StatusCode != code {
			t.Errorf("HTTP %s %s: unexpected code: %v != %v", method, path, res.StatusCode, code)
		}

		// If HEAD request, verify body is empty
		if req.Method == "HEAD" {
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}

			if body != nil && len(body) > 0 {
				t.Fatalf("HTTP %s %s: unexpected body: %v", method, path, string(body))
			}
		}
	})
}

// withContext sets up a test context with an API context wrapping a
// temporary database.
func withContext(t *testing.T, fn func(c *Context) error) {
	// Invoke tests with temporary database
	err := ditest.WithTemporaryDB(func(db *data.DB) error {
		// Build context
		c := &Context{
			db: db,
		}

		// Invoke test
		return fn(c)
	})

	// Check for errors from database setup/cleanup
	if err != nil {
		t.Fatal("ditest.WithTemporaryDB:", err)
	}
}

// withContextUser builds upon withContext, adding a mock user.
func withContextUser(t *testing.T, fn func(c *Context, user *models.User) error) {
	withContext(t, func(c *Context) error {
		// Generate mock user
		user := ditest.MockUser()
		if err := c.db.InsertUser(user); err != nil {
			return err
		}

		// Invoke test
		return fn(c, user)
	})
}
