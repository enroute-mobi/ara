package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"bitbucket.org/enroute-mobi/ara/core"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type SituationController struct {
	referential *core.Referential
}

func NewSituationController(referential *core.Referential) ControllerInterface {
	return &Controller{
		restfulResource: &SituationController{
			referential: referential,
		},
	}
}

func (controller *SituationController) findSituation(identifier string) (model.Situation, bool) {
	foundStrings := idPattern.FindStringSubmatch(identifier)
	if foundStrings != nil {
		code := model.NewCode(foundStrings[1], foundStrings[2])
		return controller.referential.Model().Situations().FindByCode(code)
	}
	return controller.referential.Model().Situations().Find(model.SituationId(identifier))
}

func (controller *SituationController) Index(response http.ResponseWriter, filters url.Values) {
	logger.Log.Debugf("Situations Index")

	jsonBytes, _ := json.Marshal(controller.referential.Model().Situations().FindAll())
	response.Write(jsonBytes)
}

func (controller *SituationController) Show(response http.ResponseWriter, identifier string) {
	situation, ok := controller.findSituation(identifier)
	if !ok {
		http.Error(response, fmt.Sprintf("Situation not found: %s", identifier), http.StatusNotFound)
		return
	}
	logger.Log.Debugf("Get situation %s", identifier)

	jsonBytes, _ := situation.MarshalJSON()
	response.Write(jsonBytes)
}

func (controller *SituationController) Delete(response http.ResponseWriter, identifier string) {
	ctx := context.Background()
	span, _ := tracer.StartSpanFromContext(ctx, "api.situations.delete")
	defer span.Finish()

	span.SetTag("situation_id", identifier)

	situation, ok := controller.findSituation(identifier)
	if !ok {
		http.Error(response, fmt.Sprintf("Situation not found: %s", identifier), http.StatusNotFound)
		return
	}
	logger.Log.Debugf("Delete situation %s", identifier)

	jsonBytes, _ := situation.MarshalJSON()
	controller.referential.Model().Situations().Delete(&situation)
	response.Write(jsonBytes)
}

func (controller *SituationController) Update(response http.ResponseWriter, identifier string, body []byte) {
	ctx := context.Background()
	span, _ := tracer.StartSpanFromContext(ctx, "api.situations.update")
	defer span.Finish()

	span.SetTag("situation", string(body))

	situation, ok := controller.findSituation(identifier)
	if !ok {
		http.Error(response, fmt.Sprintf("Situation not found: %s", identifier), http.StatusNotFound)
		return
	}

	logger.Log.Debugf("Update situation %s: %s", identifier, string(body))

	apiSituation := situation.Definition()
	err := json.Unmarshal(body, &apiSituation)
	if err != nil {
		http.Error(response, fmt.Sprintf("Invalid request: can't parse request body: %v", err), http.StatusBadRequest)
		return
	}

	if apiSituation.Id != situation.Id() {
		http.Error(response, "Invalid request (Id specified)", http.StatusBadRequest)
		return
	}

	code := model.NewCode(apiSituation.CodeSpace, apiSituation.SituationNumber)
	s, found := controller.referential.Model().Situations().FindByCode(code)
	if found && s.Id() != situation.Id() {
		apiSituation.ExistingSituationCode = true
	}

	if !apiSituation.Validate() && !apiSituation.IgnoreValidation {
		jsonBytes, _ := json.Marshal(apiSituation)
		response.WriteHeader(http.StatusBadRequest)
		response.Write(jsonBytes)
		return
	}

	situation.SetDefinition(apiSituation)
	controller.referential.Model().Situations().Save(&situation)
	jsonBytes, _ := json.Marshal(&situation)

	response.Write(jsonBytes)
}

func (controller *SituationController) Create(response http.ResponseWriter, body []byte) {
	logger.Log.Debugf("Create situation: %s", string(body))

	ctx := context.Background()
	span, _ := tracer.StartSpanFromContext(ctx, "api.situations.create")
	defer span.Finish()

	span.SetTag("situation", string(body))

	situation := controller.referential.Model().Situations().New()
	apiSituation := situation.Definition()
	err := json.Unmarshal(body, &apiSituation)

	if err != nil {
		http.Error(response, fmt.Sprintf("Invalid request: can't parse request body: %v", err), http.StatusBadRequest)
		return
	}

	if apiSituation.Id != "" {
		http.Error(response, "Invalid request", http.StatusBadRequest)
		return
	}

	code := model.NewCode(apiSituation.CodeSpace, apiSituation.SituationNumber)
	_, ok := controller.referential.Model().Situations().FindByCode(code)
	if ok {
		apiSituation.ExistingSituationCode = true
	}

	if !apiSituation.Validate() && !apiSituation.IgnoreValidation {
		jsonBytes, _ := json.Marshal(apiSituation)
		response.WriteHeader(http.StatusBadRequest)
		response.Write(jsonBytes)
		return
	}

	situation.SetDefinition(apiSituation)
	controller.referential.Model().Situations().Save(&situation)
	jsonBytes, _ := json.Marshal(&situation)

	response.Write(jsonBytes)
}
