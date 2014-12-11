package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestJSONAPIHandlerNoBody verifies that JSONAPIHandler returns correct
// results for an input function with no response body.
func TestJSONAPIHandlerNoBody(t *testing.T) {
	// emptyFn returns HTTP OK and nothing else
	emptyFn := func(r *http.Request, vars Vars) (int, []byte, error) {
		return http.StatusOK, nil, nil
	}

	// Perform test with HTTP GET and HEAD
	for _, m := range []string{"GET", "HEAD"} {
		testJSONAPIHandler(t, emptyFn, m, http.StatusOK, nil, nil)
	}
}

// TestJSONAPIHandlerBody verifies that JSONAPIHandler returns correct
// results for an input function with a response body.
func TestJSONAPIHandlerBody(t *testing.T) {
	// bodyFn returns HTTP OK and a small message body
	expBody := []byte("hello world")
	bodyFn := func(r *http.Request, vars Vars) (int, []byte, error) {
		return http.StatusOK, expBody, nil
	}

	// Perform test with HTTP GET and HEAD
	for _, m := range []string{"GET", "HEAD"} {
		testJSONAPIHandler(t, bodyFn, m, http.StatusOK, expBody, nil)
	}
}

// TestJSONAPIHandlerServerError verifies that JSONAPIHandler returns correct
// results for an input function with an internal server error.
func TestJSONAPIHandlerServerError(t *testing.T) {
	// errFn returns HTTP Internal Server Error and an error
	expErr := errors.New("a fake test error")
	errFn := func(r *http.Request, vars Vars) (int, []byte, error) {
		return JSONAPIErr(expErr)
	}

	// Perform test with HTTP GET and HEAD
	for _, m := range []string{"GET", "HEAD"} {
		testJSONAPIHandler(t, errFn, m, http.StatusInternalServerError, nil, expErr)
	}
}

// TestJSONAPIHandlerMethodNotAllowed verifies that JSONAPIHandler returns correct
// results for an input function with an unknown HTTP request type.
func TestJSONAPIHandlerMethodNotAllowed(t *testing.T) {
	testJSONAPIHandler(t, MethodNotAllowed, "CAT", http.StatusMethodNotAllowed, JSON[methodNotAllowed], nil)
}

// testJSONAPIHandler accepts input parameters and expected results for
// JSONAPIHandler, and ensures it behaves as expected.
func testJSONAPIHandler(t *testing.T, fn JSONAPIFunc, method string, code int, body []byte, expErr error) {
	// Capture log output in buffer
	buffer := bytes.NewBuffer(nil)
	log.SetOutput(buffer)

	// Create mock HTTP request
	r, err := http.NewRequest(method, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Capture output, invoke function as http.HandlerFunc
	w := httptest.NewRecorder()
	JSONAPIHandler(fn).ServeHTTP(w, r)

	// If body not empty, ensure JSON response header set
	if len(body) > 0 {
		if contentType := w.Header().Get(httpContentType); contentType != jsonContentType {
			t.Fatalf("unexpected Content-Type header: %v != %v", contentType, jsonContentType)
		}
	}

	// Verify expected code
	if w.Code != code {
		t.Fatalf("unexpected code: %v != %v", w.Code, code)
	}

	// If no error, verify expected body
	if expErr == nil && method != "HEAD" {
		if !bytes.Equal(w.Body.Bytes(), body) {
			t.Fatalf("unexpected body: %v != %v", w.Body.Bytes(), body)
		}

		return
	} else if method == "HEAD" {
		// If HEAD, there will be no body
		length := len(w.Body.Bytes())
		if length > 0 {
			t.Fatalf("non-empty body for HTTP HEAD: %v", length)
		}

		return
	}

	// Verify expected error, by unmarshaling body
	var errRes ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &errRes); err != nil {
		t.Fatal(err)
	}

	// Verify error fields
	if errRes.Error.Code != code {
		t.Fatalf("unexpected error code: %v != %v", errRes.Error.Code, code)
	}
	if errRes.Error.Message != InternalServerError {
		t.Fatalf("unexpected error message: %v != %v", errRes.Error.Message, InternalServerError)
	}

	// Verify error was logged
	if !bytes.Contains(buffer.Bytes(), []byte(expErr.Error())) {
		t.Fatalf("error not logged: %v", expErr)
	}
}
