package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"bitbucket.org/enroute-mobi/ara/core"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
)

type StopVisitController struct {
	svs model.StopVisits
}

func NewStopVisitController(referential *core.Referential) ControllerInterface {
	return &Controller{
		restfulResource: &StopVisitController{
			svs: referential.Model().StopVisits(),
		},
	}
}

func NewScheduledStopVisitController(referential *core.Referential) ControllerInterface {
	return &Controller{
		restfulResource: &StopVisitController{
			svs: referential.Model().ScheduledStopVisits(),
		},
	}
}

func (controller *StopVisitController) findStopVisit(identifier string) (*model.StopVisit, bool) {
	idRegexp := "([0-9a-zA-Z-]+):([0-9a-zA-Z-:]+)"
	pattern := regexp.MustCompile(idRegexp)
	foundStrings := pattern.FindStringSubmatch(identifier)
	if foundStrings != nil {
		objectid := model.NewObjectID(foundStrings[1], foundStrings[2])
		return controller.svs.FindByObjectId(objectid)
	}
	return controller.svs.Find(model.StopVisitId(identifier))
}

func (controller *StopVisitController) filterStopVisits(stopVisits []*model.StopVisit, filters url.Values) []*model.StopVisit {
	selectors := []model.StopVisitSelector{}
	filteredStopVisits := []*model.StopVisit{}

	for key, value := range filters {
		switch key {
		case "After":
			layout := "2006/01/02-15:04:05"
			startTime, err := time.Parse(layout, value[0])
			if err != nil {
				continue
			}
			selectors = append(selectors, model.StopVisitSelectorAfterTime(startTime))
		case "Before":
			layout := "2006/01/02-15:04:05"
			endTime, err := time.Parse(layout, value[0])
			if err != nil {
				continue
			}
			selectors = append(selectors, model.StopVisitSelectorAfterTime(endTime))
		case "StopArea":
			selectors = append(selectors, model.StopVisitSelectByStopAreaId(model.StopAreaId(value[0])))
		}
	}

	selector := model.CompositeStopVisitSelector(selectors)
	for _, sv := range stopVisits {
		if !selector(sv) {
			continue
		}
		filteredStopVisits = append(filteredStopVisits, sv)
	}

	return filteredStopVisits
}

func (controller *StopVisitController) Index(response http.ResponseWriter, filters url.Values) {
	stopVisits := controller.svs.FindAll()
	filteredStopVisits := controller.filterStopVisits(stopVisits, filters)

	logger.Log.Debugf("StopVisits Index")
	jsonBytes, _ := json.Marshal(filteredStopVisits)
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

	for _, obj := range stopVisit.ObjectIDs() {
		sv, ok := controller.svs.FindByObjectId(obj)
		if ok && sv.Id() != stopVisit.Id() {
			http.Error(response, fmt.Sprintf("Invalid request: stopVisit %v already have an objectid %v", sv.Id(), obj.String()), http.StatusBadRequest)
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

	for _, obj := range stopVisit.ObjectIDs() {
		sv, ok := controller.svs.FindByObjectId(obj)
		if ok {
			http.Error(response, fmt.Sprintf("Invalid request: stopVisit %v already have an objectid %v", sv.Id(), obj.String()), http.StatusBadRequest)
			return
		}
	}

	controller.svs.Save(stopVisit)
	jsonBytes, _ := stopVisit.MarshalJSON()
	response.Write(jsonBytes)
}
