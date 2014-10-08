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

// Test_makeAuthHandler ensures that makeAuthHandler generates the appropriate
// http.HandlerFunc for an input AuthenticateFunc and http.HandlerFunc.
func Test_makeAuthHandler(t *testing.T) {
	// Test function which returns a formatted client error
	clientErrFn := func(r *http.Request) (*models.User, *models.Session, error, error) {
		return nil, nil, &Error{
			Reason: "foo bar",
		}, nil
	}

	// Test function which returns a non-standard client error
	clientBadErrFn := func(r *http.Request) (*models.User, *models.Session, error, error) {
		return nil, nil, errors.New("some error"), nil
	}

	// Test function which returns a server error
	errServer := errors.New("internal server error")
	serverErrFn := func(r *http.Request) (*models.User, *models.Session, error, error) {
		return nil, nil, nil, errServer
	}

	// Test function which returns OK
	okFn := func(r *http.Request) (*models.User, *models.Session, error, error) {
		return nil, nil, nil, nil
	}

	// Test handler called on successful authentication
	handlerFn := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello world"))
	}

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

	// Build JSON response for client errors
	authFailJSON := func(reason string) []byte {
		return []byte(`{"error":{"code":401,"message":"authentication failed: ` + reason + `"}}`)
	}

	var tests = []struct {
		fn   AuthenticateFunc
		h    http.HandlerFunc
		code int
		body []byte
		err  error
	}{
		// Client error
		{clientErrFn, handlerFn, http.StatusUnauthorized, authFailJSON("foo bar"), nil},
		// Non-standard client error
		{clientBadErrFn, handlerFn, http.StatusUnauthorized, util.JSON[util.NotAuthorized], nil},
		// Server error
		{serverErrFn, handlerFn, http.StatusInternalServerError, util.JSON[util.InternalServerError], errServer},
		// No error
		{okFn, handlerFn, http.StatusOK, []byte("hello world"), nil},
		// No error, with context
		{contextFn, contextHandlerFn, http.StatusOK, []byte(user.Username + session.Key), nil},
	}

	for _, test := range tests {
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
		makeAuthHandler(test.fn, test.h).ServeHTTP(w, r)

		// Verify expected code
		if w.Code != test.code {
			t.Fatalf("unexpected code: %v != %v", w.Code, test.code)
		}

		// Verify expected body
		if !bytes.Equal(w.Body.Bytes(), test.body) {
			t.Fatalf("unexpected body: %v != %v", string(w.Body.Bytes()), string(test.body))
		}

		// Continue if status OK
		if w.Code == http.StatusOK {
			continue
		}

		// Check for error body from server
		if test.err != nil {
			if !bytes.Contains(buffer.Bytes(), []byte(test.err.Error())) {
				t.Fatal("error not logged:", test.err)
			}
		}
	}
}
