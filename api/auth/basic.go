package auth

import (
	"database/sql"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/mdlayher/deltaiota/data/models"
)

var (
	// errNoAuthorizationHeader is returned when an input Authorization header
	// is blank.
	errNoAuthorizationHeader = &AuthError{
		Reason: "no HTTP Authorization header",
	}

	// errNoAuthorizationType is returned when an input Authorization header
	// contains no type.
	errNoAuthorizationType = &AuthError{
		Reason: "no HTTP Authorization type",
	}

	// errNotBasicAuthorization is returned when an input Authorization header
	// is not the HTTP Basic type.
	errNotBasicAuthorization = &AuthError{
		Reason: "not HTTP Basic Authorization type",
	}

	// errInvalidBase64Authorization is returned when an input Authorization header
	// does not contain valid base64-encoded data.
	errInvalidBase64Authorization = &AuthError{
		Reason: "invalid base64 HTTP Basic Authorization header",
	}

	// errInvalidBasicCredentialPair is returned when an input Authorization header
	// does not contain a valid credential pair.
	errInvalidBasicCredentialPair = &AuthError{
		Reason: "invalid credential pair in HTTP Basic Authorization header",
	}
)

// BasicAuthHandler is a http.HandlerFunc which performs HTTP Basic authentication.
func (a *AuthContext) BasicAuthHandler(h http.HandlerFunc) http.HandlerFunc {
	return makeAuthHandler(a.basicAuthenticate, h)
}

// basicAuthenticate is a AuthenticateFunc which authenticates a user via HTTP Basic.
// On success, a user is returned (nil session is returned).  On failure, either a
// client or server error is returned.
func (a *AuthContext) basicAuthenticate(r *http.Request) (*models.User, *models.Session, error, error) {
	// Attempt to fetch username/password pair from Authorization header
	username, password, err := basicCredentials(r.Header.Get("Authorization"))
	if err != nil {
		// Return client authentication error
		return nil, nil, err, nil
	}

	// Check for blank credentials
	if username == "" {
		return nil, nil, &AuthError{
			Reason: "no username provided",
		}, nil
	}
	if password == "" {
		return nil, nil, &AuthError{
			Reason: "no password provided",
		}, nil
	}

	// Attempt to select user for authentication by username
	user, err := a.db.SelectUserByUsername(username)
	if err != nil {
		// Check for unknown user
		if err == sql.ErrNoRows {
			return nil, nil, &AuthError{
				Reason: "invalid username",
			}, nil
		}

		return nil, nil, nil, err
	}

	// Attempt authentication using input password
	if err := user.TryPassword(password); err != nil {
		// Check for invalid password
		if err == models.ErrInvalidPassword {
			return nil, nil, &AuthError{
				Reason: err.Error(),
			}, nil
		}

		return nil, nil, nil, err
	}

	// Return authenticated user
	return user, nil, nil, nil
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
	if len(pair) < 2 || pair[0] == "" || pair[1] == "" {
		return "", "", errInvalidBasicCredentialPair
	}

	return pair[0], pair[1], nil
}
