// Package util provides common API helper types and functions for the Phi Mu
// Alpha Sinfonia - Delta Iota chapter website.
package util

import (
	"encoding/json"
	"net/http"
)

const (
	// httpConnection is the name of the ConnectionHTTP header.
	httpConnection = "Connection"

	// httpContentLength is the name of the Content-Length HTTP header.
	httpContentLength = "Content-Length"

	// httpContentType is the name of the Content-Type HTTP header.
	httpContentType = "Content-Type"

	// jsonContentType is the name of the JSON HTTP Content-Type.
	jsonContentType = "application/json"
)

// JSON util, human-readable client error responses.
const (
	InternalServerError = "internal server error"
	NotAuthorized       = "not authorized"

	methodNotAllowed = "method not allowed"
)

// JSON util, map of client errors to response codes.
var Code = map[string]int{
	InternalServerError: http.StatusInternalServerError,
	NotAuthorized:       http.StatusUnauthorized,

	methodNotAllowed: http.StatusMethodNotAllowed,
}

// Generated JSON responses for various client-facing errors.
var JSON = map[string][]byte{}

// init initializes the stored JSON responses for client-facing errors.
func init() {
	// Iterate all error strings and code integers
	for k, v := range Code {
		// Generate error response with appropriate string and code
		body, err := json.Marshal(ErrRes(v, k))
		if err != nil {
			panic(err)
		}

		// Store for later use
		JSON[k] = body
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
