package diclient

import (
	"github.com/mdlayher/deltaiota/api/v0"
	"github.com/mdlayher/deltaiota/data/models"
)

// SessionsService provides access to the Sessions API.
type SessionsService struct {
	client *Client
}

// GetSession retrieves the current authenticated Session for the API.
func (s *SessionsService) GetSession() (*models.Session, *Response, error) {
	sessionsRes, res, err := s.sessionsRequest("GET", "sessions", nil)
	return sessionsRes.Session, res, err
}

// CreateSession attempts to generate a new Session for the API, using
// the input username and password.
func (s *SessionsService) CreateSession(username string, password string) (*models.Session, *Response, error) {
	// Create request for Sessions endpoint
	req, err := s.client.NewRequest("POST", "sessions", nil)
	req.SetBasicAuth(username, password)
	if err != nil {
		return nil, nil, err
	}

	// Perform request, attempt to unmarshal response into a
	// Sessions API response
	sessionRes := new(v0.SessionsResponse)
	res, err := s.client.Do(req, &sessionRes)
	if err != nil {
		return nil, res, err
	}

	// Return session from API
	return sessionRes.Session, res, nil
}

// DeleteSession attempts to destroy the current Session for the API.
func (s *SessionsService) DeleteSession() (*Response, error) {
	_, res, err := s.sessionsRequest("DELETE", "sessions", nil)
	return res, err
}

// sessionsRequest generates and performs a HTTP request to the Sessions API,
// with the exception of new session creation, due to a different authentication
// mechanism.
func (s *SessionsService) sessionsRequest(method string, endpoint string, body interface{}) (*v0.SessionsResponse, *Response, error) {
	// Create request for Sessions endpoint
	req, err := s.client.NewRequest(method, endpoint, body)
	if err != nil {
		return nil, nil, err
	}

	// Perform request, attempt to unmarshal response into a
	// Sessions API response
	sessionRes := new(v0.SessionsResponse)
	res, err := s.client.Do(req, &sessionRes)
	if err != nil {
		return nil, res, err
	}

	return sessionRes, res, nil
}
