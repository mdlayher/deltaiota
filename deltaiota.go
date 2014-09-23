// Package deltaiota provides a HTTP API for the Phi Mu Alpha Sinfonia - Delta
// Iota chapter website.
package deltaiota

import (
	"net/http"

	"github.com/gorilla/mux"
)

// NewServeMux returns a new http.Handler which contains the necessary HTTP routes
// for the deltaiota HTTP server.
func NewServeMux(db *DB) http.Handler {
	// Create new mux to be configured
	r := mux.NewRouter()

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
	db *DB
}

// Root handles the root HTTP route.
func (h *Handler) Root(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
}
