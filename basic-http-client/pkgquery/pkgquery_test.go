package pkgquery

import(
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)
// To test the fetchPackageData(), we need to create a test server
func startTestPackageServer() *httptest.Server{
	pkgData := `[
		{
			"name":"package1",
			"version":"1.1"
		},{
			"name":"package2",
			"version":"1.0"
		}
	]`

	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, pkgData)
			},
		),
	)
	return ts
}

func TestFetchPackageData(t *testing.T){
	//start the test server
	ts := startTestPackageServer()
	//close the server upon completion of objective
	defer ts.Close()

	//get the data from the server
	packages, err := fetchPackageData(ts.URL)
	//handle the error
	if err != nil{
		t.Fatal(err)
	}

	if len(packages) != 2{
		t.Fatalf("Expected 2 packages, Got back: %d", len(packages))
	}
}