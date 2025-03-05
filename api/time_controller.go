package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/logger"
)

type TimeController struct {
	server *Server
}

func NewTimeController(server *Server) *TimeController {
	return &TimeController{
		server: server,
	}
}

func (controller *TimeController) get(response http.ResponseWriter) {
	responseTime := controller.server.Clock().Now().Format(`{ "time": "2006-01-02T15:04:05.000Z07:00" }`)
	response.Write([]byte(responseTime))
}

func (controller *TimeController) advance(response http.ResponseWriter, request *http.Request) {
	if _, ok := controller.server.Clock().(clock.FakeClock); !ok {
		http.Error(response, "Invalid request: server has a real Clock", http.StatusBadRequest)
		return
	}

	body := getRequestBody(response, request)
	if body == nil {
		return
	}
	var responseBody map[string]string
	if err := json.Unmarshal(body, &responseBody); err != nil {
		http.Error(response, fmt.Sprintf("Invalid request: can't parse request body: %v", err), http.StatusBadRequest)
		return
	}
	duration, ok := responseBody["duration"]
	if !ok {
		http.Error(response, "Invalid request: can't find duration", http.StatusBadRequest)
		return
	}
	parsedDuration, err := time.ParseDuration(duration)
	if err != nil {
		http.Error(response, fmt.Sprintf("Invalid request: can't parse duration: %v", err), http.StatusBadRequest)
		return
	}
	logger.Log.Printf("Advance time by %v", parsedDuration)
	controller.server.Clock().(clock.FakeClock).Advance(parsedDuration)
	logger.Log.Printf("Time is now %v", controller.server.Clock().Now())

	controller.get(response)
}
