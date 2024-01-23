/*
	We have seen how to create a custom HTTP client

	Furthermore, we have used methods such as Get() on the Client object to make requests

	Underneath, the client is using a default request object of type http.Request that is defined in the standard library.

	Customizing the http.Request object allows you to add headers or cookies or simply set the time-out for a request.

	Creating a new request is done by calling the NewRequest() function

	The NewRequestWithContext() function has the exact same purpose, but additionally it allows passing a context to the request
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)

		The first argument to the function is the context object.

		The second parameter is the HTTP method, for which we are creating the request

		url points to the URL of the resource to which we are going to make the request

		The last argument is an io.Reader object pointing to the body, which in the case of a GET request will in most cases likely be empty

		To create a request for a POST request, with io.Reader and body, you would make the following function call:
			req, err := http.NewRequestWithContext(ctx, "POST", url, body)

		We can send body(data) through get request as well, but that practice is not recommended.

	Once you create the request object, you can then add a header using the following:
		req.Header().Add("X-AUTH-HASH", "authhash")

		This will add a header X-AUTH-HASH with the value authhash to the outgoing request.

	You can encapsulate this logic in a function that creates a custom http.Request object for making a GET request with headers:


*/
/*
	Context:


		In Go, the context package provides a powerful tool for managing concurrent operations, cancellation, and deadlines. It's crucial for several aspects of Go programming, especially when dealing with goroutines, web requests, and other asynchronous tasks.

		Components:

			context.Context interface:
				Defines the fundamental type representing a context. It holds values like deadlines, cancellation signals, and optional key-value pairs.

			Methods:
				Context offers methods for accessing deadlines, checking cancellation state, retrieving values, and creating child contexts.

			Context Propagation:
				Contexts can be passed down through function calls, ensuring all operations within a specific scope share the same information.

		Key Features:

			Cancellation:
				Allows cancelling ongoing operations gracefully, preventing them from running indefinitely.

			Deadlines:
				Sets timeouts for operations, ensuring they terminate within a specified duration.

			Value Sharing:
				Enables efficient sharing of additional information through key-value pairs within a context.

			Error Handling:
				Simplifies error handling related to cancelled contexts or missed deadlines.


		Read More: https://pkg.go.dev/context
*/
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


	When the Transport object is not specified while creating a client, a
	predefined object of type Transport, DefaultTransport is used. It is
	defined as follows (with fields omitted):
		var DefaultTransport RoundTripper = &Transport{
			..fields omitted
		}


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

package client

import (
	"context"
	"log"
	"net/http"
)

type LoggingClient struct{
	log *log.Logger
}

// To satisfy the RoundTripper interface, we implement the RoundTrip() method:
func(c LoggingClient) RoundTrip(r *http.Request)(*http.Response, error){
	c.log.Printf("Sending a %s request to %s over %s\n", r.Method, r.URL, r.Proto)

	resp, err := http.DefaultTransport.RoundTrip(r)
	c.log.Printf("Got back a response over %s\n", resp.Proto)
	return resp, err
}
 
func createHTTPGetRequest(ctx context.Context, url string, headers map[string]string) (*http.Request, error){
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	//handle the error
	if err != nil{
		return nil, err
	}

	for k, v := range headers{
		req.Header.Add(k,v)
	}
	return req, err
}

