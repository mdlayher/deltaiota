package v0

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/mdlayher/deltaiota/api/auth"
	"github.com/mdlayher/deltaiota/api/util"
	"github.com/mdlayher/deltaiota/data"
	"github.com/mdlayher/deltaiota/data/models"
	"github.com/mdlayher/deltaiota/ditest"
)

// TestListNotificationsForUserNoNotifications verifies that ListNotificationsForUser
// returns no notifications when no notifications exist for a user in the database.
func TestListNotificationsForUserNoNotifications(t *testing.T) {
	// Invoke tests with temporary database
	err := ditest.WithTemporaryDB(func(db *data.DB) {
		// Build context
		c := &Context{
			db: db,
		}

		// Generate mock users
		user := ditest.MockUser()
		if err := c.db.InsertUser(user); err != nil {
			t.Error(err)
			return
		}

		user2 := ditest.MockUser()
		if err := c.db.InsertUser(user2); err != nil {
			t.Error(err)
			return
		}

		// Add notification for second user, which should not
		// appear in results
		if err := c.db.InsertNotification(&models.Notification{
			UserID: user2.ID,
		}); err != nil {
			t.Error(err)
			return
		}

		// Generate HTTP request
		r, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Error(err)
			return
		}

		// Store mock-authenticated user
		auth.SetUser(r, user)

		// Fetch list of current notifications for user
		code, body, err := c.ListNotificationsForUser(r, util.Vars{})
		if err != nil {
			t.Error(err)
			return
		}

		// Ensure proper HTTP status code
		if code != http.StatusOK {
			t.Errorf("unexpected code: %v != %v", code, http.StatusOK)
			return
		}

		// Unmarshal response body
		var res NotificationsResponse
		if err := json.Unmarshal(body, &res); err != nil {
			t.Error(err)
			return
		}

		// Verify length of notifications slice
		if len(res.Notifications) != 0 {
			t.Errorf("notifications slice not empty: %v", res.Notifications)
			return
		}
	})

	// Check for errors from database setup/cleanup
	if err != nil {
		t.Fatal("ditest.WithTemporaryDB:", err)
	}
}

// TestListNotificationsForUserManyNotifications verifies that ListNotificationsForUser
// returns many notifications when many notifications for the user exist in the database.
func TestListNotificationsManyNotifications(t *testing.T) {
	// Invoke tests with temporary database
	err := ditest.WithTemporaryDB(func(db *data.DB) {
		// Build context
		c := &Context{
			db: db,
		}

		// Generate mock users
		user := ditest.MockUser()
		if err := c.db.InsertUser(user); err != nil {
			t.Error(err)
			return
		}

		user2 := ditest.MockUser()
		if err := c.db.InsertUser(user2); err != nil {
			t.Error(err)
			return
		}

		// Add notification for second user, which should not
		// appear in results
		if err := c.db.InsertNotification(&models.Notification{
			UserID: user2.ID,
		}); err != nil {
			t.Error(err)
			return
		}

		// Generate HTTP request
		r, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Error(err)
			return
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
				t.Error(err)
				return
			}

			notifications[i] = notification
		}

		// Fetch list of current notifications for user
		code, body, err := c.ListNotificationsForUser(r, util.Vars{})
		if err != nil {
			t.Error(err)
			return
		}

		// Ensure proper HTTP status code
		if code != http.StatusOK {
			t.Errorf("unexpected code: %v != %v", code, http.StatusOK)
			return
		}

		// Unmarshal response body
		var res NotificationsResponse
		if err := json.Unmarshal(body, &res); err != nil {
			t.Error(err)
			return
		}

		// Check length of response slice
		if len(notifications) != len(res.Notifications) {
			t.Errorf("unexpected Notifications slice length: %v != %v", len(notifications), len(res.Notifications))
			return
		}

		// Check if all generated notifications returned, verify
		// they belong to the same user
		for i := range res.Notifications {
			if res.Notifications[i].UserID != notifications[i].UserID {
				t.Errorf("unexpected Notification UserID: %v != %v", res.Notifications[i].UserID, notifications[i].UserID)
			}
		}
	})

	// Check for errors from database setup/cleanup
	if err != nil {
		t.Fatal("ditest.WithTemporaryDB:", err)
	}
}
