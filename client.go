package mondo

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type httpclient interface {
	Do(*http.Request) (*http.Response, error)
}

type auth interface {
	Get(invalidate bool, client *Client) (string, error)
}

// Client wraps Mondo-specific error-handling and authentication around an HTTP
// Client.
type Client struct {
	HTTPClient httpclient
	Auth       auth
}

// Do performs a request and returns the raw HTTP response. Any authorization
// headers on the request are overridden by what the Client's Auth provides.
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	_, authedReq := req.Header[http.CanonicalHeaderKey("Authorization")]
	if !authedReq || c.Auth == nil {
		return c.do(req)
	}
	return c.doAuth(req)
}

// DoInto performs a request and stores the result into the given object. Any
// authorization headers on the request are overridden by what the Client's
// Auth provides.
func (c *Client) DoInto(req *http.Request, target interface{}) error {
	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return WrapError(err, "Failed to read response body", req, resp)
	}

	err = json.Unmarshal(body, target)
	if err != nil {
		return WrapError(err, "Failed to parse response body", req, resp)
	}

	return nil
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return resp, WrapError(err, "HTTP request failed", req, resp)
	}

	if resp.StatusCode != 200 {
		respErr := DecodeError(req, resp)
		if _, ok := respErr.(*ResponseError); ok {
			return resp, respErr
		}
		return resp, WrapError(respErr, "Non-200 response returned and an error occured during parsing", req, resp)
	}

	return resp, nil
}

func (c *Client) doAuth(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var _err error

	for attempt := 0; attempt < 2; attempt++ {
		// Get the token. Get a fresh token in subsequent attempts.
		invalidate := attempt > 0
		token, err := c.Auth.Get(invalidate, c)
		if err != nil {
			return nil, WrapError(err, "Failed to get authentication details", req, nil)
		}

		req.Header.Set("Authorization", token)
		resp, err = c.do(req)
		if err != nil {
			if mondoErr, ok := err.(*ResponseError); ok && mondoErr.InvalidToken {
				_err = mondoErr
				continue
			}
			return resp, err
		}

		break
	}

	return resp, _err
}
