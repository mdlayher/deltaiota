// Package util provides common API helper types and functions for the Phi Mu
// Alpha Sinfonia - Delta Iota chapter website.
package util

import (
	"encoding/json"
	"log"
	"net/http"
)

const (
	// httpContentType is the name of the Content-Type HTTP header.
	httpContentType = "Content-Type"

	// jsonContentType is the name of the JSON HTTP Content-Type.
	jsonContentType = "application/json"
)

// JSON util, human-readable client error responses.
const (
	utilInternalServerError = "internal server error"
)

// JSON util, map of client errors to response codes.
var utilCode = map[string]int{
	utilInternalServerError: http.StatusInternalServerError,
}

// Generated JSON responses for various client-facing errors.
var utilJSON = map[string][]byte{}

// init initializes the stored JSON responses for client-facing errors.
func init() {
	// Iterate all error strings and code integers
	for k, v := range utilCode {
		// Generate error response with appropriate string and code
		body, err := json.Marshal(ErrRes(v, k))
		if err != nil {
			panic(err)
		}

		// Store for later use
		utilJSON[k] = body
	}
}

// ErrorResponse is the response returned whenever the API generates a client or server
// error.  It contains a nested Error object, which provides further information.
type ErrorResponse struct {
	Error *Error `json:"error"`
}

// Error contains a status code and human-readable error message, and is generated for
// client and server facing errors.
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ErrRes generates and returns an ErrorResponse using the input parameters.
func ErrRes(code int, message string) *ErrorResponse {
	return &ErrorResponse{
		Error: &Error{
			Code:    code,
			Message: message,
		},
	}
}

// JSONAPIFunc is a type which accepts an input http.Request, and responds with a HTTP
// status code, and either a JSON response body, or an error which is reported as an
// internal server error to the client.
type JSONAPIFunc func(r *http.Request) (int, []byte, error)

// JSONAPIHandler returns a http.HandlerFunc by invoking an input JSONAPIFunc, setting
// necessary headers, and writing a response to the client.
func JSONAPIHandler(fn JSONAPIFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Invoke input closure to retrieve a HTTP status, a response body, and any
		// possible errors which occurred.
		code, body, err := fn(r)
		if err != nil {
			log.Println(err)
			code = utilCode[utilInternalServerError]
			body = utilJSON[utilInternalServerError]
		}

		// If body is non-empty, set JSON content type
		if len(body) > 0 {
			w.Header().Set(httpContentType, jsonContentType)
		}

		// Write HTTP status code and body
		w.WriteHeader(code)
		w.Write(body)
	})
}
