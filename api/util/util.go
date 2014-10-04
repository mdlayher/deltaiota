// Package util provides common API helper types and functions for the Phi Mu
// Alpha Sinfonia - Delta Iota chapter website.
package util

import (
	"encoding/json"
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
