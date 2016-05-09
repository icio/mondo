package mondo

import (
	"net/http"
)

// HTTPClient wraps another HTTPClient (such as *http.Client) instances
// with Host and UserAgent overrides.
type HTTPClient struct {
	Client    httpclient
	Host      string
	UserAgent string
}

// Do performs the HTTP request as http.Client.Do, setting overrides.
func (c *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	if c.UserAgent == "" {
		req.Header.Del("User-Agent")
	} else {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	if c.Host != "" {
		req.URL.Host = c.Host
	}
	return c.Client.Do(req)
}
