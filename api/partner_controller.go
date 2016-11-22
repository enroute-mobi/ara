package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/af83/edwig/core"
	"github.com/af83/edwig/logger"
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
	partner := controller.referential.Partners().Find(core.PartnerId(identifier))
	if partner == nil {
		http.Error(response, fmt.Sprintf("Partner not found: %s", identifier), 500)
		return
	}
	logger.Log.Debugf("Get partner %s", identifier)

	jsonBytes, _ := partner.MarshalJSON()
	response.Write(jsonBytes)
}

func (controller *PartnerController) Delete(response http.ResponseWriter, identifier string) {
	partner := controller.referential.Partners().Find(core.PartnerId(identifier))
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
	partner := controller.referential.Partners().Find(core.PartnerId(identifier))
	if partner == nil {
		http.Error(response, fmt.Sprintf("Partner not found: %s", identifier), 500)
		return
	}

	logger.Log.Debugf("Update partner %s: %s", identifier, string(body))

	apiPartner := partner.Definition()
	err := json.Unmarshal(body, apiPartner)
	if err != nil {
		http.Error(response, "Invalid request: can't parse body", 400)
		return
	}

	if !apiPartner.Validate() {
		jsonBytes, _ := json.Marshal(apiPartner)
		response.WriteHeader(http.StatusBadRequest)
		response.Write(jsonBytes)
		return
	}

	partner.SetDefinition(apiPartner)
	partner.Save()
	jsonBytes, _ := partner.MarshalJSON()
	response.Write(jsonBytes)
}

func (controller *PartnerController) Create(response http.ResponseWriter, body []byte) {
	logger.Log.Debugf("Create partner: %s", string(body))

	partner := core.Partner{}
	apiPartner := partner.Definition()
	err := json.Unmarshal(body, apiPartner)
	if err != nil {
		http.Error(response, "Invalid request: can't parse body", 400)
		return
	}
	if partner.Id() != "" {
		http.Error(response, "Invalid request (Id specified)", 400)
		return
	}

	if !apiPartner.Validate() {
		jsonBytes, _ := json.Marshal(apiPartner)
		response.WriteHeader(http.StatusBadRequest)
		response.Write(jsonBytes)
		return
	}

	partner.SetDefinition(apiPartner)
	controller.referential.Partners().Save(&partner)
	jsonBytes, _ := partner.MarshalJSON()
	response.Write(jsonBytes)
}