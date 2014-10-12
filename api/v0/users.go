package v0

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/mdlayher/deltaiota/api/util"
	"github.com/mdlayher/deltaiota/data"
	"github.com/mdlayher/deltaiota/data/models"
)

// JSON Users API, human-readable client error responses.
const (
	// HTTP GET
	userInvalidID = "invalid user ID"
	userMissingID = "missing user ID"
	userNotFound  = "user not found"

	// HTTP POST
	userConflict          = "user already exists"
	userInvalidParameters = "invalid parameters"
	userJSONSyntax        = "invalid JSON request"
	userMissingParameters = "missing required parameters"
)

// JSON Users API, map of client errors to response codes.
var usersCode = map[string]int{
	// HTTP GET
	userInvalidID: http.StatusBadRequest,
	userMissingID: http.StatusBadRequest,
	userNotFound:  http.StatusNotFound,

	// HTTP POST
	userConflict:          http.StatusConflict,
	userInvalidParameters: http.StatusBadRequest,
	userJSONSyntax:        http.StatusBadRequest,
	userMissingParameters: http.StatusBadRequest,
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

// UsersAPI is a util.JSONAPIFunc, and is the single entry point for the Users API.
// This method delegates to other methods as appropriate to handle incoming requests.
func (c *Context) UsersAPI(r *http.Request, vars util.Vars) (int, []byte, error) {
	// Switch based on HTTP method
	switch r.Method {
	case "GET", "HEAD":
		// If ID present, request for single user
		if _, ok := vars["id"]; ok {
			return c.GetUser(r, vars)
		}

		// No ID, request for list of users
		return c.ListUsers(r, vars)
	case "POST":
		return c.PostUser(r, vars)
	case "PUT":
		return c.PutUser(r, vars)
	case "DELETE":
		return c.DeleteUser(r, vars)
	default:
		return util.MethodNotAllowed(r, vars)
	}
}

// ListUsers is a util.JSONAPIFunc which returns HTTP 200 and a JSON list of users
// on success, or a non-200 HTTP status code and an error response on failure.
func (c *Context) ListUsers(r *http.Request, vars util.Vars) (int, []byte, error) {
	// Fetch a list of all users from the database
	users, err := c.db.SelectAllUsers()
	if err != nil {
		return util.JSONAPIErr(err)
	}

	// Strip all passwords from output
	for i := range users {
		users[i].Password = ""
	}

	// Wrap in response and return
	body, err := json.Marshal(UsersResponse{
		Users: users,
	})
	return http.StatusOK, body, err
}

// GetUser is a util.JSONAPIFunc which returns HTTP 200 and a JSON user object
// on success, or a non-200 HTTP status code and an error response on failure.
func (c *Context) GetUser(r *http.Request, vars util.Vars) (int, []byte, error) {
	// Fetch input user ID
	strID, ok := vars["id"]
	if !ok {
		return usersCode[userMissingID], usersJSON[userMissingID], nil
	}

	// Convert string to integer
	id, err := strconv.ParseUint(strID, 10, 64)
	if err != nil {
		return usersCode[userInvalidID], usersJSON[userInvalidID], nil
	}

	// Select single user by ID from the database
	user, err := c.db.SelectUserByID(id)
	if err != nil {
		// If no results found, return HTTP not found
		if err == sql.ErrNoRows {
			return usersCode[userNotFound], usersJSON[userNotFound], nil
		}

		return util.JSONAPIErr(err)
	}

	// Strip password from output
	user.Password = ""

	// Wrap in response and return
	body, err := json.Marshal(UsersResponse{
		Users: []*models.User{user},
	})
	return http.StatusOK, body, err
}

// PostUser is a util.JSONAPIFunc which creates a User and returns HTTP 201
// and a JSON user object on success, or a non-200 HTTP status code and an
// error response on failure.
func (c *Context) PostUser(r *http.Request, vars util.Vars) (int, []byte, error) {
	// Read and validate request input into a User struct
	user, code, body, err := c.jsonToUser(r)
	if err != nil {
		return util.JSONAPIErr(err)
	}

	// If a body was written (probably client error), return now
	if body != nil {
		return code, body, nil
	}

	// No body written, all checks passed, so insert new user
	if err := c.db.InsertUser(user); err != nil {
		// Check for constraint failure, meaning user already exists
		if c.db.IsConstraintFailure(err) {
			return usersCode[userConflict], usersJSON[userConflict], nil
		}

		return util.JSONAPIErr(err)
	}

	// Strip password from output
	user.Password = ""

	// Wrap in response and return
	body, err = json.Marshal(UsersResponse{
		Users: []*models.User{user},
	})
	return http.StatusCreated, body, err
}

// PutUser is a util.JSONAPIFunc which updates a User and returns HTTP 200
// and a JSON user object on success, or a non-200 HTTP status code and an
// error response on failure.
func (c *Context) PutUser(r *http.Request, vars util.Vars) (int, []byte, error) {
	// Fetch input user ID
	strID, ok := vars["id"]
	if !ok {
		return usersCode[userMissingID], usersJSON[userMissingID], nil
	}

	// Convert string to integer
	id, err := strconv.ParseUint(strID, 10, 64)
	if err != nil {
		return usersCode[userInvalidID], usersJSON[userInvalidID], nil
	}

	// Select single user by ID from the database
	user, err := c.db.SelectUserByID(id)
	if err != nil {
		// If no results found, return HTTP not found
		if err == sql.ErrNoRows {
			return usersCode[userNotFound], usersJSON[userNotFound], nil
		}

		return util.JSONAPIErr(err)
	}

	// Read and validate request input into a User struct
	newUser, code, body, err := c.jsonToUser(r)
	if err != nil {
		return util.JSONAPIErr(err)
	}

	// If a body was written (probably client error), return now
	if body != nil {
		return code, body, nil
	}

	// No body written, all checks passed, so update existing user with
	// new fields
	//  - Email already validated in jsonToUser
	//  - Password already hashed in jsonToUser
	user.CopyFrom(newUser)
	if err := c.db.UpdateUser(user); err != nil {
		// Check for constraint failure, meaning a unique check failed
		if c.db.IsConstraintFailure(err) {
			return usersCode[userConflict], usersJSON[userConflict], nil
		}

		return util.JSONAPIErr(err)
	}

	// Strip password from output
	user.Password = ""

	// Wrap in response and return
	body, err = json.Marshal(UsersResponse{
		Users: []*models.User{user},
	})
	return http.StatusOK, body, err
}

// DeleteUser is a util.JSONAPIFunc which deletes a User and returns HTTP 204
// on success, or a non-200 HTTP status code and an error response on failure.
func (c *Context) DeleteUser(r *http.Request, vars util.Vars) (int, []byte, error) {
	// Fetch input user ID
	strID, ok := vars["id"]
	if !ok {
		return usersCode[userMissingID], usersJSON[userMissingID], nil
	}

	// Convert string to integer
	id, err := strconv.ParseUint(strID, 10, 64)
	if err != nil {
		return usersCode[userInvalidID], usersJSON[userInvalidID], nil
	}

	// Select single user by ID from the database
	user, err := c.db.SelectUserByID(id)
	if err != nil {
		// If no results found, return HTTP not found
		if err == sql.ErrNoRows {
			return usersCode[userNotFound], usersJSON[userNotFound], nil
		}

		return util.JSONAPIErr(err)
	}

	// Clear user data within a transaction
	err = c.db.WithTx(func(tx *data.Tx) error {
		// Delete all sessions for user
		if err := tx.DeleteSessionsByUserID(user.ID); err != nil {
			return err
		}

		// Delete user
		return tx.DeleteUser(user)
	})

	// Check for transaction errors
	if err != nil {
		return util.JSONAPIErr(err)
	}

	return http.StatusNoContent, nil, nil
}

// jsonToUser reads the JSON body of an incoming HTTP request, validates that
// all required fields are set, and returns a User on success.
// On failure, it will return a message body or an error, causing the caller to
// immediately send the result.
func (c *Context) jsonToUser(r *http.Request) (*models.User, int, []byte, error) {
	// Unmarshal body into a User
	user := new(models.User)
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		// Check for bad input JSON
		if _, ok := err.(*json.SyntaxError); ok || err == io.EOF || err == io.ErrUnexpectedEOF {
			return nil, usersCode[userJSONSyntax], usersJSON[userJSONSyntax], nil
		}

		return nil, http.StatusInternalServerError, nil, err
	}

	// Attempt to set password from input
	if err := user.SetPassword(user.Password); err != nil {
		// If empty password was passed, we are missing a parameter
		if emptyErr, ok := err.(*models.EmptyFieldError); ok {
			// Set code for missing parameter
			code := usersCode[userMissingParameters]

			// Return customized error object
			body, err := json.Marshal(util.ErrRes(code, emptyErr.Error()))
			return nil, code, body, err
		}

		return nil, http.StatusInternalServerError, nil, err
	}

	// Validate input for user
	if err := user.Validate(); err != nil {
		// If a required field was empty, report missing parameters
		if emptyErr, ok := err.(*models.EmptyFieldError); ok {
			// Set code for missing parameter
			code := usersCode[userMissingParameters]

			// Return customized error object
			body, err := json.Marshal(util.ErrRes(code, emptyErr.Error()))
			return nil, code, body, err
		}

		// If a field was invalid, report invalid input
		if invalidErr, ok := err.(*models.InvalidFieldError); ok {
			// Set code for missing parameter
			code := usersCode[userInvalidParameters]

			// Return customized error object
			body, err := json.Marshal(util.ErrRes(code, invalidErr.Error()))
			return nil, code, body, err
		}

		// For any other errors, report a server error
		return nil, http.StatusInternalServerError, nil, err
	}

	// All validations passed, return User with no body so processing
	// can continue in caller
	return user, http.StatusOK, nil, nil
}
