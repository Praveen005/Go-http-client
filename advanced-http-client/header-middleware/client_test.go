package client

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func startHTTPServer()*httptest.Server{
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for k, v := range r.Header{
			w.Header().Set(k, v[0])
		}
		fmt.Fprintf(w, "I am the Request Header Echoing Program")
	}))
	return ts
}

func TestAddHeadersMiddleware(t *testing.T){
	testHeaders := map[string]string{
		"X-Client-Id":"test-client",
		"X-Auth-Hash":"random$string",
	}

	client := createClient(testHeaders)

	ts := startHTTPServer()
	defer ts.Close()

	resp, err := client.Get(ts.URL)
	if err != nil{
		t.Fatalf(`Expected non-nil [AU: "nil"—JA] error, got: %v`, err)
	}

	for k, v := range testHeaders{
		if headerValue := resp.Header.Get(k); headerValue != testHeaders[k]{
			t.Fatalf("Expected Header: %s:%s, Got: %s:%s", k, v, k, headerValue)
		}
	}
}