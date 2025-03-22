package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"bitbucket.org/enroute-mobi/ara/core"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
)

type StopVisitController struct {
	svs model.StopVisits
}

func NewStopVisitController(referential *core.Referential) RestfulResource {
	return &StopVisitController{
		svs: referential.Model().StopVisits(),
	}
}

func NewScheduledStopVisitController(referential *core.Referential) RestfulResource {
	return &StopVisitController{
		svs: referential.Model().ScheduledStopVisits(),
	}
}

func (controller *StopVisitController) findStopVisit(identifier string) (*model.StopVisit, bool) {
	foundStrings := idPattern.FindStringSubmatch(identifier)
	if foundStrings != nil {
		code := model.NewCode(foundStrings[1], foundStrings[2])
		return controller.svs.FindByCode(code)
	}
	return controller.svs.Find(model.StopVisitId(identifier))
}

func (controller *StopVisitController) Index(response http.ResponseWriter, _params url.Values) {
	stopVisits := controller.svs.FindAll()

	logger.Log.Debugf("StopVisits Index")
	jsonBytes, _ := json.Marshal(stopVisits)
	response.Write(jsonBytes)
}

func (controller *StopVisitController) Show(response http.ResponseWriter, identifier string) {
	stopVisit, ok := controller.findStopVisit(identifier)
	if !ok {
		http.Error(response, fmt.Sprintf("Stop visit not found: %s", identifier), http.StatusNotFound)
		return
	}
	logger.Log.Debugf("Get stopVisit %s", identifier)

	jsonBytes, _ := stopVisit.MarshalJSON()
	response.Write(jsonBytes)
}

func (controller *StopVisitController) Delete(response http.ResponseWriter, identifier string) {
	stopVisit, ok := controller.findStopVisit(identifier)
	if !ok {
		http.Error(response, fmt.Sprintf("Stop visit not found: %s", identifier), http.StatusNotFound)
		return
	}
	logger.Log.Debugf("Delete stopVisit %s", identifier)

	jsonBytes, _ := stopVisit.MarshalJSON()
	controller.svs.Delete(stopVisit)
	response.Write(jsonBytes)
}

func (controller *StopVisitController) Update(response http.ResponseWriter, identifier string, body []byte) {
	stopVisit, ok := controller.findStopVisit(identifier)
	if !ok {
		http.Error(response, fmt.Sprintf("Stop visit not found: %s", identifier), http.StatusNotFound)
		return
	}

	logger.Log.Debugf("Update stopVisit %s: %s", identifier, string(body))

	err := json.Unmarshal(body, &stopVisit)
	if err != nil {
		http.Error(response, fmt.Sprintf("Invalid request: can't parse request body: %v", err), http.StatusBadRequest)
		return
	}

	for _, obj := range stopVisit.Codes() {
		sv, ok := controller.svs.FindByCode(obj)
		if ok && sv.Id() != stopVisit.Id() {
			http.Error(response, fmt.Sprintf("Invalid request: stopVisit %v already have a code %v", sv.Id(), obj.String()), http.StatusBadRequest)
			return
		}
	}

	controller.svs.Save(stopVisit)
	jsonBytes, _ := stopVisit.MarshalJSON()
	response.Write(jsonBytes)
}

func (controller *StopVisitController) Create(response http.ResponseWriter, body []byte) {
	stopVisit := controller.svs.New()

	err := json.Unmarshal(body, &stopVisit)
	if err != nil {
		http.Error(response, fmt.Sprintf("Invalid request: can't parse request body: %v", err), http.StatusBadRequest)
		return
	}

	if stopVisit.Id() != "" {
		http.Error(response, "Invalid request", http.StatusBadRequest)
		return
	}

	for _, obj := range stopVisit.Codes() {
		sv, ok := controller.svs.FindByCode(obj)
		if ok {
			http.Error(response, fmt.Sprintf("Invalid request: stopVisit %v already have a code %v", sv.Id(), obj.String()), http.StatusBadRequest)
			return
		}
	}

	controller.svs.Save(stopVisit)
	jsonBytes, _ := stopVisit.MarshalJSON()
	response.Write(jsonBytes)
}
