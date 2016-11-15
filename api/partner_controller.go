package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
)

type PartnerController struct {
	ControllerReferential
}

func NewPartnerController() (controller *Controller) {
	return &Controller{
		ressourceController: &PartnerController{},
	}
}

func (controller *PartnerController) Ressources() string {
	return "partners"
}

func (controller *PartnerController) Index(response http.ResponseWriter) {
	logger.Log.Debugf("Partners Index")

	jsonBytes, _ := json.Marshal(controller.referential.Partners().FindAll())
	response.Write(jsonBytes)
}

func (controller *PartnerController) Show(response http.ResponseWriter, identifier string) {
	partner := controller.referential.Partners().Find(model.PartnerId(identifier))
	if partner == nil {
		http.Error(response, fmt.Sprintf("Partner not found: %s", identifier), 500)
		return
	}
	logger.Log.Debugf("Get partner %s", identifier)

	jsonBytes, _ := partner.MarshalJSON()
	response.Write(jsonBytes)
}

func (controller *PartnerController) Delete(response http.ResponseWriter, identifier string) {
	partner := controller.referential.Partners().Find(model.PartnerId(identifier))
	if partner == nil {
		http.Error(response, fmt.Sprintf("Partner not found: %s", identifier), 500)
		return
	}
	logger.Log.Debugf("Delete partner %s", identifier)

	jsonBytes, _ := partner.MarshalJSON()
	controller.referential.Partners().Delete(partner)
	response.Write(jsonBytes)
}

func (controller *PartnerController) Update(response http.ResponseWriter, identifier string, body []byte) {
	partner := controller.referential.Partners().Find(model.PartnerId(identifier))
	if partner == nil {
		http.Error(response, fmt.Sprintf("Partner not found: %s", identifier), 500)
		return
	}

	logger.Log.Debugf("Update partner %s: %s", identifier, string(body))

	err := json.Unmarshal(body, &partner)
	if err != nil {
		http.Error(response, fmt.Sprintf("Invalid request: %v", err), 400)
		return
	}

	partner.Save()
	jsonBytes, _ := partner.MarshalJSON()
	response.Write(jsonBytes)
}

func (controller *PartnerController) Create(response http.ResponseWriter, body []byte) {
	logger.Log.Debugf("Create partner: %s", string(body))

	partner := model.Partner{}
	err := json.Unmarshal(body, &partner)
	if err != nil {
		http.Error(response, fmt.Sprintf("Invalid request: %v", err), 400)
		return
	}
	if partner.Id() != "" {
		http.Error(response, "Invalid request (Id specified)", 400)
		return
	}

	controller.referential.Partners().Save(&partner)
	jsonBytes, _ := partner.MarshalJSON()
	response.Write(jsonBytes)
}
