package siri_tests

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/siri/sxml"
)

func getXMLGeneralMessageResponse(t *testing.T) *sxml.XMLGeneralMessageResponse {
	file, err := os.Open("testdata/general-messages-response.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	response, _ := sxml.NewXMLGeneralMessageResponseFromContent(content)
	return response
}

func Test_XMLGeneralMessage(t *testing.T) {
	response := getXMLGeneralMessageResponse(t)
	generalMessage := response.XMLGeneralMessages()[0]
	content := generalMessage.Content().(sxml.IDFGeneralMessageStructure)
	lineSection := content.LineSections()[0]

	if expected := time.Date(2017, time.March, 29, 03, 30, 06, 0, generalMessage.RecordedAtTime().Location()); generalMessage.RecordedAtTime() != expected {
		t.Errorf("Wrong RecordedAtTime: \n got: %v\nwant: %v", generalMessage.RecordedAtTime(), expected)
	}

	if expected := "3477"; generalMessage.ItemIdentifier() != expected {
		t.Errorf("Wrong ItemIdentifier: \n got: %v\nwant: %v", generalMessage.ItemIdentifier(), expected)
	}

	if expected := "NINOXE:GeneralMessage:27_1"; generalMessage.InfoMessageIdentifier() != expected {
		t.Errorf("Wrong InfoMessageIdentifier: \n got: %v\nwant: %v", generalMessage.InfoMessageIdentifier(), expected)
	}

	if expected := 1; generalMessage.InfoMessageVersion() != expected {
		t.Errorf("Wrong InfoMessageVersion: \n got: %v\nwant: %v", generalMessage.InfoMessageVersion(), expected)
	}

	if expected := "Commercial"; generalMessage.InfoChannelRef() != expected {
		t.Errorf("Wrong InfoChannelRef: \n got: %v\nwant: %v", generalMessage.InfoChannelRef(), expected)
	}

	if expected := time.Date(2017, time.March, 29, 03, 30, 06, 0, generalMessage.ValidUntilTime().Location()); generalMessage.ValidUntilTime() != expected {
		t.Errorf("Wrong RecordedAtTime: \n got: %v\nwant: %v", generalMessage.ValidUntilTime(), expected)
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

func checkGeneralMessagesEquivalence(s1 *sxml.XMLGeneralMessageResponse, s2 *sxml.XMLGeneralMessageResponse, t *testing.T) {

	if s1.Address() != s2.Address() {
		t.Errorf("Wrong Address: \n got: %v\nwant: %v", s2.Address(), s1.Address())
	}

	if s1.ProducerRef() != s2.ProducerRef() {
		t.Errorf("Wrong ProducerRef: \n got: %v\nwant: %v", s2.ProducerRef(), s1.ProducerRef())
	}

	if s1.ResponseMessageIdentifier() != s2.ResponseMessageIdentifier() {
		t.Errorf("Wrong ResponseMessageIdentifier: \n got: %v\nwant: %v", s2.ResponseMessageIdentifier(), s1.ResponseMessageIdentifier())
	}

	if s1.Status() != s2.Status() {
		t.Errorf("Wrong Status: \n got: %v\nwant: %v", s2.Status(), s1.Status())
	}

	if s1.ResponseTimestamp() != s2.ResponseTimestamp() {
		t.Errorf("Wrong ResponseTimestamp: \n got: %v\nwant: %v", s2.ResponseTimestamp(), s1.ResponseTimestamp())
	}

	if len(s1.XMLGeneralMessages()) != len(s2.XMLGeneralMessages()) {
		t.Errorf("Wrong XMLGeneralMessage: \n got: %v\nwant: %v", len(s2.XMLGeneralMessages()), len(s1.XMLGeneralMessages()))
	}

	expectedGM := s1.XMLGeneralMessages()[0]
	gotGM := s2.XMLGeneralMessages()[0]

	if expectedGM.RecordedAtTime() != gotGM.RecordedAtTime() {
		t.Errorf("Wrong RecordedAtTime: \n got: %v\nwant: %v", gotGM.RecordedAtTime(), expectedGM.RecordedAtTime())
	}

	if expectedGM.ValidUntilTime() != gotGM.ValidUntilTime() {
		t.Errorf("Wrong ValidUntilTime: \n got: %v\nwant: %v", gotGM.ValidUntilTime(), expectedGM.ValidUntilTime())
	}

	if expectedGM.InfoMessageVersion() != gotGM.InfoMessageVersion() {
		t.Errorf("Wrong InfoMessageVersion: \n got: %v\nwant: %v", gotGM.InfoMessageVersion(), expectedGM.InfoMessageVersion())
	}

	if expectedGM.InfoMessageIdentifier() != gotGM.InfoMessageIdentifier() {
		t.Errorf("Wrong InfoMessageIdentifier: \n got: %v\nwant: %v", gotGM.InfoMessageIdentifier(), expectedGM.InfoMessageIdentifier())
	}

	if expectedGM.FormatRef() != gotGM.FormatRef() {
		t.Errorf("Wrong FormatRef: \n got: %v\nwant: %v", gotGM.FormatRef(), expectedGM.FormatRef())
	}

	if expectedGM.InfoChannelRef() != gotGM.InfoChannelRef() {
		t.Errorf("Wrong InfoChannelRef: \n got: %v\n want: %v", gotGM.InfoChannelRef(), expectedGM.RecordedAtTime())
	}

	expectedContent := expectedGM.Content().(sxml.IDFGeneralMessageStructure)
	gotContent := gotGM.Content().(sxml.IDFGeneralMessageStructure)

	expedtedMessages := expectedContent.Messages()[0]
	gotMessages := gotContent.Messages()[0]

	if expedtedMessages.MessageText() != gotMessages.MessageText() {
		t.Errorf("Wrong Message Content: \n got: %v\nwant: %v", expedtedMessages.MessageText(), gotMessages.MessageText())
	}

	if expedtedMessages.NumberOfLines() != gotMessages.NumberOfLines() {
		t.Errorf("Wrong Message NumberOfLines: \n got: %v\nwant: %v", expedtedMessages.NumberOfLines(), gotMessages.NumberOfLines())
	}

	expectedLineSection := expectedContent.LineSections()[0]
	gotLineSection := gotContent.LineSections()[0]

	if expectedLineSection.LineRef() != gotLineSection.LineRef() {
		t.Errorf("Wrong MessageType: \n got: %v\nwant: %v", gotLineSection.LineRef(), expectedLineSection.LineRef())
	}

	if expectedLineSection.FirstStop() != gotLineSection.FirstStop() {
		t.Errorf("Wrong MessageType: \n got: %v\nwant: %v", gotLineSection.FirstStop(), expectedLineSection.FirstStop())
	}

	if expectedLineSection.LastStop() != gotLineSection.LastStop() {
		t.Errorf("Wrong MessageType: \n got: %v\nwant: %v", gotLineSection.LastStop(), expectedLineSection.LastStop())
	}

}
