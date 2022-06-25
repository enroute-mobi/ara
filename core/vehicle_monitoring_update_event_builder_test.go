package core

import (
	"io/ioutil"
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
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	response, err := sxml.NewXMLVehicleMonitoringResponseFromContent(content)
	if err != nil {
		t.Fatal(err)
	}

	return response
}

func Test_Coordinates_Transform(t *testing.T) {
	p := NewPartner()
	builder := NewVehicleMonitoringUpdateEventBuilder(p)

	vm := getvm(t, "testdata/vm_response_soap.xml")
	va := vm.VehicleMonitoringDeliveries()[0].VehicleActivities()[0]

	coord, err := builder.handleCoordinates(va)

	if err != nil {
		t.Errorf("Error while converting: %v", err)
	}

	if e := 7.2761920740520; round(coord.X) != e {
		t.Errorf("Wrong coord longitude. Expected %v, got %v", e, round(coord.X))
	}
	if e := 43.703478618706; round(coord.Y) != e {
		t.Errorf("Wrong coord latitude. Expected %v, got %v", e, round(coord.Y))
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

	coord, err := builder.handleCoordinates(va)

	if err != nil {
		t.Errorf("Error while converting: %v", err)
	}

	if e := 1.1234; coord.X != e {
		t.Errorf("Wrong coord longitude. Expected %v, got %v", e, coord.X)
	}
	if e := 2.3456; coord.Y != e {
		t.Errorf("Wrong coord latitude. Expected %v, got %v", e, coord.Y)
	}
}

func Test_Coordinates_WithLongitude(t *testing.T) {
	p := NewPartner()
	builder := NewVehicleMonitoringUpdateEventBuilder(p)

	vm := getvm(t, "testdata/vm_response_soap.xml")
	va := vm.VehicleMonitoringDeliveries()[0].VehicleActivities()[0]

	va.SetLongitude("1.1234")

	coord, err := builder.handleCoordinates(va)

	if err != nil {
		t.Errorf("Error while converting: %v", err)
	}

	if e := 1.1234; coord.X != e {
		t.Errorf("Wrong coord longitude. Expected %v, got %v", e, coord.X)
	}
	if e := 0.0; coord.Y != e {
		t.Errorf("Wrong coord latitude. Expected %v, got %v", e, coord.Y)
	}
}

func Test_Coordinates_WithLatitude(t *testing.T) {
	p := NewPartner()
	builder := NewVehicleMonitoringUpdateEventBuilder(p)

	vm := getvm(t, "testdata/vm_response_soap.xml")
	va := vm.VehicleMonitoringDeliveries()[0].VehicleActivities()[0]

	va.SetLatitude("2.3456")

	coord, err := builder.handleCoordinates(va)

	if err != nil {
		t.Errorf("Error while converting: %v", err)
	}

	if e := 0.0; coord.X != e {
		t.Errorf("Wrong coord longitude. Expected %v, got %v", e, coord.X)
	}
	if e := 2.3456; coord.Y != e {
		t.Errorf("Wrong coord latitude. Expected %v, got %v", e, coord.Y)
	}
}

func Test_Coordinates_InvalidSRS(t *testing.T) {
	p := NewPartner()
	builder := NewVehicleMonitoringUpdateEventBuilder(p)

	vm := getvm(t, "testdata/vm_response_soap.xml")
	va := vm.VehicleMonitoringDeliveries()[0].VehicleActivities()[0]

	va.SetSRSName("invalid srs name")

	coord, err := builder.handleCoordinates(va)

	if err == nil {
		t.Errorf("Converting coordinates should return an error, got nothing and the following coordinates: %v", coord)
	}
}

func Test_Coordinates_InvalidCoordinates(t *testing.T) {
	p := NewPartner()
	builder := NewVehicleMonitoringUpdateEventBuilder(p)

	vm := getvm(t, "testdata/vm_response_soap.xml")
	va := vm.VehicleMonitoringDeliveries()[0].VehicleActivities()[0]

	va.SetCoordinates("invalid")

	coord, err := builder.handleCoordinates(va)

	if err == nil {
		t.Errorf("Converting coordinates should return an error, got nothing and the following coordinates: %v", coord)
	}
}
