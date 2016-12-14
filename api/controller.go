package api

import (
	"io/ioutil"
	"net/http"

	"github.com/af83/edwig/core"
)

var newControllerMap = map[string](func(*Server) *Controller){
	"_referentials": NewReferentialController,
}

var newWithReferentialControllerMap = map[string](func(*core.Referential) *Controller){
	"stop_areas": NewStopAreaController,
	"partners":   NewPartnerController,
	"lines":      NewLineController,
}

type RestfulRessource interface {
	Index(response http.ResponseWriter)
	Show(response http.ResponseWriter, identifier string)
	Delete(response http.ResponseWriter, identifier string)
	Update(response http.ResponseWriter, identifier string, body []byte)
	Create(response http.ResponseWriter, body []byte)
}

type Controller struct {
	restfulRessource RestfulRessource
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

func (controller *Controller) serve(response http.ResponseWriter, request *http.Request, identifier string) {
	switch {
	case request.Method == "GET":
		if identifier == "" {
			controller.restfulRessource.Index(response)
		} else {
			controller.restfulRessource.Show(response, identifier)
		}
	case request.Method == "DELETE":
		if identifier == "" {
			http.Error(response, "Invalid request", 400)
			return
		}
		controller.restfulRessource.Delete(response, identifier)
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
		controller.restfulRessource.Update(response, identifier, body)
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
		controller.restfulRessource.Create(response, body)
	}
}
