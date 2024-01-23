package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

/*

	The chan struct{} data type represents a channel in Go that can only send and receive the zero-sized struct {}.

	This type of channel is often used for signaling and synchronization, where the actual data being sent is not relevant, and only the signaling event matters.


	We create an unbuffered channel, shutdownServer, and pass it to the
	function startBadTestHTTPServerV2() as a parameter. Then, inside
	the handler of the test server, we attempt to read from the channel,
	thus creating a potential point of infinitely blocking the execution of
	the handler. Since we do not care about the value inside the channel,
	the type of the channel is the empty struct, struct{}. Replacing the
	time.Sleep() statement via a blocking read operation allows us to
	have more control over the test server operation

*/
func startBadTestServerV2(shutdownServer chan struct{}) *httptest.Server{
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-shutdownServer
		fmt.Fprint(w, "Hello World")
	}))

	return ts
}

func TestFetchRemoteResourcev2(t *testing.T){
	// We create an unbuffered channel, shutdownServer, of type struct{} —an empty struct type
	shutdownServer := make(chan struct{})
	ts := startBadTestServerV2(shutdownServer)
	defer ts.Close()

	// We create a new deferred call to an anonymous function that writes an empty struct value to the channel. This call comes after the ts.Close() call so that it is called before the ts.Close() function.
	defer func(){
		shutdownServer <-struct{}{}
	}()

	client := CreateHTTPClientWithTimeout(200 *time.Millisecond)

	_, err := FetchRemoteResource(client, ts.URL)
	if err == nil{
		t.Log("Expected non-nil Error")
		t.Fail()
	}

	if !strings.Contains(err.Error(), "context deadline exceeded"){
		t.Fatalf("Expected error to contain: context deadline exceeded, Got: %s\n", err.Error())
	}

}
/*
Order of execution:
	1. The time out of 200 millisecond is set for http client, means if the server takes more than 200ms, the request will get timed out.

	2. When fetchRemoteResource() is called, it makes a Get request to the ts server

	3. when it goes to the handler function, it gets blocked at <-shutdownServer, and further execution of this function stops, it remains so, till the timeout is hit.

	4. once timeout is hit, the error is returned to fetchRemoteResource() and it in turns returns error to TestFetchRemoteResourcev2() function.

	5. This error is handled in TestFetchRemoteResourcev2() and now there is no further statement to be executed, so now the time is for the defer statements to be executed in LIFO order,

	6. so, 2nd defer is executed, which writes to the shutdownServer channel, and it unblocks the handler function and the handler's execution resumes to print "Hello World", but this output is likely discarded as the client has already received a timeout error.

	7. at last the server is closed by 2nd defer statement, defer ts.close()



*/

/*
	t.Fatal:
		Use t.Fatal when you want to indicate that a test has failed and stop the execution of the test immediately.
		It marks the test as failed and stops further execution of the test function. Subsequent code in the test function will not be executed.

	t.Fail:
		Use t.Fail when you want to indicate that a test has failed, but you still want the test to continue running.
		It marks the test as failed but does not stop the execution of the test function. It allows the test to proceed to the end and report other failures or errors.

	t.Errorf:
		Use t.Errorf when you want to format and report an error message but still allow the test to continue running.


	Use t.Fail for non-critical failures when you want the test to continue.
	Use t.Fatal for critical failures when you want to stop the test immediately.
	Use t.Errorf for reporting non-fatal errors with formatted error messages.

*/
/*
	In Go, when you have multiple defer statements in a function, they are executed in reverse order, i.e., the last defer statement gets executed first, and the first defer statement gets executed last. This behavior is known as "defer stacking."

*/