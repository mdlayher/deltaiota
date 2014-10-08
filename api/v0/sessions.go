package v0

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/mdlayher/deltaiota/api/util"
	"github.com/mdlayher/deltaiota/data/models"
)

// SessionsResponse is the output response for the Sessions API.
type SessionsResponse struct {
	Session *models.Session `json:"session"`
}

// PostSession is a util.JSONAPIFunc which creates a new Session and returns HTTP 200
// and a JSON session object on success, or a non-200 HTTP status code and an
// error response on failure.
func (c *Context) PostSession(r *http.Request, vars util.Vars) (int, []byte, error) {
	// Retrieve authenticated user
	user := util.SessionUser(r)

	// Generate a new session for the user
	session, err := user.NewSession(time.Now().Add(7 * 24 * time.Hour))
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	// Store session for later use
	if err := c.db.InsertSession(session); err != nil {
		return http.StatusInternalServerError, nil, err
	}

	// Wrap in response and return
	body, err := json.Marshal(SessionsResponse{
		Session: session,
	})
	return http.StatusOK, body, err
}
