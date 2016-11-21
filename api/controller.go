package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/af83/edwig/core"
	"github.com/af83/edwig/logger"
)

type RessourceController interface {
	SetReferential(referential *core.Referential)
	Ressources() string
	Index(response http.ResponseWriter)
	Show(response http.ResponseWriter, identifier string)
	Delete(response http.ResponseWriter, identifier string)
	Update(response http.ResponseWriter, identifier string, body []byte)
	Create(response http.ResponseWriter, body []byte)
}

type ControllerReferential struct {
	referential *core.Referential
}

func (controller *ControllerReferential) SetReferential(referential *core.Referential) {
	controller.referential = referential
	return
}

type Controller struct {
	ressourceController RessourceController
}

func (controller *Controller) SetReferential(referential *core.Referential) {
	controller.ressourceController.SetReferential(referential)
	return
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

func (controller *Controller) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	logger.Log.Debugf("%s controller request: %s", controller.ressourceController.Ressources(), request)

	path := request.URL.Path
	ressourcePath := fmt.Sprintf("/%s(?:/([0-9a-zA-Z-]+))?", controller.ressourceController.Ressources())
	resourcePathPattern := regexp.MustCompile(ressourcePath)
	identifier := resourcePathPattern.FindStringSubmatch(path)[1]

	response.Header().Set("Content-Type", "application/json")

	switch {
	case request.Method == "GET":
		if identifier == "" {
			controller.ressourceController.Index(response)
		} else {
			controller.ressourceController.Show(response, identifier)
		}
	case request.Method == "DELETE":
		if identifier == "" {
			http.Error(response, "Invalid request", 400)
			return
		}
		controller.ressourceController.Delete(response, identifier)
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
		controller.ressourceController.Update(response, identifier, body)
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
		controller.ressourceController.Create(response, body)
	}
}
