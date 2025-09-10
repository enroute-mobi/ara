package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"bitbucket.org/enroute-mobi/ara/core"
	"bitbucket.org/enroute-mobi/ara/core/partners"
)

type PartnerTemplateController struct {
	referential *core.Referential
}

func NewPartnerTemplateController(referential *core.Referential) RestfulResource {
	return &PartnerTemplateController{
		referential: referential,
	}
}

func (controller *PartnerTemplateController) findPartnerTemplate(identifier string) *core.PartnerTemplate {
	pt := controller.referential.PartnerTemplates().FindBySlug(partners.Slug(identifier))
	if pt != nil {
		return pt
	}
	return controller.referential.PartnerTemplates().Find(partners.Id(identifier))
}

func (controller *PartnerTemplateController) Index(response http.ResponseWriter, _params url.Values) {
	jsonBytes, _ := json.Marshal(controller.referential.PartnerTemplates())
	response.Write(jsonBytes)
}

func (controller *PartnerTemplateController) Show(response http.ResponseWriter, identifier string) {
	pt := controller.findPartnerTemplate(identifier)
	if pt == nil {
		http.Error(response, fmt.Sprintf("Partner template not found: %s", identifier), http.StatusNotFound)
		return
	}

	jsonBytes, _ := pt.MarshalJSON()
	response.Write(jsonBytes)
}

func (controller *PartnerTemplateController) Delete(response http.ResponseWriter, identifier string) {
	pt := controller.findPartnerTemplate(identifier)
	if pt == nil {
		http.Error(response, fmt.Sprintf("Partner template not found: %s", identifier), http.StatusNotFound)
		return
	}

	jsonBytes, _ := pt.MarshalJSON()
	controller.referential.PartnerTemplates().Delete(pt)
	response.Write(jsonBytes)
}

func (controller *PartnerTemplateController) Update(response http.ResponseWriter, identifier string, body []byte) {
	pt := controller.findPartnerTemplate(identifier)
	if pt == nil {
		http.Error(response, fmt.Sprintf("Partner template not found: %s", identifier), http.StatusNotFound)
		return
	}

	newPt := pt.Copy()
	err := json.Unmarshal(body, newPt)
	if err != nil {
		http.Error(response, fmt.Sprintf("Invalid request: can't parse request body: %v", err), http.StatusBadRequest)
		return
	}

	if !newPt.Validate() {
		jsonBytes, _ := json.Marshal(newPt)
		response.WriteHeader(http.StatusBadRequest)
		response.Write(jsonBytes)
		return
	}

	newPt.Save()

	jsonBytes, _ := pt.MarshalJSON()
	response.Write(jsonBytes)
}

func (controller *PartnerTemplateController) Create(response http.ResponseWriter, body []byte) {
	pt := controller.referential.PartnerTemplates().New("")
	err := json.Unmarshal(body, pt)
	if err != nil {
		http.Error(response, fmt.Sprintf("Invalid request: can't parse request body: %v", err), http.StatusBadRequest)
		return
	}

	if !pt.Validate() {
		jsonBytes, _ := json.Marshal(pt)
		response.WriteHeader(http.StatusBadRequest)
		response.Write(jsonBytes)
		return
	}

	controller.referential.PartnerTemplates().Save(pt)
	jsonBytes, _ := pt.MarshalJSON()
	response.Write(jsonBytes)
}

// func (controller *PartnerTemplateController) Save(response http.ResponseWriter) {

// 	status, err := controller.referential.PartnerTemplates().SaveToDatabase()

// 	if err != nil {
// 		monitoring.ReportError(err)

// 		response.WriteHeader(status)
// 		jsonBytes, _ := json.Marshal(map[string]string{"error": err.Error()})
// 		response.Write(jsonBytes)
// 	}
// }
