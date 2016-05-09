package mondo

import (
	"errors"
	"github.com/icio/mondo/mondodomain"
	"github.com/icio/mondo/mondohttp"
	"sync"
)

// ErrNoCredentials indicates when UserAuth doesn't have credentials to auth.
var ErrNoCredentials = errors.New("mondo: No credentials for generating access token")

// UserAuth provides the access token for a single user-account, requesting
// a new access token using a refresh token or username/password when required.
// Thread-safe.
type UserAuth struct {
	ClientID     string
	ClientSecret string
	AccessToken  string
	RefreshToken string
	Lock         sync.RWMutex
}

// NewAccessTokenAuth prepares an auth with existing access token (without the Bearer prefix).
func NewAccessTokenAuth(accessToken string) *UserAuth {
	return &UserAuth{AccessToken: "Bearer " + accessToken}
}

// NewClientAccessTokenAuth prepares an auth with existing access token (with a
// Bearer prefix) and client credentials for refreshing when it expires.
func NewClientAccessTokenAuth(clientID, clientSecret, accessToken, refreshToken string) *UserAuth {
	return &UserAuth{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}

// Get returns an access token (e.g. "Bearer xyz..."), making authentication
// requests using the credentials available to it, where necessary.
func (auth *UserAuth) Get(invalidate bool, client *Client) (string, error) {
	// Attempt to reuse an existing access token.
	auth.Lock.RLock()
	if !invalidate && auth.AccessToken != "" {
		auth.Lock.RUnlock()
		return auth.AccessToken, nil
	}

	// Prepare to fetch a token.
	existingToken := auth.AccessToken
	auth.Lock.RUnlock()
	auth.Lock.Lock()
	defer auth.Lock.Unlock()

	// Check if the token was refreshed elsewhere.
	if auth.AccessToken != "" && existingToken != auth.AccessToken {
		return auth.AccessToken, nil
	}

	// Invalidate any existing access token.
	auth.AccessToken = ""
	var err error

	token := new(mondodomain.Token)
	if auth.ClientID != "" && auth.ClientSecret != "" && auth.RefreshToken != "" {
		// Refresh the access token.
		err = client.DoInto(
			mondohttp.NewRefreshAccessRequest(auth.ClientID, auth.ClientSecret, auth.RefreshToken),
			token,
		)
	} else {
		err = ErrNoCredentials
	}
	if err != nil {
		return "", err
	}

	auth.AccessToken = token.TokenType + " " + token.AccessToken
	auth.RefreshToken = token.RefreshToken

	return auth.AccessToken, nil
}
