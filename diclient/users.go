package diclient

import (
	"github.com/mdlayher/deltaiota/api/v0"
	"github.com/mdlayher/deltaiota/data/models"
)

// UsersService provides access to the Users API.
type UsersService struct {
	client *Client
}

// ListUsers returns a slice of all User objects from the API.
func (u *UsersService) ListUsers() ([]*models.User, *Response, error) {
	// Create request for Users endpoint
	req, err := u.client.NewRequest("GET", "users", nil)
	if err != nil {
		return nil, nil, err
	}

	// Perform request, attempt to unmarshal response into a
	// Users API response
	usersRes := new(v0.UsersResponse)
	res, err := u.client.Do(req, &usersRes)
	if err != nil {
		return nil, res, err
	}

	// Return users found by API
	return usersRes.Users, res, nil
}
