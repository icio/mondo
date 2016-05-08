package mondodomain

import (
	"encoding/json"
	"time"
)

// Ping mirrors the response format of /ping requests.
type Ping struct {
	Ping string `json:"ping"`
}

// Token mirrors the response format of /oauth2/token requests.
// https://getmondo.co.uk/docs/#authentication
type Token struct {
	AccessToken  string `json:"access_token"`
	ClientID     string `json:"client_id"`
	ExpiresIn    uint64 `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	UserID       string `json:"user_id"`
}

// Identity mirrors the response format of /ping/whoami requests.
type Identity struct {
	Authenticated bool   `json:"authenticated"`
	ClientID      string `json:"client_id"`
	UserID        string `json:"user_id"`
}

// AccountsResponse mirrors the response format of /accounts requests.
// https://getmondo.co.uk/docs/#list-accounts
type AccountsResponse struct {
	Accounts []Account `json:"accounts"`
}

// Account is the structure of each account listed in /accounts.
// https://getmondo.co.uk/docs/#list-accounts
type Account struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	Created     time.Time `json:"created"`
}

// Balance is the structure of an account balance, mirroring the response format of /balance requests.
// https://getmondo.co.uk/docs/#read-balance
type Balance struct {
	Balance    int    `json:"balance"`
	Currency   string `json:"currency"`
	SpendToday int    `json:"spend_today"`
}

// TransactionResponse mirrors the response format of /transaction requests,
// and utilises field hoisting to directly expose the wrapped Transaction.
type TransactionResponse struct {
	Transaction `json:"transaction"`
}

// TransactionsResponse mirrors the response format of /transactions requests.
type TransactionsResponse struct {
	Transactions []Transaction `json:"transactions"`
}

// Transaction is the structure of a single transaction in /transactions and /transaction.
// https://getmondo.co.uk/docs/#transactions
type Transaction struct {
	ID             string            `json:"id"`
	Created        string            `json:"created"`
	Amount         int               `json:"amount"`
	Currency       string            `json:"currency"`
	AccountBalance int               `json:"account_balance"`
	Merchant       *Merchant         `json:"merchant"`
	Description    string            `json:"description"`
	DeclineReason  string            `json:"decline_reason,omitempty"`
	IsLoad         bool              `json:"is_load"`
	Settled        string            `json:"settled"`
	Metadata       map[string]string `json:"metadata"`
	Notes          string            `json:"notes"`
}

// Merchant is the structure of merchant information on a Transaction. Can be
// Unmarshalled from an object; or a string which becomes the ID value with all
// other fields remaining their zero value.
type Merchant struct {
	ID       string           `json:"id"`
	Name     string           `json:"name"`
	Created  time.Time        `json:"created"`
	Address  *MerchantAddress `json:"address"`
	GroupID  string           `json:"group_id"`
	Logo     string           `json:"logo"`
	Emoji    string           `json:"emoji"`
	Category string           `json:"category"` // Is this also present on the transaction?
}

// rawUnmarshallMerchant is equivalent to Merchant except where decoding json.
//
// json.Unmarshall won't see rawUnmarshallMerchant instances as instances of
// json.Unmarshaller and will just perform its default decoding.
type rawUnmarshallMerchant Merchant

// UnmarshalJSON unmarshals a Merchant from either a JSON object, or from a
// JSON string, e.g. `"merch_123"` -> Merchant{ID: "merch_123"}.
func (m *Merchant) UnmarshalJSON(body []byte) error {
	// Try to decode as a string ID.
	if json.Unmarshal(body, &m.ID) == nil {
		return nil
	}

	// Try to decode a full merchant object.
	return json.Unmarshal(body, (*rawUnmarshallMerchant)(m))
}

// MerchantAddress is the structure of a Merchant's address. Yup.
type MerchantAddress struct {
	Address   string  `json:"address"`
	City      string  `json:"city"`
	Country   string  `json:"country"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
	Postcode  string  `json:"postcode"`
	Region    string  `json:"region"`
}
