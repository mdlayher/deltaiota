package v0

import (
	"bytes"
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/mdlayher/deltaiota/api/util"
	"github.com/mdlayher/deltaiota/data"
	"github.com/mdlayher/deltaiota/data/models"
	"github.com/mdlayher/deltaiota/ditest"
)

// TestListUsers verifies that ListUsers returns the appropriate HTTP status
// code, body, and any errors which occur.
func TestListUsers(t *testing.T) {
	// Invoke tests with temporary database
	err := ditest.WithTemporaryDB(func(db *data.DB) {
		// Build context
		c := &context{
			db: db,
		}

		// Loop twice to check for empty users, add a user, then verify
		// that the same user is returned
		var lastUser *models.User
		for i := 0; i < 1; i++ {
			// Fetch list of current users
			code, body, err := c.ListUsers(nil, util.Vars{})
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
			var res UsersResponse
			if err := json.Unmarshal(body, &res); err != nil {
				t.Error(err)
				return
			}

			// Verify length of users slice
			if len(res.Users) != i {
				t.Errorf("unexpected number of users returned: %v", res.Users)
				return
			}

			// Check if generated user returned and break loop on second run
			if i == 1 {
				if res.Users[0] != lastUser {
					t.Errorf("unexpected User: %v != %v", res.Users[0], lastUser)
					return
				}

				break
			}

			// Generate and save a mock user in the database
			lastUser = ditest.MockUser()
			if err := c.db.InsertUser(lastUser); err != nil {
				t.Error(err)
				return
			}
		}
	})

	// Check for errors from database setup/cleanup
	if err != nil {
		t.Fatalf("ditest.WithTemporaryDB:", err)
	}
}

// TestGetUser verifies that GetUser returns the appropriate HTTP status
// code, body, and any errors which occur.
func TestGetUser(t *testing.T) {
	// Invoke tests with temporary database
	err := ditest.WithTemporaryDB(func(db *data.DB) {
		// Build context
		c := &context{
			db: db,
		}

		// Generate and store mock user
		user := ditest.MockUser()
		if err := c.db.InsertUser(user); err != nil {
			t.Error(err)
			return
		}

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
				t.Error(err)
				return
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
				t.Error(err)
				return
			}

			// Ensure proper HTTP status code
			if code != test.code {
				t.Errorf("unexpected code: %v != %v", code, test.code)
				return
			}

			// If code is not HTTP OK, check error response
			if code != http.StatusOK {
				// Unmarshal error JSON into struct
				var errRes util.ErrorResponse
				if err := json.Unmarshal(body, &errRes); err != nil {
					t.Error(err)
					return
				}

				// Verify error code and message
				if errRes.Error.Code != test.code {
					t.Errorf("unexpected error code: %v != %v", errRes.Error.Code, test.code)
					return
				}
				if errRes.Error.Message != test.errMessage {
					t.Errorf("unexpected error message: %v != %v", errRes.Error.Message, test.errMessage)
					return
				}

				// Skip to next test
				continue
			}

			// Unmarshal response body
			var res UsersResponse
			if err := json.Unmarshal(body, &res); err != nil {
				t.Error(err)
				return
			}

			// Verify length of users slice
			if len(res.Users) != 1 {
				t.Errorf("unexpected number of users returned: %v", res.Users)
				return
			}

			// Verify user is the same as the mock we created earlier
			if !reflect.DeepEqual(user, res.Users[0]) {
				t.Errorf("unexpected user: %v != %v", user, res.Users[0])
			}
		}
	})

	// Check for errors from database setup/cleanup
	if err != nil {
		t.Fatalf("ditest.WithTemporaryDB:", err)
	}
}

// TestPostUser verifies that PostUser returns the appropriate HTTP status
// code, body, and any errors which occur.
func TestPostUser(t *testing.T) {
	// Invoke tests with temporary database
	err := ditest.WithTemporaryDB(func(db *data.DB) {
		// Build context
		c := &context{
			db: db,
		}

		// JSON used to generate a temporary user
		mockUserJSON := []byte(`{"id": 1, "password":"test","firstName":"test","lastName":"test","username":"test","email":"test@test.com"}`)

		// Unmarshal into mock user
		user := new(models.User)
		if err := json.Unmarshal(mockUserJSON, user); err != nil {
			t.Error(err)
			return
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
				t.Error(err)
				return
			}

			// Invoke PostUser with HTTP request
			code, body, err := c.PostUser(r, util.Vars{})
			if err != nil {
				t.Error(err)
				return
			}

			// Ensure proper HTTP status code
			if code != test.code {
				t.Errorf("unexpected code: %v != %v", code, test.code)
				return
			}

			// If code is in HTTP 400 or above, check error response
			if code >= http.StatusBadRequest {
				// Unmarshal error JSON into struct
				var errRes util.ErrorResponse
				if err := json.Unmarshal(body, &errRes); err != nil {
					t.Error(err)
					return
				}

				// Verify error code and message
				if errRes.Error.Code != test.code {
					t.Errorf("unexpected error code: %v != %v", errRes.Error.Code, test.code)
					return
				}
				if errRes.Error.Message != test.errMessage {
					t.Errorf("unexpected error message: %v != %v", errRes.Error.Message, test.errMessage)
					return
				}

				// Skip to next test
				continue
			}

			// Unmarshal response body
			var res UsersResponse
			if err := json.Unmarshal(body, &res); err != nil {
				t.Error(err)
				return
			}

			// Verify length of users slice
			if len(res.Users) != 1 {
				t.Errorf("unexpected number of users returned: %v", res.Users)
				return
			}

			// Verify user is the same as the mock we created earlier
			if !reflect.DeepEqual(user, res.Users[0]) {
				t.Errorf("unexpected user: %v != %v", user, res.Users[0])
			}
		}
	})

	// Check for errors from database setup/cleanup
	if err != nil {
		t.Fatalf("ditest.WithTemporaryDB:", err)
	}
}
