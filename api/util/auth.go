package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/mdlayher/deltaiota/data"
	"github.com/mdlayher/deltaiota/data/models"

	"github.com/gorilla/context"
)

// AuthenticateFunc is a function which may be used to authenticate a user from an
// input HTTP request.  On success, a user is returned.  On failure, either a client
// or server error is returned.
type AuthenticateFunc func(r *http.Request) (*models.User, error, error)

// AuthContext provides all shared members required for user authentication.
type AuthContext struct {
	DB *data.DB
}

// BasicAuthHandler is a http.HandlerFunc which performs HTTP Basic authentication.
func (a *AuthContext) BasicAuthHandler(h http.HandlerFunc) http.HandlerFunc {
	return makeAuthHandler(a.basicAuthenticate, h)
}

// basicAuthenticate is a AuthenticateFunc which authenticates a user via HTTP Basic.
// On success, a user is returned.  On failure, either a client or server error is
// returned.
func (a *AuthContext) basicAuthenticate(r *http.Request) (*models.User, error, error) {
	return nil, nil, errors.New("basicAuthenticate not yet implemented")
}

// AuthError is an error returned on client authentication failure.
type AuthError struct {
	Reason string
}

// Error returns the string representation of an AuthError.
func (e *AuthError) Error() string {
	return fmt.Sprintf("authentication failed: %s", e.Reason)
}

// makeAuthHandler generates a common authentication http.HandlerFunc using an input
// AuthenticateFunc and http.HandlerFunc.
func makeAuthHandler(fn AuthenticateFunc, h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Invoke authentication function, retrieve output user, client
		// error, and server error
		user, cErr, sErr := fn(r)

		// On server error, log error and return internal server error
		if sErr != nil {
			log.Println(sErr)
			w.WriteHeader(utilCode[utilInternalServerError])
			return
		}

		// On client error, return details regarding failure
		if cErr != nil {
			code := utilCode[utilNotAuthorized]
			w.WriteHeader(code)

			// If not a specific authentication error, return generic error
			authErr, ok := cErr.(*AuthError)
			if !ok {
				w.Write(utilJSON[utilNotAuthorized])
				return
			}

			// Marshal specific error to JSON
			body, err := json.Marshal(ErrRes(code, authErr.Error()))
			if err != nil {
				// On failed JSON marshal, return server error
				log.Println(err)
				w.WriteHeader(utilCode[utilInternalServerError])
				return
			}

			// If not a HEAD request, write error body
			if r.Method != "HEAD" {
				w.Write(body)
			}
			return
		}

		// Authentication succeeded, store user for later use
		context.Set(r, ctxUser, user)

		// Invoke input handler
		h.ServeHTTP(w, r)

		// Clear context after use
		context.Clear(r)
	})
}
