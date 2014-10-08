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

// JSONAPIHandler returns a http.HandlerFunc by invoking an input JSONAPIFunc.
func JSONAPIHandler(fn JSONAPIFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Invoke input closure to retrieve a HTTP status, a response body, and any
		// possible errors which occurred.
		code, body, err := fn(r, mux.Vars(r))
		if err != nil {
			log.Println(err)
			code = utilCode[utilInternalServerError]
			body = utilJSON[utilInternalServerError]
		}

		// Write HTTP status code
		w.Header().Set(httpContentType, jsonContentType)
		w.WriteHeader(code)

		// If HTTP HEAD request, write no body
		if r.Method == "HEAD" {
			w.Write(nil)
		}

		w.Write(body)
	})
}
