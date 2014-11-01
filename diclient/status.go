package diclient

import (
	"github.com/mdlayher/deltaiota/api/v0"
)

// StatusService provides access to the Status API.
type StatusService struct {
	client *Client
}

// Get returns the current API server status.
func (s *StatusService) Get() (*v0.Status, *Response, error) {
	// Get request for Status endpoint
	req, err := s.client.NewRequest("GET", "status", nil)
	if err != nil {
		return nil, nil, err
	}

	// Perform request, attempt to unmarshal response into a
	// Status API response
	sRes := new(v0.StatusResponse)
	res, err := s.client.Do(req, &sRes)
	if err != nil {
		return nil, res, err
	}

	return sRes.Status, res, nil
}
