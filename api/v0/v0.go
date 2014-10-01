// Package v0 provides the development HTTP API for the Phi Mu Alpha Sinfonia - Delta
// Iota chapter website.
package v0

import (
	"net/http"

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

	// Create a handler which stores the database connection
	h := &Handler{
		db: db,
	}

	// Set up HTTP routes
	r.HandleFunc("/", h.Root).Methods("GET")

	return r
}

// Handler stores shared members for HTTP handlers.
type Handler struct {
	db *data.DB
}

// Root handles the root HTTP route.
func (h *Handler) Root(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
}
