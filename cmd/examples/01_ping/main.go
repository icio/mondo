// Package main (01_ping) demonstrates the basic request structure.
package main

import (
	"fmt"
	"github.com/icio/mondo"
	"github.com/icio/mondo/mondodomain"
	"github.com/icio/mondo/mondohttp"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

func main() {
	// We can construct an unauthenticated client using only a http.Client.
	client := &mondo.Client{
		HTTPClient: http.DefaultClient,
	}

	// Working with raw responses: Whilst we use mondo.Client here, we could
	// use http.Client as their Do methods share the same signature.  We might
	// prefer mondo.Client for its Mondo-specific error-handling, however,
	// which we'll see below.
	resp, err := client.Do(mondohttp.NewPingRequest())
	if err != nil {
		fmt.Println("Error during ping:", err)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading body of ping:", err)
		return
	}
	fmt.Printf("Ping Response: [%d] %s\n", resp.StatusCode, body)

	// Working with unmarshalled responses: mondo.Client.DoInto combines the
	// performing of requesting with unmarshalling for convenience' sake.
	whoami := new(mondodomain.Identity)
	err = client.DoInto(mondohttp.NewWhoAmIRequest(""), whoami)
	if err != nil {
		fmt.Println("Error during whoami:", err)
		return
	}
	fmt.Printf("WhoAmI Response: %#v\n", whoami)

	// Error handling: errors are unmarshalled into mondo.Error objects.
	_, err = client.Do(mondohttp.NewAccountsRequest(""))
	fmt.Printf("Error: %s\n", err)

	// Once again, attempting to authenticate. See later examples for actual
	// means of authentication.
	_, err = client.Do(mondohttp.NewAccountsRequest("Bearer abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzab"))
	fmt.Printf("Error: %s\n", err)
	if mondoErr, ok := errors.Cause(err).(*mondo.ResponseError); ok {
		fmt.Println("Invalid token:", mondoErr.InvalidToken)
	}
}
