package mondohttp

import "testing"

func TestNewAccountsRequest(t *testing.T) {
	req := NewAccountsRequest("token")
	assertReqEquals(t, req, `GET /accounts HTTP/1.1
Host: api.getmondo.co.uk
User-Agent: Go-http-client/1.1
Authorization: token

`)
}

func TestBalanceRequest(t *testing.T) {
	req := NewBalanceRequest("token", "acc_123")
	assertReqEquals(t, req, `GET /balance?account_id=acc_123 HTTP/1.1
Host: api.getmondo.co.uk
User-Agent: Go-http-client/1.1
Authorization: token

`)
}

func TestTransactionRequest(t *testing.T) {
	req := NewTransactionRequest("token", "trans_456", false)
	assertReqEquals(t, req, `GET /transactions/trans_456 HTTP/1.1
Host: api.getmondo.co.uk
User-Agent: Go-http-client/1.1
Authorization: token

`)
}

func TestTransactionRequest_Merchants(t *testing.T) {
	req := NewTransactionRequest("token", "trans_456", true)
	assertReqEquals(t, req, `GET /transactions/trans_456?expand%5B%5D=merchant HTTP/1.1
Host: api.getmondo.co.uk
User-Agent: Go-http-client/1.1
Authorization: token

`)
}

func TestTransactionsRequest_NoPage(t *testing.T) {
	req := NewTransactionsRequest("token", "acc_123", true, "", "", 0)
	assertReqEquals(t, req, `GET /transactions?account_id=acc_123&expand%5B%5D=merchant HTTP/1.1
Host: api.getmondo.co.uk
User-Agent: Go-http-client/1.1
Authorization: token

`)
}

func TestTransactionsRequest_Page(t *testing.T) {
	req := NewTransactionsRequest("token", "acc_123", true, "start", "end", 50)
	assertReqEquals(t, req, `GET /transactions?account_id=acc_123&before=end&expand%5B%5D=merchant&limit=50&since=start HTTP/1.1
Host: api.getmondo.co.uk
User-Agent: Go-http-client/1.1
Authorization: token

`)
}

func TestAnnotateTransactionRequest(t *testing.T) {
	req := NewAnnotateTransactionRequest("token", "trans_456", map[string]string{
		"test_a": "abc",
		"test_b": "",
	})
	assertReqEquals(t, req, `PATCH /transactions/trans_456 HTTP/1.1
Host: api.getmondo.co.uk
User-Agent: Go-http-client/1.1
Content-Length: 46
Authorization: token
Content-Type: application/x-www-form-urlencoded

metadata%5Btest_a%5D=abc&metadata%5Btest_b%5D=`)
}

func TestCreateURLFeedItemRequest(t *testing.T) {
	req := NewCreateURLFeedItemRequest("token", "acc_123", "https://www.google.com", "My feed item", "http://test.com/image.png")
	assertReqEquals(t, req, `POST /feed HTTP/1.1
Host: api.getmondo.co.uk
User-Agent: Go-http-client/1.1
Content-Length: 149
Authorization: token
Content-Type: application/x-www-form-urlencoded

account_id=acc_123&params%5Bimage_url%5D=http%3A%2F%2Ftest.com%2Fimage.png&params%5Btitle%5D=My+feed+item&type=basic&url=https%3A%2F%2Fwww.google.com`)
}

func TestCreateBasicFeedItemRequest_Minimum(t *testing.T) {
	req := NewCreateBasicFeedItemRequest("token", "acc_123", "", "My feed item", "http://test.com/image.png", "You've created a feed item!", "", "", "")
	assertReqEquals(t, req, `POST /feed HTTP/1.1
Host: api.getmondo.co.uk
User-Agent: Go-http-client/1.1
Content-Length: 165
Authorization: token
Content-Type: application/x-www-form-urlencoded

account_id=acc_123&params%5Bbody%5D=You%27ve+created+a+feed+item%21&params%5Bimage_url%5D=http%3A%2F%2Ftest.com%2Fimage.png&params%5Btitle%5D=My+feed+item&type=basic`)
}

func TestCreateBasicFeedItemRequest_Maximum(t *testing.T) {
	req := NewCreateBasicFeedItemRequest("token", "acc_123", "https://override.com/", "My feed item", "http://test.com/image.png", "You've created a feed item!", "bg-color", "h1-color", "p-color")
	assertReqEquals(t, req, `POST /feed HTTP/1.1
Host: api.getmondo.co.uk
User-Agent: Go-http-client/1.1
Content-Length: 301
Authorization: token
Content-Type: application/x-www-form-urlencoded

account_id=acc_123&params%5Bbackground_color%5D=bg-color&params%5Bbody%5D=You%27ve+created+a+feed+item%21&params%5Bbody_color%5D=p-color&params%5Bimage_url%5D=http%3A%2F%2Ftest.com%2Fimage.png&params%5Btitle%5D=My+feed+item&params%5Btitle_color%5D=h1-color&type=basic&url=https%3A%2F%2Foverride.com%2F`)
}
