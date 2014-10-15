package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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

	// Perform test with HTTP GET
	if err := testJSONAPIHandler(emptyFn, "GET", http.StatusOK, nil, nil); err != nil {
		t.Fatal(err)
	}

	// Perform test with HTTP HEAD
	if err := testJSONAPIHandler(emptyFn, "HEAD", http.StatusOK, nil, nil); err != nil {
		t.Fatal(err)
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

	// Perform test with HTTP GET
	if err := testJSONAPIHandler(bodyFn, "GET", http.StatusOK, expBody, nil); err != nil {
		t.Fatal(err)
	}

	// Perform test with HTTP HEAD - no response body expected
	if err := testJSONAPIHandler(bodyFn, "HEAD", http.StatusOK, nil, nil); err != nil {
		t.Fatal(err)
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

	// Perform test with HTTP GET
	if err := testJSONAPIHandler(errFn, "GET", http.StatusInternalServerError, nil, expErr); err != nil {
		t.Fatal(err)
	}

	// Perform test with HTTP HEAD - no response body expected
	if err := testJSONAPIHandler(errFn, "HEAD", http.StatusInternalServerError, nil, expErr); err != nil {
		t.Fatal(err)
	}
}

// TestJSONAPIHandlerMethodNotAllowed verifies that JSONAPIHandler returns correct
// results for an input function with an unknown HTTP request type.
func TestJSONAPIHandlerMethodNotAllowed(t *testing.T) {
	if err := testJSONAPIHandler(MethodNotAllowed, "CAT", http.StatusMethodNotAllowed, JSON[methodNotAllowed], nil); err != nil {
		t.Fatal(err)
	}
}

// testJSONAPIHandler accepts input parameters and expected results for
// JSONAPIHandler, and ensures it behaves as expected.
func testJSONAPIHandler(fn JSONAPIFunc, method string, code int, body []byte, expErr error) error {
	// Capture log output in buffer
	buffer := bytes.NewBuffer(nil)
	log.SetOutput(buffer)

	// Create mock HTTP request
	r, err := http.NewRequest(method, "/", nil)
	if err != nil {
		return err
	}

	// Capture output, invoke function as http.HandlerFunc
	w := httptest.NewRecorder()
	JSONAPIHandler(fn).ServeHTTP(w, r)

	// If body not empty, ensure JSON response header set
	if len(body) > 0 {
		if contentType := w.Header().Get(httpContentType); contentType != jsonContentType {
			return fmt.Errorf("unexpected Content-Type header: %v != %v", contentType, jsonContentType)
		}
	}

	// Verify expected code
	if w.Code != code {
		return fmt.Errorf("unexpected code: %v != %v", w.Code, code)
	}

	// If no error, verify expected body
	if expErr == nil {
		if !bytes.Equal(w.Body.Bytes(), body) {
			return fmt.Errorf("unexpected body: %v != %v", w.Body.Bytes(), body)
		}

		return nil
	}

	// If HEAD, there will be no body
	if method == "HEAD" {
		// Verify no body
		length := len(w.Body.Bytes())
		if length > 0 {
			return fmt.Errorf("non-empty body for HTTP HEAD: %v", length)
		}

		return nil
	}

	// Verify expected error, by unmarshaling body
	var errRes ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &errRes); err != nil {
		return err
	}

	// Verify error fields
	if errRes.Error.Code != code {
		return fmt.Errorf("unexpected error code: %v != %v", errRes.Error.Code, code)
	}
	if errRes.Error.Message != InternalServerError {
		return fmt.Errorf("unexpected error message: %v != %v", errRes.Error.Message, InternalServerError)
	}

	// Verify error was logged
	if !bytes.Contains(buffer.Bytes(), []byte(expErr.Error())) {
		return fmt.Errorf("error not logged: %v", expErr)
	}

	return nil
}
