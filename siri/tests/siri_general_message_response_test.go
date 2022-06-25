package siri_tests

import (
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/siri/siri"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
)

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
	<AnswerExtension/>
</sw:GetGeneralMessageResponse>`

	response, _ := sxml.NewXMLGeneralMessageResponseFromContent([]byte(expectedXML))
	responseTimestamp := time.Date(2016, time.September, 21, 20, 14, 46, 0, time.UTC)

	request := &siri.SIRIGeneralMessageResponse{
		Address:                   "address",
		ProducerRef:               "producer",
		ResponseMessageIdentifier: "identifier",
	}

	lineSection := &siri.SIRILineSection{
		FirstStop: "NINOXE:StopPoint:SP:24:LOC",
		LastStop:  "NINOXE:StopPoint:SP:12:LOC",
		LineRef:   "NINOXE:Line::3:LOC",
	}

	gM := &siri.SIRIGeneralMessage{
		RecordedAtTime: time.Date(2016, time.September, 21, 20, 14, 46, 0, time.UTC),
		ValidUntilTime: time.Date(2016, time.September, 21, 20, 14, 46, 0, time.UTC),
		LineSections:   []*siri.SIRILineSection{lineSection},
	}

	request.Status = true
	request.ResponseTimestamp = responseTimestamp
	request.RequestMessageRef = "ref"

	request.GeneralMessages = []*siri.SIRIGeneralMessage{gM}
	request.GeneralMessages[0].Messages = append(request.GeneralMessages[0].Messages, &siri.SIRIMessage{Content: "Je suis un texte", Type: "Un Type"})
	request.GeneralMessages[0].InfoMessageVersion = 1
	request.GeneralMessages[0].InfoChannelRef = "Chan"
	request.GeneralMessages[0].FormatRef = "STIF-IDF"

	xml, err := request.BuildXML()
	if err != nil {
		t.Fatal(err)
	}

	xmlResponse, _ := sxml.NewXMLGeneralMessageResponseFromContent([]byte(xml))

	checkGeneralMessagesEquivalence(response, xmlResponse, t)
}
