package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/af83/edwig/core"
	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
)

type LineController struct {
	referential *core.Referential
}

func NewLineController(referential *core.Referential) ControllerInterface {
	return &Controller{
		restfulRessource: &LineController{
			referential: referential,
		},
	}
}

func (controller *LineController) Index(response http.ResponseWriter) {
	logger.Log.Debugf("Lines Index")

	jsonBytes, _ := json.Marshal(controller.referential.Model().Lines().FindAll())
	response.Write(jsonBytes)
}

func (controller *LineController) Show(response http.ResponseWriter, identifier string) {
	line, ok := controller.referential.Model().Lines().Find(model.LineId(identifier))
	if !ok {
		http.Error(response, fmt.Sprintf("Line not found: %s", identifier), 500)
		return
	}
	logger.Log.Debugf("Get line %s", identifier)

	jsonBytes, _ := line.MarshalJSON()
	response.Write(jsonBytes)
}

func (controller *LineController) Delete(response http.ResponseWriter, identifier string) {
	// New transaction
	tx := controller.referential.NewTransaction()
	defer tx.Close()

	line, ok := tx.Model().Lines().Find(model.LineId(identifier))
	if !ok {
		http.Error(response, fmt.Sprintf("Line not found: %s", identifier), 500)
		return
	}
	logger.Log.Debugf("Delete line %s", identifier)

	jsonBytes, _ := line.MarshalJSON()
	tx.Model().Lines().Delete(&line)
	err := tx.Commit()
	if err != nil {
		logger.Log.Debugf("Transaction error: %v", err)
		http.Error(response, "Internal error", 500)
		return
	}
	response.Write(jsonBytes)
}

func (controller *LineController) Update(response http.ResponseWriter, identifier string, body []byte) {
	// New transaction
	tx := controller.referential.NewTransaction()
	defer tx.Close()

	line, ok := tx.Model().Lines().Find(model.LineId(identifier))
	if !ok {
		http.Error(response, fmt.Sprintf("Line not found: %s", identifier), 500)
		return
	}

	logger.Log.Debugf("Update line %s: %s", identifier, string(body))

	err := json.Unmarshal(body, &line)
	if err != nil {
		http.Error(response, "Invalid request: can't parse request body", 400)
		return
	}

	tx.Model().Lines().Save(&line)
	err = tx.Commit()
	if err != nil {
		logger.Log.Debugf("Transaction error: %v", err)
		http.Error(response, "Internal error", 500)
		return
	}
	jsonBytes, _ := json.Marshal(&line)
	response.Write(jsonBytes)
}

func (controller *LineController) Create(response http.ResponseWriter, body []byte) {
	// New transaction
	tx := controller.referential.NewTransaction()
	defer tx.Close()

	logger.Log.Debugf("Create line: %s", string(body))

	line := tx.Model().Lines().New()

	err := json.Unmarshal(body, &line)
	if err != nil {
		http.Error(response, "Invalid request: can't parse request body", 400)
		return
	}

	if line.Id() != "" {
		http.Error(response, "Invalid request", 400)
		return
	}

	tx.Model().Lines().Save(&line)
	err = tx.Commit()
	if err != nil {
		logger.Log.Debugf("Transaction error: %v", err)
		http.Error(response, "Internal error", 500)
		return
	}
	jsonBytes, _ := json.Marshal(&line)
	response.Write(jsonBytes)
}
