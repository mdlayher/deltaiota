package v0

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/mdlayher/deltaiota/api/auth"
	"github.com/mdlayher/deltaiota/api/util"
	"github.com/mdlayher/deltaiota/data/models"
	"github.com/mdlayher/deltaiota/ditest"
)

// TestNotificationsAPI verifies that NotificationsAPI correctly routes requests to
// other Notifications API handlers, using the input HTTP request.
func TestNotificationsAPI(t *testing.T) {
	withContextUser(t, func(c *Context, user *models.User) error {
		var tests = []struct {
			method string
			code   int
		}{
			// ListNotificationsForUser
			{"GET", http.StatusOK},
			{"HEAD", http.StatusOK},
			// Method not allowed
			{"POST", http.StatusMethodNotAllowed},
			{"PUT", http.StatusMethodNotAllowed},
			{"PATCH", http.StatusMethodNotAllowed},
			{"DELETE", http.StatusMethodNotAllowed},
			{"CAT", http.StatusMethodNotAllowed},
		}

		for _, test := range tests {
			// Generate HTTP request
			r, err := http.NewRequest(test.method, "/", nil)
			if err != nil {
				return err
			}

			// Store mock-authenticated user
			auth.SetUser(r, user)

			// Delegate to appropriate handler
			code, _, err := c.NotificationsAPI(r, util.Vars{})
			if err != nil {
				return err
			}

			// Ensure proper HTTP status code
			if code != test.code {
				return fmt.Errorf("unexpected code: %v != %v", code, test.code)
			}
		}

		return nil
	})
}

// TestListNotificationsForUserNoNotifications verifies that ListNotificationsForUser
// returns no notifications when no notifications exist for a user in the database.
func TestListNotificationsForUserNoNotifications(t *testing.T) {
	withContextUser(t, func(c *Context, user *models.User) error {
		// Generate another mock user
		user2 := ditest.MockUser()
		if err := c.db.InsertUser(user2); err != nil {
			return err
		}

		// Add notification for second user, which should not
		// appear in results
		if err := c.db.InsertNotification(&models.Notification{
			UserID: user2.ID,
		}); err != nil {
			return err
		}

		// Generate HTTP request
		r, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			return err
		}

		// Store mock-authenticated user
		auth.SetUser(r, user)

		// Fetch list of current notifications for user
		code, body, err := c.ListNotificationsForUser(r, util.Vars{})
		if err != nil {
			return err
		}

		// Ensure proper HTTP status code
		if code != http.StatusOK {
			return fmt.Errorf("unexpected code: %v != %v", code, http.StatusOK)
		}

		// Unmarshal response body
		var res NotificationsResponse
		if err := json.Unmarshal(body, &res); err != nil {
			return err
		}

		// Verify length of notifications slice
		if len(res.Notifications) != 0 {
			return fmt.Errorf("notifications slice not empty: %v", res.Notifications)
		}

		return nil
	})
}

// TestListNotificationsForUserManyNotifications verifies that ListNotificationsForUser
// returns many notifications when many notifications for the user exist in the database.
func TestListNotificationsManyNotifications(t *testing.T) {
	withContextUser(t, func(c *Context, user *models.User) error {
		// Generate another mock user
		user2 := ditest.MockUser()
		if err := c.db.InsertUser(user2); err != nil {
			return err
		}

		// Add notification for second user, which should not
		// appear in results
		if err := c.db.InsertNotification(&models.Notification{
			UserID: user2.ID,
		}); err != nil {
			return err
		}

		// Generate HTTP request
		r, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			return err
		}

		// Store mock-authenticated user
		auth.SetUser(r, user)

		// Generate and save a mock notifications in the database
		notifications := make([]*models.Notification, 100)
		for i := range notifications {
			notification := &models.Notification{
				UserID: user.ID,
			}
			if err := c.db.InsertNotification(notification); err != nil {
				return err
			}

			notifications[i] = notification
		}

		// Fetch list of current notifications for user
		code, body, err := c.ListNotificationsForUser(r, util.Vars{})
		if err != nil {
			return err
		}

		// Ensure proper HTTP status code
		if code != http.StatusOK {
			return fmt.Errorf("unexpected code: %v != %v", code, http.StatusOK)
		}

		// Unmarshal response body
		var res NotificationsResponse
		if err := json.Unmarshal(body, &res); err != nil {
			return err
		}

		// Check length of response slice
		if len(notifications) != len(res.Notifications) {
			return fmt.Errorf("unexpected Notifications slice length: %v != %v", len(notifications), len(res.Notifications))
		}

		// Check if all generated notifications returned, verify
		// they belong to the same user
		for i := range res.Notifications {
			if res.Notifications[i].UserID != notifications[i].UserID {
				return fmt.Errorf("unexpected Notification UserID: %v != %v", res.Notifications[i].UserID, notifications[i].UserID)
			}
		}

		return nil
	})
}
