package v0

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mdlayher/deltaiota/data/models"
)

// UsersResponse is the output response for the Users API
type UsersResponse struct {
	Users []models.User `json:"users"`
}

// ListUsers is a util.JSONAPIFunc which returns HTTP 200 and a JSON list of users
// on success, or a non-200 HTTP status code and an error response on failure.
func (c *context) ListUsers(r *http.Request) (int, []byte) {
	// Fetch a list of all users from the database
	users, err := c.db.FetchAllUsers()
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, nil
	}

	// Marshal list of users to JSON
	body, err := json.Marshal(users)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, nil
	}

	return http.StatusOK, body
}
