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
	generalMessage := response.XMLGeneralMessages()[0]
	content := generalMessage.Content().(IDFGeneralMessageStructure)
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

func checkGeneralMessagesEquivalence(s1 *XMLGeneralMessageResponse, s2 *XMLGeneralMessageResponse, t *testing.T) {

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

	expectedContent := expectedGM.Content().(IDFGeneralMessageStructure)
	gotContent := gotGM.Content().(IDFGeneralMessageStructure)

	expedtedMessages := expectedContent.Messages()[0]
	gotMessages := gotContent.Messages()[0]

	if expedtedMessages.messageText != gotMessages.messageText {
		t.Errorf("Wrong Message Content: \n got: %v\nwant: %v", expedtedMessages.messageText, gotMessages.messageText)
	}

	if expedtedMessages.numberOfLines != gotMessages.numberOfLines {
		t.Errorf("Wrong Message NumberOfLines: \n got: %v\nwant: %v", expedtedMessages.numberOfLines, gotMessages.numberOfLines)
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

func Test_SIRIGeneralMessageResponse_BuildXML(t *testing.T) {
	expectedXML := `<sw:GetGeneralMessageResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
	<ServiceDeliveryInfo>
		<siri:ResponseTimestamp>2016-09-21T20:14:46.000Z</siri:ResponseTimestamp>
		<siri:ProducerRef>producer</siri:ProducerRef>
		<siri:Address>address</siri:Address>
		<siri:ResponseMessageIdentifier>identifier</siri:ResponseMessageIdentifier>
		<siri:RequestMessageRef>ref</siri:RequestMessageRef>
	</ServiceDeliveryInfo>
	<Answer>
		<siri:GeneralMessageDelivery version="2.0:FR-IDF-2.4">
			<siri:ResponseTimestamp>2016-09-21T20:14:46.000Z</siri:ResponseTimestamp>
			<siri:Status>true</siri:Status>
			<siri:GeneralMessage formatRef="STIF-IDF">
				<siri:RecordedAtTime>2016-09-21T20:14:46.000Z</siri:RecordedAtTime>
				<siri:ValidUntilTime>2016-09-21T20:14:46.000Z</siri:ValidUntilTime>
				<siri:InfoMessageVersion>1</siri:InfoMessageVersion>
				<siri:InfoChannelRef>Chan</siri:InfoChannelRef>
				<siri:Content xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
				xsi:type="ns9:IDFGeneralMessageStructure">
					<Message>
						<MessageType>Je suis de type texte</MessageType>
						<MessageText xml:lang="NL">Je suis un texte</MessageText>
					</Message>
					<LineSection>
					  <FirstStop>NINOXE:StopPoint:SP:24:LOC</FirstStop>
					  <LastStop>NINOXE:StopPoint:SP:12:LOC</LastStop>
					  <LineRef>NINOXE:Line::3:LOC</LineRef>
					</LineSection>
				</siri:Content>
			</siri:GeneralMessage>{{end}}
		</siri:GeneralMessageDelivery>
	</Answer>
</sw:GetGeneralMessageResponse>`

	response, _ := NewXMLGeneralMessageResponseFromContent([]byte(expectedXML))
	responseTimestamp := time.Date(2016, time.September, 21, 20, 14, 46, 0, time.UTC)

	request := &SIRIGeneralMessageResponse{
		Address:                   "address",
		ProducerRef:               "producer",
		ResponseMessageIdentifier: "identifier",
	}

	lineSection := &SIRILineSection{
		FirstStop: "NINOXE:StopPoint:SP:24:LOC",
		LastStop:  "NINOXE:StopPoint:SP:12:LOC",
		LineRef:   "NINOXE:Line::3:LOC",
	}

	gM := &SIRIGeneralMessage{
		RecordedAtTime: time.Date(2016, time.September, 21, 20, 14, 46, 0, time.UTC),
		ValidUntilTime: time.Date(2016, time.September, 21, 20, 14, 46, 0, time.UTC),
		LineSections:   []*SIRILineSection{lineSection},
	}

	request.Status = true
	request.ResponseTimestamp = responseTimestamp
	request.RequestMessageRef = "ref"

	request.GeneralMessages = []*SIRIGeneralMessage{gM}
	request.GeneralMessages[0].Messages = append(request.GeneralMessages[0].Messages, &SIRIMessage{Content: "Je suis un texte", Type: "Un Type"})
	request.GeneralMessages[0].InfoMessageVersion = 1
	request.GeneralMessages[0].InfoChannelRef = "Chan"
	request.GeneralMessages[0].FormatRef = "STIF-IDF"

	xml, err := request.BuildXML()
	if err != nil {
		t.Fatal(err)
	}

	xmlResponse, _ := NewXMLGeneralMessageResponseFromContent([]byte(xml))

	checkGeneralMessagesEquivalence(response, xmlResponse, t)
}
