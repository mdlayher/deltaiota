// Package diclient implements a client library for the HTTP API of the
// Phi Mu Alpha Sinfonia - Delta Iota chapter website.
//
// This package is inspired by Google's go-github library, which can
// be found here: https://github.com/google/go-github.
package diclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	// version is the API version which this client implements
	version = "v0"

	// jsonContentType is the content type for JSON data
	jsonContentType = "application/json"
)

// Client provides a client interface for the HTTP API of the Phi Mu Alpha
// Sinfonia - Delta Iota chapter website.
type Client struct {
	client *http.Client
	url    *url.URL

	userAgent string

	Users *UsersService
}

// NewClient creates a new Client for the HTTP API at the specified host.
// Optionally, a custom http.Client may be specified.  If no http.Client is
// specified, http.DefaultClient will be used.
func NewClient(host string, client *http.Client) (*Client, error) {
	// Parse input host for a valid URL
	u, err := url.Parse(host)
	if err != nil {
		return nil, err
	}

	// If input client is nil, use http.DefaultClient
	if client == nil {
		client = http.DefaultClient
	}

	// Set up basic client
	c := &Client{
		client: client,
		url:    u,

		userAgent: "diclient",
	}

	// Set up individual services within client
	c.Users = &UsersService{
		client: c,
	}

	return c, nil
}

// NewRequest creates a new HTTP request, using the specified HTTP method and API endpoint.
// Optionally, a request body may be sent.
func (c *Client) NewRequest(method string, endpoint string, body interface{}) (*http.Request, error) {
	// Generate relative URL using API root, version, and endpoint
	rel, err := url.Parse(fmt.Sprintf("api/%s/%s", version, endpoint))
	if err != nil {
		return nil, err
	}

	// Resolve relative URL to base, using input host
	u := c.url.ResolveReference(rel)

	// If a body object was specified, encode it to JSON
	buf := bytes.NewBuffer(nil)
	if body != nil {
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	// Generate new HTTP request for appropriate URL, with optional body
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	// Set headers to indicate proper content type
	req.Header.Add("Accept", jsonContentType)
	req.Header.Add("Content-Type", jsonContentType)

	// Identify the client
	req.Header.Add("User-Agent", c.userAgent)

	return req, nil
}

// Response is a wrapped http.Response.  It may be expanded upon later.
type Response struct {
	*http.Response
}

// Do invokes the input HTTP request, and attempts to unmarshal or stream any
// JSON response into the object passed by the second parameter.
func (c *Client) Do(req *http.Request, v interface{}) (*Response, error) {
	// Invoke request using underlying HTTP client
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Wrap underlying response in custom type
	wrapRes := &Response{
		Response: res,
	}

	// If no second parameter was passed, do not attempt to handle response
	if v == nil {
		return wrapRes, nil
	}

	// Attempt to unmarshal JSON
	switch vt := v.(type) {
	// If the input is a io.Writer, attempt to stream response body into writer
	case io.Writer:
		_, err = io.Copy(vt, res.Body)
	// For other cases, attempt to decode response body into object
	default:
		err = json.NewDecoder(res.Body).Decode(v)
	}

	return wrapRes, err
}
