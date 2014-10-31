package diclient

import (
	"github.com/mdlayher/deltaiota/api/v0"
)

// StatusService provides access to the Status API.
type StatusService struct {
	client *Client
}

// GetStatus returns the current API server status.
func (u *StatusService) GetStatus() (*v0.StatusResponse, *Response, error) {
	// Get request for Status endpoint
	req, err := u.client.NewRequest("GET", "status", nil)
	if err != nil {
		return nil, nil, err
	}

	// Perform request, attempt to unmarshal response into a
	// Status API response
	statusRes := new(v0.StatusResponse)
	res, err := u.client.Do(req, &statusRes)
	if err != nil {
		return nil, res, err
	}

	return statusRes, res, nil
}
