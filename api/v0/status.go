package v0

import (
	"encoding/json"
	"net/http"
	"os"
	"runtime"

	"github.com/mdlayher/deltaiota/api/util"
)

// StatusResponse is the output response for the Status API
type StatusResponse struct {
	Status *Status `json:"status"`
}

// Status contains information about the running server process.
type Status struct {
	Architecture string `json:"architecture"`
	Hostname     string `json:"hostname"`
	NumCPU       int    `json:"numCpu"`
	NumGoroutine int    `json:"numGoroutine"`
	PID          int    `json:"pid"`
	Platform     string `json:"platform"`
}

// StatusAPI is a util.JSONAPIFunc, and is the single entry point for the Status API.
// This method delegates to other methods as appropriate to handle incoming requests.
func (c *Context) StatusAPI(r *http.Request, vars util.Vars) (int, []byte, error) {
	// Switch based on HTTP method
	switch r.Method {
	case "GET", "HEAD":
		return c.GetStatus(r, vars)
	default:
		return util.MethodNotAllowed(r, vars)
	}
}

// GetStatus is a util.JSONAPIFunc which returns HTTP 200 and current server status
// on success, or a non-200 HTTP status code and an error response on failure.
func (c *Context) GetStatus(r *http.Request, vars util.Vars) (int, []byte, error) {
	// Fetch hostname
	hostname, err := os.Hostname()
	if err != nil {
		return util.JSONAPIErr(err)
	}

	// Wrap in response
	body, err := json.Marshal(StatusResponse{
		Status: &Status{
			Architecture: runtime.GOARCH,
			Hostname:     hostname,
			NumCPU:       runtime.NumCPU(),
			NumGoroutine: runtime.NumGoroutine(),
			PID:          os.Getpid(),
			Platform:     runtime.GOOS,
		},
	})
	return http.StatusOK, body, err
}
