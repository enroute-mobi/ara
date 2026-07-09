package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/model/schedules"
)

// PassingTimeChronology controller is created on StopVisits, but called on VehicleJourneys
func NewPassingTimeChronologyController(sc *SelectControl) (controller, error) {
	if sc.Hook.String != "AfterAllStopVisitSave" {
		return nil, errors.New("'unexpected' controller must be defined AfterAllStopVisitSave")
	}
	if sc.ModelType.String != "StopVisit" {
		return nil, fmt.Errorf("don't know how to handle model type %s in 'unexpected' controller", sc.ModelType.String)
	}

	var remoteCodeSpace string
	attrs := make(map[string]string)

	err := json.Unmarshal([]byte(sc.Attributes.String), &attrs)
	if err == nil {
		remoteCodeSpace = attrs["code_space"]
	}

	return func(mi ModelInstance) error {
		vj, ok := mi.(*VehicleJourney)
		if !ok {
			return errors.New("PassingTimeChronologyController called on non VehicleJourney Model")
		}
		vjName := vj.Name
		if vjName == "" && remoteCodeSpace != "" {
			code, _ := vj.Code(remoteCodeSpace)
			vjName = code.Value()
		}

		svs := vj.model.StopVisits().FindByVehicleJourneyId(vj.id)

		sort.Slice(svs, func(i int, j int) bool {
			return svs[i].PassageOrder < svs[j].PassageOrder
		})

		for i := range svs {
			// Control if arrival time if after departure time for each schedule type
			for _, j := range schedules.ScheduleOrderArray {
				if !svs[i].Schedules.Schedule(j).ArrivalTime().IsZero() &&
					!svs[i].Schedules.Schedule(j).DepartureTime().IsZero() &&
					svs[i].Schedules.Schedule(j).ArrivalTime().After(svs[i].Schedules.Schedule(j).DepartureTime()) {

					// Attributes: vj name, departure time, arrival time
					messageAttribute := fmt.Sprintf("%s,%s,%s",
						vjName,
						svs[i].Schedules.Schedule(j).DepartureTime().Format("2006-01-02T15:04:05.000Z07:00"),
						svs[i].Schedules.Schedule(j).ArrivalTime().Format("2006-01-02T15:04:05.000Z07:00"),
					)
					m := &audit.BigQueryControlEvent{
						Criticity:                        sc.Criticity.String,
						ControlType:                      "PassingTimeChronology",
						InternalCode:                     sc.InternalCode.String,
						TargetModelClass:                 sc.ModelType.String,
						TargetModelUUID:                  string(svs[i].ModelId()),
						TranslationInfoMessageKey:        fmt.Sprintf("%s_arrival_before_departure", j),
						TranslationInfoMessageAttributes: messageAttribute,
					}

					audit.CurrentBigQuery(sc.ReferentialSlug).WriteEvent(m)
				}

				// Control if next Stop Visit arrival time is before the saved Stop Visit departure time
				if i != len(svs)-1 &&
					!svs[i].Schedules.Schedule(j).DepartureTime().IsZero() &&
					!svs[i+1].Schedules.Schedule(j).ArrivalTime().IsZero() &&
					svs[i].Schedules.Schedule(j).DepartureTime().After(svs[i+1].Schedules.Schedule(j).ArrivalTime()) {

					// Attributes: vj name, current departure time, next arrival time
					messageAttribute := fmt.Sprintf("%s,%s,%s",
						vjName,
						svs[i].Schedules.Schedule(j).DepartureTime().Format("2006-01-02T15:04:05.000Z07:00"),
						svs[i+1].Schedules.Schedule(j).ArrivalTime().Format("2006-01-02T15:04:05.000Z07:00"),
					)
					m := &audit.BigQueryControlEvent{
						Criticity:                        sc.Criticity.String,
						ControlType:                      "PassingTimeChronology",
						InternalCode:                     sc.InternalCode.String,
						TargetModelClass:                 sc.ModelType.String,
						TargetModelUUID:                  string(svs[i].ModelId()),
						TranslationInfoMessageKey:        fmt.Sprintf("%s_departure_after_next_arrival", j),
						TranslationInfoMessageAttributes: messageAttribute,
					}

					audit.CurrentBigQuery(sc.ReferentialSlug).WriteEvent(m)
				}
			}
		}

		return nil
	}, nil
}
