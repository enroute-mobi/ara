package model

import (
	"io"
	"os"
	"reflect"
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/model/schedules"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
	"github.com/stretchr/testify/assert"
)

func Test_UpdateManager_UpdateVehicle_WithNextStopVisitOrderExisting(t *testing.T) {
	assert := assert.New(t)

	model := NewTestMemoryModel()
	code := NewCode("codeSpace", "value")

	sa := model.StopAreas().New()
	sa.SetCode(code)
	sa.Save()

	l := model.Lines().New()
	l.SetCode(code)
	l.Save()

	vj := model.VehicleJourneys().New()
	vj.SetCode(code)
	vj.LineId = l.Id()
	vj.Save()

	stopVisit := model.StopVisits().New()
	stopVisit.SetCode(code)
	stopVisit.VehicleJourneyId = vj.Id()
	stopVisit.StopAreaId = sa.Id()
	stopVisit.PassageOrder = 5
	stopVisit.Save()

	vehicle := model.Vehicles().New()
	vehicle.SetCode(code)
	vehicle.LineId = l.Id()
	vehicle.StopAreaId = sa.Id()
	vehicle.VehicleJourneyId = vj.Id()
	vehicle.Save()

	manager := newUpdateManager(model)

	event := &VehicleUpdateEvent{
		Code:               code,
		StopAreaCode:       code,
		VehicleJourneyCode: code,
		NextStopPointOrder: 5,
	}

	manager.Update(event)

	updatedVehicle, _ := model.vehicles.Find(vehicle.Id())

	assert.Equal(stopVisit.Id(), updatedVehicle.NextStopVisitId)
}

func Test_UpdateManager_UpdateVehicle_WithNextStopVisitOrderNotExisting(t *testing.T) {
	assert := assert.New(t)

	model := NewTestMemoryModel()
	code := NewCode("codeSpace", "value")

	sa := model.StopAreas().New()
	sa.SetCode(code)
	sa.Save()

	l := model.Lines().New()
	l.SetCode(code)
	l.Save()

	vj := model.VehicleJourneys().New()
	vj.SetCode(code)
	vj.LineId = l.Id()
	vj.Save()

	stopVisit := model.StopVisits().New()
	stopVisit.SetCode(code)
	stopVisit.VehicleJourneyId = vj.Id()
	stopVisit.StopAreaId = sa.Id()
	stopVisit.PassageOrder = 6
	stopVisit.Save()

	vehicle := model.Vehicles().New()
	vehicle.SetCode(code)
	vehicle.LineId = l.Id()
	vehicle.StopAreaId = sa.Id()
	vehicle.VehicleJourneyId = vj.Id()
	vehicle.Save()

	manager := newUpdateManager(model)

	event := &VehicleUpdateEvent{
		Code:               code,
		StopAreaCode:       code,
		VehicleJourneyCode: code,
		NextStopPointOrder: 5,
	}

	manager.Update(event)

	updatedVehicle, _ := model.vehicles.Find(vehicle.Id())

	assert.Equal(StopVisitId(""), updatedVehicle.NextStopVisitId)
}

func Test_UpdateManager_UpdateVehicle_WithNextStop_WithoutORder_With_One_StopVisit(t *testing.T) {
	assert := assert.New(t)

	model := NewTestMemoryModel()
	code := NewCode("codeSpace", "value")

	sa := model.StopAreas().New()
	sa.SetCode(code)
	sa.Save()

	l := model.Lines().New()
	l.SetCode(code)
	l.Save()

	vj := model.VehicleJourneys().New()
	vj.SetCode(code)
	vj.LineId = l.Id()
	vj.Save()

	stopVisit := model.StopVisits().New()
	stopVisit.SetCode(code)
	stopVisit.VehicleJourneyId = vj.Id()
	stopVisit.StopAreaId = sa.Id()
	stopVisit.PassageOrder = 6
	stopVisit.Save()

	vehicle := model.Vehicles().New()
	vehicle.SetCode(code)
	vehicle.LineId = l.Id()
	vehicle.StopAreaId = sa.Id()
	vehicle.VehicleJourneyId = vj.Id()
	vehicle.Save()

	manager := newUpdateManager(model)

	event := &VehicleUpdateEvent{
		Code:               code,
		StopAreaCode:       code,
		VehicleJourneyCode: code,
	}

	manager.Update(event)

	updatedVehicle, _ := model.vehicles.Find(vehicle.Id())

	assert.Equal(stopVisit.Id(), updatedVehicle.NextStopVisitId)
}

func Test_UpdateManager_UpdateVehicle_WithNextStop_WithoutOrder_With_More_Than_One_StopVisit(t *testing.T) {
	assert := assert.New(t)

	model := NewTestMemoryModel()
	code := NewCode("codeSpace", "value")

	sa := model.StopAreas().New()
	sa.SetCode(code)
	sa.Save()

	l := model.Lines().New()
	l.SetCode(code)
	l.Save()

	vj := model.VehicleJourneys().New()
	vj.SetCode(code)
	vj.LineId = l.Id()
	vj.Save()

	stopVisit := model.StopVisits().New()
	stopVisit.SetCode(code)
	stopVisit.VehicleJourneyId = vj.Id()
	stopVisit.StopAreaId = sa.Id()
	stopVisit.PassageOrder = 6
	stopVisit.Save()

	stopVisit1 := model.StopVisits().New()
	stopVisit1.SetCode(code)
	stopVisit1.VehicleJourneyId = vj.Id()
	stopVisit1.StopAreaId = sa.Id()
	stopVisit1.PassageOrder = 7
	stopVisit1.Save()

	vehicle := model.Vehicles().New()
	vehicle.SetCode(code)
	vehicle.LineId = l.Id()
	vehicle.StopAreaId = sa.Id()
	vehicle.VehicleJourneyId = vj.Id()
	vehicle.Save()

	manager := newUpdateManager(model)

	event := &VehicleUpdateEvent{
		Code:               code,
		StopAreaCode:       code,
		VehicleJourneyCode: code,
	}

	manager.Update(event)

	updatedVehicle, _ := model.vehicles.Find(vehicle.Id())

	assert.Equal(StopVisitId(""), updatedVehicle.NextStopVisitId)
}

func Test_UpdateManager_CreateStopVisit(t *testing.T) {
	model := NewTestMemoryModel()
	code := NewCode("codeSpace", "value")
	sa := model.StopAreas().New()
	sa.SetCode(code)
	sa.Save()

	l := model.Lines().New()
	l.SetCode(code)
	l.Save()

	vj := model.VehicleJourneys().New()
	vj.SetCode(code)
	vj.LineId = l.Id()
	vj.Save()

	manager := newUpdateManager(model)

	event := &StopVisitUpdateEvent{
		Code:               code,
		StopAreaCode:       code,
		VehicleJourneyCode: code,
		DepartureStatus:    STOP_VISIT_DEPARTURE_CANCELLED,
		ArrivalStatus:      STOP_VISIT_ARRIVAL_ONTIME,
		Schedules:          schedules.NewStopVisitSchedules(),
	}

	manager.Update(event)
	updatedStopVisit, ok := model.StopVisits().FindByCode(code)
	if !ok {
		t.Fatalf("StopVisit should be created")
	}
	if updatedStopVisit.DepartureStatus != STOP_VISIT_DEPARTURE_CANCELLED {
		t.Errorf("StopVisit DepartureStatus should be updated")
	}
	if updatedStopVisit.ArrivalStatus != STOP_VISIT_ARRIVAL_ONTIME {
		t.Errorf("StopVisit ArrivalStatus should be updated")
	}
	if !updatedStopVisit.IsCollected() {
		t.Errorf("StopVisit ArrivalStatus should be collected")
	}
	updatedStopArea, _ := model.StopAreas().Find(sa.Id())
	if !updatedStopArea.LineIds.Contains(l.Id()) {
		t.Errorf("StopArea LineIds should be updated")
	}
}

func Test_UpdateManager_UpdateStopVisit(t *testing.T) {
	model := NewTestMemoryModel()
	code := NewCode("codeSpace", "value")
	sa := model.StopAreas().New()
	sa.SetCode(code)
	sa.Save()

	l := model.Lines().New()
	l.SetCode(code)
	l.Save()

	vj := model.VehicleJourneys().New()
	vj.SetCode(code)
	vj.LineId = l.Id()
	vj.Save()

	stopVisit := model.StopVisits().New()
	stopVisit.SetCode(code)
	stopVisit.Save()

	manager := newUpdateManager(model)

	event := &StopVisitUpdateEvent{
		Code:               code,
		StopAreaCode:       code,
		VehicleJourneyCode: code,
		DepartureStatus:    STOP_VISIT_DEPARTURE_CANCELLED,
		ArrivalStatus:      STOP_VISIT_ARRIVAL_ONTIME,
		Schedules:          schedules.NewStopVisitSchedules(),
	}

	manager.Update(event)
	updatedStopVisit, _ := model.StopVisits().Find(stopVisit.Id())
	if updatedStopVisit.DepartureStatus != STOP_VISIT_DEPARTURE_CANCELLED {
		t.Errorf("StopVisit DepartureStatus should be updated")
	}
	if updatedStopVisit.ArrivalStatus != STOP_VISIT_ARRIVAL_ONTIME {
		t.Errorf("StopVisit ArrivalStatus should be updated")
	}
	if !updatedStopVisit.IsCollected() {
		t.Errorf("StopVisit ArrivalStatus should be collected")
	}
	updatedStopArea, _ := model.StopAreas().Find(sa.Id())
	if !updatedStopArea.LineIds.Contains(l.Id()) {
		t.Errorf("StopArea LineIds should be updated")
	}
}

func Test_UpdateManager_CreateStopVisit_NoStopAreaId(t *testing.T) {
	emptyCode := NewCode("codeSpace", "")

	model := NewTestMemoryModel()
	code := NewCode("codeSpace", "value")
	sa := model.StopAreas().New()
	sa.SetCode(code)
	sa.Save()

	l := model.Lines().New()
	l.SetCode(code)
	l.Save()

	vj := model.VehicleJourneys().New()
	vj.SetCode(code)
	vj.LineId = l.Id()
	vj.Save()

	manager := newUpdateManager(model)

	event := &StopVisitUpdateEvent{
		Code:               code,
		StopAreaCode:       emptyCode,
		VehicleJourneyCode: code,
		DepartureStatus:    STOP_VISIT_DEPARTURE_CANCELLED,
		ArrivalStatus:      STOP_VISIT_ARRIVAL_ONTIME,
		Schedules:          schedules.NewStopVisitSchedules(),
	}

	manager.Update(event)
	_, ok := model.StopVisits().FindByCode(code)
	if ok {
		t.Fatalf("StopVisit should not be created")
	}
}

func Test_UpdateManager_UpdateStopVisit_NoStopAreaId(t *testing.T) {
	emptyCode := NewCode("codeSpace", "")

	model := NewTestMemoryModel()
	code := NewCode("codeSpace", "value")
	sa := model.StopAreas().New()
	sa.SetCode(code)
	sa.Save()

	l := model.Lines().New()
	l.SetCode(code)
	l.Save()

	vj := model.VehicleJourneys().New()
	vj.SetCode(code)
	vj.LineId = l.Id()
	vj.Save()

	stopVisit := model.StopVisits().New()
	stopVisit.SetCode(code)
	stopVisit.StopAreaId = sa.Id()
	stopVisit.Save()

	manager := newUpdateManager(model)

	event := &StopVisitUpdateEvent{
		Code:               code,
		StopAreaCode:       emptyCode,
		VehicleJourneyCode: code,
		DepartureStatus:    STOP_VISIT_DEPARTURE_CANCELLED,
		ArrivalStatus:      STOP_VISIT_ARRIVAL_ONTIME,
		Schedules:          schedules.NewStopVisitSchedules(),
	}

	manager.Update(event)
	updatedStopVisit, _ := model.StopVisits().Find(stopVisit.Id())
	if updatedStopVisit.DepartureStatus != STOP_VISIT_DEPARTURE_CANCELLED {
		t.Errorf("StopVisit DepartureStatus should be updated")
	}
	if updatedStopVisit.ArrivalStatus != STOP_VISIT_ARRIVAL_ONTIME {
		t.Errorf("StopVisit ArrivalStatus should be updated")
	}
	if !updatedStopVisit.IsCollected() {
		t.Errorf("StopVisit ArrivalStatus should be collected")
	}
	updatedStopArea, _ := model.StopAreas().Find(sa.Id())
	if !updatedStopArea.LineIds.Contains(l.Id()) {
		t.Errorf("StopArea LineIds should be updated")
	}
}

func Test_UpdateManager_UpdateStatus(t *testing.T) {
	model := NewTestMemoryModel()
	manager := newUpdateManager(model)

	sa := model.StopAreas().New()
	sa.Name = "Parent"
	sa.Save()

	sa2 := model.StopAreas().New()
	sa2.Name = "Son"
	sa2.ParentId = sa.id
	sa2.Save()

	sa3 := model.StopAreas().New()
	sa3.Name = "Grandson"
	sa3.ParentId = sa2.id
	sa3.Save()

	event := NewStatusUpdateEvent(sa3.Id(), "test_origin", true)
	manager.Update(event)

	stopArea, _ := model.StopAreas().Find(sa.Id())
	if status, ok := stopArea.Origins.Origin("test_origin"); !ok || !status {
		t.Errorf("Parent StopArea status should have been updated, got found origin: %v and status: %v", ok, status)
	}

	stopArea2, _ := model.StopAreas().Find(sa2.Id())
	if status, ok := stopArea2.Origins.Origin("test_origin"); !ok || !status {
		t.Errorf("StopArea status should have been updated, got found origin: %v and status: %v", ok, status)
	}

	stopArea3, _ := model.StopAreas().Find(sa3.Id())
	if status, ok := stopArea3.Origins.Origin("test_origin"); !ok || !status {
		t.Errorf("StopArea status should have been updated, got found origin: %v and status: %v", ok, status)
	}
}

func Test_UpdateManager_UpdateNotCollected(t *testing.T) {
	assert := assert.New(t)

	model := NewTestMemoryModel()
	manager := newUpdateManager(model)
	code := NewCode("codeSpace", "value")

	sa := model.StopAreas().New()
	sa.SetCode(code)
	sa.Save()

	l := model.Lines().New()
	l.SetCode(code)
	l.Save()

	stopVisit := model.StopVisits().New()
	stopVisit.SetCode(code)
	stopVisit.StopAreaId = sa.Id()
	stopVisit.collected = true
	stopVisit.Save()

	time := time.Now()

	manager.Update(NewNotCollectedUpdateEvent(code, time))
	updatedStopVisit, _ := model.StopVisits().Find(stopVisit.Id())

	assert.Equal(updatedStopVisit.ArrivalStatus, STOP_VISIT_ARRIVAL_ARRIVED)
	assert.Equal(updatedStopVisit.DepartureStatus, STOP_VISIT_DEPARTURE_DEPARTED)

	assert.False(updatedStopVisit.collected)

	assert.Equal(time, updatedStopVisit.Schedules.Schedule(schedules.Actual).ArrivalTime())
	assert.Equal(time, updatedStopVisit.Schedules.Schedule(schedules.Actual).DepartureTime())
}

func Test_UpdateManager_UpdateFreshVehicleJourney(t *testing.T) {
	InitTestDb(t)
	defer CleanTestDb(t)

	// Insert Data in the test db
	databaseVehicleJourney := DatabaseVehicleJourney{
		Id:              "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
		ReferentialSlug: "referential",
		ModelName:       "2017-01-01",
		Name:            "vehicleJourney",
		Codes:           `{"internal":"value"}`,
		LineId:          "c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
		Attributes:      "{}",
		References:      `{}`,
	}

	Database.AddTableWithName(databaseVehicleJourney, "vehicle_journeys")
	err := Database.Insert(&databaseVehicleJourney)
	if err != nil {
		t.Fatal(err)
	}

	// Fetch data from the db
	model := NewTestMemoryModel()
	model.date = Date{
		Year:  2017,
		Month: time.January,
		Day:   1,
	}
	vehicleJourneys := model.VehicleJourneys().(*MemoryVehicleJourneys)
	err = vehicleJourneys.Load("referential")
	if err != nil {
		t.Fatal(err)
	}

	vehicleJourneyId := VehicleJourneyId(databaseVehicleJourney.Id)
	_, ok := vehicleJourneys.Find(vehicleJourneyId)
	if !ok {
		t.Fatal("Loaded VehicleJourneys should be found")
	}

	file, err := os.Open("testdata/stopmonitoring-response-soap.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	response, _ := sxml.NewXMLStopMonitoringResponseFromContent(content)

	manager := newUpdateManager(model)

	code := NewCode("internal", "value")
	event := &VehicleJourneyUpdateEvent{
		CodeSpace: "internal",
		Code:      code,
		SiriXML:   &response.StopMonitoringDeliveries()[0].XMLMonitoredStopVisits()[0].XMLMonitoredVehicleJourney,
	}

	manager.Update(event)

	updatedVehicleJourney, _ := vehicleJourneys.Find(vehicleJourneyId)
	if updatedVehicleJourney.Attributes.IsEmpty() {
		t.Fatal("Attributes shouldn't be empty after update")
	}

	if updatedVehicleJourney.References.IsEmpty() {
		t.Fatalf("References shouldn't be empty after update: %v", updatedVehicleJourney.References)
	}
}

func Test_SituationUpdateManager_Update(t *testing.T) {
	assert := assert.New(t)
	code := NewCode("codeSpace", "value")
	testTime := time.Now()

	model := NewTestMemoryModel()
	situation := model.Situations().New()
	situation.SetCode(code)
	situation.SetCode(NewCode("_default", code.HashValue()))
	model.Situations().Save(situation)

	manager := newUpdateManager(model)
	event := completeEvent(code, testTime)

	manager.Update(event)

	updatedSituation, _ := model.Situations().Find(situation.Id())

	assert.True(checkSituation(*updatedSituation, code, testTime))

}

func Test_SituationUpdateManager_SameRecordedAtAndSameVersion(t *testing.T) {
	assert := assert.New(t)
	code := NewCode("codeSpace", "value")
	testTime := time.Now()

	model := NewTestMemoryModel()
	situation := model.Situations().New()
	situation.SetCode(code)
	situation.RecordedAt = testTime
	situation.Version = 1
	model.Situations().Save(situation)

	manager := newUpdateManager(model)
	event := completeEvent(code, testTime)

	manager.Update(event)

	updatedSituation, _ := model.Situations().Find(situation.Id())

	assert.False(checkSituation(*updatedSituation, code, testTime), "Situation should not be updated")
}

func completeEvent(code Code, testTime time.Time) (event *SituationUpdateEvent) {
	period := &TimeRange{EndTime: testTime}

	event = &SituationUpdateEvent{
		RecordedAt:      testTime,
		SituationCode:   code,
		Version:         1,
		ProducerRef:     "Ara",
		ValidityPeriods: []*TimeRange{period},
		Keywords:        []string{"channel"},
	}

	summary := &TranslatedString{
		DefaultValue: "Message Text",
	}

	event.Summary = summary
	event.Format = "format"

	return
}

func checkSituation(situation Situation, code Code, testTime time.Time) bool {
	summary := &TranslatedString{
		DefaultValue: "Message Text",
	}

	period := &TimeRange{EndTime: testTime}

	testSituation := Situation{
		id:              situation.id,
		Summary:         summary,
		RecordedAt:      testTime,
		ValidityPeriods: []*TimeRange{period},
		Format:          "format",
		Keywords:        []string{"channel"},
		ProducerRef:     "Ara",
		Version:         1,
	}

	testSituation.model = situation.model
	testSituation.codes = make(Codes)
	testSituation.SetCode(code)
	testSituation.SetCode(NewCode("_default", code.HashValue()))

	return reflect.DeepEqual(situation, testSituation)
}
