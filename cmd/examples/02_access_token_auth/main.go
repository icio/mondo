// Package main (02_access_token_auth) demonstrates how to make authenticated
// requests to Mondo. To follow this example:
//
// 1. Log into https://developers.getmondo.co.uk/api/playground;
// 2. Copy the "Access token" shown on;
// 3. Run the this program with:
//    MONDO_ACCESS_TOKEN=<paste> go run main.go
//
// (Note that the environment variables a program runs with can be read, so be
// wary of the security implications of this approach.)
package main

import (
	"github.com/icio/mondo"
	"github.com/icio/mondo/mondodomain"
	"github.com/icio/mondo/mondohttp"
	"log"
	"net/http"
	"os"
)

func main() {
	client := &mondo.Client{
		// The mondohttp package assumes the Production API host and uses Go's
		// default User-Agent, but we can override these before each request is
		// made when using mondo.HTTPClient.
		HTTPClient: &mondo.HTTPClient{
			Client:    http.DefaultClient,
			Host:      "api.getmondo.co.uk",
			UserAgent: "example/0.1 (+https://github.com/icio/mondo)",
		},
		// When we give mondo.Client a mondo.auth, any Authorization headers
		// present in requests are updated before the request is made.
		Auth: mondo.NewAccessTokenAuth(os.Getenv("MONDO_ACCESS_TOKEN")),
	}

	// mondo.Client is going to override the Authorization header, so we just
	// provide an empty string for the access token when creating the request.
	accounts := new(mondodomain.AccountCollection)
	err := client.DoInto(mondohttp.NewAccountsRequest(""), accounts)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Accounts: %#v\n", accounts)
}
