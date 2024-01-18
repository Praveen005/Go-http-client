/*

	Downloading from an overloaded server.

	We will create an always overloaded test HTTP server where every response is delayed by 60 seconds:

*/

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func fetchRemoteResource(url string) ([]byte, error){
	r, err := http.Get(url)
	if err != nil{
		return nil, err
	}
	//close the response
	defer r.Body.Close()

	return io.ReadAll(r.Body) //io.ReadAll returns a slice of byte and an error
}

func main(){
	if len(os.Args) != 2{
		fmt.Fprint(os.Stdout, "Must specify the URL to pull the data from")
		os.Exit(1)
	}
	body, err := fetchRemoteResource(os.Args[1])
	// handle the error
	if err != nil{
		fmt.Fprintf(os.Stdout, "%#v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "The Required data is: %s\n", body)
}

