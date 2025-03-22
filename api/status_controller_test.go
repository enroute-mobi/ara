package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func statusPrepareRequest(method string, t *testing.T) (server *Server, request *http.Request, responseRecorder *httptest.ResponseRecorder) {
	require := require.New(t)
	server = &Server{}
	var err error
	request, err = http.NewRequest(method, "/_status", nil)
	require.NoError(err)

	responseRecorder = httptest.NewRecorder()
	return
}

func Test_status_check(t *testing.T) {
	assert := assert.New(t)

	server, request, responseRecorder := statusPrepareRequest("GET", t)
	server.handleStatus(responseRecorder, request)

	assert.Equal(http.StatusOK, responseRecorder.Code)
}

func Test_status_check_withApiKey(t *testing.T) {
	assert := assert.New(t)

	server, request, responseRecorder := statusPrepareRequest("GET", t)
	server.apiKey = "dummy"
	server.handleStatus(responseRecorder, request)

	assert.Equal(http.StatusOK, responseRecorder.Code)
}

func Test_status_check_version(t *testing.T) {
	assert := assert.New(t)

	server, request, responseRecorder := statusPrepareRequest("GET", t)
	server.handleStatus(responseRecorder, request)

	assert.Equal(http.StatusOK, responseRecorder.Code)

	status := &Status{}
	body, _ := io.ReadAll(responseRecorder.Body)
	err := json.Unmarshal(body, status)
	assert.NoError(err)
	assert.Equal("ok", status.Status)
	assert.NotZero(status.Version)
}
