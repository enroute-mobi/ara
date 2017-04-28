package siri

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/af83/edwig/model"
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
	generalMessage := response.XMLGeneralMessage()[0]
	content := generalMessage.Content().(IDFGeneralMessageStructure)
	lineSection := content.LineSection()

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

	if len(s1.XMLGeneralMessage()) != len(s2.XMLGeneralMessage()) {
		t.Errorf("Wrong XMLGeneralMessage: \n got: %v\nwant: %v", len(s2.XMLGeneralMessage()), len(s1.XMLGeneralMessage()))
	}

	expectedGM := s1.XMLGeneralMessage()[0]
	gotGM := s2.XMLGeneralMessage()[0]

	if expectedGM.RecordedAtTime() != gotGM.RecordedAtTime() {
		t.Errorf("Wrong RecordedAtTime: \n got: %v\nwant: %v", gotGM.RecordedAtTime(), expectedGM.RecordedAtTime())
	}

	if expectedGM.ValidUntilTime() != gotGM.ValidUntilTime() {
		t.Errorf("Wrong ValidUntilTime: \n got: %v\nwant: %v", gotGM.ValidUntilTime(), expectedGM.ValidUntilTime())
	}

	if expectedGM.InfoMessageIdentifier() != gotGM.InfoMessageIdentifier() {
		t.Errorf("Wrong InfoMessageIdentifier: \n got: %v\nwant: %v", gotGM.InfoMessageIdentifier(), expectedGM.InfoMessageIdentifier())
	}

	if expectedGM.InfoMessageVersion() != gotGM.InfoMessageVersion() {
		t.Errorf("Wrong InfoMessageVersion: \n got: %v\nwant: %v", gotGM.InfoMessageVersion(), expectedGM.InfoMessageVersion())
	}

	if expectedGM.InfoMessageIdentifier() != gotGM.InfoMessageIdentifier() {
		t.Errorf("Wrong InfoMessageIdentifier: \n got: %v\nwant: %v", gotGM.InfoMessageIdentifier(), expectedGM.InfoMessageIdentifier())
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

	expectedLineSection := expectedContent.LineSection()
	gotLineSection := gotContent.LineSection()

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
	expectedXML := `<ns8:GetGeneralMessageResponse xmlns:ns3="http://www.siri.org.uk/siri"
															 xmlns:ns4="http://www.ifopt.org.uk/acsb"
															 xmlns:ns5="http://www.ifopt.org.uk/ifopt"
															 xmlns:ns6="http://datex2.eu/schema/2_0RC1/2_0"
															 xmlns:ns7="http://scma/siri"
															 xmlns:ns8="http://wsdl.siri.org.uk"
															 xmlns:ns9="http://wsdl.siri.org.uk/siri">
	<ServiceDeliveryInfo>
		<ns3:ResponseTimestamp>2016-09-21T20:14:46.000Z</ns3:ResponseTimestamp>
		<ns3:ProducerRef>producer</ns3:ProducerRef>
		<ns3:Address>address</ns3:Address>
		<ns3:ResponseMessageIdentifier>identifier</ns3:ResponseMessageIdentifier>
		<ns3:RequestMessageRef>ref</ns3:RequestMessageRef>
	</ServiceDeliveryInfo>
	<Answer>
		<ns3:GeneralMessageDelivery version="2.0:FR-IDF-2.4">
			<ns3:ResponseTimestamp>2016-09-21T20:14:46.000Z</ns3:ResponseTimestamp>
			<ns3:Status>true</ns3:Status>
			<ns3:GeneralMessage formatRef="FRANCE">
				<ns3:RecordedAtTime>2016-09-21T20:14:46.000Z</ns3:RecordedAtTime>
				<ns3:ValidUntilTime>2016-09-21T20:14:46.000Z</ns3:ValidUntilTime>
				<ns3:InfoMessageVersion>1</ns3:InfoMessageVersion>
				<ns3:InfoChannelRef>Chan</ns3:InfoChannelRef>
				<ns3:Content xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
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
				</ns3:Content>
			</ns3:GeneralMessage>{{end}}
		</ns3:GeneralMessageDelivery>
	</Answer>
	<AnswerExtension/>
</ns8:GetGeneralMessageResponse>`

	response, _ := NewXMLGeneralMessageResponseFromContent([]byte(expectedXML))
	responseTimestamp := time.Date(2016, time.September, 21, 20, 14, 46, 0, time.UTC)

	request := &SIRIGeneralMessageResponse{
		Address:                   "address",
		ProducerRef:               "producer",
		ResponseMessageIdentifier: "identifier",
	}

	gM := &SIRIGeneralMessage{
		RecordedAtTime: time.Date(2016, time.September, 21, 20, 14, 46, 0, time.UTC),
		ValidUntilTime: time.Date(2016, time.September, 21, 20, 14, 46, 0, time.UTC),
		FirstStop:      "NINOXE:StopPoint:SP:24:LOC",
		LastStop:       "NINOXE:StopPoint:SP:12:LOC",
		LineRef:        "NINOXE:Line::3:LOC",
	}

	request.Status = true
	request.ResponseTimestamp = responseTimestamp
	request.RequestMessageRef = "ref"

	request.GeneralMessages = []*SIRIGeneralMessage{gM}
	request.GeneralMessages[0].Messages = append(request.GeneralMessages[0].Messages, &model.Message{Content: "Je suis un texte", Type: "Un Type"})
	request.GeneralMessages[0].InfoMessageVersion = 1
	request.GeneralMessages[0].InfoChannelRef = "Chan"

	xml, err := request.BuildXML()
	if err != nil {
		t.Fatal(err)
	}

	xmlResponse, _ := NewXMLGeneralMessageResponseFromContent([]byte(xml))

	checkGeneralMessagesEquivalence(response, xmlResponse, t)
}
