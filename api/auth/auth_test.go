package auth

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mdlayher/deltaiota/api/util"
	"github.com/mdlayher/deltaiota/data/models"
	"github.com/mdlayher/deltaiota/ditest"
)

// Test_makeAuthHandlerClientError verifies that makeAuthHandler generates
// appropriate output for a client error.
func Test_makeAuthHandlerClientError(t *testing.T) {
	// Test function which returns a formatted client error
	reason := "foo bar"
	clientErrFn := func(r *http.Request) (*models.User, *models.Session, error, error) {
		return nil, nil, &Error{
			Reason: reason,
		}, nil
	}

	// Build client failure JSON
	authFailJSON := []byte(`{"error":{"code":401,"message":"authentication failed: ` + reason + `"}}`)

	test_makeAuthHandler(t, clientErrFn, okHandler(), http.StatusUnauthorized, authFailJSON, nil)
}

// Test_makeAuthHandlerClientNonStandardError verifies that makeAuthHandler
// generates a generic error when the wrapped Error type is not used.
func Test_makeAuthHandlerClientNonStandardError(t *testing.T) {
	// Test function which returns a non-standard client error
	clientBadErrFn := func(r *http.Request) (*models.User, *models.Session, error, error) {
		return nil, nil, errors.New("some error"), nil
	}

	test_makeAuthHandler(t, clientBadErrFn, okHandler(), http.StatusUnauthorized, util.JSON[util.NotAuthorized], nil)
}

// Test_makeAuthHandlerServerError verifies that makeAuthHandler
// generates appropriate server error output and logging when an internal
// server error occurs.
func Test_makeAuthHandlerServerError(t *testing.T) {
	// Test function which returns a server error
	errServer := errors.New("internal server error")
	serverErrFn := func(r *http.Request) (*models.User, *models.Session, error, error) {
		return nil, nil, nil, errServer
	}

	test_makeAuthHandler(t, serverErrFn, okHandler(), http.StatusInternalServerError, util.JSON[util.InternalServerError], errServer)
}

// Test_makeAuthHandlerNoError verifies that makeAuthHandler allows a
// HTTP handler to proceed with no authentication context.
func Test_makeAuthHandlerNoError(t *testing.T) {
	// Test function which returns OK
	okFn := func(r *http.Request) (*models.User, *models.Session, error, error) {
		return nil, nil, nil, nil
	}

	test_makeAuthHandler(t, okFn, okHandler(), http.StatusOK, []byte("hello world"), nil)
}

// Test_makeAuthHandlerUserSessionContext verifies that makeAuthHandler allows
// a HTTP handler to proceed with a user and session authentication context.
func Test_makeAuthHandlerUserSessionContext(t *testing.T) {
	// Test function which stores some request context
	user := ditest.MockUser()
	session, err := user.NewSession(time.Now().Add(1 * time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	contextFn := func(r *http.Request) (*models.User, *models.Session, error, error) {
		return user, session, nil, nil
	}

	// Test handler which retrieves data from request context
	contextHandlerFn := func(w http.ResponseWriter, r *http.Request) {
		cUser := User(r)
		cSession := Session(r)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(cUser.Username + cSession.Key))
	}

	test_makeAuthHandler(t, contextFn, contextHandlerFn, http.StatusOK, []byte(user.Username+session.Key), nil)
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
		// Invalid credential pair (no colon)
		{"Basic dGVzdA==", "", "", errInvalidBasicCredentialPair},
		// Valid pair
		{"Basic dGVzdDp0ZXN0", "test", "test", nil},
	}

	for _, test := range tests {
		// Split input header into credentials
		username, password, err := basicCredentials(test.input)
		if err != nil && err != test.err {
			t.Fatalf("unexpected err: %v != %v", err, test.err)
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

// test_makeAuthHandler accepts input parameters and expected results for
// makeAuthHandler, and ensures it behaves as expected.
func test_makeAuthHandler(t *testing.T, fn AuthenticateFunc, h http.HandlerFunc, code int, body []byte, expErr error) {
	// Capture log output in buffer
	buffer := bytes.NewBuffer(nil)
	log.SetOutput(buffer)

	// Create mock request
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Invoke auth handler and capture output
	w := httptest.NewRecorder()
	makeAuthHandler(fn, h).ServeHTTP(w, r)

	// Verify expected code
	if w.Code != code {
		t.Errorf("unexpected code: %v != %v", w.Code, code)
	}

	// Verify expected body
	if !bytes.Equal(w.Body.Bytes(), body) {
		t.Errorf("unexpected body: %v != %v", string(w.Body.Bytes()), string(body))
	}

	// Check for error body from server
	if expErr != nil {
		if !bytes.Contains(buffer.Bytes(), []byte(expErr.Error())) {
			t.Errorf("error not logged: %v", expErr)
		}
	}
}

// Test handler called on successful authentication
func okHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello world"))
	}
}
