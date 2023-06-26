package remote

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createSIRILiteServer(t *testing.T, returnedFile string, opts ...int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// force status code when provided
		if len(opts) != 0 {
			w.WriteHeader(opts[0])
		}
		file, err := os.Open(fmt.Sprintf("testdata/%s.json", returnedFile))
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()
		io.Copy(w, file)
	}))
}

func Test_SIRILiteClient_StopMonitoringDelivery(t *testing.T) {
	assert := assert.New(t)

	// Create a test http server
	ts := createSIRILiteServer(t, "stopmonitoring-lite-delivery")
	defer ts.Close()

	// Create and send request
	httpClient := NewHTTPClient(HTTPClientOptions{Urls: HTTPClientUrls{Url: ts.URL}})
	client := httpClient.SIRILiteClient()

	stopArea := "STIF:StopPoint:Q:41178:"

	dest, err := client.StopMonitoring(stopArea)
	assert.Nil(err)
	assert.Equal("STIF:StopPoint:Q:41178:", dest.Siri.ServiceDelivery.StopMonitoringDelivery[0].MonitoredStopVisit[0].MonitoringRef)
}

func Test_SIRILiteClient_StopMonitoringDelivery_With_Error_400(t *testing.T) {
	assert := assert.New(t)

	// Create a test http server with return code 400 and error payload
	ts := createSIRILiteServer(t, "stopmonitoring-lite-delivery-error", 400)
	defer ts.Close()

	// Create and send request
	httpClient := NewHTTPClient(HTTPClientOptions{Urls: HTTPClientUrls{Url: ts.URL}})
	client := httpClient.SIRILiteClient()
	stopArea := "STIF:StopPoint:Q:41178:"

	_, err := client.StopMonitoring(stopArea)
	assert.EqualError(err, "request failed with status 400: La requête contient des identifiants qui sont inconnus", err.Error())
}

func Test_SIRILiteClient_StopMonitoringDelivery_With_Error_400_With_unmarshallable_Error(t *testing.T) {
	assert := assert.New(t)

	// Create a test http server with return code 400 and error payload
	ts := createSIRILiteServer(t, "stop-monitoring-lite-delivery-error-unmarshallable", 400)
	defer ts.Close()

	// Create and send request
	httpClient := NewHTTPClient(HTTPClientOptions{Urls: HTTPClientUrls{Url: ts.URL}})
	client := httpClient.SIRILiteClient()
	stopArea := "STIF:StopPoint:Q:41178:"
	expectedErrorMessage := "request failed with status 400: \"ERROR\""

	_, err := client.StopMonitoring(stopArea)
	assert.EqualError(err, expectedErrorMessage, err.Error())
}
