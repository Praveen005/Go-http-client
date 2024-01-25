/*
	The term middleware (or interceptor) is used for custom code that can be configured to be executed along with the core operation in a network server or client application.

	In a server application, it would be code that gets executed when the server is processing a request from a client.

	In a client application, it would be code that is executed when making an HTTP request to a server application

	In the context of web development, a middleware is a piece of software that sits between a client and a server (or between different components of an application) and performs various functions to enhance, modify, or facilitate the request-response cycle. Middleware functions are executed in a sequential order, allowing each middleware to process the request before it reaches the final handler or the next middleware in the chain.
	Examples:
		API gateways: act as a single entry point for managing API access, authorization, and data transformation for backend services.

		Message brokers: facilitate communication between different applications by routing and delivering messages using asynchronous messaging patterns.

*/
/*

	Understanding the RoundTripper Interface:

	The http.Client struct defines a field, Transport, as follows:

		type Client struct {
			Transport RoundTripper
			...Other fields
		}

	The RoundTripper interface defined in the net/http package defines a
	type that will carry an HTTP request from the client to the remote
	server and carry the response back to the client. The only method
	this type needs to implement is RoundTrip():

		type RoundTripper interface {
			RoundTrip(*Request) (*Response, error)
		}

	[In simple term, RoundTripper ek interface hai, jo ki RoundTrip() method implement karti hai.
	Ab RoundTrip() method kya karti hai, ye yahan Padhen, really nicely explained: https://pkg.go.dev/net/http#RoundTripper]

	When the Transport object is not specified while creating a client, a
	predefined object of type Transport, DefaultTransport is used. It is
	defined as follows (with fields omitted):
		var DefaultTransport RoundTripper = &Transport{
			..fields omitted
		}

		var DefaultTransport RoundTripper: Declares a variable named DefaultTransport of type RoundTripper

		= &Transport{...}: Assigns a new Transport struct to DefaultTransport. This struct contains configuration for how HTTP requests should be made.

	The Transport type defined in the net/http package implements the
	RoundTrip() method as required by the RoundTripper interface. It is
	responsible for creating and managing the underlying Transmission
	Control Protocol (TCP) connections over which an HTTP request-response transaction occurs:

		1. You create a Client object.
		2. You create an HTTP Request .
		3. The HTTP request is then “carried over” the RoundTripper
		implementation (for instance, over a TCP connection) to the
		server, and the response is carried back.
		4. If you make more than one request with the same client, step 2
		and step 3 will be repeated.


	To implement a client middleware, we will write a custom type that
	will encapsulate the DefaultTransport's RoundTripper
	implementation.

	For control over proxies, TLS configuration, keep-alives, compression, and other settings, create a Transport.

	Clients and Transports are safe for concurrent use by multiple goroutines and for efficiency should only be created once and re-used.

	Refer: https://pkg.go.dev/net/http#hdr-Clients_and_Transports  for more

*/
/*

	The first middleware that you will write will log a message before sending a request. It will log another message when a response is received.




*/

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)


type LoggingClient struct{
	log *log.Logger
}

// To satisfy the RoundTripper interface, we implement the RoundTrip() method:
//The RoundTripper must not mutate(modify) the request or response or return an error. 
// If there is an error during the round-trip (e.g., network error), it should be communicated through the Response object rather than returning an error. 
func(c LoggingClient) RoundTrip(r *http.Request)(*http.Response, error){
	c.log.Printf("Sending a %s request to %s over %s\n", r.Method, r.URL, r.Proto)

	// http.DefaultTransport, is the standard RoundTripper used by Go's http.Client
	resp, err := http.DefaultTransport.RoundTrip(r)
	c.log.Printf("Got back a response over %s\n", resp.Proto)
	return resp, err
}

func createHTTPClientWithTimeout(d time.Duration) *http.Client{
	client := http.Client{Timeout: d}
	return &client
}

func fetchRemoteResource(client *http.Client, url string)([]byte, error){
	r, err := client.Get(url)
	if err != nil{
		return nil, err
	}
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil{
		return nil, err
	}

	return body, nil
}

func main(){
	if len(os.Args) != 2{
		fmt.Fprintf(os.Stdout, "Must specify the URL to get data from\n")
		os.Exit(1)
	}

	myTransport := LoggingClient{}
	l := log.New(os.Stdout, "", log.LstdFlags)  //new log.Logger object
	myTransport.log = l

	client := createHTTPClientWithTimeout(15 *time.Second)
	client.Transport =&myTransport

	body, err := fetchRemoteResource(client, os.Args[1])
	if err != nil{
		fmt.Fprintf(os.Stdout, "%#v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "Bytes in Response: %d\n", len(body))
}
 
/*
	we create a new
	LoggingClient object. Then, we create a new log.Logger object by
	calling the log.New() function. 
	
	The first parameter to the function is
	a io.Writer object to which the logs will be written. Here we use
	os.Stdout. 
	
	The second parameter to the function is the prefix string
	to add to each log statement—an empty string is specified here. 
	
	The last parameter to the function is a flag—text to prefix to each log line.
	Here we use log.LstdFlags, which will display the date and time. We
	then assign the log.Logger object to the l field of the myTransport
	object. 
	
	Finally, we set client.Transport to &myTransport .



*/

