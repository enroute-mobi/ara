package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/clock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func timePrepareRequest(method string, advance bool, body []byte, t *testing.T) (responseRecorder *httptest.ResponseRecorder, server *Server) {
	require := require.New(t)

	server = &Server{}
	server.SetClock(clock.NewFakeClock())

	// Create a request
	var address string
	if advance {
		address = "/_time/advance"
	} else {
		address = "/_time"
	}

	request, err := http.NewRequest(method, address, bytes.NewReader(body))
	require.NoError(err)

	// Create a ResponseRecorder
	responseRecorder = httptest.NewRecorder()

	if advance {
		server.handleTimeAdvance(responseRecorder, request)
	} else {
		server.handleTimeGet(responseRecorder, request)
	}

	return
}

func Test_TimeController_Get(t *testing.T) {
	assert := assert.New(t)

	// Send request
	responseRecorder, _ := timePrepareRequest("GET", false, nil, t)

	assert.Equal(http.StatusOK, responseRecorder.Code)
	assert.Equal("application/json", responseRecorder.Header().Get("Content-Type"))
	assert.Equal(`{ "time": "1984-04-04T00:00:00.000Z" }`, responseRecorder.Body.String())
}

func Test_TimeController_Advance(t *testing.T) {
	assert := assert.New(t)

	// Send request
	body := []byte(`{ "duration": "1s" }`)
	responseRecorder, server := timePrepareRequest("POST", true, body, t)

	assert.Equal(http.StatusOK, responseRecorder.Code)
	assert.Equal("application/json", responseRecorder.Header().Get("Content-Type"))
	assert.Equal(`{ "time": "1984-04-04T00:00:01.000Z" }`, responseRecorder.Body.String())

	expected := time.Date(1984, time.April, 4, 0, 0, 1, 0, time.UTC)
	assert.True(server.Clock().Now().Equal(expected))
}
