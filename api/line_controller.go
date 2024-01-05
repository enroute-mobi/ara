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

type LineController struct {
	referential *core.Referential
}

func NewLineController(referential *core.Referential) ControllerInterface {
	return &Controller{
		restfulResource: &LineController{
			referential: referential,
		},
	}
}

func (controller *LineController) findLine(identifier string) (*model.Line, bool) {
	foundStrings := idPattern.FindStringSubmatch(identifier)
	if foundStrings != nil {
		code := model.NewCode(foundStrings[1], foundStrings[2])
		return controller.referential.Model().Lines().FindByCode(code)
	}
	return controller.referential.Model().Lines().Find(model.LineId(identifier))
}

func (controller *LineController) Index(response http.ResponseWriter, filters url.Values) {
	logger.Log.Debugf("Lines Index")

	jsonBytes, _ := json.Marshal(controller.referential.Model().Lines().FindAll())
	response.Write(jsonBytes)
}

func (controller *LineController) Show(response http.ResponseWriter, identifier string) {
	line, ok := controller.findLine(identifier)
	if !ok {
		http.Error(response, fmt.Sprintf("Line not found: %s", identifier), http.StatusNotFound)
		return
	}
	logger.Log.Debugf("Get line %s", identifier)

	jsonBytes, _ := line.MarshalJSON()
	response.Write(jsonBytes)
}

func (controller *LineController) Delete(response http.ResponseWriter, identifier string) {
	line, ok := controller.findLine(identifier)
	if !ok {
		http.Error(response, fmt.Sprintf("Line not found: %s", identifier), http.StatusNotFound)
		return
	}
	logger.Log.Debugf("Delete line %s", identifier)

	jsonBytes, _ := line.MarshalJSON()
	controller.referential.Model().Lines().Delete(line)

	response.Write(jsonBytes)
}

func (controller *LineController) Update(response http.ResponseWriter, identifier string, body []byte) {
	line, ok := controller.findLine(identifier)
	if !ok {
		http.Error(response, fmt.Sprintf("Line not found: %s", identifier), http.StatusNotFound)
		return
	}

	logger.Log.Debugf("Update line %s: %s", identifier, string(body))

	err := json.Unmarshal(body, &line)
	if err != nil {
		http.Error(response, fmt.Sprintf("Invalid request: can't parse request body: %v", err), http.StatusBadRequest)
		return
	}

	for _, obj := range line.Codes() {
		l, ok := controller.referential.Model().Lines().FindByCode(obj)
		if ok && l.Id() != line.Id() {
			http.Error(response, fmt.Sprintf("Invalid request: line %v already have a code %v", l.Id(), obj.String()), http.StatusBadRequest)
			return
		}
	}

	controller.referential.Model().Lines().Save(line)
	jsonBytes, _ := json.Marshal(&line)
	response.Write(jsonBytes)
}

func (controller *LineController) Create(response http.ResponseWriter, body []byte) {
	logger.Log.Debugf("Create line: %s", string(body))

	line := controller.referential.Model().Lines().New()

	err := json.Unmarshal(body, &line)
	if err != nil {
		http.Error(response, fmt.Sprintf("Invalid request: can't parse request body: %v", err), http.StatusBadRequest)
		return
	}

	if line.Id() != "" {
		http.Error(response, "Invalid request", http.StatusBadRequest)
		return
	}

	for _, obj := range line.Codes() {
		l, ok := controller.referential.Model().Lines().FindByCode(obj)
		if ok {
			http.Error(response, fmt.Sprintf("Invalid request: line %v already have a code %v", l.Id(), obj.String()), http.StatusBadRequest)
			return
		}
	}

	controller.referential.Model().Lines().Save(line)
	jsonBytes, _ := json.Marshal(&line)
	response.Write(jsonBytes)
}
