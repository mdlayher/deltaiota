package util

import (
	"log"
	"net/http"
)

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
