package v0

import (
	"log"
	"net/http"
	"strconv"

	"github.com/mdlayher/deltaiota/api/util"
	"github.com/mdlayher/deltaiota/data/models"

	"github.com/gorilla/mux"
)

// UsersResponse is the output response for the Users API
type UsersResponse struct {
	Users []*models.User `json:"users"`
}

// ListUsers is a util.JSONAPIFunc which returns HTTP 200 and a JSONable list of users
// on success, or a non-200 HTTP status code and an error response on failure.
func (c *context) ListUsers(r *http.Request) (int, util.JSONable) {
	// Fetch a list of all users from the database
	users, err := c.db.FetchAllUsers()
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, nil
	}

	// Wrap in response and return
	return http.StatusOK, UsersResponse{
		Users: users,
	}
}

// GetUser is a util.JSONAPIFunc which returns HTTP 200 and a JSONable user object
// on success, or a non-200 HTTP status code and an error response on failure.
func (c *context) GetUser(r *http.Request) (int, util.JSONable) {
	// Fetch input user ID
	strID, ok := mux.Vars(r)["id"]
	if !ok {
		return http.StatusBadRequest, nil
	}

	// Convert string to integer
	id, err := strconv.ParseInt(strID, 10, 64)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, nil
	}

	// Select single user by ID from the database
	user, err := c.db.SelectUserByID(id)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, nil
	}

	// Wrap in response and return
	return http.StatusOK, UsersResponse{
		Users: []*models.User{user},
	}
}
