package api

import (
	"io/ioutil"
	"net/http"

	"github.com/af83/edwig/core"
)

var newControllerMap = map[string](func(*Server) ControllerInterface){
	"_referentials": NewReferentialController,
	"_time":         NewTimeController,
	"_status":       NewStatusController,
}

var newWithReferentialControllerMap = map[string](func(*core.Referential) ControllerInterface){
	"stop_areas":       NewStopAreaController,
	"partners":         NewPartnerController,
	"lines":            NewLineController,
	"stop_visits":      NewStopVisitController,
	"vehicle_journeys": NewVehicleJourneyController,
	"situations":       NewSituationController,
}

type RestfulRessource interface {
	Index(response http.ResponseWriter)
	Show(response http.ResponseWriter, identifier string)
	Delete(response http.ResponseWriter, identifier string)
	Update(response http.ResponseWriter, identifier string, body []byte)
	Create(response http.ResponseWriter, body []byte)
}

type ActionResource interface {
	Action(response http.ResponseWriter, requestData *RequestData)
}

type ControllerInterface interface {
	serve(response http.ResponseWriter, request *http.Request, requestData *RequestData)
}

type Controller struct {
	restfulRessource RestfulRessource
}

func getRequestBody(response http.ResponseWriter, request *http.Request) []byte {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(response, "Invalid request: Can't read request body", 400)
		return nil
	}
	if len(body) == 0 {
		http.Error(response, "Invalid request: Empty body", 400)
		return nil
	}
	return body
}

func (controller *Controller) serve(response http.ResponseWriter, request *http.Request, requestData *RequestData) {

	if requestData.Action != "" {
		if actionResource, ok := controller.restfulRessource.(ActionResource); ok {
			actionResource.Action(response, requestData)
			return
		}
	}

	switch request.Method {
	case "GET":
		if requestData.Id == "" {
			controller.restfulRessource.Index(response)
		} else {
			controller.restfulRessource.Show(response, requestData.Id)
		}
	case "DELETE":
		if requestData.Id == "" {
			http.Error(response, "Invalid request", 400)
			return
		}
		controller.restfulRessource.Delete(response, requestData.Id)
	case "PUT":
		if requestData.Id == "" {
			http.Error(response, "Invalid request", 400)
			return
		}
		body := getRequestBody(response, request)
		if body == nil {
			return
		}
		controller.restfulRessource.Update(response, requestData.Id, body)
	case "POST":
		if requestData.Id != "" {
			http.Error(response, "Invalid request", 400)
			return
		}
		body := getRequestBody(response, request)
		if body == nil {
			return
		}
		controller.restfulRessource.Create(response, body)
	}
}
