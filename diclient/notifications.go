package diclient

import (
	"github.com/mdlayher/deltaiota/api/v0"
	"github.com/mdlayher/deltaiota/data/models"
)

// NotificationsService provides access to the Notifications API.
type NotificationsService struct {
	client *Client
}

// ListNotifications attempts to return a list of current notifications
// for the active user.
func (n *NotificationsService) ListNotifications() ([]*models.Notification, *Response, error) {
	notificationsRes, res, err := n.notificationsRequest("GET", "notifications", nil)

	// Check for empty notifications
	if notificationsRes == nil || notificationsRes.Notifications == nil {
		return nil, res, err
	}

	return notificationsRes.Notifications, res, err
}

// notificationsRequest generates and performs a HTTP request to the Notifications API.
func (n *NotificationsService) notificationsRequest(method string, endpoint string, body interface{}) (*v0.NotificationsResponse, *Response, error) {
	// Create request for Notifications endpoint
	req, err := n.client.NewRequest(method, endpoint, body)
	if err != nil {
		return nil, nil, err
	}

	// Perform request, attempt to unmarshal response into a
	// Notifications API response
	notificationsRes := new(v0.NotificationsResponse)
	res, err := n.client.Do(req, &notificationsRes)
	if err != nil {
		return nil, res, err
	}

	return notificationsRes, res, nil
}
