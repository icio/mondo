package mondohttp

import (
	"bytes"
	"net/http"
	"strings"
	"testing"
)

func assertReqEquals(t *testing.T, req *http.Request, expected string) {
	buf := new(bytes.Buffer)
	req.Write(buf)
	actual := buf.String()

	expected = strings.Replace(expected, "\n", "\r\n", -1)

	if actual != expected {
		t.Logf("Actual: %q", actual)
		t.Logf("Expect: %q", expected)
		t.Error("Actual != Expected")
	}
}
