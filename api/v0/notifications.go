package v0

import (
	"encoding/json"
	"net/http"

	"github.com/mdlayher/deltaiota/api/auth"
	"github.com/mdlayher/deltaiota/api/util"
	"github.com/mdlayher/deltaiota/data/models"
)

// NotificationsResponse is the output response for the Notifications API
type NotificationsResponse struct {
	Notifications []*models.Notification `json:"notifications"`
}

// ListNotificationsForUser is a util.JSONAPIFunc which returns HTTP 200 and a
// JSON list of notifications for the authenticated user on success, or a non-200
// HTTP status code and an error response on failure.
func (c *Context) ListNotificationsForUser(r *http.Request, vars util.Vars) (int, []byte, error) {
	// Fetch a list of notifications for this user from the database
	notifications, err := c.db.SelectNotificationsByUserID(auth.User(r).ID)
	if err != nil {
		return util.JSONAPIErr(err)
	}

	// Wrap in response
	body, err := json.Marshal(NotificationsResponse{
		Notifications: notifications,
	})
	return http.StatusOK, body, err
}
