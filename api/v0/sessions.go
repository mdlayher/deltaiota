package v0

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/mdlayher/deltaiota/api/auth"
	"github.com/mdlayher/deltaiota/api/util"
	"github.com/mdlayher/deltaiota/data/models"
)

// SessionsResponse is the output response for the Sessions API.
type SessionsResponse struct {
	Session *models.Session `json:"session"`
}

// SessionsAPI is a util.JSONAPIFunc, and is the single entry point for all non-POST
// methods for the Sessions API.  The POST endpoint is separate due to using password
// authentication, rather than key authentication.
// This method delegates to other methods as appropriate to handle incoming requests.
func (c *Context) SessionsAPI(r *http.Request, vars util.Vars) (int, []byte, error) {
	// Switch based on HTTP method
	switch r.Method {
	case "GET":
		return c.GetSession(r, vars)
	case "DELETE":
		return c.DeleteSession(r, vars)
	default:
		return util.MethodNotAllowed(r, vars)
	}
}

// GetSession is a util.JSONAPIFunc which returns the current Session and a HTTP 200
// on success, or a non-200 HTTP status code and an error response on failure.
func (c *Context) GetSession(r *http.Request, vars util.Vars) (int, []byte, error) {
	body, err := json.Marshal(SessionsResponse{
		Session: auth.Session(r),
	})
	return http.StatusOK, body, err
}

// PostSession is a util.JSONAPIFunc which creates a new Session and returns HTTP 200
// and a JSON session object on success, or a non-200 HTTP status code and an
// error response on failure.
func (c *Context) PostSession(r *http.Request, vars util.Vars) (int, []byte, error) {
	// Retrieve authenticated user
	user := auth.User(r)

	// Generate a new session for the user
	session, err := user.NewSession(time.Now().Add(auth.SessionDuration))
	if err != nil {
		return util.JSONAPIErr(err)
	}

	// Store session for later use
	if err := c.db.InsertSession(session); err != nil {
		return util.JSONAPIErr(err)
	}

	// Wrap in response and return
	body, err := json.Marshal(SessionsResponse{
		Session: session,
	})
	return http.StatusOK, body, err
}

// DeleteSession is a util.JSONAPIFunc which deletes an existing Session and returns
// HTTP 204 on success, or a non-200 HTTP status code and an error response on failure.
func (c *Context) DeleteSession(r *http.Request, vars util.Vars) (int, []byte, error) {
	// Retrieve authenticated session
	session := auth.Session(r)

	// Delete session now
	if err := c.db.DeleteSession(session); err != nil {
		return util.JSONAPIErr(err)
	}

	return http.StatusNoContent, nil, nil
}
