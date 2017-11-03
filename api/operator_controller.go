package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"

	"github.com/af83/edwig/core"
	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
)

type OperatorController struct {
	referential *core.Referential
}

func NewOperatorController(referential *core.Referential) ControllerInterface {
	return &Controller{
		restfulRessource: &OperatorController{
			referential: referential,
		},
	}
}

func (controller *OperatorController) findOperator(tx *model.Transaction, identifier string) (model.Operator, bool) {
	idRegexp := "([0-9a-zA-Z-]+):([0-9a-zA-Z-]+)"
	pattern := regexp.MustCompile(idRegexp)
	foundStrings := pattern.FindStringSubmatch(identifier)
	if foundStrings != nil {
		objectid := model.NewObjectID(foundStrings[1], foundStrings[2])
		return tx.Model().Operators().FindByObjectId(objectid)
	}
	return tx.Model().Operators().Find(model.OperatorId(identifier))
}

func (controller *OperatorController) Index(response http.ResponseWriter, filters url.Values) {
	tx := controller.referential.NewTransaction()
	defer tx.Close()

	logger.Log.Debugf("Operators Index")

	jsonBytes, _ := json.Marshal(tx.Model().Operators().FindAll())
	response.Write(jsonBytes)
}

func (controller *OperatorController) Show(response http.ResponseWriter, identifier string) {
	tx := controller.referential.NewTransaction()
	defer tx.Close()

	operator, ok := controller.findOperator(tx, identifier)
	if !ok {
		http.Error(response, fmt.Sprintf("Operator not found: %s", identifier), 404)
		return
	}
	logger.Log.Debugf("Get operator %s", identifier)

	jsonBytes, _ := operator.MarshalJSON()
	response.Write(jsonBytes)
}

func (controller *OperatorController) Delete(response http.ResponseWriter, identifier string) {
	// New transaction
	tx := controller.referential.NewTransaction()
	defer tx.Close()

	operator, ok := controller.findOperator(tx, identifier)
	if !ok {
		http.Error(response, fmt.Sprintf("Operator not found: %s", identifier), 404)
		return
	}
	logger.Log.Debugf("Delete operator %s", identifier)

	jsonBytes, _ := operator.MarshalJSON()
	tx.Model().Operators().Delete(&operator)
	err := tx.Commit()
	if err != nil {
		logger.Log.Debugf("Transaction error: %v", err)
		http.Error(response, "Internal error", 500)
		return
	}
	response.Write(jsonBytes)
}

func (controller *OperatorController) Update(response http.ResponseWriter, identifier string, body []byte) {
	// New transaction
	tx := controller.referential.NewTransaction()
	defer tx.Close()

	operator, ok := controller.findOperator(tx, identifier)
	if !ok {
		http.Error(response, fmt.Sprintf("Operator not found: %s", identifier), 404)
		return
	}

	logger.Log.Debugf("Update operator %s: %s", identifier, string(body))

	err := json.Unmarshal(body, &operator)
	if err != nil {
		http.Error(response, fmt.Sprintf("Invalid request: can't parse request body: %v", err), 400)
		return
	}

	tx.Model().Operators().Save(&operator)
	err = tx.Commit()
	if err != nil {
		logger.Log.Debugf("Transaction error: %v", err)
		http.Error(response, "Internal error", 500)
		return
	}
	jsonBytes, _ := json.Marshal(&operator)
	response.Write(jsonBytes)
}

func (controller *OperatorController) Create(response http.ResponseWriter, body []byte) {
	// New transaction
	tx := controller.referential.NewTransaction()
	defer tx.Close()

	logger.Log.Debugf("Create operator: %s", string(body))

	operator := tx.Model().Operators().New()

	err := json.Unmarshal(body, &operator)
	if err != nil {
		http.Error(response, fmt.Sprintf("Invalid request: can't parse request body: %v", err), 400)
		return
	}

	if operator.Id() != "" {
		http.Error(response, "Invalid request", 400)
		return
	}

	tx.Model().Operators().Save(&operator)
	err = tx.Commit()
	if err != nil {
		logger.Log.Debugf("Transaction error: %v", err)
		http.Error(response, "Internal error", 500)
		return
	}
	jsonBytes, _ := json.Marshal(&operator)
	response.Write(jsonBytes)
}
