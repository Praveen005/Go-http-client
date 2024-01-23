package pkgregisterdata

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
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

// After POST req. we read the response and unmarshal it into the pkgRegisterResult object.
type packageRegisterResult struct{
	Id				string		`json:"id"`
	Filename		string		`json:"filename"`
	Size			int64		`json:"size"`
}


func createMultipartMessage(data pkgData) ([]byte, string, error){
	// bytes.Buffer is a dynamic byte slice that can grow as needed. It provides an efficient way to handle the variable-sized content of a multipart message without needing to pre-allocate a fixed-size buffer.
	var b bytes.Buffer  
	var err error
	var fw io.Writer

	// The multipart.NewWriter function is writes data to a buffer in memory, rather than directly to a destination like a file or network connection. 
	mw := multipart.NewWriter(&b)

	// the method mw.CreateFormField("name") to create a form field object, fw, with the field name "name".
	// bytes.Buffer objects implement io.Writer for writing to in-memory buffers.
	// io.Writer is a powerful interface that enables you to write data to diverse destinations
	fw, err = mw.CreateFormField("name")
	if err != nil{
		return nil, "", err
	}
	// fw is passed to Fprintf to ensure that the formatted string is written to the correct part of the multipart form (the "name" field).
	// fw acts as a mediator between the bytes.Buffer and the multipart.Writer, ensuring that data is written to the correct field within the multipart form context while maintaining proper formatting and structure.
	fmt.Fprint(fw, data.Name)

	fw, err = mw.CreateFormField("version")
	if err != nil{
		return nil, "", err
	}
	fmt.Fprintf(fw, data.Version)

	fw, err = mw.CreateFormFile("filedata", data.Filename)
	if err != nil{
		return nil, "", err
	}
	// The underscore (_) is used to discard the number of bytes written, and the error is stored in the variable err.
	_, err = io.Copy(fw, data.Bytes)
	if err != nil{
		return nil, "", err
	}
	err =mw.Close()
	if err != nil{
		return nil, "", err
	}
	contentType := mw.FormDataContentType()
	return b.Bytes(), contentType, nil

}

