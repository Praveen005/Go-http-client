package pkgregister

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type pkgData struct{
	Name 		string		`json:"name"`
	Version		string		`json:"version"`
}

type pkgRegisterResult struct{
	ID string	`json:"id"`
}

func registerPackageData(url string, data pkgData) (pkgRegisterResult, error){
	//make an instance of pkgRegisterResult
	p := pkgRegisterResult{}
	//serialize the data as JSON to send it to the server as request body
	b, err := json.Marshal(data)
	//handle the error
	if err != nil{
		return p, err
	}
	//we will create an io.Reader object for this byte slice(returned by json.Marshal) using the NewReader() function from the bytes package.
	//The io.Reader interface represents a stream of data that can be read.
	//By accepting an io.Reader, the Post function allows you to provide the request body from different sources: a file, a network connection, or any other source that implements the io.Reader interface.
	//Using an io.Reader allows for streaming data, which is useful when dealing with large datasets. You can read and send the data in chunks without loading the entire payload into memory.
	//If you have a large payload, creating a bytes.Reader from the byte slice allows you to avoid loading the entire payload into memory at once.
	reader := bytes.NewReader(b)
	response, err := http.Post(url, "application/json", reader)
	//handle the error
	if err != nil{
		return p, err
	}

	//close the response body, after the completion of objective
	defer response.Body.Close()

	responseData, err := io.ReadAll(response.Body)
	if err != nil{
		return p, err
	}

	if response.StatusCode != http.StatusOK{
		return p, errors.New(string(responseData))
	}

	err = json.Unmarshal(responseData, &p)
	return p, err

}