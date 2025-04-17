package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"

	"bitbucket.org/enroute-mobi/ara/core"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
)

type StopAreaController struct {
	referential *core.Referential
}

func NewStopAreaController(referential *core.Referential) RestfulResource {
	return &StopAreaController{
		referential: referential,
	}
}

func (controller *StopAreaController) findStopArea(identifier string) (*model.StopArea, bool) {
	foundStrings := idPattern.FindStringSubmatch(identifier)
	if foundStrings != nil {
		code := model.NewCode(foundStrings[1], foundStrings[2])
		return controller.referential.Model().StopAreas().FindByCode(code)
	}
	return controller.referential.Model().StopAreas().Find(model.StopAreaId(identifier))
}

func (controller *StopAreaController) Index(response http.ResponseWriter, params url.Values) {
	logger.Log.Debugf("StopAreas Index")

	allStopAreas := controller.referential.Model().StopAreas().FindAll()
	direction := params.Get("direction")
	switch direction {
	case "desc":
		sort.Slice(allStopAreas, func(i, j int) bool {
			return allStopAreas[i].Name > allStopAreas[j].Name
		})
	case "asc", "":
		sort.Slice(allStopAreas, func(i, j int) bool {
			return allStopAreas[i].Name < allStopAreas[j].Name
		})
	default:
		http.Error(response, fmt.Sprintf("invalid request: query parameter \"direction\": %s", params.Get("direction")), http.StatusBadRequest)
		return
	}

	paginatedStopAreas, err := paginate(allStopAreas, params)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	jsonBytes, _ := json.Marshal(paginatedStopAreas)
	response.Write(jsonBytes)
}

func (controller *StopAreaController) Show(response http.ResponseWriter, identifier string) {
	stopArea, ok := controller.findStopArea(identifier)
	if !ok {
		http.Error(response, fmt.Sprintf("Stop area not found: %s", identifier), http.StatusNotFound)
		return
	}
	logger.Log.Debugf("Get stopArea %s", identifier)

	jsonBytes, _ := stopArea.MarshalJSON()
	response.Write(jsonBytes)
}

func (controller *StopAreaController) Delete(response http.ResponseWriter, identifier string) {
	stopArea, ok := controller.findStopArea(identifier)
	if !ok {
		http.Error(response, fmt.Sprintf("Stop area not found: %s", identifier), http.StatusNotFound)
		return
	}
	logger.Log.Debugf("Delete stopArea %s", identifier)

	jsonBytes, _ := stopArea.MarshalJSON()
	controller.referential.Model().StopAreas().Delete(stopArea)
	response.Write(jsonBytes)
}

func (controller *StopAreaController) Update(response http.ResponseWriter, identifier string, body []byte) {
	stopArea, ok := controller.findStopArea(identifier)
	if !ok {
		http.Error(response, fmt.Sprintf("Stop area not found: %s", identifier), http.StatusNotFound)
		return
	}

	logger.Log.Debugf("Update stopArea %s: %s", identifier, string(body))

	err := json.Unmarshal(body, &stopArea)
	if err != nil {
		http.Error(response, fmt.Sprintf("Invalid request: can't parse request body: %v", err), http.StatusBadRequest)
		return
	}

	for _, obj := range stopArea.Codes() {
		sa, ok := controller.referential.Model().StopAreas().FindByCode(obj)
		if ok && sa.Id() != stopArea.Id() {
			http.Error(response, fmt.Sprintf("Invalid request: stopArea %v already have a code %v", sa.Id(), obj.String()), http.StatusBadRequest)
			return
		}
	}

	controller.referential.Model().StopAreas().Save(stopArea)
	jsonBytes, _ := stopArea.MarshalJSON()
	response.Write(jsonBytes)
}

func (controller *StopAreaController) Create(response http.ResponseWriter, body []byte) {
	logger.Log.Debugf("Create stopArea: %s", string(body))

	stopArea := controller.referential.Model().StopAreas().New()

	err := json.Unmarshal(body, &stopArea)
	if err != nil {
		http.Error(response, fmt.Sprintf("Invalid request: can't parse request body: %v", err), http.StatusBadRequest)
		return
	}

	if stopArea.Id() != "" {
		http.Error(response, "Invalid request", http.StatusBadRequest)
		return
	}

	for _, obj := range stopArea.Codes() {
		sa, ok := controller.referential.Model().StopAreas().FindByCode(obj)
		if ok {
			http.Error(response, fmt.Sprintf("Invalid request: stopArea %v already have a code %v", sa.Id(), obj.String()), http.StatusBadRequest)
			return
		}
	}

	controller.referential.Model().StopAreas().Save(stopArea)
	jsonBytes, _ := stopArea.MarshalJSON()
	response.Write(jsonBytes)
}
