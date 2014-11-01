package diclient

import (
	"github.com/mdlayher/deltaiota/api/v0"
	"github.com/mdlayher/deltaiota/data/models"
)

// SessionsService provides access to the Sessions API.
type SessionsService struct {
	client *Client
}

// Get retrieves the current authenticated Session for the API.
func (s *SessionsService) Get() (*models.Session, *Response, error) {
	sRes, res, err := s.request("GET", "sessions", nil)
	return sRes.Session, res, err
}

// Create attempts to generate a new Session for the API, using
// the input username and password.
func (s *SessionsService) Create(username string, password string) (*models.Session, *Response, error) {
	// Create request for Sessions endpoint
	req, err := s.client.NewRequest("POST", "sessions", nil)
	req.SetBasicAuth(username, password)
	if err != nil {
		return nil, nil, err
	}

	// Perform request, attempt to unmarshal response into a
	// Sessions API response
	sRes := new(v0.SessionsResponse)
	res, err := s.client.Do(req, &sRes)
	if err != nil {
		return nil, res, err
	}

	// Return session from API
	return sRes.Session, res, nil
}

// Delete attempts to destroy the current Session for the API.
func (s *SessionsService) Delete() (*Response, error) {
	_, res, err := s.request("DELETE", "sessions", nil)
	return res, err
}

// request generates and performs a HTTP request to the Sessions API,
// with the exception of new session creation, due to a different authentication
// mechanism.
func (s *SessionsService) request(method string, endpoint string, body interface{}) (*v0.SessionsResponse, *Response, error) {
	// Create request for Sessions endpoint
	req, err := s.client.NewRequest(method, endpoint, body)
	if err != nil {
		return nil, nil, err
	}

	// Perform request, attempt to unmarshal response into a
	// Sessions API response
	sRes := new(v0.SessionsResponse)
	res, err := s.client.Do(req, &sRes)
	if err != nil {
		return nil, res, err
	}

	return sRes, res, nil
}
