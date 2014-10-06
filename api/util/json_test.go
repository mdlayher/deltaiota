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

// TestJSONAPIHandler verifies that JSONAPIHandler generates the appropriate
// http.HandlerFunc for one or more input JSONAPIFunc.
func TestJSONAPIHandler(t *testing.T) {
	// emptyFn returns HTTP OK and nothing else
	emptyFn := func(r *http.Request) (int, []byte, error) {
		return http.StatusOK, nil, nil
	}

	// bodyFn returns HTTP OK and a small message body
	expBody := []byte("hello world")
	bodyFn := func(r *http.Request) (int, []byte, error) {
		return http.StatusOK, expBody, nil
	}

	// errFn returns HTTP Internal Server Error and an error
	expErr := errors.New("a fake test error")
	errFn := func(r *http.Request) (int, []byte, error) {
		return http.StatusInternalServerError, nil, expErr
	}

	// Table of test JSONAPIFunc and expected results
	var tests = []struct {
		fn   JSONAPIFunc
		code int
		body []byte
		err  error
	}{
		// Empty function
		{emptyFn, http.StatusOK, nil, nil},
		// Body function
		{bodyFn, http.StatusOK, expBody, nil},
		// Error function
		{errFn, http.StatusInternalServerError, nil, expErr},
	}

	// Iterate and run all tests
	for _, test := range tests {
		// Capture log output in buffer
		buffer := bytes.NewBuffer(nil)
		log.SetOutput(buffer)

		// Capture output, invoke function as http.HandlerFunc
		w := httptest.NewRecorder()
		JSONAPIHandler(test.fn).ServeHTTP(w, nil)

		// If body not empty, ensure JSON response header set
		if len(test.body) > 0 {
			if contentType := w.Header().Get(httpContentType); contentType != jsonContentType {
				t.Fatalf("unexpected Content-Type header: %v != %v", contentType, jsonContentType)
			}
		}

		// Verify expected code
		if w.Code != test.code {
			t.Fatal("unexpected code: %v != %v", w.Code, test.code)
		}

		// If no error, verify expected body
		if test.err == nil {
			if !bytes.Equal(w.Body.Bytes(), test.body) {
				t.Fatal("unexpected body: %v != %v", w.Body.Bytes(), test.body)
			}

			continue
		}

		// Verify expected error, by unmarshaling body
		var errRes ErrorResponse
		if err := json.Unmarshal(w.Body.Bytes(), &errRes); err != nil {
			t.Fatal(err)
		}

		// Verify error fields
		if errRes.Error.Code != test.code {
			t.Fatal("unexpected error code: %v != %v", errRes.Error.Code, test.code)
		}
		if errRes.Error.Message != utilInternalServerError {
			t.Fatal("unexpected error message: %v != %v", errRes.Error.Message, utilInternalServerError)
		}

		// Verify error was logged
		if !bytes.Contains(buffer.Bytes(), []byte(test.err.Error())) {
			t.Fatal("error not logged:", test.err)
		}
	}
}