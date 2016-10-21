package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
)

type StopAreaController struct {
	memoryStopAreas *model.MemoryStopAreas
}

func NewStopAreaController() (controller *StopAreaController) {
	return &StopAreaController{memoryStopAreas: model.NewMemoryStopAreas()}
}

func (controller *StopAreaController) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	logger.Log.Debugf("StopAreaController request: %s", request)

	path := request.URL.Path
	resourcePathPattern := regexp.MustCompile("/stop_areas(?:/([0-9a-zA-Z-]+))?")
	identifier := model.StopAreaId(resourcePathPattern.FindStringSubmatch(path)[1])

	response.Header().Set("Content-Type", "application/json")

	switch {
	case request.Method == "GET":
		if identifier == "" {
			controller.Index(response)
		} else {
			controller.Show(response, identifier)
		}
	case request.Method == "DELETE":
		if identifier == "" {
			http.Error(response, "Invalid request", 400)
			return
		}
		controller.Delete(response, identifier)
	case request.Method == "PUT":
		if identifier == "" {
			http.Error(response, "Invalid request", 400)
			return
		}
		body := getRequestBody(response, request)
		if body == nil {
			http.Error(response, "Invalid request", 400)
			return
		}
		controller.Update(response, identifier, body)
	case request.Method == "POST":
		if identifier != "" {
			http.Error(response, "Invalid request", 400)
			return
		}
		body := getRequestBody(response, request)
		if body == nil {
			http.Error(response, "Invalid request", 400)
			return
		}
		controller.Create(response, body)
	}
}

func getRequestBody(response http.ResponseWriter, request *http.Request) []byte {
	if request.Body == nil {
		http.Error(response, "Invalid request: Empty body", 400)
		return nil
	}
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(response, "Invalid request: Can't read request body", 400)
		return nil
	}
	return body
}

func (controller *StopAreaController) Index(response http.ResponseWriter) {
	logger.Log.Debugf("StopAreas Index")

	jsonBytes, _ := json.Marshal(controller.memoryStopAreas.FindAll())
	response.Write(jsonBytes)
}

func (controller *StopAreaController) Show(response http.ResponseWriter, identifier model.StopAreaId) {
	stopArea, ok := controller.memoryStopAreas.Find(identifier)
	if !ok {
		http.Error(response, fmt.Sprintf("Stop area not found: %s", identifier), 500)
		return
	}
	logger.Log.Debugf("Get stopArea %s", identifier)

	jsonBytes, _ := stopArea.MarshalJSON()
	response.Write(jsonBytes)
}

func (controller *StopAreaController) Delete(response http.ResponseWriter, identifier model.StopAreaId) {
	stopArea, ok := controller.memoryStopAreas.Find(identifier)
	if !ok {
		http.Error(response, fmt.Sprintf("Stop area not found: %s", identifier), 500)
		return
	}
	logger.Log.Debugf("Delete stopArea %s", identifier)

	jsonBytes, _ := stopArea.MarshalJSON()
	controller.memoryStopAreas.Delete(&stopArea)
	response.Write(jsonBytes)
}

func (controller *StopAreaController) Update(response http.ResponseWriter, identifier model.StopAreaId, body []byte) {
	stopArea, ok := controller.memoryStopAreas.Find(identifier)
	if !ok {
		http.Error(response, fmt.Sprintf("Stop area not found: %s", identifier), 500)
		return
	}

	logger.Log.Debugf("Update stopArea %s: %s", identifier, string(body))

	err := json.Unmarshal(body, &stopArea)
	if err != nil {
		http.Error(response, "Invalid request: can't parse request body", 400)
		return
	}

	controller.memoryStopAreas.Save(&stopArea)
	jsonBytes, _ := stopArea.MarshalJSON()
	response.Write(jsonBytes)
}

func (controller *StopAreaController) Create(response http.ResponseWriter, body []byte) {
	logger.Log.Debugf("Create stopArea: %s", string(body))

	stopArea := controller.memoryStopAreas.New()

	err := json.Unmarshal(body, &stopArea)
	if err != nil {
		http.Error(response, "Invalid request: can't parse request body", 400)
		return
	}
	if stopArea.Id() != "" {
		http.Error(response, "Invalid request", 400)
		return
	}

	controller.memoryStopAreas.Save(&stopArea)
	jsonBytes, _ := stopArea.MarshalJSON()
	response.Write(jsonBytes)
}
