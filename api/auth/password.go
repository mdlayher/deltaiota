package auth

import (
	"database/sql"
	"net/http"

	"github.com/mdlayher/deltaiota/data/models"
)

var (
	// errNoPassword is returned when no password is provided for authentication.
	errNoPassword = &Error{
		Reason: "no password provided",
	}

	// errInvalidPassword is returned when an invalid password is provided for
	// authentication.
	errInvalidPassword = &Error{
		Reason: models.ErrInvalidPassword.Error(),
	}
)

// PasswordAuthHandler is a http.HandlerFunc which performs HTTP Basic authentication
// using a username and password pair from an Authorization header.
func (a *Context) PasswordAuthHandler(h http.HandlerFunc) http.HandlerFunc {
	return makeAuthHandler(a.passwordAuthenticate, h)
}

// passwordAuthenticate is a AuthenticateFunc which authenticates a user via HTTP Basic,
// using a username and password pair from an Authorization header.
// On success, a user is returned (nil session is returned).  On failure, either a
// client or server error is returned.
func (a *Context) passwordAuthenticate(r *http.Request) (*models.User, *models.Session, error, error) {
	// Attempt to fetch username/password pair from Authorization header
	username, password, err := basicCredentials(r.Header.Get("Authorization"))
	if err != nil {
		// Return client authentication error
		return nil, nil, err, nil
	}

	// Check for blank credentials
	if username == "" {
		return nil, nil, errNoUsername, nil
	}
	if password == "" {
		return nil, nil, errNoPassword, nil
	}

	// Attempt to select user for authentication by username
	user, err := a.db.SelectUserByUsername(username)
	if err != nil {
		// Check for unknown user
		if err == sql.ErrNoRows {
			return nil, nil, errInvalidUsername, nil
		}

		return nil, nil, nil, err
	}

	// Attempt authentication using input password
	if err := user.TryPassword(password); err != nil {
		// Check for invalid password
		if err == models.ErrInvalidPassword {
			return nil, nil, errInvalidPassword, nil
		}

		return nil, nil, nil, err
	}

	// Return authenticated user
	return user, nil, nil, nil
}
