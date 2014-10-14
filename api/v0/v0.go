// Package v0 provides the development HTTP API for the Phi Mu Alpha Sinfonia - Delta
// Iota chapter website.
package v0

import (
	"net/http"

	"github.com/mdlayher/deltaiota/api/auth"
	"github.com/mdlayher/deltaiota/api/util"
	"github.com/mdlayher/deltaiota/data"

	"github.com/gorilla/mux"
)

const (
	// APIPrefix is the string which prefixes all routes for this API.
	APIPrefix = "/api/v0"
)

// NewServeMux returns a new http.Handler which contains the necessary HTTP routes
// for the development deltaiota HTTP server.
func NewServeMux(db *data.DB) http.Handler {
	// Create new mux to be configured
	r := mux.NewRouter().StrictSlash(true).PathPrefix(APIPrefix).Subrouter()

	// Create a context which stores any shared members
	c := &Context{
		db: db,
	}

	// Set up authentication context
	ac := auth.NewContext(db)

	// Set up HTTP routes

	// Notifications API
	r.Handle("/notifications", ac.KeyAuthHandler(util.JSONAPIHandler(c.NotificationsAPI)))

	// Sessions API
	r.Handle("/sessions", ac.PasswordAuthHandler(util.JSONAPIHandler(c.PostSession))).Methods("POST")
	r.Handle("/sessions", ac.KeyAuthHandler(util.JSONAPIHandler(c.SessionsAPI))).Methods("GET", "HEAD", "PUT", "PATCH", "DELETE")

	// Users API
	r.Handle("/users", ac.KeyAuthHandler(util.JSONAPIHandler(c.UsersAPI)))
	r.Handle("/users/{id}", ac.KeyAuthHandler(util.JSONAPIHandler(c.UsersAPI)))

	return r
}

// Context stores shared members for API v0 HTTP handlers.
type Context struct {
	db *data.DB
}
