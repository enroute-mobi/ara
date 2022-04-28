package api

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/core"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

func checkStopVisitResponseStatus(responseRecorder *httptest.ResponseRecorder, t *testing.T) {
	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code:\n got %v\n want %v",
			status, http.StatusOK)
	}

	if contentType := responseRecorder.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Handler returned wrong Content-Type:\n got: %v\n want: %v",
			contentType, "application/json")
	}
}

func prepareStopVisitRequest(method string, sendIdentifier bool, body []byte, t *testing.T) (stopVisit *model.StopVisit, responseRecorder *httptest.ResponseRecorder, referential *core.Referential) {
	// Create a referential
	referentials := core.NewMemoryReferentials()
	server := &Server{}
	server.SetReferentials(referentials)
	referential = referentials.New("default")
	referential.Tokens = []string{"testToken"}
	referential.Save()

	// Set the fake UUID generator
	uuid.SetDefaultUUIDGenerator(uuid.NewFakeUUIDGenerator())
	// Save a new stopVisit
	timeLayout := "2006/01/02-15:04:05"
	stopVisit = model.NewStopVisit(referential.Model())

	svTime, _ := time.Parse(timeLayout, "2007/01/02-15:04:05")
	stopVisit.Schedules.SetArrivalTime("actual", svTime)
	referential.Model().StopVisits().Save(stopVisit)

	stopVisit2 := model.NewStopVisit(referential.Model())
	svTime, _ = time.Parse(timeLayout, "2005/01/02-15:04:05")
	stopVisit2.Schedules.SetArrivalTime("actual", svTime)
	referential.Model().StopVisits().Save(stopVisit2)

	url := "/default/stop_visits"
	if method == "GET" && !sendIdentifier {
		url = url + "?After=2006/01/02-15:04:05"
	}
	address := []byte(url)
	if sendIdentifier {
		address = append(address, fmt.Sprintf("/%s", stopVisit.Id())...)
	}
	request, err := http.NewRequest(method, string(address), bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	request.Header.Set("Authorization", "Token token=testToken")

	// Create a ResponseRecorder
	responseRecorder = httptest.NewRecorder()

	// Call HandleFlow method and pass in our Request and ResponseRecorder.
	server.HandleFlow(responseRecorder, request)

	return
}

func Test_StopVisitController_Delete(t *testing.T) {
	// Send request
	stopVisit, responseRecorder, referential := prepareStopVisitRequest("DELETE", true, nil, t)

	// Test response
	checkStopVisitResponseStatus(responseRecorder, t)

	//Test Results
	_, ok := referential.Model().StopVisits().Find(stopVisit.Id())
	if ok {
		t.Errorf("StopVisit shouldn't be found after DELETE request")
	}
	if expected, _ := stopVisit.MarshalJSON(); responseRecorder.Body.String() != string(expected) {
		t.Errorf("Wrong body for DELETE response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}

func Test_StopVisitController_Update(t *testing.T) {
	// Prepare and send request
	body := []byte(`{ "ObjectIDs": { "reflex": "FR:77491:ZDE:34004:STIF" } }`)
	stopVisit, responseRecorder, referential := prepareStopVisitRequest("PUT", true, body, t)

	// Check response
	checkStopVisitResponseStatus(responseRecorder, t)

	// Test Results
	updatedStopVisit, ok := referential.Model().StopVisits().Find(stopVisit.Id())
	if !ok {
		t.Errorf("StopVisit should be found after PUT request")
	}

	if expected, _ := updatedStopVisit.MarshalJSON(); responseRecorder.Body.String() != string(expected) {
		t.Errorf("Wrong body for PUT response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}

func Test_StopVisitController_Show(t *testing.T) {
	// Send request
	stopVisit, responseRecorder, _ := prepareStopVisitRequest("GET", true, nil, t)

	// Test response
	checkStopVisitResponseStatus(responseRecorder, t)

	//Test Results
	if expected, _ := stopVisit.MarshalJSON(); responseRecorder.Body.String() != string(expected) {
		t.Errorf("Wrong body for GET (show) response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}

func Test_StopVisitController_Create(t *testing.T) {
	// Prepare and send request
	body := []byte(`{ "ObjectIDs": { "reflex": "FR:77491:ZDE:34004:STIF" } }`)
	_, responseRecorder, referential := prepareStopVisitRequest("POST", false, body, t)

	// Check response
	checkStopVisitResponseStatus(responseRecorder, t)

	// Test Results
	// Using the fake uuid generator, the uuid of the created
	// stopVisit should be 6ba7b814-9dad-11d1-2-00c04fd430c8
	stopVisit, ok := referential.Model().StopVisits().Find("6ba7b814-9dad-11d1-2-00c04fd430c8")
	if !ok {
		t.Errorf("StopVisit should be found after POST request")
	}
	if expected, _ := stopVisit.MarshalJSON(); responseRecorder.Body.String() != string(expected) {
		t.Errorf("Wrong body for POST response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}

func Test_StopVisitController_Index(t *testing.T) {
	// Send request
	_, responseRecorder, _ := prepareStopVisitRequest("GET", false, nil, t)

	// Test response
	checkStopVisitResponseStatus(responseRecorder, t)

	//Test Results

	expected := `[{"Id":"6ba7b814-9dad-11d1-0-00c04fd430c8","Collected":false,"Schedules":[{"ArrivalTime":"2007-01-02T15:04:05Z","Kind":"actual"}],"Origin":"","VehicleAtStop":false}]`
	if responseRecorder.Body.String() != string(expected) {
		t.Errorf("Wrong body for GET (index) response request:\n got: %v\n want: %v", responseRecorder.Body.String(), string(expected))
	}
}

func Test_StopVisitController_FindStopVisit(t *testing.T) {
	ref := core.NewMemoryReferentials().New("test")
	stopVisit := ref.Model().StopVisits().New()
	objectid := model.NewObjectID("kind", "stif:value")
	stopVisit.SetObjectID(objectid)
	ref.Model().StopVisits().Save(&stopVisit)

	controller := &StopVisitController{
		referential: ref,
	}

	_, ok := controller.findStopVisit("kind:stif:value")
	if !ok {
		t.Error("Can't find StopVisit by ObjectId")
	}

	_, ok = controller.findStopVisit(string(stopVisit.Id()))
	if !ok {
		t.Error("Can't find StopVisit by Id")
	}
}

func benchmarkStopVisitsIndex(sv int, b *testing.B) {
	referentials := core.NewMemoryReferentials()
	server := &Server{}
	server.SetReferentials(referentials)
	referential := referentials.New("default")
	referential.Tokens = []string{"testToken"}
	referential.Save()

	for i := 0; i != sv; i++ {
		line := referential.Model().Lines().New()
		line.Save()

		vehicleJourney := referential.Model().VehicleJourneys().New()
		vehicleJourney.LineId = line.Id()
		vehicleJourney.Save()

		stopVisit := referential.Model().StopVisits().New()
		stopVisit.VehicleJourneyId = vehicleJourney.Id()
		stopVisit.Save()
	}

	for n := 0; n < b.N; n++ {
		request, err := http.NewRequest("GET", "/default/stop_visits", bytes.NewReader(nil))
		if err != nil {
			b.Fatal(err)
		}

		request.Header.Set("Authorization", "Token token=testToken")

		responseRecorder := httptest.NewRecorder()

		server.HandleFlow(responseRecorder, request)
	}
}

func BenchmarkStopVisitsIndex10(b *testing.B)    { benchmarkStopVisitsIndex(10, b) }
func BenchmarkStopVisitsIndex50(b *testing.B)    { benchmarkStopVisitsIndex(50, b) }
func BenchmarkStopVisitsIndex100(b *testing.B)   { benchmarkStopVisitsIndex(100, b) }
func BenchmarkStopVisitsIndex1000(b *testing.B)  { benchmarkStopVisitsIndex(1000, b) }
func BenchmarkStopVisitsIndex5000(b *testing.B)  { benchmarkStopVisitsIndex(5000, b) }
func BenchmarkStopVisitsIndex10000(b *testing.B) { benchmarkStopVisitsIndex(10000, b) }
