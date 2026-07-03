package core

import (
	"io"
	"os"
	"testing"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	s "bitbucket.org/enroute-mobi/ara/core/settings"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
	"bitbucket.org/enroute-mobi/ara/uuid"
	"github.com/stretchr/testify/assert"
)

func setupVehicleMonitoringModel(referential *Referential) *model.Vehicle {
	lineCode := model.NewCode("internal", "line-ref")
	line := referential.model.Lines().New()
	line.SetCode(lineCode)
	line.Save()

	vjCode := model.NewCode("internal", "vj-ref")
	vj := referential.model.VehicleJourneys().New()
	vj.SetCode(vjCode)
	vj.LineId = line.Id()
	vj.Save()

	vehicleCode := model.NewCode("internal", "vehicle-ref")
	vehicle := referential.model.Vehicles().New()
	vehicle.SetCode(vehicleCode)
	vehicle.LineId = line.Id()
	vehicle.VehicleJourneyId = vj.Id()
	vehicle.RawAttributes[siri_attributes.VehicleActivityNote] = "A vehicle note"
	vehicle.Save()

	return vehicle
}

func readVehicleMonitoringRequest(t *testing.T) *sxml.XMLGetVehicleMonitoring {
	t.Helper()
	file, err := os.Open("testdata/vehiclemonitoring-request-soap.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	request, err := sxml.NewXMLGetVehicleMonitoringFromContent(content)
	if err != nil {
		t.Fatal(err)
	}
	return request
}

func Test_SIRIVehicleMonitoringRequestBroadcaster_RequestVehicles_IgnoreNotes(t *testing.T) {
	assert := assert.New(t)

	_, referential := newTestReferential(t)
	partner := referential.Partners().New("partner")
	partner.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, map[string]string{
		"remote_code_space":           "internal",
		s.BROADCAST_SIRI_IGNORE_NOTES: "true",
	})

	connector := NewSIRIVehicleMonitoringRequestBroadcaster(partner)
	connector.SetClock(clock.NewFakeClock())
	connector.Start()

	setupVehicleMonitoringModel(referential)

	request := readVehicleMonitoringRequest(t)
	response := connector.RequestVehicles(request, &audit.BigQueryMessage{})

	if len(response.SIRIVehicleMonitoringDelivery.VehicleActivity) != 1 {
		t.Fatalf("Expected 1 VehicleActivity, got %d", len(response.SIRIVehicleMonitoringDelivery.VehicleActivity))
	}
	assert.Empty(response.SIRIVehicleMonitoringDelivery.VehicleActivity[0].VehicleActivityNote,
		"VehicleActivityNote should be empty when ignore_notes is true")
}

func Test_SIRIVehicleMonitoringRequestBroadcaster_RequestVehicles_DoNotIgnoreNotes(t *testing.T) {
	assert := assert.New(t)

	_, referential := newTestReferential(t)
	partner := referential.Partners().New("partner")
	partner.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, map[string]string{
		"remote_code_space": "internal",
	})

	connector := NewSIRIVehicleMonitoringRequestBroadcaster(partner)
	connector.SetClock(clock.NewFakeClock())
	connector.Start()

	setupVehicleMonitoringModel(referential)

	request := readVehicleMonitoringRequest(t)
	response := connector.RequestVehicles(request, &audit.BigQueryMessage{})

	if len(response.SIRIVehicleMonitoringDelivery.VehicleActivity) != 1 {
		t.Fatalf("Expected 1 VehicleActivity, got %d", len(response.SIRIVehicleMonitoringDelivery.VehicleActivity))
	}
	assert.Equal("A vehicle note", response.SIRIVehicleMonitoringDelivery.VehicleActivity[0].VehicleActivityNote)
}
