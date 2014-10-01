// Package util provides common API helper types and functions for the Phi Mu
// Alpha Sinfonia - Delta Iota chapter website.
package util

import (
	"net/http"
)

const (
	// httpContentType is the name of the Content-Type HTTP header.
	httpContentType = "Content-Type"

	// jsonContentType is the name of the JSON HTTP Content-Type.
	jsonContentType = "application/json"
)

// JSONAPIFunc is a type which accepts an input http.Request, and responds with a HTTP
// status code and a slice of bytes containing the response body.
type JSONAPIFunc func(r *http.Request) (int, []byte)

// JSONAPIHandler returns a http.HandlerFunc by invoking an input JSONAPIFunc, setting
// necessary headers, and writing a response to the client.
func JSONAPIHandler(fn JSONAPIFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Invoke input closure to retrieve a HTTP status and JSON body
		code, body := fn(r)

		// If body is non-empty, set JSON content type
		if len(body) > 0 {
			w.Header().Set(httpContentType, jsonContentType)
		}

		// Write HTTP status code and body
		w.WriteHeader(code)
		w.Write(body)
	})
}
