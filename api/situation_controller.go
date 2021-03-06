package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"

	"bitbucket.org/enroute-mobi/ara/core"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
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

func (controller *SituationController) findSituation(tx *model.Transaction, identifier string) (model.Situation, bool) {
	idRegexp := "([0-9a-zA-Z-]+):([0-9a-zA-Z-]+)"
	pattern := regexp.MustCompile(idRegexp)
	foundStrings := pattern.FindStringSubmatch(identifier)
	if foundStrings != nil {
		objectid := model.NewObjectID(foundStrings[1], foundStrings[2])
		return tx.Model().Situations().FindByObjectId(objectid)
	}
	return tx.Model().Situations().Find(model.SituationId(identifier))
}

func (controller *SituationController) Index(response http.ResponseWriter, filters url.Values) {
	tx := controller.referential.NewTransaction()
	defer tx.Close()

	logger.Log.Debugf("Situations Index")

	jsonBytes, _ := json.Marshal(tx.Model().Situations().FindAll())
	response.Write(jsonBytes)
}

func (controller *SituationController) Show(response http.ResponseWriter, identifier string) {
	tx := controller.referential.NewTransaction()
	defer tx.Close()

	situation, ok := controller.findSituation(tx, identifier)
	if !ok {
		http.Error(response, fmt.Sprintf("Situation not found: %s", identifier), http.StatusNotFound)
		return
	}
	logger.Log.Debugf("Get situation %s", identifier)

	jsonBytes, _ := situation.MarshalJSON()
	response.Write(jsonBytes)
}

func (controller *SituationController) Delete(response http.ResponseWriter, identifier string) {
	// New transaction
	tx := controller.referential.NewTransaction()
	defer tx.Close()

	situation, ok := controller.findSituation(tx, identifier)
	if !ok {
		http.Error(response, fmt.Sprintf("Situation not found: %s", identifier), http.StatusNotFound)
		return
	}
	logger.Log.Debugf("Delete situation %s", identifier)

	jsonBytes, _ := situation.MarshalJSON()
	tx.Model().Situations().Delete(&situation)
	err := tx.Commit()
	if err != nil {
		logger.Log.Debugf("Transaction error: %v", err)
		http.Error(response, "Internal error", http.StatusInternalServerError)
		return
	}
	response.Write(jsonBytes)
}

func (controller *SituationController) Update(response http.ResponseWriter, identifier string, body []byte) {
	// New transaction
	tx := controller.referential.NewTransaction()
	defer tx.Close()

	situation, ok := controller.findSituation(tx, identifier)
	if !ok {
		http.Error(response, fmt.Sprintf("Situation not found: %s", identifier), http.StatusNotFound)
		return
	}

	logger.Log.Debugf("Update situation %s: %s", identifier, string(body))

	err := json.Unmarshal(body, &situation)
	if err != nil {
		http.Error(response, fmt.Sprintf("Invalid request: can't parse request body: %v", err), http.StatusBadRequest)
		return
	}

	tx.Model().Situations().Save(&situation)
	err = tx.Commit()
	if err != nil {
		logger.Log.Debugf("Transaction error: %v", err)
		http.Error(response, "Internal error", http.StatusInternalServerError)
		return
	}
	jsonBytes, _ := json.Marshal(&situation)
	response.Write(jsonBytes)
}

func (controller *SituationController) Create(response http.ResponseWriter, body []byte) {
	// New transaction
	tx := controller.referential.NewTransaction()
	defer tx.Close()

	logger.Log.Debugf("Create situation: %s", string(body))

	situation := tx.Model().Situations().New()

	err := json.Unmarshal(body, &situation)
	if err != nil {
		http.Error(response, fmt.Sprintf("Invalid request: can't parse request body: %v", err), http.StatusBadRequest)
		return
	}

	if situation.Id() != "" {
		http.Error(response, "Invalid request", http.StatusBadRequest)
		return
	}

	tx.Model().Situations().Save(&situation)
	err = tx.Commit()
	if err != nil {
		logger.Log.Debugf("Transaction error: %v", err)
		http.Error(response, "Internal error", http.StatusInternalServerError)
		return
	}
	jsonBytes, _ := json.Marshal(&situation)
	response.Write(jsonBytes)
}
