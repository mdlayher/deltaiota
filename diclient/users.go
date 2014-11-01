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

// List returns a slice of all User objects from the API.
func (u *UsersService) List() ([]*models.User, *Response, error) {
	uRes, res, err := u.request("GET", "users", nil)

	// Check for empty users
	if uRes == nil || uRes.Users == nil {
		return nil, res, err
	}

	return uRes.Users, res, err
}

// Get returns a single User object with the input ID from the API.
func (u *UsersService) Get(id uint64) (*models.User, *Response, error) {
	uRes, res, err := u.request("GET", fmt.Sprintf("users/%d", id), nil)

	// Check for no user found
	if uRes == nil || uRes.Users == nil || len(uRes.Users) == 0 {
		return nil, res, err
	}

	return uRes.Users[0], res, err
}

// Create generates an API user using the input User object.
func (u *UsersService) Create(user *models.User) (*Response, error) {
	_, res, err := u.request("POST", "users", user)
	return res, err
}

// Update updates an existing API user using the input User object.
func (u *UsersService) Update(user *models.User) (*Response, error) {
	_, res, err := u.request("PUT", fmt.Sprintf("users/%d", user.ID), user)
	return res, err
}

// request generates and performs a HTTP request to the Users API.
func (u *UsersService) request(method string, endpoint string, body interface{}) (*v0.UsersResponse, *Response, error) {
	// Create request for Users endpoint
	req, err := u.client.NewRequest(method, endpoint, body)
	if err != nil {
		return nil, nil, err
	}

	// Perform request, attempt to unmarshal response into a
	// Users API response
	uRes := new(v0.UsersResponse)
	res, err := u.client.Do(req, &uRes)
	if err != nil {
		return nil, res, err
	}

	return uRes, res, nil
}
