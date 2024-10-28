package api

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/core"
	"bitbucket.org/enroute-mobi/ara/model"
	"github.com/stretchr/testify/assert"
)

func prepareMultipart(t *testing.T, values map[string]io.Reader) (request *http.Request) {
	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range values {
		var fw io.Writer
		var err error
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		// Add an image file
		if x, ok := r.(*os.File); ok {
			if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
				t.Fatal(err)
			}
		} else {
			// Add other fields
			if fw, err = w.CreateFormField(key); err != nil {
				t.Fatal(err)
			}
		}
		if _, err = io.Copy(fw, r); err != nil {
			t.Fatal(err)
		}

	}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	request, err := http.NewRequest("POST", "/test/import", &b)
	if err != nil {
		t.Fatal(err)
	}
	// Don't forget to set the content type, this will contain the boundary.
	request.Header.Set("Content-Type", w.FormDataContentType())
	request.Header.Set("Authorization", "Token token=testToken")

	return request
}

func mustOpen(t *testing.T, f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func Test_Serve(t *testing.T) {
	assert := assert.New(t)

	model.InitTestDb(t)
	defer model.CleanTestDb(t)

	clock.SetDefaultClock(clock.NewFakeClockAt(time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)))

	// Initialize referential manager
	referentials := core.NewMemoryReferentials()
	// Save a new referential
	referential := referentials.New("test")
	referential.Tokens = []string{"testToken"}
	referentials.Save(referential)

	server := &Server{}
	server.SetReferentials(referentials)
	// Create a request
	//prepare the reader instances to encode
	values := map[string]io.Reader{
		"data":    mustOpen(t, "testdata/import.csv"),
		"request": strings.NewReader("{\"force\": false}"),
	}
	request := prepareMultipart(t, values)

	// Create a ResponseRecorder
	responseRecorder := httptest.NewRecorder()

	// Call HandleFlow method and pass in our Request and ResponseRecorder.
	server.HandleFlow(responseRecorder, request)

	// Test the result
	assert.Equal(http.StatusOK, responseRecorder.Code)
	assert.Equal("application/json", responseRecorder.Header().Get("Content-Type"))

	expectedBody := `{"Import":{"Total":13,"line":2,"line_group":1,"operator":2,"stop_area":2,` +
		`"stop_area_group":1,"stop_visit":2,"vehicle_journey":3},"Errors":{}}`

	assert.Equal(expectedBody, responseRecorder.Body.String())

	referential.ReloadModel()

	_, ok := referential.Model().Operators().Find("03eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
	assert.True(ok)

	_, ok = referential.Model().StopAreas().Find("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
	assert.True(ok)

	_, ok = referential.Model().StopAreaGroups().Find("cf3e1970-7a7e-4379-ae67-a67abe1c7c1b")
	assert.True(ok)

	_, ok = referential.Model().Lines().Find("f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
	assert.True(ok)

	_, ok = referential.Model().LineGroups().Find("59208069-3cad-4108-968f-349c5d50a351")
	assert.True(ok)

	_, ok = referential.Model().VehicleJourneys().Find("01eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
	assert.True(ok)

	// FIXME: We don't reload SV for now
	// _, ok = referential.Model().StopVisits().Find("02eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
	// if !ok {
	// 	t.Errorf("StopVisit should be found after the reload")
	// }

	// Send more requests to test force param
	values2 := map[string]io.Reader{
		"data":    mustOpen(t, "testdata/import.csv"),
		"request": strings.NewReader("{\"force\": false}"),
	}
	request2 := prepareMultipart(t, values2)

	responseRecorder2 := httptest.NewRecorder()
	server.HandleFlow(responseRecorder2, request2)

	// Test the result
	assert.Equal(http.StatusOK, responseRecorder2.Code)

	result := model.Result{}
	jsonDecoder := json.NewDecoder(responseRecorder2.Body)
	jsonDecoder.Decode(&result)

	assert.Equal(int64(7), result.Import["Errors"])

	values3 := map[string]io.Reader{
		"data":    mustOpen(t, "testdata/import.csv"),
		"request": strings.NewReader("{\"force\": true}"),
	}
	request3 := prepareMultipart(t, values3)

	responseRecorder3 := httptest.NewRecorder()
	server.HandleFlow(responseRecorder3, request3)

	// Test the result
	assert.Equal(http.StatusOK, responseRecorder3.Code)

	result3 := model.Result{}
	jsonDecoder = json.NewDecoder(responseRecorder3.Body)
	jsonDecoder.Decode(&result3)
	assert.Zero(result3.Import["Errors"])
}

func Test_Serve_With_NoToken(t *testing.T) {
	model.InitTestDb(t)
	defer model.CleanTestDb(t)

	clock.SetDefaultClock(clock.NewFakeClockAt(time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)))

	// Initialize referential manager
	referentials := core.NewMemoryReferentials()
	// Save a new referential
	referential := referentials.New("test")
	referentials.Save(referential)

	server := &Server{}
	server.SetReferentials(referentials)
	// Create a request
	//prepare the reader instances to encode
	values := map[string]io.Reader{
		"data":    mustOpen(t, "testdata/import.csv"),
		"request": strings.NewReader("{\"force\": false}"),
	}
	request := prepareMultipart(t, values)

	// Create a ResponseRecorder
	responseRecorder := httptest.NewRecorder()

	// Call HandleFlow method and pass in our Request and ResponseRecorder.
	server.HandleFlow(responseRecorder, request)

	// Test the result
	if status := responseRecorder.Code; status == http.StatusOK {
		t.Errorf("Handler returned wrong status code:\n got %v\n want %v", http.StatusUnauthorized, http.StatusOK)
	}
}

func Test_Serve_With_ImportToken(t *testing.T) {
	assert := assert.New(t)

	model.InitTestDb(t)
	defer model.CleanTestDb(t)

	clock.SetDefaultClock(clock.NewFakeClockAt(time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)))

	// Initialize referential manager
	referentials := core.NewMemoryReferentials()
	// Save a new referential
	referential := referentials.New("test")
	referential.ImportTokens = []string{"testToken"}
	referentials.Save(referential)

	server := &Server{}
	server.SetReferentials(referentials)
	// Create a request
	//prepare the reader instances to encode
	values := map[string]io.Reader{
		"data":    mustOpen(t, "testdata/import.csv"),
		"request": strings.NewReader("{\"force\": false}"),
	}
	request := prepareMultipart(t, values)

	// Create a ResponseRecorder
	responseRecorder := httptest.NewRecorder()

	// Call HandleFlow method and pass in our Request and ResponseRecorder.
	server.HandleFlow(responseRecorder, request)

	// Test the result
	assert.Equal(http.StatusOK, responseRecorder.Code)
	assert.Equal("application/json", responseRecorder.Header().Get("Content-Type"))

	expectedBody := `{"Import":{"Total":13,"line":2,"line_group":1,"operator":2,"stop_area":2,` +
		`"stop_area_group":1,"stop_visit":2,"vehicle_journey":3},"Errors":{}}`

	assert.Equal(expectedBody, responseRecorder.Body.String())
}
