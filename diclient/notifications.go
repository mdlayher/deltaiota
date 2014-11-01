package diclient

import (
	"github.com/mdlayher/deltaiota/api/v0"
	"github.com/mdlayher/deltaiota/data/models"
)

// NotificationsService provides access to the Notifications API.
type NotificationsService struct {
	client *Client
}

// List attempts to return a list of current notification for the active user.
func (n *NotificationsService) List() ([]*models.Notification, *Response, error) {
	nRes, res, err := n.request("GET", "notifications", nil)

	// Check for empty notifications
	if nRes == nil || nRes.Notifications == nil {
		return nil, res, err
	}

	return nRes.Notifications, res, err
}

// request generates and performs a HTTP request to the Notifications API.
func (n *NotificationsService) request(method string, endpoint string, body interface{}) (*v0.NotificationsResponse, *Response, error) {
	// Create request for Notifications endpoint
	req, err := n.client.NewRequest(method, endpoint, body)
	if err != nil {
		return nil, nil, err
	}

	// Perform request, attempt to unmarshal response into a
	// Notifications API response
	nRes := new(v0.NotificationsResponse)
	res, err := n.client.Do(req, &nRes)
	if err != nil {
		return nil, res, err
	}

	return nRes, res, nil
}
