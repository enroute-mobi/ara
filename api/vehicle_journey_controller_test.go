package api

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"bitbucket.org/enroute-mobi/ara/core"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func checkVehicleJourneyResponseStatus(responseRecorder *httptest.ResponseRecorder, t *testing.T) {
	require := require.New(t)

	require.Equal(http.StatusOK, responseRecorder.Code)
	require.Equal("application/json", responseRecorder.Header().Get("Content-Type"))
}

func prepareVehicleJourneyRequest(method string, sendIdentifier bool, body []byte, t *testing.T) (vehicleJourney *model.VehicleJourney, responseRecorder *httptest.ResponseRecorder, referential *core.Referential) {
	// Create a referential
	referentials := core.NewMemoryReferentials()
	server := &Server{}
	server.SetReferentials(referentials)
	referential = referentials.New("default")
	referential.Tokens = []string{"testToken"}
	referential.Save()

	// Set the fake UUID generator
	uuid.SetDefaultUUIDGenerator(uuid.NewFakeUUIDGenerator())
	// Save a new vehicleJourney
	vehicleJourney = referential.Model().VehicleJourneys().New()
	referential.Model().VehicleJourneys().Save(vehicleJourney)

	// Create a request
	address := []byte("/default/vehicle_journeys")
	if sendIdentifier {
		address = append(address, fmt.Sprintf("/%s", vehicleJourney.Id())...)
	}
	request, err := http.NewRequest(method, string(address), bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Authorization", "Token token=testToken")

	// Create a ResponseRecorder
	responseRecorder = httptest.NewRecorder()

	request.SetPathValue("referential_slug", string(referential.Slug()))
	request.SetPathValue("model", "vehicle_journeys")
	switch method {
	case "GET":
		if sendIdentifier == false && body == nil {
			server.handleReferentialModelIndex(responseRecorder, request)
		} else {
			request.SetPathValue("id", string(vehicleJourney.Id()))
			server.handleReferentialModelShow(responseRecorder, request)
		}
	case "POST":
		server.handleReferentialModelCreate(responseRecorder, request)
	case "PUT":
		request.SetPathValue("id", string(vehicleJourney.Id()))
		server.handleReferentialModelUpdate(responseRecorder, request)
	case "DELETE":
		request.SetPathValue("id", string(vehicleJourney.Id()))
		server.handleReferentialModelDelete(responseRecorder, request)
	default:
		t.Fatalf("Unknown method: %s", method)
	}

	return
}

func Test_VehicleJourneyController_Delete(t *testing.T) {
	assert := assert.New(t)

	vehicleJourney, responseRecorder, referential := prepareVehicleJourneyRequest("DELETE", true, nil, t)

	// Test response
	checkVehicleJourneyResponseStatus(responseRecorder, t)

	//Test Results
	_, ok := referential.Model().VehicleJourneys().Find(vehicleJourney.Id())
	assert.False(ok, "VehicleJourney shouldn't be found after DELETE request")

	expected, err := vehicleJourney.MarshalJSON()
	assert.NoError(err)
	assert.JSONEq(string(expected), responseRecorder.Body.String())
}

func Test_VehicleJourneyController_Update(t *testing.T) {
	assert := assert.New(t)

	// Prepare and send request
	body := []byte(`{ "Codes": { "reflex": "FR:77491:ZDE:34004:STIF" } }`)
	vehicleJourney, responseRecorder, referential := prepareVehicleJourneyRequest("PUT", true, body, t)

	// Check response
	checkVehicleJourneyResponseStatus(responseRecorder, t)

	// Test Results
	updatedVehicleJourney, ok := referential.Model().VehicleJourneys().Find(vehicleJourney.Id())
	assert.True(ok, "VehicleJourney should be found after PUT request")

	expected, err := updatedVehicleJourney.MarshalJSON()
	assert.NoError(err)
	assert.JSONEq(string(expected), responseRecorder.Body.String())
}

func Test_VehicleJourneyController_Show(t *testing.T) {
	assert := assert.New(t)

	// Send request
	vehicleJourney, responseRecorder, _ := prepareVehicleJourneyRequest("GET", true, nil, t)

	// Test response
	checkVehicleJourneyResponseStatus(responseRecorder, t)

	//Test Results
	expectedVehicleJourney, _ := vehicleJourney.MarshalJSON()
	assert.JSONEq(string(expectedVehicleJourney), responseRecorder.Body.String())
}

func Test_VehicleJourneyController_Create(t *testing.T) {
	assert := assert.New(t)

	// Prepare and send request
	body := []byte(`{ "Codes": { "reflex": "FR:77491:ZDE:34004:STIF" } }`)
	_, responseRecorder, referential := prepareVehicleJourneyRequest("POST", false, body, t)

	// Check response
	checkVehicleJourneyResponseStatus(responseRecorder, t)

	// Test Results
	// Using the fake uuid generator, the uuid of the created
	// vehicleJourney should be 6ba7b814-9dad-11d1-1-00c04fd430c8
	vehicleJourney, ok := referential.Model().VehicleJourneys().Find("6ba7b814-9dad-11d1-1-00c04fd430c8")
	assert.True(ok, "VehicleJourney should be found after POST request")

	expected := `{"Codes":{"reflex":"FR:77491:ZDE:34004:STIF"},"HasCompleteStopSequence":false,"Id":"6ba7b814-9dad-11d1-1-00c04fd430c8","Monitored":false}`
	vehicleJourneyMarshal, err := vehicleJourney.MarshalJSON()
	assert.NoError(err)
	assert.JSONEq(expected, string(vehicleJourneyMarshal))
}

func Test_VehicleJourneyController_Index(t *testing.T) {
	assert := assert.New(t)

	// Send request
	_, responseRecorder, _ := prepareVehicleJourneyRequest("GET", false, nil, t)

	// Test response
	checkVehicleJourneyResponseStatus(responseRecorder, t)

	//Test Results
	expected := `[{"Monitored":false,"HasCompleteStopSequence":false,"Id":"6ba7b814-9dad-11d1-0-00c04fd430c8"}]`
	assert.JSONEq(expected, responseRecorder.Body.String())
}

func Test_VehicleJourneyController_FindVehicleJourney(t *testing.T) {
	assert := assert.New(t)

	ref := core.NewMemoryReferentials().New("test")

	vehicleJourney := ref.Model().VehicleJourneys().New()
	code := model.NewCode("codeSpace", "value")
	vehicleJourney.SetCode(code)
	ref.Model().VehicleJourneys().Save(vehicleJourney)

	controller := &VehicleJourneyController{
		referential: ref,
	}

	_, ok := controller.findVehicleJourney("codeSpace:value")
	assert.True(ok, "Can't find VehicleJourney by Code")

	_, ok = controller.findVehicleJourney(string(vehicleJourney.Id()))
	assert.True(ok, "Can't find VehicleJourney by Id")
}
