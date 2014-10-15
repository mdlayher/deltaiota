package auth

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mdlayher/deltaiota/api/util"
	"github.com/mdlayher/deltaiota/data/models"
	"github.com/mdlayher/deltaiota/ditest"
)

// Test_makeAuthHandler_ClientError verifies that makeAuthHandler generates
// appropriate output for a client error.
func Test_makeAuthHandler_ClientError(t *testing.T) {
	// Test function which returns a formatted client error
	reason := "foo bar"
	clientErrFn := func(r *http.Request) (*models.User, *models.Session, error, error) {
		return nil, nil, &Error{
			Reason: reason,
		}, nil
	}

	// Build client failure JSON
	authFailJSON := []byte(`{"error":{"code":401,"message":"authentication failed: ` + reason + `"}}`)

	// Perform test
	if err := makeAuthHandlerTest(clientErrFn, okHandler(), http.StatusUnauthorized, authFailJSON, nil); err != nil {
		t.Fatal(err)
	}
}

// Test_makeAuthHandler_ClientNonStandardError verifies that makeAuthHandler
// generates a generic error when the wrapped Error type is not used.
func Test_makeAuthHandler_ClientNonStandardError(t *testing.T) {
	// Test function which returns a non-standard client error
	clientBadErrFn := func(r *http.Request) (*models.User, *models.Session, error, error) {
		return nil, nil, errors.New("some error"), nil
	}

	// Perform test
	if err := makeAuthHandlerTest(clientBadErrFn, okHandler(), http.StatusUnauthorized, util.JSON[util.NotAuthorized], nil); err != nil {
		t.Fatal(err)
	}
}

// Test_makeAuthHandler_ServerError verifies that makeAuthHandler
// generates appropriate server error output and logging when an internal
// server error occurs.
func Test_makeAuthHandler_ServerError(t *testing.T) {
	// Test function which returns a server error
	errServer := errors.New("internal server error")
	serverErrFn := func(r *http.Request) (*models.User, *models.Session, error, error) {
		return nil, nil, nil, errServer
	}

	// Perform test
	if err := makeAuthHandlerTest(serverErrFn, okHandler(), http.StatusInternalServerError, util.JSON[util.InternalServerError], errServer); err != nil {
		t.Fatal(err)
	}
}

// Test_makeAuthHandler_NoError verifies that makeAuthHandler allows a
// HTTP handler to proceed with no authentication context.
func Test_makeAuthHandler_NoError(t *testing.T) {
	// Test function which returns OK
	okFn := func(r *http.Request) (*models.User, *models.Session, error, error) {
		return nil, nil, nil, nil
	}

	// Perform test
	if err := makeAuthHandlerTest(okFn, okHandler(), http.StatusOK, []byte("hello world"), nil); err != nil {
		t.Fatal(err)
	}
}

// Test_makeAuthHandler_UserSessionContext verifies that makeAuthHandler allows
// a HTTP handler to proceed with a user and session authentication context.
func Test_makeAuthHandler_UserSessionContext(t *testing.T) {
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

	// Perform test
	if err := makeAuthHandlerTest(contextFn, contextHandlerFn, http.StatusOK, []byte(user.Username+session.Key), nil); err != nil {
		t.Fatal(err)
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

// makeAuthHandlerTest accepts input parameters and expected results for
// makeAuthHandler, and ensures it behaves as expected.
func makeAuthHandlerTest(fn AuthenticateFunc, h http.HandlerFunc, code int, body []byte, expErr error) error {
	// Capture log output in buffer
	buffer := bytes.NewBuffer(nil)
	log.SetOutput(buffer)

	// Create mock request
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		return err
	}

	// Invoke auth handler and capture output
	w := httptest.NewRecorder()
	makeAuthHandler(fn, h).ServeHTTP(w, r)

	// Verify expected code
	if w.Code != code {
		return fmt.Errorf("unexpected code: %v != %v", w.Code, code)
	}

	// Verify expected body
	if !bytes.Equal(w.Body.Bytes(), body) {
		return fmt.Errorf("unexpected body: %v != %v", string(w.Body.Bytes()), string(body))
	}

	// End test if status OK
	if w.Code == http.StatusOK {
		return nil
	}

	// Check for error body from server
	if expErr != nil {
		if !bytes.Contains(buffer.Bytes(), []byte(expErr.Error())) {
			return fmt.Errorf("error not logged: %v", expErr)
		}
	}

	return nil
}

// Test handler called on successful authentication
func okHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello world"))
	}
}
