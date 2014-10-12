package diclient

import (
	"github.com/mdlayher/deltaiota/api/v0"
	"github.com/mdlayher/deltaiota/data/models"
)

// SessionsService provides access to the Sessions API.
type SessionsService struct {
	client *Client
}

// CreateSession attempts to generate a new Session for the API, using
// the input username and password.
func (u *SessionsService) CreateSession(username string, password string) (*models.Session, *Response, error) {
	// Create request for Sessions endpoint
	req, err := u.client.NewRequest("POST", "sessions", nil)
	req.SetBasicAuth(username, password)
	if err != nil {
		return nil, nil, err
	}

	// Perform request, attempt to unmarshal response into a
	// Sessions API response
	sessionRes := new(v0.SessionsResponse)
	res, err := u.client.Do(req, &sessionRes)
	if err != nil {
		return nil, res, err
	}

	// Return session from API
	return sessionRes.Session, res, nil
}

// DeleteSession attempts to destroy the current Session for the API.
func (u *SessionsService) DeleteSession() (*Response, error) {
	// Delete request for Sessions endpoint
	req, err := u.client.NewRequest("DELETE", "sessions", nil)
	if err != nil {
		return nil, err
	}

	// Perform request
	res, err := u.client.Do(req, nil)
	return res, err
}
