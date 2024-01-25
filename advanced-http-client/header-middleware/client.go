/*
	Let's implement a middleware that will add
	one or more HTTP headers to every outgoing request. we will likely
	end up needing this functionality in various scenarios—sending an
	authentication header, propagating a request ID, and so on.

*/

package client

import "net/http"


type AddHeadersMiddleware struct{
	headers		map[string]string
}

/*

	Encapsulation: Middleware should ideally operate on a copy of the request to maintain encapsulation and avoid unintended side effects on the original request object. This promotes modularity and reusability of middleware components.

	Concurrent Safety: As RoundTrip can be called concurrently, modifying the request in-place could lead to race conditions if multiple goroutines try to access or modify the same request object simultaneously.

	Context Propagation: The r.Context() argument ensures that the context associated with the original request is also copied to the new request, maintaining context-specific values and deadlines.

	Header Addition: The headers are added to the reqCopy object's header, ensuring the original request remains unchanged.

	return http.DefaultTransport.RoundTrip(reqCopy): The modified copy is then passed to the default transport for execution, ensuring the actual HTTP request sent contains the added headers without affecting the original request object.

	This middleware will modify the original request by adding headers
	to it. However, instead of modifying it in place, we clone the request
	using the Clone() method and add headers to it. We then call the
	DefaultTransport's RoundTrip() implementation with the new
	request.


*/

func (h AddHeadersMiddleware) RoundTrip(r *http.Request)(*http.Response, error){
	reqCopy := r.Clone(r.Context())
	for k, v := range h.headers{
		reqCopy.Header.Add(k,v)
	}
	return http.DefaultTransport.RoundTrip(reqCopy)
}


func createClient(headers map[string]string) *http.Client{
	h := AddHeadersMiddleware{
		headers: headers,
	}
	client := http.Client{
		Transport: &h,   // ab naye header ke sath transport hoga request, yayy!!
	}

	return &client
}