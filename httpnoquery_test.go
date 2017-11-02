package httpnoquery_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/lag13/httpnoquery"
)

// TestClientDo tests that sending the request works as expected.
func TestClientDo(t *testing.T) {
	want := "hello"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(want))
	}))
	defer server.Close()
	tests := []struct {
		name       string
		httpClient *http.Client
	}{
		{
			name:       "default http client",
			httpClient: nil,
		},
		{
			name:       "specific http client",
			httpClient: &http.Client{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := httpnoquery.Client{test.httpClient}
			req, err := http.NewRequest(http.MethodGet, server.URL, nil)
			if err != nil {
				panic(err)
			}
			resp, err := client.Do(httpnoquery.Request{req})
			if err != nil {
				t.Errorf("got non-nil error: %v", err)
			}
			gotBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				panic(err)
			}
			if got := string(gotBody); got != want {
				t.Errorf("got response body %s, wanted %s", got, want)
			}
		})
	}
}

// TestClientDoNoQueryStr tests that when using the client to send a
// request, the error does not contain any query string parameters.
func TestClientDoNoQueryStr(t *testing.T) {
	client := httpnoquery.Client{&http.Client{}}
	queryStr := "?user=hello&password=super-secret-password"
	req, err := http.NewRequest(http.MethodGet, "/path"+queryStr, nil)
	if err != nil {
		panic(err)
	}
	_, err = client.Do(httpnoquery.Request{req})
	if got, want := fmt.Sprintf("%v", err), "sending request: Get /path:"; !strings.Contains(got, want) {
		t.Errorf("error string %s should contain the message %s", got, want)
	}
}
