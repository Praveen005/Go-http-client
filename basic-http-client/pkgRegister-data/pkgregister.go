package pkgregisterdata

import (
	"io"
)

//A multipart message typically consists of multiple parts, each delineated by a boundary string.
//The boundary string is a unique identifier that separates the individual parts within the message.
//In the context of HTTP, a multipart message is commonly used in form submissions that include files (file uploads). When a form contains file input fields, the browser sends a multipart/form-data request to the server. This request format enables the transmission of both text fields and binary files in a single HTTP request.

//Given an object of type pkgData, we can create a multipart message to “package” the data.
type pkgData struct{
	Name		string
	Version		string
	Filename	string  		//The Filename field will store the filename of the package
	Bytes		io.Reader		//pointing to the opened file, The io.Reader interface represents a stream of data that can be read.
}
