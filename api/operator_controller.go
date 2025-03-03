package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"bitbucket.org/enroute-mobi/ara/core"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
)

type OperatorController struct {
	referential *core.Referential
}

func NewOperatorController(referential *core.Referential) RestfulResource {
	return &OperatorController{
		referential: referential,
	}
}

func (controller *OperatorController) findOperator(identifier string) (*model.Operator, bool) {
	foundStrings := idPattern.FindStringSubmatch(identifier)
	if foundStrings != nil {
		code := model.NewCode(foundStrings[1], foundStrings[2])
		return controller.referential.Model().Operators().FindByCode(code)
	}
	return controller.referential.Model().Operators().Find(model.OperatorId(identifier))
}

func (controller *OperatorController) Index(response http.ResponseWriter) {
	logger.Log.Debugf("Operators Index")

	jsonBytes, _ := json.Marshal(controller.referential.Model().Operators().FindAll())
	response.Write(jsonBytes)
}

func (controller *OperatorController) Show(response http.ResponseWriter, identifier string) {
	operator, ok := controller.findOperator(identifier)
	if !ok {
		http.Error(response, fmt.Sprintf("Operator not found: %s", identifier), http.StatusNotFound)
		return
	}
	logger.Log.Debugf("Get operator %s", identifier)

	jsonBytes, _ := operator.MarshalJSON()
	response.Write(jsonBytes)
}

func (controller *OperatorController) Delete(response http.ResponseWriter, identifier string) {
	operator, ok := controller.findOperator(identifier)
	if !ok {
		http.Error(response, fmt.Sprintf("Operator not found: %s", identifier), http.StatusNotFound)
		return
	}
	logger.Log.Debugf("Delete operator %s", identifier)

	jsonBytes, _ := operator.MarshalJSON()
	controller.referential.Model().Operators().Delete(operator)
	response.Write(jsonBytes)
}

func (controller *OperatorController) Update(response http.ResponseWriter, identifier string, body []byte) {
	operator, ok := controller.findOperator(identifier)
	if !ok {
		http.Error(response, fmt.Sprintf("Operator not found: %s", identifier), http.StatusNotFound)
		return
	}

	logger.Log.Debugf("Update operator %s: %s", identifier, string(body))

	err := json.Unmarshal(body, &operator)
	if err != nil {
		http.Error(response, fmt.Sprintf("Invalid request: can't parse request body: %v", err), http.StatusBadRequest)
		return
	}

	for _, obj := range operator.Codes() {
		o, ok := controller.referential.Model().Operators().FindByCode(obj)
		if ok && o.Id() != operator.Id() {
			http.Error(response, fmt.Sprintf("Invalid request: operator %v already have a code %v", o.Id(), obj.String()), http.StatusBadRequest)
			return
		}
	}

	controller.referential.Model().Operators().Save(operator)
	jsonBytes, _ := json.Marshal(&operator)
	response.Write(jsonBytes)
}

func (controller *OperatorController) Create(response http.ResponseWriter, body []byte) {
	logger.Log.Debugf("Create operator: %s", string(body))

	operator := controller.referential.Model().Operators().New()

	err := json.Unmarshal(body, &operator)
	if err != nil {
		http.Error(response, fmt.Sprintf("Invalid request: can't parse request body: %v", err), http.StatusBadRequest)
		return
	}

	if operator.Id() != "" {
		http.Error(response, "Invalid request", http.StatusBadRequest)
		return
	}

	for _, obj := range operator.Codes() {
		o, ok := controller.referential.Model().Operators().FindByCode(obj)
		if ok {
			http.Error(response, fmt.Sprintf("Invalid request: operator %v already have a code %v", o.Id(), obj.String()), http.StatusBadRequest)
			return
		}
	}

	controller.referential.Model().Operators().Save(operator)
	jsonBytes, _ := json.Marshal(&operator)
	response.Write(jsonBytes)
}
