package api

import (
	"io/ioutil"
	"net/http"
	"net/url"

	"bitbucket.org/enroute-mobi/edwig/core"
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
	"operators":        NewOperatorController,
	"vehicles":         NewVehicleController,
	"import":           NewImportController,
}

type RestfulResource interface {
	Index(response http.ResponseWriter, filters url.Values)
	Show(response http.ResponseWriter, identifier string)
	Delete(response http.ResponseWriter, identifier string)
	Update(response http.ResponseWriter, identifier string, body []byte)
	Create(response http.ResponseWriter, body []byte)
}

type ActionResource interface {
	Action(response http.ResponseWriter, requestData *RequestData)
}

type Savable interface {
	Save(response http.ResponseWriter)
}

type ControllerInterface interface {
	serve(response http.ResponseWriter, request *http.Request, requestData *RequestData)
}

type Controller struct {
	restfulResource RestfulResource
}

func getRequestBody(response http.ResponseWriter, request *http.Request) []byte {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(response, "Invalid request: Can't read request body", http.StatusBadRequest)
		return nil
	}
	if len(body) == 0 {
		http.Error(response, "Invalid request: Empty body", http.StatusBadRequest)
		return nil
	}
	return body
}

func (controller *Controller) serve(response http.ResponseWriter, request *http.Request, requestData *RequestData) {
	// Check request body
	if requestData.Method == "PUT" || (requestData.Method == "POST" && requestData.Id != "save") {
		requestData.Body = getRequestBody(response, request)
		if requestData.Body == nil {
			http.Error(response, "Invalid request", http.StatusBadRequest)
			return
		}
	}

	// Check request Id
	if (requestData.Method == "DELETE" || requestData.Method == "PUT") && requestData.Id == "" {
		http.Error(response, "Invalid request", http.StatusBadRequest)
		return
	}

	if requestData.Action != "" {
		if actionResource, ok := controller.restfulResource.(ActionResource); ok {
			actionResource.Action(response, requestData)
			return
		}
	}

	switch requestData.Method {
	case "GET":
		if requestData.Id == "" {
			controller.restfulResource.Index(response, requestData.Filters)
			return
		}
		controller.restfulResource.Show(response, requestData.Id)
	case "DELETE":
		controller.restfulResource.Delete(response, requestData.Id)
	case "PUT":
		controller.restfulResource.Update(response, requestData.Id, requestData.Body)
	case "POST":
		if requestData.Id == "save" {
			if savableResource, ok := controller.restfulResource.(Savable); ok {
				savableResource.Save(response)
				return
			}
		}
		if requestData.Id != "" {
			http.Error(response, "Invalid request", http.StatusBadRequest)
			return
		}
		controller.restfulResource.Create(response, requestData.Body)
	}
}
