// Package httpnoquery defines a http request and client. The
// difference between this client and the standard http.Client is that
// the error returned from a failed request will not contain any query
// string parameters. It was created because some API's, for better or
// worse, include sensitive information in query string parameters.
package httpnoquery

import (
	"fmt"
	"net/http"
	"net/url"
)

// Request is sent by this client. It was created to make it less
// possible to mess up and use a regular http.Client to send data.
type Request struct {
	request *http.Request
}

// NewRequest creates a Request.
func NewRequest(r *http.Request) Request {
	return Request{r}
}

// Client sends http requests and modifies the returned error so it
// contains no query strings parameters.
type Client struct {
	HTTPClient *http.Client
}

func (c Client) httpClient() *http.Client {
	if c.HTTPClient == nil {
		return http.DefaultClient
	}
	return c.HTTPClient
}

// Do sends the request and removes any query string from the returned
// error message.
func (c Client) Do(r Request) (*http.Response, error) {
	resp, err := c.httpClient().Do(r.request)
	if urlErr, ok := err.(*url.Error); ok {
		urlNoQuery := *r.request.URL
		urlNoQuery.RawQuery = ""
		return resp, fmt.Errorf("sending request: %s %s: %v", urlErr.Op, urlNoQuery.String(), urlErr.Err)
	}
	return resp, err
}
