package siri_tests

import (
	"io"
	"os"
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/siri/sxml"
	"github.com/stretchr/testify/require"
)

func getXMLGeneralMessageResponse(t *testing.T) *sxml.XMLGeneralMessageResponse {
	file, err := os.Open("testdata/general-messages-response.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := io.ReadAll(file)
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
	require := require.New(t)

	require.Equal(s1.Address(), s2.Address())
	require.Equal(s1.ProducerRef(), s2.ProducerRef())
	require.Equal(s1.ResponseMessageIdentifier(), s2.ResponseMessageIdentifier())
	require.Equal(s1.Status(), s2.Status())
	require.Equal(s1.ResponseTimestamp(), s2.ResponseTimestamp())
	require.Len(s1.XMLGeneralMessages(), len(s2.XMLGeneralMessages()))

	expectedGM := s1.XMLGeneralMessages()[0]
	gotGM := s2.XMLGeneralMessages()[0]

	require.Equal(expectedGM.RecordedAtTime(), gotGM.RecordedAtTime())
	require.Equal(expectedGM.ValidUntilTime(), gotGM.ValidUntilTime())
	require.Equal(expectedGM.InfoMessageVersion(), gotGM.InfoMessageVersion())
	require.Equal(expectedGM.InfoMessageIdentifier(), gotGM.InfoMessageIdentifier())
	require.Equal(expectedGM.FormatRef(), gotGM.FormatRef())
	require.Equal(expectedGM.InfoChannelRef(), gotGM.InfoChannelRef())

	expectedContent := expectedGM.Content().(sxml.IDFGeneralMessageStructure)
	gotContent := gotGM.Content().(sxml.IDFGeneralMessageStructure)

	expedtedMessages := expectedContent.Messages()[0]
	gotMessages := gotContent.Messages()[0]

	require.Equal(expedtedMessages.MessageTexts(), gotMessages.MessageTexts())

	expectedLineSection := expectedContent.LineSections()[0]
	gotLineSection := gotContent.LineSections()[0]

	require.Equal(expectedLineSection.LineRef(), gotLineSection.LineRef())
	require.Equal(expectedLineSection.FirstStop(), gotLineSection.FirstStop())
	require.Equal(expectedLineSection.LastStop(), gotLineSection.LastStop())
}
