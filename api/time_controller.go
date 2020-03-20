package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
)

type TimeController struct {
	server *Server
}

func NewTimeController(server *Server) ControllerInterface {
	return &TimeController{
		server: server,
	}
}

func (controller *TimeController) serve(response http.ResponseWriter, request *http.Request, requestData *RequestData) {
	switch {
	case request.Method == "GET":
		if requestData.Resource != "" {
			http.Error(response, "Invalid request", http.StatusBadRequest)
			return
		}
		controller.get(response)
	case request.Method == "POST":
		if _, ok := controller.server.Clock().(model.FakeClock); !ok {
			http.Error(response, "Invalid request: server has a real Clock", http.StatusBadRequest)
			return
		}
		if requestData.Resource != "advance" {
			http.Error(response, "Invalid request: invalid action", http.StatusBadRequest)
			return
		}
		body := getRequestBody(response, request)
		if body == nil {
			return
		}
		controller.advance(response, body)
	default:
		http.Error(response, "Invalid request", http.StatusBadRequest)
		return
	}
}

func (controller *TimeController) get(response http.ResponseWriter) {
	responseTime := controller.server.Clock().Now().Format(`{ "time": "2006-01-02T15:04:05.000Z07:00" }`)
	response.Write([]byte(responseTime))
}

func (controller *TimeController) advance(response http.ResponseWriter, body []byte) {
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
	controller.server.Clock().(model.FakeClock).Advance(parsedDuration)
	logger.Log.Printf("Time is now %v", controller.server.Clock().Now())

	controller.get(response)
}
