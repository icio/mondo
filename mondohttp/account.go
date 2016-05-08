package mondohttp

import (
	"net/http"
	"net/url"
	"strings"
)

// NewAccountsRequest creates a request for a listing of the user's accounts.
// https://getmondo.co.uk/docs/#list-accounts.
func NewAccountsRequest(accessToken string) *http.Request {
	req, _ := http.NewRequest("GET", ProductionAPI+"accounts", nil)
	req.Header.Set(auth(accessToken))
	return req
}

// NewBalanceRequest creates a request for an account's current balance.
// https://getmondo.co.uk/docs/#read-balance.
func NewBalanceRequest(accessToken, accountID string) *http.Request {
	req, _ := http.NewRequest("GET", ProductionAPI+"balance?account_id="+url.QueryEscape(accountID), nil)
	req.Header.Set(auth(accessToken))
	return req
}

// NewTransactionRequest creates a request for a single transactions.
// https://getmondo.co.uk/docs/#retrieve-transaction
func NewTransactionRequest(accessToken, transactionID string, expandMerchants bool) *http.Request {
	var query string
	if expandMerchants {
		query = "?expand%5B%5D=merchant"
	}

	req, _ := http.NewRequest("GET", ProductionAPI+"transactions/"+transactionID+query, nil)
	req.Header.Set(auth(accessToken))
	return req
}

// NewTransactionsRequest creates a request for a series of account transactions.
// https://getmondo.co.uk/docs/#list-transactions
func NewTransactionsRequest(accessToken, accountID string, expandMerchants bool, since, before string, limit int) *http.Request {
	query := &url.Values{
		"account_id": {accountID},
	}
	appendPaginationParams(query, since, before, limit)
	if expandMerchants {
		query.Add("expand[]", "merchant")
	}

	req, _ := http.NewRequest("GET", ProductionAPI+"transactions?"+query.Encode(), nil)
	req.Header.Set(auth(accessToken))
	return req
}

// NewAnnotateTransactionRequest creates a request for updating annotations on a transactions.
// https://getmondo.co.uk/docs/#annotate-transaction
func NewAnnotateTransactionRequest(accessToken, transactionID string, metadata map[string]string) *http.Request {
	body := &url.Values{}
	appendQueryMap(body, "metadata[", "]", metadata)

	req, _ := http.NewRequest("PATCH", ProductionAPI+"transactions/"+transactionID, strings.NewReader(body.Encode()))
	req.Header.Set(formContentType())
	req.Header.Set(auth(accessToken))
	return req
}

// NewCreateFeedItemRequest creates a request for adding a feed item to an account.
// https://getmondo.co.uk/docs/#create-feed-item
func NewCreateFeedItemRequest(accessToken, accountID, itemType, itemURL string, params map[string]string) *http.Request {
	body := &url.Values{
		"account_id": {accountID},
		"type":       {itemType},
	}
	if itemURL != "" {
		body.Set("url", itemURL)
	}
	appendQueryMap(body, "params[", "]", params)

	req, _ := http.NewRequest("POST", ProductionAPI+"feed", strings.NewReader(body.Encode()))
	req.Header.Set(formContentType())
	req.Header.Set(auth(accessToken))
	return req
}

// NewCreateURLFeedItemRequest is shorthand for creating a request to add a
// basic account feed item with only a URL, title, and image.
// https://getmondo.co.uk/docs/#create-feed-item
func NewCreateURLFeedItemRequest(accessToken, accountID, url, title, imageURL string) *http.Request {
	return NewCreateFeedItemRequest(accessToken, accountID, "basic", url, map[string]string{"image_url": imageURL, "title": title})
}

// NewCreateBasicFeedItemRequest is shorthand for creating a request to add a
// basic account feed item.
// https://getmondo.co.uk/docs/#create-feed-item
func NewCreateBasicFeedItemRequest(accessToken, accountID, url, title, imageURL, body, backgroundColor, titleColor, bodyColor string) *http.Request {
	params := map[string]string{
		"title":     title,
		"image_url": imageURL,
	}
	if body != "" {
		params["body"] = body
	}
	if backgroundColor != "" {
		params["background_color"] = backgroundColor
	}
	if titleColor != "" {
		params["title_color"] = titleColor
	}
	if bodyColor != "" {
		params["body_color"] = bodyColor
	}

	return NewCreateFeedItemRequest(accessToken, accountID, "basic", url, params)
}

// TODO: https://getmondo.co.uk/docs/#webhooks
// TODO: https://getmondo.co.uk/docs/#attachments
