package pkgregister

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func packageRegHandler(w http.ResponseWriter, r* http.Request){
	if r.Method =="POST"{
		//incoming package data
		p := pkgData{}

		//package registration response
		d := pkgRegisterResult{}
		defer r.Body.Close()

		data, err := io.ReadAll(r.Body)
		//handle error
		if err != nil{
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = json.Unmarshal(data, &p)
		if err != nil || len(p.Name)==0 ||len(p.Version)==0{
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		d.ID =p.Name+ "-" + p.Version
		jsonData, err := json.Marshal(d)
		//handle error
		if err != nil{
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		//if marshalling does not throw any error, set the header of the response
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(jsonData))
	}else{
		http.Error(w, "Invalid HTTP Method specified", http.StatusMethodNotAllowed)
		return
	}
}

//make a test server to check if or not the post req. goes through properly
//above function handles the post request, the function below, returns a server
func startTestPackageServer() *httptest.Server{
	//creating a server instance
	ts := httptest.NewServer(http.HandlerFunc(packageRegHandler))
	return ts
}

func TestRegisterPackageData(t *testing.T){
	ts := startTestPackageServer()

	defer ts.Close()

	p := pkgData{
		Name: "MyPackage",
		Version: "0.1",
	}

	resp, err := registerPackageData(ts.URL, p)
	//handle the error
	if err != nil{
		t.Fatal(err)
	}

	if resp.ID != "MyPackage-0.1"{
		t.Errorf("Expected package id to be MyPackage-0.1, Got: %s", resp.ID)
	}
}

func TestRegisterEmptyPackageData(t *testing.T){
	ts := startTestPackageServer()
	defer ts.Close()

	p := pkgData{}

	resp, err := registerPackageData(ts.URL, p)
	if err == nil{
		t.Fatal("Expected error to be non-nil, got nil")
	}

	if len(resp.ID) != 0{
		t.Errorf("Expected package ID to be empty, got : %s", resp.ID)
	}
}