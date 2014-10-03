// Package v0 provides the development HTTP API for the Phi Mu Alpha Sinfonia - Delta
// Iota chapter website.
package v0

import (
	"net/http"

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
	c := &context{
		db: db,
	}

	// Set up HTTP routes

	// Users API
	r.Handle("/users", util.JSONAPIHandler(c.ListUsers)).Methods("GET")
	r.Handle("/users/{id}", util.JSONAPIHandler(c.GetUser)).Methods("GET")

	return r
}

// context stores shared members for API v0 HTTP handlers.
type context struct {
	db *data.DB
}
