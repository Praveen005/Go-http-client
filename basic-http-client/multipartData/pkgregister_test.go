package pkgregisterdata

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)


func packageRegHandler(w http.ResponseWriter, r *http.Request){
	if r.Method == "POST"{
		// Response that is to be returned upon successful package registration
		d := packageRegisterResult{}
		// parses the incoming request as multipart form data.
		// The 5000 argument sets a maximum memory limit(in bytes) for parsing
		err := r.ParseMultipartForm(5000)
		if err != nil{
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// accesses the parsed form data.
		mForm := r.MultipartForm
		// Get File data: retrieves the first uploaded file from the "filedata" field.
		f := mForm.File["filedata"][0]

		// Constructs a package ID using values from the "name" and "version" fields.
		d.Id = fmt.Sprintf("%s-%s", mForm.Value["name"][0], mForm.Value["version"][0])
		//stores the uploaded file's name
		d.Filename = f.Filename
		// stores the uploaded file's size.
		d.Size = f.Size

		// marshals the packageRegisterResult struct (containing ID, filename, and size) into JSON format.
		jsonData, err := json.Marshal(d)
		if err != nil{
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(jsonData))
		
	}else{
		http.Error(w, "Invalid HTTP method specified", http.StatusMethodNotAllowed)
	}
}

// Starting a test server to test resgistration of package and handling of the multipart message

func startTestPackageServer() *httptest.Server{
	ts := httptest.NewServer(http.HandlerFunc(packageRegHandler))
	return ts
}

func TestRegisterPackageData(t * testing.T){
	//start a server instance
	ts := startTestPackageServer()

	//close the server upon the completion of objective
	defer ts.Close()

	p := pkgData{
		Name:		"mypackage",
		Version: 	"0.1",
		Filename: 	"mypackage-0.1.tar.gz",
		Bytes:		strings.NewReader("data"),
	}

	pResult, err := registerPackageData(ts.URL, p)
	//handle the error
	if err != nil{
		t.Fatal(err)
	}

	if pResult.Id != fmt.Sprintf("%s-%s", p.Name, p.Version){
		t.Errorf("Expected package ID to be %s-%s, Got: %s", p.Name, p.Version, pResult.Id)
	}

	if pResult.Filename != p.Filename{
		t.Errorf("Expected package filename to be %s, Got: %s", p.Filename, pResult.Filename)
	}

	if pResult.Size != 4{
		t.Errorf("Expected package size to be 4, Got: %d", pResult.Size)
	}
}