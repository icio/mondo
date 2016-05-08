package mondohttp

import (
	"net/http"
	"net/url"
	"strings"
)

// NewAuthCodeAccessRequest creates a request for exchanging authorization codes.
// https://getmondo.co.uk/docs/#exchange-the-authorization-code
func NewAuthCodeAccessRequest(clientID, clientSecret, redirectURI, authCode string) *http.Request {
	body := strings.NewReader(url.Values{
		"grant_type":    {"authorization_code"},
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"redirect_uri":  {redirectURI},
		"code":          {authCode},
	}.Encode())

	req, _ := http.NewRequest("POST", ProductionAPI+"oauth2/token", body)
	req.Header.Set(formContentType())
	return req
}

// NewRefreshAccessRequest creates a request for refreshing an access token.
// https://getmondo.co.uk/docs/#refreshing-access
func NewRefreshAccessRequest(clientID, clientSecret, refreshToken string) *http.Request {
	body := strings.NewReader(url.Values{
		"grant_type":    {"refresh_token"},
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"refresh_token": {refreshToken},
	}.Encode())

	req, _ := http.NewRequest("POST", ProductionAPI+"oauth2/token", body)
	req.Header.Set(formContentType())
	return req
}

// NewWhoAmIRequest creates a request for verifying the authenticated identity.
// https://getmondo.co.uk/docs/#authenticating-requests
func NewWhoAmIRequest(accessToken string) *http.Request {
	req, _ := http.NewRequest("GET", ProductionAPI+"ping/whoami", nil)
	req.Header.Set(auth(accessToken))
	return req
}

// NewPingRequest creates a request for pinging the API service.
// (Not documented, but is an example in the API Playground.)
func NewPingRequest() *http.Request {
	req, _ := http.NewRequest("GET", ProductionAPI+"ping", nil)
	return req
}
