package mondo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// InvalidToken is the Mondo error type representing an expired access token.
const InvalidToken string = "invalid_token"

// Error is a generic library error with additional details about the context.
type Error struct {
	Request  *http.Request
	Response *http.Response
	Message  string
	cause    error
}

// WrapError creates a library Error from a causal error request context.
func WrapError(cause error, message string, req *http.Request, resp *http.Response) *Error {
	return &Error{
		Request:  req,
		Response: resp,
		Message:  message,
		cause:    cause,
	}
}

func (err *Error) Error() string {
	return fmt.Sprintf(
		"mondo: %s (caused by: %s) during {%s}",
		err.Message,
		err.cause,
		formatReqResp(err.Request, err.Response),
	)
}

// Cause returns the error that resulted in this Error being generated. This
// method satisfies the causer interface of github.com/pkg/errors for
// error-message unravelling.
func (err *Error) Cause() error {
	return err.cause
}

// ResponseError respresents error messages successfully returned from Mondo.
type ResponseError struct {
	Request  *http.Request
	Response *http.Response

	// Values read from an HTTP error response.
	TraceID   string            `json:"request_id"`
	Code      string            `json:"code"`
	ErrorType string            `json:"error"`
	Message   string            `json:"message"`
	Params    map[string]string `json:"params"`

	// InvalidToken usually indicates whether a token has expired.
	InvalidToken bool
}

// DecodeError parses ResponseErrors from Mondo API calls.
func DecodeError(req *http.Request, resp *http.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return WrapError(err, "Failed to read response body", req, resp)
	}

	mondoErr := &ResponseError{
		Request:  req,
		Response: resp,
		TraceID:  resp.Header.Get("Trace-ID"),
	}
	err = json.Unmarshal(body, mondoErr)
	if err != nil {
		return WrapError(err, "Failed to decode error response", req, resp)
	}

	mondoErr.InvalidToken = mondoErr.ErrorType == InvalidToken
	return mondoErr
}

func (err *ResponseError) Error() string {
	errType := ""
	if err.ErrorType != "" {
		errType = fmt.Sprintf(" (%s)", err.ErrorType)
	}
	return fmt.Sprintf(
		"mondo: %s%s during {%s} (trace: %s)",
		err.Message,
		errType,
		formatReqResp(err.Request, err.Response),
		err.TraceID,
	)
}

func (err *Error) String() string {
	return err.Error()
}

func formatReqResp(req *http.Request, resp *http.Response) string {
	desc := ""
	if req != nil {
		desc += formatReq(req)
	}
	if resp != nil {
		if resp.Request != nil && resp.Request != req {
			if req != nil {
				desc += " -> "
			}
			desc += formatReq(resp.Request)
		}
		desc += " => " + formatResp(resp)
	}
	return desc
}

func formatReq(req *http.Request) string {
	return fmt.Sprintf("%s %s", req.Method, req.URL)
}

func formatResp(resp *http.Response) string {
	return resp.Status
}
