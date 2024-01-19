/*

	As we saw, the TestFetchBadRemoteResource test took 60 seconds to run. In fact, if the bad server were to sleep for 600 seconds before it sent back the response, our client code in fetchRemoteResource() would wait the same amount of time. As you can imagine, this will lead to a very bad user experience.

	Let's improve our data downloader function so that it doesn't wait for the response if the server takes more than a specified duration.


	The answer to making our data downloader wait only for a specified maximum period of time is to use a custom HTTP client.

	When we used the http.Get() function, we implicitly used a default HTTP client that is defined in the net/http package

*/

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func CreateHTTPClientWithTimeout(d time.Duration) *http.Client{
	client := http.Client{Timeout: d}
	return &client
}

func FetchRemoteResource(client *http.Client, url string)([]byte, error){
	//make the Get request to the url
	r, err := client.Get(url)
	//handle the error
	
	if err != nil{
		return nil, err
	}
	//close the response body
	// But why? Here's from http.Client documentation:
	/*

	If the returned error is nil, the Response will contain a non-nil Body which the user is expected to close. If the Body is not both read to EOF and closed, the Client's underlying RoundTripper (typically Transport) may not be able to re-use a persistent TCP connection to the server for a subsequent "keep-alive" request

	*/
	defer r.Body.Close()

	return io.ReadAll(r.Body)
}


func main(){
	if len(os.Args) != 2{
		fmt.Fprint(os.Stdout, "Please enter the URL to pull the data from")
		os.Exit(1)
	}
	client := CreateHTTPClientWithTimeout(15 *time.Second)
	data, err := FetchRemoteResource(client ,os.Args[1])
	if err != nil{
		fmt.Fprintf(os.Stdout, "%#v", err) // # produces output under double-quotes
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "%s\n", data)
}