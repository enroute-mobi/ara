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

type FacilityController struct {
	referential *core.Referential
}

func NewFacilityController(referential *core.Referential) RestfulResource {
	return &FacilityController{
		referential: referential,
	}
}

func (controller *FacilityController) findFacility(identifier string) (*model.Facility, bool) {
	foundStrings := idPattern.FindStringSubmatch(identifier)
	if foundStrings != nil {
		code := model.NewCode(foundStrings[1], foundStrings[2])
		return controller.referential.Model().Facilities().FindByCode(code)
	}
	return controller.referential.Model().Facilities().Find(model.FacilityId(identifier))
}

func (controller *FacilityController) Index(response http.ResponseWriter, _params url.Values) {
	logger.Log.Debugf("Facilities Index")

	jsonBytes, _ := json.Marshal(controller.referential.Model().Facilities().FindAll())
	response.Write(jsonBytes)
}

func (controller *FacilityController) Show(response http.ResponseWriter, identifier string, _params url.Values) {
	facility, ok := controller.findFacility(identifier)
	if !ok {
		http.Error(response, fmt.Sprintf("Facility not found: %s", identifier), http.StatusNotFound)
		return
	}
	logger.Log.Debugf("Get facility %s", identifier)

	jsonBytes, _ := facility.MarshalJSON()
	response.Write(jsonBytes)
}

func (controller *FacilityController) Delete(response http.ResponseWriter, identifier string) {
	facility, ok := controller.findFacility(identifier)
	if !ok {
		http.Error(response, fmt.Sprintf("Facility not found: %s", identifier), http.StatusNotFound)
		return
	}
	logger.Log.Debugf("Delete facility %s", identifier)

	jsonBytes, _ := facility.MarshalJSON()
	controller.referential.Model().Facilities().Delete(facility)
	response.Write(jsonBytes)
}

func (controller *FacilityController) Update(response http.ResponseWriter, identifier string, body []byte) {
	facility, ok := controller.findFacility(identifier)
	if !ok {
		http.Error(response, fmt.Sprintf("Facility not found: %s", identifier), http.StatusNotFound)
		return
	}

	logger.Log.Debugf("Update facility %s: %s", identifier, string(body))

	err := json.Unmarshal(body, &facility)
	if err != nil {
		http.Error(response, fmt.Sprintf("Invalid request: can't parse request body: %v", err), http.StatusBadRequest)
		return
	}

	for _, obj := range facility.Codes() {
		o, ok := controller.referential.Model().Facilities().FindByCode(obj)
		if ok && o.Id() != facility.Id() {
			http.Error(response, fmt.Sprintf("Invalid request: facility %v already have a code %v", o.Id(), obj.String()), http.StatusBadRequest)
			return
		}
	}

	controller.referential.Model().Facilities().Save(facility)
	jsonBytes, _ := json.Marshal(&facility)
	response.Write(jsonBytes)
}

func (controller *FacilityController) Create(response http.ResponseWriter, body []byte) {
	logger.Log.Debugf("Create facility: %s", string(body))

	facility := controller.referential.Model().Facilities().New()

	err := json.Unmarshal(body, &facility)
	if err != nil {
		http.Error(response, fmt.Sprintf("Invalid request: can't parse request body: %v", err), http.StatusBadRequest)
		return
	}

	if facility.Id() != "" {
		http.Error(response, "Invalid request", http.StatusBadRequest)
		return
	}

	for _, obj := range facility.Codes() {
		o, ok := controller.referential.Model().Facilities().FindByCode(obj)
		if ok {
			http.Error(response, fmt.Sprintf("Invalid request: facility %v already have a code %v", o.Id(), obj.String()), http.StatusBadRequest)
			return
		}
	}

	controller.referential.Model().Facilities().Save(facility)
	jsonBytes, _ := json.Marshal(&facility)
	response.Write(jsonBytes)
}
