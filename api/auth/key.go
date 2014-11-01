package auth

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/mdlayher/deltaiota/data/models"
)

var (
	// errNoKey is returned when no API key is provided for authentication.
	errNoKey = &Error{
		Reason: "no API key provided",
	}

	// errInvalidKey is returned when an invalid key is provided for
	// authentication.
	errInvalidKey = &Error{
		Reason: "invalid API key",
	}

	// errExpiredKey is returned when an expired key is provided for
	// authentication.
	errExpiredKey = &Error{
		Reason: "expired API key",
	}
)

// KeyAuthHandler is a http.HandlerFunc which performs API Key authentication.
func (a *Context) KeyAuthHandler(h http.HandlerFunc) http.HandlerFunc {
	return makeAuthHandler(a.keyAuthenticate, h)
}

// keyAuthenticate is a AuthenticateFunc which authenticates a user via API key,
// using HTTP Basic to pass the credentials.
// On success, a user and session are returned.  On failure, either a
// client or server error is returned.
func (a *Context) keyAuthenticate(r *http.Request) (*models.User, *models.Session, error, error) {
	// Attempt to fetch username/key pair from Authorization header
	username, key, err := basicCredentials(r.Header.Get("Authorization"))
	if err != nil {
		// Return client authentication error
		return nil, nil, err, nil
	}

	// Check for blank credentials
	if username == "" {
		return nil, nil, errNoUsername, nil
	}
	if key == "" {
		return nil, nil, errNoKey, nil
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

	// Attempt to select session for authentication by key
	session, err := a.db.SelectSessionByKey(key)
	if err != nil {
		// Check for unknown session
		if err == sql.ErrNoRows {
			return nil, nil, errInvalidKey, nil
		}

		return nil, nil, nil, err
	}

	// Verify key belongs to this user
	if user.ID != session.UserID {
		return nil, nil, errInvalidKey, nil
	}

	// Verify key is not expired
	if session.IsExpired() {
		// Delete expired key
		if err := a.db.DeleteSession(session); err != nil {
			return nil, nil, nil, err
		}

		// Return expired key error
		return nil, nil, errExpiredKey, nil
	}

	// Update expire time, since authentication succeeded
	session.SetExpire(time.Now().Add(SessionDuration))
	if err := a.db.UpdateSession(session); err != nil {
		// If database is readonly, ignore error
		if !a.db.IsReadonly(err) {
			return nil, nil, nil, err
		}
	}

	// Return authenticated user and session
	return user, session, nil, nil
}
