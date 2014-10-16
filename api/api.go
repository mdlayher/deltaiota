// Package api provides a HTTP API for the Phi Mu Alpha Sinfonia - Delta
// Iota chapter website.
package api

import (
	"net/http"

	"github.com/mdlayher/deltaiota/api/util"
	"github.com/mdlayher/deltaiota/api/v0"
	"github.com/mdlayher/deltaiota/data"

	"github.com/gorilla/mux"
)

// NewServeMux returns a new http.Handler which contains the necessary HTTP routes
// for all versions of the deltaiota HTTP server.
func NewServeMux(db *data.DB) http.Handler {
	// Create new mux to be configured
	r := mux.NewRouter().StrictSlash(true)

	// Create a handler for all v0 API routes
	r.PathPrefix(v0.APIPrefix).Handler(util.LogHandler{v0.NewServeMux(db)})

	return r
}
