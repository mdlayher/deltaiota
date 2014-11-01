package diclient

import (
	"fmt"

	"github.com/mdlayher/deltaiota/api/v0"
	"github.com/mdlayher/deltaiota/data/models"
)

// UsersService provides access to the Users API.
type UsersService struct {
	client *Client
}

// ListUsers returns a slice of all User objects from the API.
func (u *UsersService) ListUsers() ([]*models.User, *Response, error) {
	usersRes, res, err := u.usersRequest("GET", "users", nil)
	return usersRes.Users, res, err
}

// GetUser returns a single User object with the input ID from the API.
func (u *UsersService) GetUser(id uint64) (*models.User, *Response, error) {
	usersRes, res, err := u.usersRequest("GET", fmt.Sprintf("users/%d", id), nil)

	// Check for no user found
	if usersRes == nil || usersRes.Users == nil || len(usersRes.Users) == 0 {
		return nil, res, err
	}

	return usersRes.Users[0], res, err
}

// CreateUser generates an API user using the input User object.
func (u *UsersService) CreateUser(user *models.User) (*Response, error) {
	_, res, err := u.usersRequest("POST", "users", user)
	return res, err
}

// UpdateUser updates an existing API user using the input User object.
func (u *UsersService) UpdateUser(user *models.User) (*Response, error) {
	_, res, err := u.usersRequest("PUT", fmt.Sprintf("users/%d", user.ID), user)
	return res, err
}

// usersRequest generates and performs a HTTP request to the Users API.
func (u *UsersService) usersRequest(method string, endpoint string, body interface{}) (*v0.UsersResponse, *Response, error) {
	// Create request for Users endpoint
	req, err := u.client.NewRequest(method, endpoint, body)
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

	return usersRes, res, nil
}
