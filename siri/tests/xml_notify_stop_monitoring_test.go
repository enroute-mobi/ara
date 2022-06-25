package siri_tests

import (
	"io/ioutil"
	"os"
	"testing"

	"bitbucket.org/enroute-mobi/ara/siri/sxml"
)

func getXMLNotifyStopMonitoring(t *testing.T) *sxml.XMLNotifyStopMonitoring {
	file, err := os.Open("testdata/notify_stop_monitoring.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	xmlStructure, _ := sxml.NewXMLNotifyStopMonitoringFromContent(content)
	return xmlStructure
}

func Test_XMLNotifyStopMonitoring_Deliveries(t *testing.T) {
	xmlStructure := getXMLNotifyStopMonitoring(t)

	if expected := 2; len(xmlStructure.StopMonitoringDeliveries()) != expected {
		t.Errorf("Wrong StopMonitoringDeliveries count:\n got: %v\nwant: %v", xmlStructure.StopMonitoringDeliveries(), expected)
	}

	if len(xmlStructure.StopMonitoringDeliveries()[0].XMLMonitoredStopVisitCancellations()) != 1 {
		t.Errorf("No MonitoredStopVisitCancellation found in first StopMonitoringDelivery")
	}

	cancellation := xmlStructure.StopMonitoringDeliveries()[0].XMLMonitoredStopVisitCancellations()[0]
	if expected := "SIRI:43745132"; cancellation.ItemRef() != expected {
		t.Errorf("Wrong ItemRef in MonitoredStopVisitCancellation:\n got: %v\nwant: %v", cancellation.ItemRef(), expected)
	}

	if len(xmlStructure.StopMonitoringDeliveries()[1].XMLMonitoredStopVisits()) != 1 {
		t.Errorf("No MonitoredStopVisit found in second StopMonitoringDelivery")
	}

	monitoredStopVisit := xmlStructure.StopMonitoringDeliveries()[1].XMLMonitoredStopVisits()[0]
	if expected := "SIRI:43771841"; monitoredStopVisit.ItemIdentifier() != expected {
		t.Errorf("Wrong ItemIdentifier in MonitoredStopVisit:\n got: %v\nwant: %v", monitoredStopVisit.ItemIdentifier(), expected)
	}
}
