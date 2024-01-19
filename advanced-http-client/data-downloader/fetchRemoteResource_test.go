package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func startBadTestServer() *httptest.Server{
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		time.Sleep(60 *time.Second)
		fmt.Fprint(w, "Hello World")
	}))

	return ts
}

func TestFetchRemoteResource(t *testing.T){
	// starting an instance of the test server
	ts := startBadTestServer()
	// close the server afater the program execution happens
	defer ts.Close()

	expected := "Hello World"
	data, err := fetchRemoteResource(ts.URL)
	if err != nil{
		t.Fatal(err)
	}
	if expected != string(data){
		t.Errorf("Expected: %s, Got: %s", expected, data)
	}
}
