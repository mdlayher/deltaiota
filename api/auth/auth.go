package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/mdlayher/deltaiota/api/util"
	"github.com/mdlayher/deltaiota/data"
	"github.com/mdlayher/deltaiota/data/models"

	"github.com/gorilla/context"
)

const (
	// ctxSession is the named key used to fetch a Session from gorilla/context.
	ctxSession = "session"

	// ctxUser is the named key used to fetch a User from gorilla/context.
	ctxUser = "user"
)

// AuthenticateFunc is a function which may be used to authenticate a user from an
// input HTTP request.  On success, a user is returned.  On failure, either a client
// or server error is returned.
type AuthenticateFunc func(r *http.Request) (*models.User, *models.Session, error, error)

// AuthContext provides all shared members required for user authentication.
type AuthContext struct {
	db *data.DB
}

// NewContext initializes a new AuthContext with the input parameters.
func NewContext(db *data.DB) *AuthContext {
	return &AuthContext{
		db: db,
	}
}

// AuthError is an error returned on client authentication failure.
type AuthError struct {
	Reason string
}

// Error returns the string representation of an AuthError.
func (e *AuthError) Error() string {
	return fmt.Sprintf("authentication failed: %s", e.Reason)
}

// Session returns the gorilla/context Session for the input http.Request.
// This function will panic if the user is not properly authenticated, and
// should only be used in handlers which are always authenticated.
func Session(r *http.Request) *models.Session {
	return context.Get(r, ctxSession).(*models.Session)
}

// User returns the gorilla/context User for the input http.Request.
// This function will panic if the user is not properly authenticated, and
// should only be used in handlers which are always authenticated.
func User(r *http.Request) *models.User {
	return context.Get(r, ctxUser).(*models.User)
}

// makeAuthHandler generates a common authentication http.HandlerFunc using an input
// AuthenticateFunc and http.HandlerFunc.
func makeAuthHandler(fn AuthenticateFunc, h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Invoke authentication function, retrieve output user, session,
		// client error, and server error
		user, session, cErr, sErr := fn(r)

		// On server error, log error and return internal server error
		if sErr != nil {
			log.Println(sErr)
			w.WriteHeader(util.Code[util.InternalServerError])
			w.Write(util.JSON[util.InternalServerError])
			return
		}

		// On client error, return details regarding failure
		if cErr != nil {
			code := util.Code[util.NotAuthorized]
			w.WriteHeader(code)

			// If not a specific authentication error, return generic error
			authErr, ok := cErr.(*AuthError)
			if !ok {
				w.Write(util.JSON[util.NotAuthorized])
				return
			}

			// Marshal specific error to JSON
			body, err := json.Marshal(util.ErrRes(code, authErr.Error()))
			if err != nil {
				// On failed JSON marshal, return server error
				log.Println(err)
				w.WriteHeader(util.Code[util.InternalServerError])
				return
			}

			// If not a HEAD request, write error body
			if r.Method != "HEAD" {
				w.Write(body)
			}
			return
		}

		// Authentication succeeded, store user and session for later use
		context.Set(r, ctxUser, user)
		context.Set(r, ctxUser, session)

		// Invoke input handler
		h.ServeHTTP(w, r)

		// Clear context after use
		context.Clear(r)
	})
}
