package siri

import (
	"testing"
	"time"
)

var r string

func BenchmarkTemplates(b *testing.B) {
	response := &SIRIStopMonitoringResponse{
		Address:                   "address",
		ProducerRef:               "producer",
		ResponseMessageIdentifier: "identifier",
	}
	response.RequestMessageRef = "ref"
	response.Status = true
	response.ResponseTimestamp = time.Date(2016, time.September, 21, 20, 14, 46, 0, time.UTC)
	response.MonitoringRef = "MonitoringRef"

	siriMonitoredStopVisit := &SIRIMonitoredStopVisit{
		ItemIdentifier:        "itemId",
		MonitoringRef:         "monitoringRef",
		StopPointRef:          "stopPointRef",
		StopPointName:         "stopPointName",
		LineRef:               "lineRef",
		PublishedLineName:     "lineName",
		DepartureStatus:       "depStatus",
		ArrivalStatus:         "arrStatus",
		VehicleJourneyName:    "NameOfVj",
		VehicleAtStop:         true,
		Order:                 1,
		Monitored:             true,
		RecordedAt:            time.Date(2015, time.September, 21, 20, 14, 46, 0, time.UTC),
		DataFrameRef:          "2016-09-21",
		AimedArrivalTime:      time.Date(2017, time.September, 21, 20, 14, 46, 0, time.UTC),
		ActualArrivalTime:     time.Date(2018, time.September, 21, 20, 14, 46, 0, time.UTC),
		AimedDepartureTime:    time.Date(2019, time.September, 21, 20, 14, 46, 0, time.UTC),
		ExpectedDepartureTime: time.Date(2020, time.September, 21, 20, 14, 46, 0, time.UTC),
		Attributes:            make(map[string]map[string]string),
		References:            make(map[string]map[string]string),
	}

	siriMonitoredStopVisit.Attributes["StopVisitAttributes"] = make(map[string]string)
	siriMonitoredStopVisit.References["VehicleJourney"] = make(map[string]string)
	siriMonitoredStopVisit.References["StopVisitReferences"] = make(map[string]string)
	siriMonitoredStopVisit.Attributes["VehicleJourneyAttributes"] = make(map[string]string)
	siriMonitoredStopVisit.Attributes["VehicleJourneyAttributes"]["Delay"] = "30"
	siriMonitoredStopVisit.DatedVehicleJourneyRef = "vehicleJourney#Code"
	siriMonitoredStopVisit.References["StopVisitReferences"]["OperatorRef"] = "OperatorRef"
	siriMonitoredStopVisit.References["VehicleJourney"]["DestinationRef"] = "NINOXE:StopPoint:SP:62:LOC"

	response.MonitoredStopVisits = []*SIRIMonitoredStopVisit{siriMonitoredStopVisit}

	for n := 0; n < b.N; n++ {
		xml, err := response.BuildXML()
		if err != nil {
			b.Fatal(err)
		}
		r = xml
	}
}
