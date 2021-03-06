package v0

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/mdlayher/deltaiota/api/util"
	"github.com/mdlayher/deltaiota/data/models"
	"github.com/mdlayher/deltaiota/ditest"
)

// TestUsersAPI verifies that UsersAPI correctly routes requests to
// other Users API handlers, using the input HTTP request.
func TestUsersAPI(t *testing.T) {
	withContext(t, func(c *Context) error {
		var tests = []struct {
			method string
			vars   util.Vars
			code   int
		}{
			// ListUsers
			{"GET", util.Vars{}, http.StatusOK},
			{"HEAD", util.Vars{}, http.StatusOK},
			// GetUser
			{"GET", util.Vars{"id": "1"}, http.StatusNotFound},
			{"HEAD", util.Vars{"id": "1"}, http.StatusNotFound},
			// PostUser
			{"POST", util.Vars{}, http.StatusBadRequest},
			// PutUser
			{"PUT", util.Vars{}, http.StatusBadRequest},
			// DeleteUser
			{"DELETE", util.Vars{}, http.StatusBadRequest},
			// Unknown method
			{"CAT", util.Vars{}, http.StatusMethodNotAllowed},
		}

		for _, test := range tests {
			// Generate HTTP request
			r, err := http.NewRequest(test.method, "/", nil)
			if err != nil {
				return err
			}

			// Delegate to appropriate handler
			code, _, err := c.UsersAPI(r, test.vars)
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

// TestListUsersNoUsers verifies that ListUsers returns no users when no
// users exist in the database.
func TestListUsersNoUsers(t *testing.T) {
	withContext(t, func(c *Context) error {
		// Fetch list of current users
		code, body, err := c.ListUsers(nil, util.Vars{})
		if err != nil {
			return err
		}

		// Ensure proper HTTP status code
		if code != http.StatusOK {
			return fmt.Errorf("unexpected code: %v != %v", code, http.StatusOK)
		}

		// Unmarshal response body
		var res UsersResponse
		if err := json.Unmarshal(body, &res); err != nil {
			return err
		}

		// Verify length of users slice
		if len(res.Users) != 0 {
			return fmt.Errorf("users slice not empty: %v", res.Users)
		}

		return nil
	})
}

// TestListUserManyUsers verifies that ListUsers returns many users when many users
// users exist in the database.
func TestListUsersManyUsers(t *testing.T) {
	withContext(t, func(c *Context) error {
		// Generate and save a mock users in the database
		users := make([]*models.User, 100)
		for i := range users {
			user := ditest.MockUser()
			if err := c.db.InsertUser(user); err != nil {
				return err
			}

			users[i] = user
		}

		// Fetch list of current users
		code, body, err := c.ListUsers(nil, util.Vars{})
		if err != nil {
			return err
		}

		// Ensure proper HTTP status code
		if code != http.StatusOK {
			return fmt.Errorf("unexpected code: %v != %v", code, http.StatusOK)
		}

		// Unmarshal response body
		var res UsersResponse
		if err := json.Unmarshal(body, &res); err != nil {
			return err
		}

		// Check length of response slice
		if len(users) != len(res.Users) {
			return fmt.Errorf("unexpected Users slice length: %v != %v", len(users), len(res.Users))
		}

		// Check if all generated users returned
		for i := range res.Users {
			if res.Users[i].Username != users[i].Username {
				return fmt.Errorf("unexpected User username: %v != %v", res.Users[i].Username, users[i].Username)
			}

			if res.Users[i].FirstName != users[i].FirstName {
				return fmt.Errorf("unexpected User first name: %v != %v", res.Users[i].FirstName, users[i].FirstName)
			}

			if res.Users[i].LastName != users[i].LastName {
				return fmt.Errorf("unexpected User last name: %v != %v", res.Users[i].LastName, users[i].LastName)
			}

			if res.Users[i].Email != users[i].Email {
				return fmt.Errorf("unexpected User email: %v != %v", res.Users[i].Email, users[i].Email)
			}
		}

		return nil
	})
}

// TestGetUser verifies that GetUser returns the appropriate HTTP status
// code, body, and any errors which occur.
func TestGetUser(t *testing.T) {
	withContextUser(t, func(c *Context, user *models.User) error {
		// Table of tests to iterate
		var tests = []struct {
			id         string
			code       int
			errMessage string
		}{
			// Empty ID
			{"", http.StatusBadRequest, userMissingID},
			// Bad ID
			{"-1", http.StatusBadRequest, userInvalidID},
			{"test", http.StatusBadRequest, userInvalidID},
			// ID not found
			{"2", http.StatusNotFound, userNotFound},
			// Existing user
			{"1", http.StatusOK, ""},
		}

		// Iterate and run tests
		for _, test := range tests {
			// Generate HTTP request
			r, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				return err
			}

			// Set path variables, unless ID is missing
			vars := util.Vars{}
			if test.id != "" {
				vars["id"] = test.id
			}

			// Invoke GetUser with HTTP request, manually injecting
			// path variables from test
			code, body, err := c.GetUser(r, vars)
			if err != nil {
				return err
			}

			// Ensure proper HTTP status code
			if code != test.code {
				return fmt.Errorf("unexpected code: %v != %v", code, test.code)
			}

			// If code is not HTTP OK, check error response
			if code != http.StatusOK {
				// Unmarshal error JSON into struct
				var errRes util.ErrorResponse
				if err := json.Unmarshal(body, &errRes); err != nil {
					return err
				}

				// Verify error code and message
				if errRes.Error.Code != test.code {
					return fmt.Errorf("unexpected error code: %v != %v", errRes.Error.Code, test.code)
				}
				if errRes.Error.Message != test.errMessage {
					return fmt.Errorf("unexpected error message: %v != %v", errRes.Error.Message, test.errMessage)
				}

				// Skip to next test
				continue
			}

			// Unmarshal response body
			var res UsersResponse
			if err := json.Unmarshal(body, &res); err != nil {
				return err
			}

			// Verify length of users slice
			if len(res.Users) != 1 {
				return fmt.Errorf("unexpected number of users returned: %v", res.Users)
			}

			// Verify user is the same as the mock we created earlier
			if res.Users[0].Username != user.Username {
				return fmt.Errorf("unexpected User username: %v != %v", res.Users[0].Username, user.Username)
			}

			if res.Users[0].FirstName != user.FirstName {
				return fmt.Errorf("unexpected User first name: %v != %v", res.Users[0].FirstName, user.FirstName)
			}

			if res.Users[0].LastName != user.LastName {
				return fmt.Errorf("unexpected User last name: %v != %v", res.Users[0].LastName, user.LastName)
			}

			if res.Users[0].Email != user.Email {
				return fmt.Errorf("unexpected User email: %v != %v", res.Users[0].Email, user.Email)
			}
		}

		return nil
	})
}

// TestPostUser verifies that PostUser returns the appropriate HTTP status
// code, body, and any errors which occur.
func TestPostUser(t *testing.T) {
	withContext(t, func(c *Context) error {
		// JSON used to generate a temporary user
		mockUserJSON := []byte(`{"id": 1, "password":"test","firstName":"test","lastName":"test","username":"test","email":"test@test.com"}`)

		// Unmarshal into mock user
		user := new(models.User)
		if err := json.Unmarshal(mockUserJSON, user); err != nil {
			return err
		}

		// Table of tests to iterate
		var tests = []struct {
			code       int
			errMessage string
			body       []byte
		}{
			// Empty body
			{http.StatusBadRequest, userJSONSyntax, nil},
			// Bad JSON
			{http.StatusBadRequest, userJSONSyntax, []byte(`{`)},
			// No password key
			{http.StatusBadRequest, "empty field: password", []byte(`{}`)},
			// Missing password
			{http.StatusBadRequest, "empty field: password", []byte(`{"password":""}`)},
			// Missing username
			{http.StatusBadRequest, "empty field: username", []byte(`{"password":"test","firstName":"test","lastName":"test","email":"test@test.com"}`)},
			// Missing first name
			{http.StatusBadRequest, "empty field: firstName", []byte(`{"password":"test","lastName":"test","username":"test","email":"test@test.com"}`)},
			// Missing last name
			{http.StatusBadRequest, "empty field: lastName", []byte(`{"password":"test","firstName":"test","username":"test","email":"test@test.com"}`)},
			// Missing email
			{http.StatusBadRequest, "empty field: email", []byte(`{"password":"test","firstName":"test","lastName":"test","username":"test"}`)},
			// Invalid email
			{http.StatusBadRequest, "invalid field: email (could not parse valid email address)", []byte(`{"password":"test","firstName":"test","lastName":"test","username":"test","email":"test"}`)},
			// Valid request
			{http.StatusCreated, "", mockUserJSON},
			// Duplicate username
			{http.StatusConflict, userConflict, []byte(`{"password":"test2","firstName":"test2","lastName":"test2","username":"test","email":"test2@test.com"}`)},
			// Duplicate email
			{http.StatusConflict, userConflict, []byte(`{"password":"test2","firstName":"test2","lastName":"test2","username":"test2","email":"test@test.com"}`)},
		}

		// Iterate and run tests
		for _, test := range tests {
			// Generate HTTP request
			r, err := http.NewRequest("POST", "/", bytes.NewReader(test.body))
			if err != nil {
				return err
			}

			// Invoke PostUser with HTTP request
			code, body, err := c.PostUser(r, util.Vars{})
			if err != nil {
				return err
			}

			// Ensure proper HTTP status code
			if code != test.code {
				return fmt.Errorf("unexpected code: %v != %v", code, test.code)
			}

			// If code is in HTTP 400 or above, check error response
			if code >= http.StatusBadRequest {
				// Unmarshal error JSON into struct
				var errRes util.ErrorResponse
				if err := json.Unmarshal(body, &errRes); err != nil {
					return err
				}

				// Verify error code and message
				if errRes.Error.Code != test.code {
					return fmt.Errorf("unexpected error code: %v != %v", errRes.Error.Code, test.code)
				}
				if errRes.Error.Message != test.errMessage {
					return fmt.Errorf("unexpected error message: %v != %v", errRes.Error.Message, test.errMessage)
				}

				// Skip to next test
				continue
			}

			// Unmarshal response body
			var res UsersResponse
			if err := json.Unmarshal(body, &res); err != nil {
				return err
			}

			// Verify length of users slice
			if len(res.Users) != 1 {
				return fmt.Errorf("unexpected number of users returned: %v", res.Users)
			}

			// Strip password for comparison
			user.Password = ""

			// Verify user is the same as the mock we created earlier
			if !reflect.DeepEqual(user, res.Users[0]) {
				return fmt.Errorf("unexpected user: %v != %v", user, res.Users[0])
			}
		}

		return nil
	})
}

// TestPutUser verifies that PutUser returns the appropriate HTTP status
// code, body, and any errors which occur.
func TestPutUser(t *testing.T) {
	withContext(t, func(c *Context) error {
		// JSON used to generate a temporary user
		mockUserJSON := []byte(`{"id": 1, "password":"test","firstName":"test","lastName":"test","username":"test","email":"test@test.com"}`)

		// Unmarshal into mock user
		user := new(models.User)
		if err := json.Unmarshal(mockUserJSON, user); err != nil {
			return err
		}

		// Save user in database, to be updated later
		if err := c.db.InsertUser(user); err != nil {
			return err
		}

		// Save user in database for conflicting username and email
		conflictUser := &models.User{
			Username: "conflict",
			Email:    "conflict@conflict.com",
		}
		conflictUser.SetPassword("conflict")
		if err := c.db.InsertUser(conflictUser); err != nil {
			return err
		}

		// Table of tests to iterate
		var tests = []struct {
			id         string
			code       int
			errMessage string
			body       []byte
		}{
			// Empty ID
			{"", http.StatusBadRequest, userMissingID, nil},
			// Bad ID
			{"-1", http.StatusBadRequest, userInvalidID, nil},
			{"test", http.StatusBadRequest, userInvalidID, nil},
			// ID not found
			{"3", http.StatusNotFound, userNotFound, nil},
			// Empty body
			{"1", http.StatusBadRequest, userJSONSyntax, nil},
			// Bad JSON
			{"1", http.StatusBadRequest, userJSONSyntax, []byte(`{`)},
			// No password key
			{"1", http.StatusBadRequest, "empty field: password", []byte(`{}`)},
			// Missing password
			{"1", http.StatusBadRequest, "empty field: password", []byte(`{"password":""}`)},
			// Missing username
			{"1", http.StatusBadRequest, "empty field: username", []byte(`{"password":"test","firstName":"test","lastName":"test","email":"test@test.com"}`)},
			// Missing first name
			{"1", http.StatusBadRequest, "empty field: firstName", []byte(`{"password":"test","lastName":"test","username":"test","email":"test@test.com"}`)},
			// Missing last name
			{"1", http.StatusBadRequest, "empty field: lastName", []byte(`{"password":"test","firstName":"test","username":"test","email":"test@test.com"}`)},
			// Missing email
			{"1", http.StatusBadRequest, "empty field: email", []byte(`{"password":"test","firstName":"test","lastName":"test","username":"test"}`)},
			// Invalid email
			{"1", http.StatusBadRequest, "invalid field: email (could not parse valid email address)", []byte(`{"password":"test","firstName":"test","lastName":"test","username":"test","email":"test"}`)},
			// Valid request
			{"1", http.StatusOK, "", mockUserJSON},
			// Duplicate username
			{"1", http.StatusConflict, userConflict, []byte(`{"password":"test","firstName":"test","lastName":"test","username":"conflict","email":"test@test.com"}`)},
			// Duplicate email
			{"1", http.StatusConflict, userConflict, []byte(`{"password":"test","firstName":"test","lastName":"test","username":"test","email":"conflict@conflict.com"}`)},
		}

		// Iterate and run tests
		for _, test := range tests {
			// Generate HTTP request
			r, err := http.NewRequest("PUT", "/", bytes.NewReader(test.body))
			if err != nil {
				return err
			}

			// Set path variables, unless ID is missing
			vars := util.Vars{}
			if test.id != "" {
				vars["id"] = test.id
			}

			// Invoke PutUser with HTTP request, manually injecting
			// path variables from test
			code, body, err := c.PutUser(r, vars)
			if err != nil {
				return err
			}

			// Ensure proper HTTP status code
			if code != test.code {
				return fmt.Errorf("unexpected code: %v != %v", code, test.code)
			}

			// If code is in HTTP 400 or above, check error response
			if code >= http.StatusBadRequest {
				// Unmarshal error JSON into struct
				var errRes util.ErrorResponse
				if err := json.Unmarshal(body, &errRes); err != nil {
					return err
				}

				// Verify error code and message
				if errRes.Error.Code != test.code {
					return fmt.Errorf("unexpected error code: %v != %v", errRes.Error.Code, test.code)
				}
				if errRes.Error.Message != test.errMessage {
					return fmt.Errorf("unexpected error message: %v != %v", errRes.Error.Message, test.errMessage)
				}

				// Skip to next test
				continue
			}

			// Unmarshal response body
			var res UsersResponse
			if err := json.Unmarshal(body, &res); err != nil {
				return err
			}

			// Verify length of users slice
			if len(res.Users) != 1 {
				return fmt.Errorf("unexpected number of users returned: %v", res.Users)
			}

			// Strip password for comparison
			user.Password = ""

			// Verify user is the same as the mock we created earlier
			if !reflect.DeepEqual(user, res.Users[0]) {
				return fmt.Errorf("unexpected user: %v != %v", user, res.Users[0])
			}
		}

		return nil
	})
}

// TestDeleteUser verifies that DeleteUser returns the appropriate HTTP status
// code, body, and any errors which occur.
func TestDeleteUser(t *testing.T) {
	withContextUser(t, func(c *Context, user *models.User) error {
		// Table of tests to iterate
		var tests = []struct {
			id         string
			code       int
			errMessage string
		}{
			// Empty ID
			{"", http.StatusBadRequest, userMissingID},
			// Bad ID
			{"-1", http.StatusBadRequest, userInvalidID},
			{"test", http.StatusBadRequest, userInvalidID},
			// ID not found
			{"2", http.StatusNotFound, userNotFound},
			// Existing user
			{"1", http.StatusNoContent, ""},
		}

		// Iterate and run tests
		for _, test := range tests {
			// Generate HTTP request
			r, err := http.NewRequest("DELETE", "/", nil)
			if err != nil {
				return err
			}

			// Set path variables, unless ID is missing
			vars := util.Vars{}
			if test.id != "" {
				vars["id"] = test.id
			}

			// Invoke DeleteUser with HTTP request, manually injecting
			// path variables from test
			code, body, err := c.DeleteUser(r, vars)
			if err != nil {
				return err
			}

			// Ensure proper HTTP status code
			if code != test.code {
				return fmt.Errorf("unexpected code: %v != %v", code, test.code)
			}

			// If code is HTTP No Content, ensure user was deleted
			if code == http.StatusNoContent {
				if _, err := c.db.SelectUserByID(user.ID); err != sql.ErrNoRows {
					return fmt.Errorf("called DeleteUser, but user still exists: %v", user)
				}

				continue
			}

			// Unmarshal error JSON into struct
			var errRes util.ErrorResponse
			if err := json.Unmarshal(body, &errRes); err != nil {
				return err
			}

			// Verify error code and message
			if errRes.Error.Code != test.code {
				return fmt.Errorf("unexpected error code: %v != %v", errRes.Error.Code, test.code)
			}
			if errRes.Error.Message != test.errMessage {
				return fmt.Errorf("unexpected error message: %v != %v", errRes.Error.Message, test.errMessage)
			}
		}

		return nil
	})
}
