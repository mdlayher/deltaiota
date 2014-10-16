package util

import (
	"log"
	"net/http"
)

// LogHandler provides basic logging of actions which are passed through a
// http.Handler.
type LogHandler struct {
	http.Handler
}

// ServeHTTP allows LogHandler to be used as a http.Handler, and captures
// information regarding the client request and server response.
func (l LogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Wrap writer in custom logged writer
	w = &logResponseWriter{
		ResponseWriter: w,
	}

	// Call underlying handler
	l.Handler.ServeHTTP(w, r)

	// Type-assert back to logged response writer
	lw, ok := w.(*logResponseWriter)
	if ok {
		// Log information about the request and response
		log.Printf("[%s] HTTP %d: %s %s", r.RemoteAddr, lw.Status, r.Method, r.URL)
	}
}

// logResponseWriter captures the HTTP status code sent to a client.
type logResponseWriter struct {
	http.ResponseWriter
	Status int
}

// WriteHeader captures the HTTP status code sent to a client.
func (w *logResponseWriter) WriteHeader(s int) {
	w.ResponseWriter.WriteHeader(s)
	w.Status = s
}
