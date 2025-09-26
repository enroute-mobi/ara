package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"

	"bitbucket.org/enroute-mobi/ara/core"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
)

type VehicleJourneyController struct {
	referential *core.Referential
}

func NewVehicleJourneyController(referential *core.Referential) RestfulResource {
	return &VehicleJourneyController{
		referential: referential,
	}
}

func (controller *VehicleJourneyController) findVehicleJourney(identifier string) (*model.VehicleJourney, bool) {
	foundStrings := idPattern.FindStringSubmatch(identifier)
	if foundStrings != nil {
		code := model.NewCode(foundStrings[1], foundStrings[2])
		return controller.referential.Model().VehicleJourneys().FindByCode(code)
	}
	return controller.referential.Model().VehicleJourneys().Find(model.VehicleJourneyId(identifier))
}

func (controller *VehicleJourneyController) Index(response http.ResponseWriter, params url.Values) {
	logger.Log.Debugf("VehicleJourneys Index")

	allVehicleJourneys := controller.referential.Model().VehicleJourneys().FindAll()
	direction := params.Get("direction")
	switch direction {
	case "desc":
		sort.Slice(allVehicleJourneys, func(i, j int) bool {
			return allVehicleJourneys[i].Name > allVehicleJourneys[j].Name
		})
	case "asc", "":
		sort.Slice(allVehicleJourneys, func(i, j int) bool {
			return allVehicleJourneys[i].Name < allVehicleJourneys[j].Name
		})
	default:
		http.Error(response, fmt.Sprintf("invalid request: query parameter \"direction\": %s", params.Get("direction")), http.StatusBadRequest)
		return
	}

	paginatedVehicleJourneys, err := paginate(allVehicleJourneys, params)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	jsonBytes, _ := json.Marshal(paginatedVehicleJourneys)
	response.Write(jsonBytes)
}

func (controller *VehicleJourneyController) Show(response http.ResponseWriter, identifier string, params url.Values) {
	vehicleJourney, ok := controller.findVehicleJourney(identifier)
	if !ok {
		http.Error(response, fmt.Sprintf("Vehicle journey not found: %s", identifier), http.StatusNotFound)
		return
	}
	logger.Log.Debugf("Get vehicleJourney %s", identifier)

	withDetailedStpVisits := params.Get("with_detailed_stop_visits")

	if withDetailedStpVisits != "" {
		ok, err := strconv.ParseBool(withDetailedStpVisits)
		if err != nil {
			http.Error(response, fmt.Sprintf("invalid request: query parameter \"with_detailed_stop_visits\": %s", params.Get("with_detailed_stop_visits")), http.StatusBadRequest)
		}

		if ok {

			var stopVisitsWithDetails []model.DetailedStopVisit
			svs := controller.referential.Model().StopVisits().FindByVehicleJourneyId(vehicleJourney.Id())
			for i := range svs {
				sa, ok := controller.referential.Model().StopAreas().Find(svs[i].StopAreaId)
				if !ok {
					continue
				}

				stopVisit := &model.DetailedStopVisit{}
				stopVisit.Order = svs[i].PassageOrder
				stopVisit.StopAreaId = svs[i].StopAreaId
				stopVisit.StopAreaName = sa.Name
				stopVisit.ArrivalStatus = svs[i].ArrivalStatus
				stopVisit.DepartureStatus = svs[i].DepartureStatus
				stopVisit.Schedules = svs[i].Schedules.ToSlice()
				stopVisit.CollectedAt = svs[i].CollectedAt()
				stopVisitsWithDetails = append(stopVisitsWithDetails, *stopVisit)

			}

			vehicleJourney.DetailedStopVisits = stopVisitsWithDetails
		}
	}

	jsonBytes, _ := vehicleJourney.MarshalJSON()
	response.Write(jsonBytes)
}

func (controller *VehicleJourneyController) Delete(response http.ResponseWriter, identifier string) {
	vehicleJourney, ok := controller.findVehicleJourney(identifier)
	if !ok {
		http.Error(response, fmt.Sprintf("Vehicle journey not found: %s", identifier), http.StatusNotFound)
		return
	}
	logger.Log.Debugf("Delete vehicleJourney %s", identifier)

	jsonBytes, _ := vehicleJourney.MarshalJSON()
	controller.referential.Model().VehicleJourneys().Delete(vehicleJourney)
	response.Write(jsonBytes)
}

func (controller *VehicleJourneyController) Update(response http.ResponseWriter, identifier string, body []byte) {
	vehicleJourney, ok := controller.findVehicleJourney(identifier)
	if !ok {
		http.Error(response, fmt.Sprintf("Vehicle journey not found: %s", identifier), http.StatusNotFound)
		return
	}

	logger.Log.Debugf("Update vehicleJourney %s: %s", identifier, string(body))

	err := json.Unmarshal(body, &vehicleJourney)
	if err != nil {
		http.Error(response, fmt.Sprintf("Invalid request: can't parse request body: %v", err), http.StatusBadRequest)
		return
	}

	for _, obj := range vehicleJourney.Codes() {
		vj, ok := controller.referential.Model().VehicleJourneys().FindByCode(obj)
		if ok && vj.Id() != vehicleJourney.Id() {
			http.Error(response, fmt.Sprintf("Invalid request: vehicleJourney %v already have a code %v", vj.Id(), obj.String()), http.StatusBadRequest)
			return
		}
	}

	controller.referential.Model().VehicleJourneys().Save(vehicleJourney)
	jsonBytes, _ := vehicleJourney.MarshalJSON()
	response.Write(jsonBytes)
}

func (controller *VehicleJourneyController) Create(response http.ResponseWriter, body []byte) {
	logger.Log.Debugf("Create vehicleJourney: %s", string(body))

	vehicleJourney := controller.referential.Model().VehicleJourneys().New()

	err := json.Unmarshal(body, &vehicleJourney)
	if err != nil {
		http.Error(response, fmt.Sprintf("Invalid request: can't parse request body: %v", err), http.StatusBadRequest)
		return
	}

	if vehicleJourney.Id() != "" {
		http.Error(response, "Invalid request", http.StatusBadRequest)
		return
	}

	for _, obj := range vehicleJourney.Codes() {
		vj, ok := controller.referential.Model().VehicleJourneys().FindByCode(obj)
		if ok {
			http.Error(response, fmt.Sprintf("Invalid request: vehicleJourney %v already have a code %v", vj.Id(), obj.String()), http.StatusBadRequest)
			return
		}
	}

	controller.referential.Model().VehicleJourneys().Save(vehicleJourney)
	jsonBytes, _ := vehicleJourney.MarshalJSON()
	response.Write(jsonBytes)
}
