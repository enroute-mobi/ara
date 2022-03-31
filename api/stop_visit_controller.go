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
	referential *core.Referential
}

func NewStopVisitController(referential *core.Referential) ControllerInterface {
	return &Controller{
		restfulResource: &StopVisitController{
			referential: referential,
		},
	}
}

func (controller *StopVisitController) findStopVisit(tx *model.Transaction, identifier string) (model.StopVisit, bool) {
	idRegexp := "([0-9a-zA-Z-]+):([0-9a-zA-Z-:]+)"
	pattern := regexp.MustCompile(idRegexp)
	foundStrings := pattern.FindStringSubmatch(identifier)
	if foundStrings != nil {
		objectid := model.NewObjectID(foundStrings[1], foundStrings[2])
		return tx.Model().StopVisits().FindByObjectId(objectid)
	}
	return tx.Model().StopVisits().Find(model.StopVisitId(identifier))
}

func (controller *StopVisitController) filterStopVisits(stopVisits []model.StopVisit, filters url.Values) []model.StopVisit {
	selectors := []model.StopVisitSelector{}
	filteredStopVisits := []model.StopVisit{}

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
		if !selector(&sv) {
			continue
		}
		filteredStopVisits = append(filteredStopVisits, sv)
	}

	return filteredStopVisits
}

func (controller *StopVisitController) Index(response http.ResponseWriter, filters url.Values) {
	tx := controller.referential.NewTransaction()
	defer tx.Close()

	stopVisits := tx.Model().StopVisits().FindAll()
	filteredStopVisits := controller.filterStopVisits(stopVisits, filters)

	logger.Log.Debugf("StopVisits Index")
	jsonBytes, _ := json.Marshal(filteredStopVisits)
	response.Write(jsonBytes)
}

func (controller *StopVisitController) Show(response http.ResponseWriter, identifier string) {
	tx := controller.referential.NewTransaction()
	defer tx.Close()

	stopVisit, ok := controller.findStopVisit(tx, identifier)
	if !ok {
		http.Error(response, fmt.Sprintf("Stop visit not found: %s", identifier), http.StatusNotFound)
		return
	}
	logger.Log.Debugf("Get stopVisit %s", identifier)

	jsonBytes, _ := stopVisit.MarshalJSON()
	response.Write(jsonBytes)
}

func (controller *StopVisitController) Delete(response http.ResponseWriter, identifier string) {
	// New transaction
	tx := controller.referential.NewTransaction()
	defer tx.Close()

	stopVisit, ok := controller.findStopVisit(tx, identifier)
	if !ok {
		http.Error(response, fmt.Sprintf("Stop visit not found: %s", identifier), http.StatusNotFound)
		return
	}
	logger.Log.Debugf("Delete stopVisit %s", identifier)

	jsonBytes, _ := stopVisit.MarshalJSON()
	tx.Model().StopVisits().Delete(&stopVisit)
	err := tx.Commit()
	if err != nil {
		logger.Log.Debugf("Transaction error: %v", err)
		http.Error(response, "Internal error", http.StatusInternalServerError)
		return
	}
	response.Write(jsonBytes)
}

func (controller *StopVisitController) Update(response http.ResponseWriter, identifier string, body []byte) {
	// New transaction
	tx := controller.referential.NewTransaction()
	defer tx.Close()

	stopVisit, ok := controller.findStopVisit(tx, identifier)
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
		sv, ok := tx.Model().StopVisits().FindByObjectId(obj)
		if ok && sv.Id() != stopVisit.Id() {
			http.Error(response, fmt.Sprintf("Invalid request: stopVisit %v already have an objectid %v", sv.Id(), obj.String()), http.StatusBadRequest)
			return
		}
	}

	tx.Model().StopVisits().Save(&stopVisit)
	err = tx.Commit()
	if err != nil {
		logger.Log.Debugf("Transaction error: %v", err)
		http.Error(response, "Internal error", http.StatusInternalServerError)
		return
	}

	jsonBytes, _ := stopVisit.MarshalJSON()
	response.Write(jsonBytes)
}

func (controller *StopVisitController) Create(response http.ResponseWriter, body []byte) {
	// New transaction
	tx := controller.referential.NewTransaction()
	defer tx.Close()

	//too verbose
	//logger.Log.Debugf("Create stopVisit: %s", string(body))

	stopVisit := tx.Model().StopVisits().New()

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
		sv, ok := tx.Model().StopVisits().FindByObjectId(obj)
		if ok {
			http.Error(response, fmt.Sprintf("Invalid request: stopVisit %v already have an objectid %v", sv.Id(), obj.String()), http.StatusBadRequest)
			return
		}
	}

	tx.Model().StopVisits().Save(&stopVisit)
	err = tx.Commit()
	if err != nil {
		logger.Log.Debugf("Transaction error: %v", err)
		http.Error(response, "Internal error", http.StatusInternalServerError)
		return
	}
	jsonBytes, _ := stopVisit.MarshalJSON()
	response.Write(jsonBytes)
}
