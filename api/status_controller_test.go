package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func statusCheckResponseStatus(responseRecorder *httptest.ResponseRecorder, t *testing.T) {
	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code:\n got %v\n want %v",
			status, http.StatusOK)
	}
}

func statusPrepareRequest(method string, t *testing.T) (responseRecorder *httptest.ResponseRecorder) {
	server := &Server{}
	request, err := http.NewRequest(method, "/_status", nil)
	if err != nil {
		t.Fatal(err)

	}

	responseRecorder = httptest.NewRecorder()

	server.HandleFlow(responseRecorder, request)

	return
}

func Test_status_check(t *testing.T) {

	responseRecorder := statusPrepareRequest("GET", t)

	statusCheckResponseStatus(responseRecorder, t)

}

func Test_status_check_version(t *testing.T) {

	responseRecorder := statusPrepareRequest("GET", t)

	statusCheckResponseStatus(responseRecorder, t)

	status := &Status{}
	body, _ := ioutil.ReadAll(responseRecorder.Body)
	json.Unmarshal(body, status)
	if status.Status != "ok" && status.Version != "" {
		t.Errorf("Status should be ok and vesion not empty  | Status == %v, Version %v", status.Status, status.Version)
	}

}
