package pkgquery

import (
	"encoding/json"
	"io"
	"net/http"
)

type pkgData struct{
	Name 	string 		`json:"name"`
	Version string		`json:"version"`
}


func fetchPackageData(url string) ([]pkgData, error){
	//make an instance of pkgData struct
	var packages []pkgData

	//get the response from the remote url
	resp, err := http.Get(url)
	//error handling
	if err != nil{
		return nil, err
	}

	defer resp.Body.Close()

	//check if the content type mentioned in response header is Json or not
	if resp.Header.Get("Content-Type") != "application/json"{
		return packages, nil
	}

	//now if your program execution is here means, you got a json response
	//so, let's read the body of the response
	data, err := io.ReadAll(resp.Body)
	//handle error
	if err != nil{
		return packages, err
	}

	//now that you have the data, deserialize it to required data-structure, which is pkgData struct here
	//unmarshall function only returns an error
	err = json.Unmarshal(data, &packages)

	return packages, err
}