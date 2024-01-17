package pkgregisterdata

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func registerPackageData(url string, data pkgData)(packageRegisterResult, error){
	//create an instance of response data
	p := packageRegisterResult{}
	payload, contentType, err := createMultipartMessage(data)
	//handle error
	if err != nil{
		return p, err
	}
	reader := bytes.NewReader(payload)

	resp, err := http.Post(url, contentType, reader)
	//handle error
	if err != nil{
		return p, err
	}

	defer resp.Body.Close()
	respData, err := io.ReadAll(resp.Body)
	//handle error
	if err != nil{
		return p, err
	}

	err= json.Unmarshal(respData, &p)
	return p, err
}
