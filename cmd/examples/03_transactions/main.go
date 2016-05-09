// Package main (03_transactions) demonstrates how to enumerate all of an
// accounts transactions. As with the second example, you need to set the
// MONDO_ACCESS_TOKEN environment variable.
package main

import (
	"encoding/json"
	"fmt"
	"github.com/icio/mondo"
	"github.com/icio/mondo/mondodomain"
	"github.com/icio/mondo/mondohttp"
	"log"
	"net/http"
	"os"
)

func main() {
	// Our mondo client.
	client := &mondo.Client{
		HTTPClient: &mondo.HTTPClient{
			Client:    http.DefaultClient,
			UserAgent: "example/0.1 (+https://github.com/icio/mondo)",
		},
		Auth: mondo.NewAccessTokenAuth(os.Getenv("MONDO_ACCESS_TOKEN")),
	}

	// Get the first account.
	accounts := new(mondodomain.AccountCollection)
	err := client.DoInto(mondohttp.NewAccountsRequest(""), accounts)
	if err != nil {
		log.Fatal(err)
	}
	account := accounts.Accounts[0].ID

	// The trans channel will recieve all of the transactions from the account
	// until all transactions have been listed, or we signal to stop.
	trans := make(chan mondodomain.Transaction)

	// We can use the stop channel to signal we want no more transactions. (If
	// our application always iterates through all transactions we can pass nil
	// into IterTransactions.)
	stop := make(chan bool)

	// The iteration is done from a goroutine we create ourselves:
	// mondo.IterTransactions will block until it has written all transactions
	// from the API to the channel we give it, or until it receives a stop
	// signal.
	go func() {
		defer close(trans)
		defer close(stop)

		// Load 30 transactions per page from the API, from the beginning of
		// the account's history.
		err := client.IterTransactions(trans, stop, "", account, true, "", "", 30)
		if err != nil {
			log.Fatal(err)
		}
	}()

	n := 0
	for tran := range trans {
		// Dump the transaction.
		enc, err := json.Marshal(tran)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%d: %s\n", n, enc)

		// Stop after the first 120 results.
		n++
		if n == 120 {
			stop <- true
		}
	}
}
