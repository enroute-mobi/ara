package core

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri"
)

func getMonitoredStopVisit(t *testing.T) *siri.XMLMonitoredStopVisit {
	file, err := os.Open("testdata/stopmonitoring-response-soap.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	response, _ := siri.NewXMLStopMonitoringResponseFromContent(content)
	return response.StopMonitoringDeliveries()[0].XMLMonitoredStopVisits()[0]
}

func Test_SIRIStopVisitUpdateAttributes_StopVisitAttributes(t *testing.T) {
	xmlStopVisit := getMonitoredStopVisit(t)
	stopVisitUpdateAttributes := NewSIRIStopVisitUpdateAttributes(xmlStopVisit, "objectidKind")
	stopVisitAttributes := stopVisitUpdateAttributes.StopVisitAttributes()

	expected := map[string]string{"kind": "objectidKind", "value": "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3"}
	if stopVisitAttributes.ObjectId.Kind() != expected["kind"] || stopVisitAttributes.ObjectId.Value() != expected["value"] {
		t.Errorf("Wrong ObjectId:\n expected: kind: %v value: %v\n got: kind: %v value: %v", expected["kind"], expected["value"], stopVisitAttributes.ObjectId.Kind(), stopVisitAttributes.ObjectId.Value())
	}
	expected["value"] = "NINOXE:StopPoint:Q:50:LOC"
	if stopVisitAttributes.StopAreaObjectId.Kind() != expected["kind"] || stopVisitAttributes.StopAreaObjectId.Value() != expected["value"] {
		t.Errorf("Wrong StopAreaObjectId:\n expected: kind: %v value: %v\n got: kind: %v value: %v", expected["kind"], expected["value"], stopVisitAttributes.StopAreaObjectId.Kind(), stopVisitAttributes.StopAreaObjectId.Value())
	}
	expected["value"] = "NINOXE:VehicleJourney:201"
	if stopVisitAttributes.VehicleJourneyObjectId.Kind() != expected["kind"] || stopVisitAttributes.VehicleJourneyObjectId.Value() != expected["value"] {
		t.Errorf("Wrong VehicleJourneyObjectId:\n expected: kind: %v value: %v\n got: kind: %v value: %v", expected["kind"], expected["value"], stopVisitAttributes.VehicleJourneyObjectId.Kind(), stopVisitAttributes.VehicleJourneyObjectId.Value())
	}
	if expected := 4; stopVisitAttributes.PassageOrder != expected {
		t.Errorf("Wrong PassageOrder:\n expected: %v\n got: %v", expected, stopVisitAttributes.PassageOrder)
	}
	if expected := model.STOP_VISIT_DEPARTURE_UNDEFINED; stopVisitAttributes.DepartureStatus != expected {
		t.Errorf("Wrong DepartureStatus:\n expected: %v\n got: %v", expected, stopVisitAttributes.DepartureStatus)
	}
	if expected := model.STOP_VISIT_ARRIVAL_ARRIVED; stopVisitAttributes.ArrivalStatus != expected {
		t.Errorf("Wrong ArrivalStatus:\n expected: %v\n got: %v", expected, stopVisitAttributes.ArrivalStatus)
	}
	if expected := time.Date(2016, time.September, 22, 9, 54, 00, 000000000, time.UTC); !stopVisitAttributes.Schedules.Schedule(model.STOP_VISIT_SCHEDULE_AIMED).DepartureTime().IsZero() || stopVisitAttributes.Schedules.Schedule(model.STOP_VISIT_SCHEDULE_AIMED).ArrivalTime().Equal(expected) {
		t.Errorf("Wrong Aimed Schedule:\n expected: departure: %v arrival: %v\n got: departure: %v arrival: %v", time.Time{}, expected, stopVisitAttributes.Schedules.Schedule(model.STOP_VISIT_SCHEDULE_AIMED).DepartureTime(), stopVisitAttributes.Schedules.Schedule(model.STOP_VISIT_SCHEDULE_AIMED).ArrivalTime())
	}
	if expected := time.Date(2016, time.September, 22, 9, 54, 00, 000000000, time.UTC); !stopVisitAttributes.Schedules.Schedule(model.STOP_VISIT_SCHEDULE_ACTUAL).DepartureTime().IsZero() || stopVisitAttributes.Schedules.Schedule(model.STOP_VISIT_SCHEDULE_ACTUAL).ArrivalTime().Equal(expected) {
		t.Errorf("Wrong Actual Schedule:\n expected: departure: %v arrival: %v\n got: departure: %v arrival: %v", time.Time{}, expected, stopVisitAttributes.Schedules.Schedule(model.STOP_VISIT_SCHEDULE_ACTUAL).DepartureTime(), stopVisitAttributes.Schedules.Schedule(model.STOP_VISIT_SCHEDULE_ACTUAL).ArrivalTime())
	}
	if schedule := stopVisitAttributes.Schedules.Schedule(model.STOP_VISIT_SCHEDULE_EXPECTED); !schedule.DepartureTime().IsZero() || !schedule.ArrivalTime().IsZero() {
		t.Errorf("Expected Schedule shouldn't be created, got: departure: %v arrival: %v", stopVisitAttributes.Schedules.Schedule(model.STOP_VISIT_SCHEDULE_EXPECTED).DepartureTime(), stopVisitAttributes.Schedules.Schedule(model.STOP_VISIT_SCHEDULE_EXPECTED).ArrivalTime())
	}
}

func Test_SIRIStopVisitUpdateAttributes_VehicleJourneyAttributes(t *testing.T) {
	xmlStopVisit := getMonitoredStopVisit(t)
	stopVisitUpdateAttributes := NewSIRIStopVisitUpdateAttributes(xmlStopVisit, "objectidKind")
	vehicleJourneyAttributes := stopVisitUpdateAttributes.VehicleJourneyAttributes()

	expected := map[string]string{"kind": "objectidKind", "value": "NINOXE:VehicleJourney:201"}
	if vehicleJourneyAttributes.ObjectId.Kind() != expected["kind"] || vehicleJourneyAttributes.ObjectId.Value() != expected["value"] {
		t.Errorf("Wrong ObjectId:\n expected: kind: %v value: %v\n got: kind: %v value: %v", expected["kind"], expected["value"], vehicleJourneyAttributes.ObjectId.Kind(), vehicleJourneyAttributes.ObjectId.Value())
	}
	expected["value"] = "NINOXE:Line:3:LOC"
	if vehicleJourneyAttributes.LineObjectId.Kind() != expected["kind"] || vehicleJourneyAttributes.LineObjectId.Value() != expected["value"] {
		t.Errorf("Wrong LineObjectId:\n expected: kind: %v value: %v\n got: kind: %v value: %v", expected["kind"], expected["value"], vehicleJourneyAttributes.LineObjectId.Kind(), vehicleJourneyAttributes.LineObjectId.Value())
	}
}

func Test_SIRIStopVisitUpdateAttributes_LineAttributes(t *testing.T) {
	xmlStopVisit := getMonitoredStopVisit(t)
	stopVisitUpdateAttributes := NewSIRIStopVisitUpdateAttributes(xmlStopVisit, "objectidKind")
	lineAttributes := stopVisitUpdateAttributes.LineAttributes()

	expected := map[string]string{"kind": "objectidKind", "value": "NINOXE:Line:3:LOC"}
	if lineAttributes.ObjectId.Kind() != expected["kind"] || lineAttributes.ObjectId.Value() != expected["value"] {
		t.Errorf("Wrong ObjectId:\n expected: kind: %v value: %v\n got: kind: %v value: %v", expected["kind"], expected["value"], lineAttributes.ObjectId.Kind(), lineAttributes.ObjectId.Value())
	}
	if expected := "Ligne 3 Metro"; lineAttributes.Name != expected {
		t.Errorf("Wrong Name:\n expected: %v\n got: %v", expected, lineAttributes.Name)
	}
}
