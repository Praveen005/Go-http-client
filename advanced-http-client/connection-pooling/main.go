package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/http/httptrace"
	"os"
	"time"

)

/*

	In a nutshell, the program is trying to create an HTTP GET request with additional tracing information related to DNS resolution and connection establishment. The tracing functions print information about DNS resolution and connection details during the execution of the request.

	The `net/http/httptrace` package will help us delve into the internals of connection pooling. One of the things that we can see using this package is whether a connection is being reused or whether a new one was established for making an HTTP request.


*/

func createHTTPGetRequestWithTrace(ctx context.Context, url string)(*http.Request, error){
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil{
		return nil, err
	}
	// This allows you to trace events during an HTTP request.
	/*
		The httptrace.ClientTrace struct type defines functions that will be
		called when certain events in a request's life cycle happen.

		We are interested in two events here:

		The DNSDone event happens when the DNS lookup for a
		hostname has completed.

		The GotConn event happens when a connection has been
		obtained and now you can send the request over
	*/
	// Read More:  https://pkg.go.dev/net/http/httptrace#WithClientTrace
	trace := &httptrace.ClientTrace{
		// This function is called when DNS resolution is completed. It prints information about the DNS resolution, such as the IP addresses associated with the resolved domain.
		// DNS resolution refers to the process of translating a human-readable domain name (like www.example.com) into an IP address.
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			fmt.Printf("DNS Info: %+v\n", dnsInfo)
		},
		GotConn: func(connInfo httptrace.GotConnInfo) {
			fmt.Printf("Got Conn: %+v\n", connInfo)
		},
		TLSHandshakeStart: func() {
			fmt.Printf("TLS HandShake Start\n")
		},
		TLSHandshakeDone: func(connState tls.ConnectionState, err error) {
			fmt.Printf("TLS HandShake Done\n")
		},

		PutIdleConn: func(err error) {
			fmt.Printf("Put Idle Conn Error: %+v\n", err)
		},
		ConnectStart: func(network, addr string) {
			fmt.Printf("Connecting to %s network at address %s\n", network, addr)
		},
		ConnectDone: func(network, addr string, err error){
			if err != nil{
				fmt.Printf("Connection to %s network at address %s failed with error: %s\n", network, addr, err)
			}else{
				fmt.Printf("Connection to %s network at address %s succeeded\n", network, addr)
			}
		},
	}
	// creating a new context (ctxTrace) derived from the original context associated with the HTTP request (req.Context()), but with an additional client trace (trace) attached to it.
	// This new context (ctxTrace) is then used in subsequent operations related to the HTTP request.
	ctxTrace := httptrace.WithClientTrace(req.Context(), trace)

	req = req.WithContext(ctxTrace)
	return req, nil
}


func createClientWithTimeout(d time.Duration)(*http.Client){
	client := http.Client{Timeout: d}
	return &client
}

func main(){
	if len(os.Args) != 2{
		fmt.Println("Please enter an URL to pull the data from.")
		os.Exit(1)
	}

	d := 5 *time.Second
	ctx := context.Background()
	client := createClientWithTimeout(d)
	/*

		With the above configuration, an idle connection will be kept around for a maximum of 10 seconds. Hence, if you make two requests using a client with an interval of 11 seconds in between them, the second request will trigger a new connection to be created.

	*/
	transport := &http.Transport{
		IdleConnTimeout: 10 *time.Second,
		MaxIdleConns: 15,
		MaxIdleConnsPerHost: 3,
	}
	client.Transport = transport
	
	req, err := createHTTPGetRequestWithTrace(ctx, os.Args[1])

	if err != nil{
		log.Fatal(err)
	}

	// Use ctlr + c to terminate this infinite loop
	for{
		client.Do(req)
		time.Sleep(1 *time.Second)
		fmt.Println("----------------------")
	}
}
/*

	what happens to our connection pool implementation when an underlying
	IP address with which a connection was established is no longer
	available? Will the connection pool implementation realize that it is
	no longer available and create a new connection to the new IP
	address?

	Yes, in fact, it will. When you attempt to make a new
	request, a new connection will be opened to the remote server.

*/