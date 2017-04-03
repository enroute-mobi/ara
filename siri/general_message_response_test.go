package siri

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func getXMLGeneralMessageResponse(t *testing.T) *XMLGeneralMessageResponse {
	file, err := os.Open("testdata/general-messages-response.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	response, _ := NewXMLGeneralMessageResponseFromContent(content)
	return response
}

func Test_XMLGeneralMessage(t *testing.T) {
	response := getXMLGeneralMessageResponse(t)
	contents := response.Contents().([]*IDFGeneralMessageStructure)
	content := contents[0]
	lineSection := content.LineSection()

	if expected := time.Date(2017, time.March, 29, 03, 30, 06, 0, response.RecordedAtTime().Location()); response.RecordedAtTime() != expected {
		t.Errorf("Wrong RecordedAtTime: \n got: %v\nwant: %v", response.RecordedAtTime(), expected)
	}

	if expected := "3477"; response.ItemIdentifier() != expected {
		t.Errorf("Wrong ItemIdentifier: \n got: %v\nwant: %v", response.ItemIdentifier(), expected)
	}

	if expected := "NINOXE:GeneralMessage:27_1"; response.InfoMessageIdentifier() != expected {
		t.Errorf("Wrong InfoMessageIdentifier: \n got: %v\nwant: %v", response.InfoMessageIdentifier(), expected)
	}

	if expected := "1"; response.InfoMessageVersion() != expected {
		t.Errorf("Wrong InfoMessageVersion: \n got: %v\nwant: %v", response.InfoMessageVersion(), expected)
	}

	if expected := "Commercial"; response.InfoChannelRef() != expected {
		t.Errorf("Wrong InfoChannelRef: \n got: %v\nwant: %v", response.InfoChannelRef(), expected)
	}

	if expected := time.Date(2017, time.March, 29, 03, 30, 06, 0, response.ValidUntilTime().Location()); response.ValidUntilTime() != expected {
		t.Errorf("Wrong RecordedAtTime: \n got: %v\nwant: %v", response.ValidUntilTime(), expected)
	}

	if expected := "longMessage"; content.MessageType() != expected {
		t.Errorf("Wrong MessageType: \n got: %v\nwant: %v", content.MessageType(), expected)
	}

	if expected := `La nouvelle carte d'abonnement est disponible au points de vente du réseau`; content.MessageText() != expected {
		t.Errorf("Wrong MessageText: \n got: %v\nwant: %v", content.MessageText(), expected)
	}

	if expected := "longMessage"; content.MessageType() != expected {
		t.Errorf("Wrong MessageType: \n got: %v\nwant: %v", content.MessageType(), expected)
	}

	if expected := "NINOXE:StopPoint:SP:24:LOC"; lineSection.FirstStop() != expected {
		t.Errorf("Wrong MessageType: \n got: %v\nwant: %v", lineSection.FirstStop(), expected)
	}

	if expected := "NINOXE:StopPoint:SP:12:LOC"; lineSection.LastStop() != expected {
		t.Errorf("Wrong lastStop: \n got: %v\nwant: %v", lineSection.LastStop(), expected)
	}

	if expected := "NINOXE:Line::3:LOC"; lineSection.LineRef() != expected {
		t.Errorf("Wrong lineRef: \n got: %v\nwant: %v", lineSection.LineRef(), expected)
	}

}
