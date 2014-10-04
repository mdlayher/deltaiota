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

// JSONAPIFunc is a type which accepts an input http.Request, and responds with a HTTP
// status code and a struct which can be marshaled to JSON, containing the response body.
type JSONAPIFunc func(r *http.Request) (int, JSONable)

// JSONable is an interface which structs implement if they can be marshaled to JSON
// and returned via JSONAPIHandler.
type JSONable interface{}

// JSONAPIHandler returns a http.HandlerFunc by invoking an input JSONAPIFunc, setting
// necessary headers, and writing a response to the client.
func JSONAPIHandler(fn JSONAPIFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Invoke input closure to retrieve a HTTP status and struct which can
		// be marshaled to JSON
		code, out := fn(r)

		// Marshal struct to JSON
		body, err := json.Marshal(out)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
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
