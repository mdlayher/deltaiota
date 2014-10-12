package util

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/gorilla/mux"
)

// Vars is a map of route variables, typically injected by gorilla/mux; though they
// can also be manually injected for testing handlers.
type Vars map[string]string

// JSONAPIFunc is a type which accepts an input http.Request and a map of route
// variables, and responds with a HTTP status code, and either a JSON response body,
// or an error which is reported as an internal server error to the client.
type JSONAPIFunc func(r *http.Request, vars Vars) (int, []byte, error)

// JSONAPIHandler returns a http.HandlerFunc by invoking an input JSONAPIFunc.
func JSONAPIHandler(fn JSONAPIFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Invoke input closure to retrieve a HTTP status, a response body, and any
		// possible errors which occurred.
		code, body, err := fn(r, mux.Vars(r))
		if err != nil {
			// Generate string with request information
			reqLog := fmt.Sprintf("[%s: %s %s]", r.RemoteAddr, r.Method, r.URL.Path)

			// Check for internal error, with more debugging information
			intErr, ok := err.(*InternalError)
			if !ok {
				// If not wrapped error, print basic information
				log.Printf("%s %s", reqLog, err.Error())
			} else {
				// If wrapped error, print advanced information
				// In the future, additional error hooks could be added here
				log.Printf("%s internal error: [file: %s:%d, err: %s]", reqLog, filepath.Base(intErr.File), intErr.Line, intErr.Err)
			}
		}

		// Write HTTP status code
		w.Header().Set(httpContentType, jsonContentType)
		w.Header().Set(httpConnection, "close")
		w.WriteHeader(code)

		// If HTTP HEAD request, write no body
		if r.Method == "HEAD" {
			w.Header().Set(httpContentLength, "0")
			w.Write(nil)
			return
		}

		w.Header().Set(httpContentLength, strconv.Itoa(len(body)))
		w.Write(body)
	})
}

// JSONAPIErr accepts an internal error, wraps it in useful information for
// debugging, and generates the appropriate JSONAPIFunc return signature,
// for convenience and reduced code reptition.
func JSONAPIErr(err error) (int, []byte, error) {
	// Attempt to retrieve information about the calling function
	if _, file, line, ok := runtime.Caller(1); ok {
		// Wrap error with useful debugging information
		err = &InternalError{
			File: file,
			Line: line,
			Err:  err,
		}
	}

	// Return HTTP code, body, and (wrapped, if possible) error
	return Code[InternalServerError], JSON[InternalServerError], err
}

// MethodNotAllowed is a JSONAPIFunc which returns HTTP 405, as well as the
// appropriate JSON response body.
func MethodNotAllowed(r *http.Request, vars Vars) (int, []byte, error) {
	return Code[methodNotAllowed], JSON[methodNotAllowed], nil
}
