package model

import "bitbucket.org/enroute-mobi/ara/logger"

func NewStopVisitSetVehicleJourneyCancellationUpdater() (updater, error) {
	return func(mi ModelInstance) error {
		sv := mi.(*StopVisit)
		if !sv.IsCancelled() {
			return nil
		}

		vj, ok := sv.model.VehicleJourneys().Find(sv.VehicleJourneyId)
		if !ok {
			return nil
		}

		if vj.IsCancelled() {
			return nil
		}

		if vj.IsComplete() {
			svs := vj.model.StopVisits().FindByVehicleJourneyId(vj.Id())
			for i := range svs {
				if svs[i].Id() == sv.Id() {
					continue
				}
				if !svs[i].IsCancelled() {
					return nil
				}
			}
			vj.Cancellation = true
			logger.Log.Printf("VehicleJourney %v Cancelled coming from StopVisit %v, StopArea %v", vj.Id(), sv.Id(), sv.StopAreaId)
			sv.model.VehicleJourneys().Save(vj)
		}
		return nil
	}, nil
}
