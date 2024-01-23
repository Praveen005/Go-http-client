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
		time.Sleep(15 *time.Second) //mimicking, heavy server load, which takes 15s to execute
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
	client := CreateHTTPClientWithTimeout(7 *time.Second)
	data, err := FetchRemoteResource(client, ts.URL)
	if err != nil{
		t.Fatal(err)
	}
	if expected != string(data){
		t.Errorf("Expected: %s, Got: %s", expected, data)
	}
}

/*
	Output:

	=== RUN   TestFetchRemoteResource
    fetchRemoteResource_test.go:30: Get "http://127.0.0.1:58510": context deadline exceeded (Client.Timeout exceeded while awaiting headers)
	2024/01/18 16:15:48 httptest.Server blocked in Close after 5 seconds, waiting for connections:
	*net.TCPConn 0xc000056080 127.0.0.1:58511 in state active
	--- FAIL: TestFetchRemoteResource (15.01s)
	FAIL
	exit status 1
	FAIL    github.com/Praveen005/Go-http-client/tree/main/advanced-http-client/data-downloader-timeout     18.829s

	we can see that the test fails but the execution takes slighly more than 15 secs.
	why?

		we have a call
		to the Close() function of the test server in a deferred call(defer ts.close()). After the
		test function completes execution, the Close() function is called to
		shut down the test server cleanly. However, this function checks to
		see if there are any active requests before shutting down. Hence, it
		only returns when the bad handler returns the response after 60
		seconds.
		-- means jab server respond karne me 7 sec se jyada lagta hai, to timeout ho jata hai, but hamne defer.ts.close() call kiya hai, means jab tak function TestFetchRemoteResource() ka kaam pura nahi hota, ye `ts` server close nahi hoga, aur TestFetchRemoteResource() ka kaam tab pura hoga jab startBadTestServer() response bhejega, aur ye apna response 15 sec ke baad hi bhejega, that is why execution is taking more than 15sec. 

		What can we do about it?
			-- We can rewrite our bad test server as in fetchRemoteResourceNew_test.go
*/