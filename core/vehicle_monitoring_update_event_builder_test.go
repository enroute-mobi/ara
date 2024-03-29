package core

import (
	"io"
	"math"
	"os"
	"testing"

	"bitbucket.org/enroute-mobi/ara/siri/sxml"
)

func getvm(t *testing.T, filePath string) *sxml.XMLVehicleMonitoringResponse {
	file, err := os.Open(filePath)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	response, err := sxml.NewXMLVehicleMonitoringResponseFromContent(content)
	if err != nil {
		t.Fatal(err)
	}

	return response
}

func Test_Vehicle_Code_With_VehicleRef(t *testing.T) {
	vm := getvm(t, "testdata/vm_response_soap.xml")
	mvj := vm.VehicleMonitoringDeliveries()[0].VehicleActivities()[0].XMLMonitoredVehicleJourney

	if vehicleRef := "RLA290"; mvj.VehicleRef() != vehicleRef {
		t.Errorf("Wrong VehicleRef. Expected %v, got %v", vehicleRef, mvj.VehicleRef())
	}
}

func Test_Vehicle_Code_Without_VehicleRef_With_VehicleMonitoringRef(t *testing.T) {
	vm := getvm(t, "testdata/vm_response_soap2.xml")
	mvj := vm.VehicleMonitoringDeliveries()[0].VehicleActivities()[0].XMLMonitoredVehicleJourney

	if vehicleRef := "TRANSDEV:Vehicle::7658:LOC"; mvj.VehicleRef() != vehicleRef {
		t.Errorf("Wrong VehicleRef. Expected %v, got %v", vehicleRef, mvj.VehicleRef())
	}
}

func Test_Coordinates_Transform(t *testing.T) {
	p := NewPartner()
	builder := NewVehicleMonitoringUpdateEventBuilder(p)

	vm := getvm(t, "testdata/vm_response_soap.xml")
	va := vm.VehicleMonitoringDeliveries()[0].VehicleActivities()[0]

	lon, lat, err := builder.handleCoordinates(va)

	if err != nil {
		t.Errorf("Error while converting: %v", err)
	}

	if e := 7.2761920740520; round(lon) != e {
		t.Errorf("Wrong coord longitude. Expected %v, got %v", e, round(lon))
	}
	if e := 43.703478617764; round(lat) != e {
		t.Errorf("Wrong coord latitude. Expected %v, got %v", e, round(lat))
	}
}

func round(n float64) float64 {
	r := math.Pow(10, 12)
	return math.Round((n * r)) / r
}

func Test_Coordinates_WithLonLat(t *testing.T) {
	p := NewPartner()
	builder := NewVehicleMonitoringUpdateEventBuilder(p)

	vm := getvm(t, "testdata/vm_response_soap.xml")
	va := vm.VehicleMonitoringDeliveries()[0].VehicleActivities()[0]

	va.SetLongitude("1.1234")
	va.SetLatitude("2.3456")

	lon, lat, err := builder.handleCoordinates(va)

	if err != nil {
		t.Errorf("Error while converting: %v", err)
	}

	if e := 1.1234; lon != e {
		t.Errorf("Wrong coord longitude. Expected %v, got %v", e, lon)
	}
	if e := 2.3456; lat != e {
		t.Errorf("Wrong coord latitude. Expected %v, got %v", e, lat)
	}
}

func Test_Coordinates_WithLongitude(t *testing.T) {
	p := NewPartner()
	builder := NewVehicleMonitoringUpdateEventBuilder(p)

	vm := getvm(t, "testdata/vm_response_soap.xml")
	va := vm.VehicleMonitoringDeliveries()[0].VehicleActivities()[0]

	va.SetLongitude("1.1234")

	lon, lat, err := builder.handleCoordinates(va)

	if err != nil {
		t.Errorf("Error while converting: %v", err)
	}

	if e := 1.1234; lon != e {
		t.Errorf("Wrong coord longitude. Expected %v, got %v", e, lon)
	}
	if e := 0.0; lat != e {
		t.Errorf("Wrong coord latitude. Expected %v, got %v", e, lat)
	}
}

func Test_Coordinates_WithLatitude(t *testing.T) {
	p := NewPartner()
	builder := NewVehicleMonitoringUpdateEventBuilder(p)

	vm := getvm(t, "testdata/vm_response_soap.xml")
	va := vm.VehicleMonitoringDeliveries()[0].VehicleActivities()[0]

	va.SetLatitude("2.3456")

	lon, lat, err := builder.handleCoordinates(va)

	if err != nil {
		t.Errorf("Error while converting: %v", err)
	}

	if e := 0.0; lon != e {
		t.Errorf("Wrong coord longitude. Expected %v, got %v", e, lon)
	}
	if e := 2.3456; lat != e {
		t.Errorf("Wrong coord latitude. Expected %v, got %v", e, lat)
	}
}

func Test_Coordinates_InvalidSRS(t *testing.T) {
	p := NewPartner()
	builder := NewVehicleMonitoringUpdateEventBuilder(p)

	vm := getvm(t, "testdata/vm_response_soap.xml")
	va := vm.VehicleMonitoringDeliveries()[0].VehicleActivities()[0]

	va.SetSRSName("invalid srs name")

	lon, lat, err := builder.handleCoordinates(va)

	if err == nil {
		t.Errorf("Converting coordinates should return an error, got nothing and the following coordinates: %v, %v", lon, lat)
	}
}

func Test_Coordinates_InvalidCoordinates(t *testing.T) {
	p := NewPartner()
	builder := NewVehicleMonitoringUpdateEventBuilder(p)

	vm := getvm(t, "testdata/vm_response_soap.xml")
	va := vm.VehicleMonitoringDeliveries()[0].VehicleActivities()[0]

	va.SetCoordinates("invalid")

	lon, lat, err := builder.handleCoordinates(va)

	if err == nil {
		t.Errorf("Converting coordinates should return an error, got nothing and the following coordinates: %v, %v", lon, lat)
	}
}
