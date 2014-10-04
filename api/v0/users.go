package v0

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/mdlayher/deltaiota/api/util"
	"github.com/mdlayher/deltaiota/data/models"

	"github.com/gorilla/mux"
)

// JSON Users API, human-readable client error responses.
const (
	userInvalidID = "invalid user ID"
	userMissingID = "missing user ID"
	userNotFound  = "user not found"
)

// JSON Users API, map of client errors to response codes.
var usersCode = map[string]int{
	userInvalidID: http.StatusBadRequest,
	userMissingID: http.StatusBadRequest,
	userNotFound:  http.StatusNotFound,
}

// Generated JSON responses for various client-facing errors.
var usersJSON = map[string][]byte{}

// init initializes the stored JSON responses for client-facing errors.
func init() {
	// Iterate all error strings and code integers
	for k, v := range usersCode {
		// Generate error response with appropriate string and code
		body, err := json.Marshal(util.ErrRes(v, k))
		if err != nil {
			panic(err)
		}

		// Store for later use
		usersJSON[k] = body
	}
}

// UsersResponse is the output response for the Users API
type UsersResponse struct {
	Users []*models.User `json:"users"`
}

// ListUsers is a util.JSONAPIFunc which returns HTTP 200 and a JSON list of users
// on success, or a non-200 HTTP status code and an error response on failure.
func (c *context) ListUsers(r *http.Request) (int, []byte, error) {
	// Fetch a list of all users from the database
	users, err := c.db.SelectAllUsers()
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	// Wrap in response and return
	body, err := json.Marshal(UsersResponse{
		Users: users,
	})
	return http.StatusOK, body, err
}

// GetUser is a util.JSONAPIFunc which returns HTTP 200 and a JSON user object
// on success, or a non-200 HTTP status code and an error response on failure.
func (c *context) GetUser(r *http.Request) (int, []byte, error) {
	// Fetch input user ID
	strID, ok := mux.Vars(r)["id"]
	if !ok {
		return usersCode[userMissingID], usersJSON[userMissingID], nil
	}

	// Convert string to integer
	id, err := strconv.ParseUint(strID, 10, 64)
	if err != nil {
		return usersCode[userInvalidID], usersJSON[userInvalidID], nil
	}

	// Select single user by ID from the database
	user, err := c.db.SelectUserByID(int64(id))
	if err != nil {
		// If no results found, return HTTP not found
		if err == sql.ErrNoRows {
			return usersCode[userNotFound], usersJSON[userNotFound], nil
		}

		return http.StatusInternalServerError, nil, err
	}

	// Wrap in response and return
	body, err := json.Marshal(UsersResponse{
		Users: []*models.User{user},
	})
	return http.StatusOK, body, err
}
