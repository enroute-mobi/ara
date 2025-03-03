package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"bitbucket.org/enroute-mobi/ara/core"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/monitoring"
)

type PartnerController struct {
	referential *core.Referential
}

func NewPartnerController(referential *core.Referential) RestfulResource {
	return &PartnerController{
		referential: referential,
	}
}

func (controller *PartnerController) getActionId(action, url string) (string, bool) {
	index := strings.LastIndex(url, action) + len(action) + 1
	if index >= len(url)-1 {
		return "", false
	}
	id := url[index:]

	sz := len(id)
	if sz > 0 && id[sz-1] == '/' {
		id = id[:sz-1]
	}
	return id, true
}

func (controller *PartnerController) subscriptionsIndex(response http.ResponseWriter, requestData *RequestData) {
	partner := controller.findPartner(requestData.Id)
	if partner == nil {
		http.Error(response, fmt.Sprintf("Partner not found: %s", requestData.Id), http.StatusInternalServerError)
		return
	}
	logger.Log.Debugf("Get partner %s for Subscriptions", requestData.Id)

	subscriptions := partner.Subscriptions()
	jsonBytes, _ := json.Marshal(subscriptions.FindAll())
	response.Write(jsonBytes)
}

func (controller *PartnerController) subscriptionsDelete(response http.ResponseWriter, requestData *RequestData) {
	partner := controller.findPartner(requestData.Id)
	if partner == nil {
		http.Error(response, fmt.Sprintf("Partner not found: %s", requestData.Id), http.StatusInternalServerError)
		return
	}
	logger.Log.Debugf("Get partner %s for Subscriptions", requestData.Id)

	id, ok := controller.getActionId(requestData.Action, requestData.Url)
	if !ok {
		http.Error(response, "Invalid request, id can't be nil", http.StatusBadRequest)
		return
	}
	partner.Subscriptions().DeleteById(core.SubscriptionId(id))
}

func (controller *PartnerController) subscriptionsCreate(response http.ResponseWriter, requestData *RequestData) {
	logger.Log.Debugf("Create Subscription: %s", string(requestData.Body))

	partner := controller.findPartner(requestData.Id)
	if partner == nil {
		http.Error(response, fmt.Sprintf("Partner not found: %s", requestData.Id), http.StatusInternalServerError)
		return
	}

	subscription := partner.Subscriptions().New("")
	apiSubscription := core.APISubscription{}

	err := json.Unmarshal(requestData.Body, &apiSubscription)
	if err != nil {
		http.Error(response, fmt.Sprintf("Invalid request: can't parse request body: %v", err), http.StatusBadRequest)
		return
	}

	subscription.SetDefinition(&apiSubscription)
	subscription.Save()

	// Subscribe the Resources for test immediately
	if apiSubscription.SubscribeResourcesNow {
		resources := subscription.Resources(subscription.Clock().Now())
		// sort the map to ensure consistency with the tests
		sort.Slice(resources, func(i, j int) bool {
			return resources[i].Reference.Code.Value() < resources[j].Reference.Code.Value()
		})

		for k, resource := range resources {
			// delay each Resource
			resource.Subscribed(subscription.Clock().Now().Add(time.Duration(k*40) * time.Second))
		}
	}

	jsonBytes, _ := subscription.MarshalJSON()
	response.Write(jsonBytes)
}

func (controller *PartnerController) subscriptions(response http.ResponseWriter, requestData *RequestData) {
	switch requestData.Method {
	case "GET":
		controller.subscriptionsIndex(response, requestData)
	case "POST":
		controller.subscriptionsCreate(response, requestData)
	case "DELETE":
		controller.subscriptionsDelete(response, requestData)
	}
}

func (controller *PartnerController) Action(response http.ResponseWriter, requestData *RequestData) {
	if requestData.Action == "subscriptions" {
		controller.subscriptions(response, requestData)
		return
	}
	http.Error(response, fmt.Sprintf("Action not supported: %s", requestData.Action), http.StatusInternalServerError)
}

func (controller *PartnerController) findPartner(identifier string) *core.Partner {
	partner, ok := controller.referential.Partners().FindBySlug(core.PartnerSlug(identifier))
	if ok {
		return partner
	}
	return controller.referential.Partners().Find(core.PartnerId(identifier))
}

func (controller *PartnerController) Index(response http.ResponseWriter) {
	logger.Log.Debugf("Partners Index")

	jsonBytes, _ := json.Marshal(controller.referential.Partners().FindAll())
	response.Write(jsonBytes)
}

func (controller *PartnerController) Show(response http.ResponseWriter, identifier string) {
	partner := controller.findPartner(identifier)
	if partner == nil {
		http.Error(response, fmt.Sprintf("Partner not found: %s", identifier), http.StatusNotFound)
		return
	}
	logger.Log.Debugf("Get partner %s", identifier)

	jsonBytes, _ := partner.MarshalJSON()
	response.Write(jsonBytes)
}

func (controller *PartnerController) Delete(response http.ResponseWriter, identifier string) {
	partner := controller.findPartner(identifier)
	if partner == nil {
		http.Error(response, fmt.Sprintf("Partner not found: %s", identifier), http.StatusNotFound)
		return
	}
	logger.Log.Debugf("Delete partner %s", identifier)

	jsonBytes, _ := partner.MarshalJSON()
	partner.Stop()
	controller.referential.Partners().Delete(partner)
	response.Write(jsonBytes)
}

func (controller *PartnerController) Update(response http.ResponseWriter, identifier string, body []byte) {
	partner := controller.findPartner(identifier)
	if partner == nil {
		http.Error(response, fmt.Sprintf("Partner not found: %s", identifier), http.StatusNotFound)
		return
	}

	logger.Log.Debugf("Update partner %s: %s", identifier, string(body))

	apiPartner := partner.Definition()
	err := json.Unmarshal(body, apiPartner)
	if err != nil {
		http.Error(response, fmt.Sprintf("Invalid request: can't parse request body: %v", err), http.StatusBadRequest)
		return
	}
	if apiPartner.Id != partner.Id() {
		http.Error(response, "Invalid request (Id specified)", http.StatusBadRequest)
		return
	}

	if !apiPartner.Validate() {
		jsonBytes, _ := json.Marshal(apiPartner)
		response.WriteHeader(http.StatusBadRequest)
		response.Write(jsonBytes)
		return
	}

	partner.Stop()
	partner.SetDefinition(apiPartner)
	partner.Save()
	partner.Start()

	jsonBytes, _ := partner.MarshalJSON()
	response.Write(jsonBytes)
}

func (controller *PartnerController) Create(response http.ResponseWriter, body []byte) {
	logger.Log.Debugf("Create partner: %s", string(body))

	partner := controller.referential.Partners().New("")
	apiPartner := partner.Definition()
	err := json.Unmarshal(body, apiPartner)
	if err != nil {
		http.Error(response, fmt.Sprintf("Invalid request: can't parse request body: %v", err), http.StatusBadRequest)
		return
	}
	if apiPartner.Id != "" {
		http.Error(response, "Invalid request (Id specified)", http.StatusBadRequest)
		return
	}

	if !apiPartner.Validate() {
		jsonBytes, _ := json.Marshal(apiPartner)
		response.WriteHeader(http.StatusBadRequest)
		response.Write(jsonBytes)
		return
	}

	partner.SetDefinition(apiPartner)
	controller.referential.Partners().Save(partner)
	partner.Start()
	jsonBytes, _ := partner.MarshalJSON()
	response.Write(jsonBytes)
}

func (controller *PartnerController) Save(response http.ResponseWriter) {
	logger.Log.Debugf("Saving partners to database")

	status, err := controller.referential.Partners().SaveToDatabase()

	if err != nil {
		monitoring.ReportError(err)

		response.WriteHeader(status)
		jsonBytes, _ := json.Marshal(map[string]string{"error": err.Error()})
		response.Write(jsonBytes)
	}
}
