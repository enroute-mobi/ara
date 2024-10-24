package api

import (
	"io"
	"net/http"
	"net/url"
	"regexp"

	"bitbucket.org/enroute-mobi/ara/core"
)

var idPattern = regexp.MustCompile("([0-9a-zA-Z-]+):([0-9a-zA-Z-:]+)")

var newControllerMap = map[string](func(*Server) ControllerInterface){
	"_referentials": NewReferentialController,
	"_time":         NewTimeController,
	"_status":       NewStatusController,
}

var newWithReferentialControllerMap = map[string](func(*core.Referential) ControllerInterface){
	"stop_areas":            NewStopAreaController,
	"partners":              NewPartnerController,
	"lines":                 NewLineController,
	"line_groups":           NewLineGroupsController,
	"stop_visits":           NewStopVisitController,
	"scheduled_stop_visits": NewScheduledStopVisitController,
	"vehicle_journeys":      NewVehicleJourneyController,
	"situations":            NewSituationController,
	"operators":             NewOperatorController,
	"vehicles":              NewVehicleController,
	"import":                NewImportController,
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
	if request.Body == nil {
		http.Error(response, "Invalid request: Can't read request body", http.StatusBadRequest)
		return nil
	}
	body, err := io.ReadAll(request.Body)
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
	if requestData.Method == "PUT" || (requestData.Method == "POST" && requestData.Id != "save" && requestData.Action != "reload") {
		requestData.Body = getRequestBody(response, request)
		if requestData.Body == nil {
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
