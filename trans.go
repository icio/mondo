package mondo

import (
	"github.com/icio/mondo/mondodomain"
	"github.com/icio/mondo/mondohttp"
)

// IterTransactions paginates the Mondo API's transactions endpoint, writing
// the results to the iter channel. The user can interrupt the pagination by
// signalling the kill channel.
func (client *Client) IterTransactions(
	iter chan<- mondodomain.Transaction,
	kill <-chan bool,
	accessToken string,
	accountID string,
	expandMerchants bool,
	since string,
	before string,
	pageLimit int,
) error {
	for {
		// Ensure no kill signal received before making the next request.
		select {
		case <-kill:
			return nil
		default:
		}

		// Request the next set of transactions.
		transColn := new(mondodomain.TransactionsResponse)
		err := client.DoInto(mondohttp.NewTransactionsRequest(accessToken, accountID, expandMerchants, since, before, pageLimit), transColn)
		if err != nil {
			return err
		}

		// Forward the transactions to the output channel.
		for _, tran := range transColn.Transactions {
			select {
			case iter <- tran:
			case <-kill:
				return nil
			}
		}

		// Stop if we've exhausted the available transactions.
		if len(transColn.Transactions) == 0 {
			return nil
		}

		// Update our pointer to the beginning of the next batch.
		since = transColn.Transactions[len(transColn.Transactions)-1].ID
	}
}
