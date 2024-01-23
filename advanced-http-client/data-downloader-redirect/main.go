/*

When a server issues an HTTP redirect, the default HTTP client
automatically and silently follows the redirect up to 10 times, after
which it terminates. What if you wanted to change that to, say, follow
no redirects at all, or at least let you know that it is following a
redirect



*/

package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

)

/*

	The custom function to implement the redirect policy must satisfy the following signature:
	func (req *http.Request, via []*http.Request) error

	The first argument, req, is the request to follow the redirect response that it got back from the server

	The slice, via, contains the requests that have been made so far, with the oldest request (your original request) the first element of this slice.


	1. The HTTP client sends a request to the original URL, url .
	2. The server responds with a redirect to, say, url1 .
	3. redirectPolicyFunc is now called with (url1, []{url}) .
	4. If the function returns a nil error, it will follow the redirect and send a new request for url1 .
	5. If there is another redirect to url2, the redirectPolicyFunc function is then called with ( url2, []{url, url1}) .
	6. Steps 3, 4, and 5 are repeated until the redirectPolicyFunc returns a non- nil error
*/

//want to hit this? Take a link, shorten it, and again shorten the shortened one
func redirectPolicyFunc(req *http.Request, via []*http.Request) error {
	fmt.Println("Redirected to: ", req.URL)
	if len(via) > 1 {
		return errors.New("stopped after 1 redirect")
	}
	return nil
}

// How do you hook redirectPolicyFunc() up with a custom HTTP client?

func createHTTPClientWithTimeout(d time.Duration) *http.Client{
	client := http.Client{Timeout: d, CheckRedirect: redirectPolicyFunc}

	return &client
}

func fetchRemoteResource(client *http.Client ,url string)([]byte, error){
	r, err := client.Get(url)
	//handle the error
	if err != nil{
		return nil, err
	}
	// we need to close the response body
	// But why? Here's from http.Client documentation:
	/*

	If the returned error is nil, the Response will contain a non-nil Body which the user is expected to close. If the Body is not both read to EOF and closed, the Client's underlying RoundTripper (typically Transport) may not be able to re-use a persistent TCP connection to the server for a subsequent "keep-alive" request

	*/

	defer r.Body.Close()

	return io.ReadAll(r.Body)  //io.Readall returns []byte and error
}

// Was getting a very large slice of byte in response, hence wrote this function to truncate it
func truncateByteSlice(original []byte, length int) []byte {
    // Check if the original slice is already shorter than the desired length
    if len(original) <= length {
        return original
    }

    // Create a new byte slice with the desired length
    truncated := make([]byte, length)

    // Copy content from the original slice to the truncated slice
    copy(truncated, original)

    return truncated
}

func main(){
	if len(os.Args) != 2{
		fmt.Fprint(os.Stderr, "Please enter a valid URL to pull the data from")
		os.Exit(1)
	}

	client := createHTTPClientWithTimeout(15 *time.Second)
	body, err := fetchRemoteResource(client, os.Args[1])
	//handle the error
	if err != nil{
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	body = truncateByteSlice(body, 20)
	fmt.Fprintf(os.Stdout, "%#v\n", body)
}