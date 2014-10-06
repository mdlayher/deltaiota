package util

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Vars is a map of route variables, typically injected by gorilla/mux; though they
// can also be manually injected for testing handlers.
type Vars map[string]string

// JSONAPIFunc is a type which accepts an input http.Request and a map of route
// variables, and responds with a HTTP status code, and either a JSON response body,
// or an error which is reported as an internal server error to the client.
type JSONAPIFunc func(r *http.Request, vars Vars) (int, []byte, error)

// JSONAPIHandler returns a http.HandlerFunc by invoking input a chain of JSONAPIFunc
// in order, until a response is written.
func JSONAPIHandler(functions ...JSONAPIFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Iterate each input function in order
		for _, fn := range functions {
			// Invoke input closure to retrieve a HTTP status, a response body, and any
			// possible errors which occurred.
			code, body, err := fn(r, mux.Vars(r))
			if err != nil {
				log.Println(err)
				code = utilCode[utilInternalServerError]
				body = utilJSON[utilInternalServerError]
			}

			// If body is empty, keep looping through chained functions until body is written
			if body == nil {
				continue
			}

			// Write HTTP status code and body
			w.Header().Set(httpContentType, jsonContentType)
			w.WriteHeader(code)
			w.Write(body)
			return
		}
	})
}
