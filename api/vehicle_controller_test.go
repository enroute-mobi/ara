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

func checkVehicleResponseStatus(responseRecorder *httptest.ResponseRecorder, t *testing.T) {
	require := require.New(t)

	require.Equal(http.StatusOK, responseRecorder.Code)
	require.Equal("application/json", responseRecorder.Header().Get("Content-Type"))
}

func prepareVehicleRequest(method string, sendIdentifier bool, body []byte, t *testing.T) (vehicle *model.Vehicle, responseRecorder *httptest.ResponseRecorder, referential *core.Referential) {
	// Create a referential
	referentials := core.NewMemoryReferentials()
	server := &Server{}
	server.SetReferentials(referentials)
	referential = referentials.New("default")
	referential.Tokens = []string{"testToken"}
	referential.Save()

	// Set the fake UUID generator
	uuid.SetDefaultUUIDGenerator(uuid.NewFakeUUIDGenerator())
	// Save a new vehicle
	vehicle = referential.Model().Vehicles().New()
	vehicle.Longitude = 1.2
	vehicle.Latitude = 3.4
	vehicle.Bearing = 5.6
	referential.Model().Vehicles().Save(vehicle)

	// Create a request
	address := []byte("/default/vehicles")
	if sendIdentifier {
		address = append(address, fmt.Sprintf("/%s", vehicle.Id())...)
	}
	request, err := http.NewRequest(method, string(address), bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	request.Header.Set("Authorization", "Token token=testToken")
	// Create a ResponseRecorder
	responseRecorder = httptest.NewRecorder()

	request.SetPathValue("referential_slug", string(referential.Slug()))
	request.SetPathValue("model", "vehicles")
	switch method {
	case "GET":
		if sendIdentifier == false && body == nil {
			server.handleReferentialModelIndex(responseRecorder, request)
		} else {
			request.SetPathValue("id", string(vehicle.Id()))
			server.handleReferentialModelShow(responseRecorder, request)
		}
	case "POST":
		server.handleReferentialModelCreate(responseRecorder, request)
	case "PUT":
		request.SetPathValue("id", string(vehicle.Id()))
		server.handleReferentialModelUpdate(responseRecorder, request)
	case "DELETE":
		request.SetPathValue("id", string(vehicle.Id()))
		server.handleReferentialModelDelete(responseRecorder, request)
	default:
		t.Fatalf("Unknown method: %s", method)
	}

	return
}

func Test_VehicleController_Delete(t *testing.T) {
	assert := assert.New(t)

	// Send request
	vehicle, responseRecorder, referential := prepareVehicleRequest("DELETE", true, nil, t)

	// Test response
	checkVehicleResponseStatus(responseRecorder, t)

	_, ok := referential.Model().Vehicles().Find(vehicle.Id())
	assert.False(ok, "Vehicle shouldn't be found after DELETE request")

	expectedVehicle, _ := vehicle.MarshalJSON()
	assert.JSONEq(string(expectedVehicle), responseRecorder.Body.String())
}

func Test_VehicleController_Update(t *testing.T) {
	assert := assert.New(t)

	// Prepare and send request
	body := []byte(`{ "Bearing": 13.3 }`)
	vehicle, responseRecorder, referential := prepareVehicleRequest("PUT", true, body, t)

	// Check response
	checkVehicleResponseStatus(responseRecorder, t)

	// Test Results
	updatedVehicle, ok := referential.Model().Vehicles().Find(vehicle.Id())
	assert.True(ok, "Vehicle should be found after PUT request")
	assert.Equal(13.3, updatedVehicle.Bearing)

	expectedVehicle, _ := updatedVehicle.MarshalJSON()
	assert.JSONEq(string(expectedVehicle), responseRecorder.Body.String())
}

func Test_VehicleController_Show(t *testing.T) {
	assert := assert.New(t)

	// Send request
	vehicle, responseRecorder, _ := prepareVehicleRequest("GET", true, nil, t)

	// Test response
	checkVehicleResponseStatus(responseRecorder, t)

	//Test Results
	expectedVehicle, _ := vehicle.MarshalJSON()
	assert.JSONEq(string(expectedVehicle), responseRecorder.Body.String())
}

func Test_VehicleController_Create(t *testing.T) {
	assert := assert.New(t)

	// Prepare and send request
	body := []byte(`{ "Bearing": 1.234 }`)
	_, responseRecorder, referential := prepareVehicleRequest("POST", false, body, t)

	// Check response
	checkVehicleResponseStatus(responseRecorder, t)

	// Test Results
	// Using the fake uuid generator, the uuid of the created
	// vehicle should be 6ba7b814-9dad-11d1-1-00c04fd430c8
	vehicle, ok := referential.Model().Vehicles().Find("6ba7b814-9dad-11d1-1-00c04fd430c8")
	assert.True(ok, "Vehicle should be found after POST request")
	assert.Equal(1.234, vehicle.Bearing)

	expectedVehicle, _ := vehicle.MarshalJSON()
	assert.JSONEq(string(expectedVehicle), responseRecorder.Body.String())
}

func Test_VehicleController_Index(t *testing.T) {
	assert := assert.New(t)

	// Send request
	_, responseRecorder, _ := prepareVehicleRequest("GET", false, nil, t)

	// Test response
	checkVehicleResponseStatus(responseRecorder, t)

	//Test Results
	expected := `[{"RecordedAtTime":"0001-01-01T00:00:00Z","ValidUntilTime":"0001-01-01T00:00:00Z","Longitude":1.2,"Latitude":3.4,"Bearing":5.6,"Id":"6ba7b814-9dad-11d1-0-00c04fd430c8"}]`
	assert.JSONEq(expected, responseRecorder.Body.String())
}

func Test_VehicleController_FindVehicle(t *testing.T) {
	assert := assert.New(t)

	ref := core.NewMemoryReferentials().New("test")

	vehicle := ref.Model().Vehicles().New()
	code := model.NewCode("codeSpace", "value")
	vehicle.SetCode(code)
	ref.Model().Vehicles().Save(vehicle)

	controller := &VehicleController{
		referential: ref,
	}

	_, ok := controller.findVehicle("codeSpace:value")
	assert.True(ok, "Can't find Vehicle by Code")

	_, ok = controller.findVehicle(string(vehicle.Id()))
	assert.True(ok, "Can't find Vehicle by Id")
}
