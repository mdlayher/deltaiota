// Package auth provides HTTP API authentication handlers for the Phi Mu Alpha
// Sinfonia - Delta Iota chapter website.
package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

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

var (
	// errNoUsername is returned when no username is provided for authentication.
	errNoUsername = &Error{
		Reason: "no username provided",
	}

	// errInvalidUsername is returned when an invalid username is provided for
	// authentication.
	errInvalidUsername = &Error{
		Reason: "invalid username",
	}

	// errNoAuthorizationHeader is returned when an input Authorization header
	// is blank.
	errNoAuthorizationHeader = &Error{
		Reason: "no HTTP Authorization header",
	}

	// errNoAuthorizationType is returned when an input Authorization header
	// contains no type.
	errNoAuthorizationType = &Error{
		Reason: "no HTTP Authorization type",
	}

	// errNotBasicAuthorization is returned when an input Authorization header
	// is not the HTTP Basic type.
	errNotBasicAuthorization = &Error{
		Reason: "not HTTP Basic Authorization type",
	}

	// errInvalidBase64Authorization is returned when an input Authorization header
	// does not contain valid base64-encoded data.
	errInvalidBase64Authorization = &Error{
		Reason: "invalid base64 HTTP Basic Authorization header",
	}

	// errInvalidBasicCredentialPair is returned when an input Authorization header
	// does not contain a valid credential pair.
	errInvalidBasicCredentialPair = &Error{
		Reason: "invalid credential pair in HTTP Basic Authorization header",
	}
)

var (
	// SessionDuration is the duration which a session may exist before
	// it expires, and is deleted.
	SessionDuration = time.Duration(7 * 24 * time.Hour)
)

// AuthenticateFunc is a function which may be used to authenticate a user from an
// input HTTP request.  On success, a user is returned.  On failure, either a client
// or server error is returned.
type AuthenticateFunc func(r *http.Request) (*models.User, *models.Session, error, error)

// Context provides all shared members required for user authentication.
type Context struct {
	db *data.DB
}

// NewContext initializes a new Context with the input parameters.
func NewContext(db *data.DB) *Context {
	return &Context{
		db: db,
	}
}

// Error is an error returned on client authentication failure.
type Error struct {
	Reason string
}

// Error returns the string representation of an Error.
func (e *Error) Error() string {
	return fmt.Sprintf("authentication failed: %s", e.Reason)
}

// SetSession sets a gorilla/context Session for the input http.Request.
func SetSession(r *http.Request, s *models.Session) {
	context.Set(r, ctxSession, s)
}

// Session returns the gorilla/context Session for the input http.Request.
// This function will panic if the user is not properly authenticated, and
// should only be used in handlers which are always authenticated.
func Session(r *http.Request) *models.Session {
	return context.Get(r, ctxSession).(*models.Session)
}

// SetUser sets a gorilla/context User for the input http.Request.
func SetUser(r *http.Request, s *models.User) {
	context.Set(r, ctxUser, s)
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
			authErr, ok := cErr.(*Error)
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
		SetUser(r, user)
		SetSession(r, session)

		// Invoke input handler
		h.ServeHTTP(w, r)

		// Clear context after use
		context.Clear(r)
	})
}

// basicCredentials returns HTTP Basic authentication credentials from an input header
// in the form: base64(user + ':' + password).
func basicCredentials(header string) (string, string, error) {
	// No headed provided
	if header == "" {
		return "", "", errNoAuthorizationHeader
	}

	// Ensure 2 elements
	basic := strings.Split(header, " ")
	if len(basic) != 2 {
		return "", "", errNoAuthorizationType
	}

	// Ensure valid format
	if basic[0] != "Basic" {
		return "", "", errNotBasicAuthorization
	}

	// Decode base64'd username:password pair
	buf, err := base64.URLEncoding.DecodeString(basic[1])
	if err != nil {
		return "", "", errInvalidBase64Authorization
	}

	// Split into username/password
	pair := strings.SplitN(string(buf), ":", 2)
	if len(pair) < 2 {
		return "", "", errInvalidBasicCredentialPair
	}

	return pair[0], pair[1], nil
}
