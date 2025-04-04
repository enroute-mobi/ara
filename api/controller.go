package api

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	"bitbucket.org/enroute-mobi/ara/core"
	"bitbucket.org/enroute-mobi/ara/model"
)

var idPattern = regexp.MustCompile("([0-9a-zA-Z-]+):([0-9a-zA-Z-:]+)")

var newWithReferentialControllerMap = map[string](func(*core.Referential) RestfulResource){
	"stop_areas":            NewStopAreaController,
	"partners":              NewPartnerController,
	"lines":                 NewLineController,
	"line_groups":           NewLineGroupsController,
	"stop_area_groups":      NewStopAreaGroupsController,
	"stop_visits":           NewStopVisitController,
	"scheduled_stop_visits": NewScheduledStopVisitController,
	"vehicle_journeys":      NewVehicleJourneyController,
	"situations":            NewSituationController,
	"operators":             NewOperatorController,
	"vehicles":              NewVehicleController,
}

const (
	DEFAULT_PER_PAGE = 30
)

type RestfulResource interface {
	Index(response http.ResponseWriter, params url.Values)
	Show(response http.ResponseWriter, identifier string)
	Delete(response http.ResponseWriter, identifier string)
	Update(response http.ResponseWriter, identifier string, body []byte)
	Create(response http.ResponseWriter, body []byte)
}

type SubscriptionResource interface {
	SubscriptionsIndex(response http.ResponseWriter, identifier string)
	SubscriptionsCreate(response http.ResponseWriter, identifier string, body []byte)
	SubscriptionsDelete(response http.ResponseWriter, identifier string, subscriptionId string)
}

type Savable interface {
	Save(response http.ResponseWriter)
}

type ControllerInterface interface {
	serve(response http.ResponseWriter, request *http.Request)
}

func getRequestBody(response http.ResponseWriter, request *http.Request) []byte {
	if request.Body == nil {
		http.Error(response, "Invalid request: Can't read request body", http.StatusBadRequest)
		return nil
	}
	body, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(response, "Invalid request: Can't read request body", http.StatusBadRequest)
		return nil
	}
	if len(body) == 0 {
		http.Error(response, "Invalid request: Empty body", http.StatusBadRequest)
		return nil
	}
	return body
}

type Paginable interface {
	model.Situation | model.StopArea
}

func paginate[P Paginable](p []*P, params url.Values) ([]*P, error) {
	if len(params) == 0 {
		return p, nil
	}

	page, err := strconv.Atoi(params.Get("page"))
	if err != nil {
		return nil, fmt.Errorf("invalid request: query parameter \"page\": %s", params.Get("page"))
	}

	var per_page int
	if params.Get("per_page") != "" {
		per_page, err = strconv.Atoi(params.Get("per_page"))
		if page != 0 && err != nil {
			return nil, fmt.Errorf("invalid request: query parameter \"per_page\": %s", params.Get("per_page"))
		}
	}

	if page == 0 && per_page == 0 {
		return p, nil
	}

	if per_page == 0 || per_page > DEFAULT_PER_PAGE {
		per_page = DEFAULT_PER_PAGE
	}

	start, end := paginateSlice(page, per_page, len(p))
	pagedSlice := p[start:end]

	return pagedSlice, nil
}

func paginateSlice(pageNum int, pageSize int, sliceLength int) (int, int) {
	firstEntry := (pageNum - 1) * pageSize
	lastEntry := firstEntry + pageSize

	if lastEntry > sliceLength {
		lastEntry = sliceLength
	}

	return firstEntry, lastEntry
}
