package mondohttp

import "testing"

func TestAuthCodeAccessRequest(t *testing.T) {
	req := NewAuthCodeAccessRequest("client_id_123", "client_sec_abc", "http://myapp/return", "mondo_auth_code")
	assertReqEquals(t, req, `POST /oauth2/token HTTP/1.1
Host: api.getmondo.co.uk
User-Agent: Go-http-client/1.1
Content-Length: 144
Content-Type: application/x-www-form-urlencoded

client_id=client_id_123&client_secret=client_sec_abc&code=mondo_auth_code&grant_type=authorization_code&redirect_uri=http%3A%2F%2Fmyapp%2Freturn`)
}

func TestRefreshAccessRequest(t *testing.T) {
	req := NewRefreshAccessRequest("client_id_123", "client_sec_abc", "ref_xyz")
	assertReqEquals(t, req, `POST /oauth2/token HTTP/1.1
Host: api.getmondo.co.uk
User-Agent: Go-http-client/1.1
Content-Length: 99
Content-Type: application/x-www-form-urlencoded

client_id=client_id_123&client_secret=client_sec_abc&grant_type=refresh_token&refresh_token=ref_xyz`)
}

func TestWhoAmIRequest(t *testing.T) {
	req := NewWhoAmIRequest("token")
	assertReqEquals(t, req, `GET /ping/whoami HTTP/1.1
Host: api.getmondo.co.uk
User-Agent: Go-http-client/1.1
Authorization: token

`)
}

func TestPingRequest(t *testing.T) {
	assertReqEquals(t, NewPingRequest(), `GET /ping HTTP/1.1
Host: api.getmondo.co.uk
User-Agent: Go-http-client/1.1

`)
}
