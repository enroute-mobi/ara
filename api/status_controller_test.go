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

func statusPrepareRequest(method string, t *testing.T) (server *Server, request *http.Request, responseRecorder *httptest.ResponseRecorder) {
	server = &Server{}
	var err error
	request, err = http.NewRequest(method, "/_status", nil)
	if err != nil {
		t.Fatal(err)
	}

	responseRecorder = httptest.NewRecorder()
	return
}

func Test_status_check(t *testing.T) {
	server, request, responseRecorder := statusPrepareRequest("GET", t)
	server.HandleFlow(responseRecorder, request)
	statusCheckResponseStatus(responseRecorder, t)
}

func Test_status_check_withApiKey(t *testing.T) {
	server, request, responseRecorder := statusPrepareRequest("GET", t)
	server.apiKey = "dummy"
	server.HandleFlow(responseRecorder, request)
	statusCheckResponseStatus(responseRecorder, t)
}

func Test_status_check_version(t *testing.T) {
	server, request, responseRecorder := statusPrepareRequest("GET", t)
	server.HandleFlow(responseRecorder, request)
	statusCheckResponseStatus(responseRecorder, t)

	status := &Status{}
	body, _ := ioutil.ReadAll(responseRecorder.Body)
	json.Unmarshal(body, status)
	if status.Status != "ok" && status.Version != "" {
		t.Errorf("Status should be ok and vesion not empty  | Status == %v, Version %v", status.Status, status.Version)
	}

}
