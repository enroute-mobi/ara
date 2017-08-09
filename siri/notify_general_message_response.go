package siri

import (
	"bytes"
	"html/template"
)

type SIRINotifyGeneralMessage struct {
	SIRIGeneralMessageDelivery

	Address                   string
	ProducerRef               string
	ResponseMessageIdentifier string
	Status                    bool

	SubscriberRef          string
	SubscriptionIdentifier string
}

const generalMessageNotifyTemplate = `<ns8:NotifyGeneralMessage xmlns:ns3="http://www.siri.org.uk/siri"
															 xmlns:ns5="http://www.ifopt.org.uk/ifopt"
                               xmlns:ns4="http://www.ifopt.org.uk/acsb"
															 xmlns:ns6="http://datex2.eu/schema/2_0RC1/2_0"
															 xmlns:ns7="http://scma/siri"
															 xmlns:ns8="http://wsdl.siri.org.uk"
															 xmlns:ns9="http://wsdl.siri.org.uk/siri">
	<ServiceDeliveryInfo>
		<ns3:ResponseTimestamp>{{ .ResponseTimestamp.Format "2006-01-02T15:04:05.000Z07:00" }}</ns3:ResponseTimestamp>
		<ns3:ProducerRef>{{ .ProducerRef }}</ns3:ProducerRef>{{ if .Address }}
		<ns3:Address>{{ .Address }}</ns3:Address>{{ end }}
		<ns3:ResponseMessageIdentifier>{{ .ResponseMessageIdentifier }}</ns3:ResponseMessageIdentifier>
		<ns3:RequestMessageRef>{{ .RequestMessageRef }}</ns3:RequestMessageRef>
	</ServiceDeliveryInfo>
	<Notification xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
	  <ns3:GeneralMessageDelivery version="2.0:FR-IDF-2.4">
	    <ns3:ResponseTimestamp>{{ .ResponseTimestamp.Format "2006-01-02T15:04:05.000Z07:00" }}</ns3:ResponseTimestamp>
	    <ns5:RequestMessageRef>{{.RequestMessageRef}}</ns5:RequestMessageRef>
	    <ns5:SubscriberRef>{{.SubscriberRef}}</ns5:SubscriberRef>
	    <ns5:SubscriptionRef>{{.SubscriptionIdentifier}}</ns5:SubscriptionRef>
			<ns3:Status>{{ .Status }}</ns3:Status>{{range .GeneralMessages}}
	  	<ns3:GeneralMessage>
	  		<ns3:formatRef>{{ .FormatRef }}</ns3:formatRef>
	  		<ns3:RecordedAtTime>{{ .RecordedAtTime.Format "2006-01-02T15:04:05.000Z07:00" }}</ns3:RecordedAtTime>
	  		<ns3:ItemIdentifier>{{ .ItemIdentifier }}</ns3:ItemIdentifier>
	  		<ns3:InfoMessageIdentifier>{{ .InfoMessageIdentifier }}</ns3:InfoMessageIdentifier>
	  		<ns3:InfoMessageVersion>{{ .InfoMessageVersion }}</ns3:InfoMessageVersion>
	  		<ns3:InfoChannelRef>{{ .InfoChannelRef }}</ns3:InfoChannelRef>
	  		<ns3:ValidUntilTime>{{ .ValidUntilTime.Format "2006-01-02T15:04:05.000Z07:00" }}</ns3:ValidUntilTime>
	  		<ns3:Content xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
	  		xsi:type="ns9:IDFLineSectionStructure">{{range .Messages}}
	  			<Message>{{if .Type}}
	  				<MessageType>{{ .Type }}</MessageType>{{end}}{{if .Content }}
	  				<MessageText>{{ .Content }}</MessageText>{{end}}{{if .NumberOfLines }}
	  				<NumberOfLines>{{ .NumberOfLines }}</NumberOfLines>{{end}}{{if .NumberOfCharPerLine }}
	  				<NumberOfCharPerLine>{{ .NumberOfCharPerLine }}</NumberOfCharPerLine>{{end}}
	  			</Message>{{end}}{{ if or .FirstStop .LastStop .LineRef }}
	  			<LineSection>{{ if .FirstStop }}
	  				<FirstStop>{{ .FirstStop }}</FirstStop>{{end}}{{if .LastStop }}
	  			  <LastStop>{{ .LastStop }}</LastStop>{{end}}{{if .LineRef }}
	  			  <LineRef>{{ .LineRef }}</LineRef>{{end}}
	  			</LineSection>{{end}}
	  		</ns3:Content>
	  	</ns3:GeneralMessage>{{end}}
	   </ns3:GeneralMessageDelivery>
		</Notification>
 <NotifyExtension xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri"/>
</ns8:NotifyGeneralMessage>`

func (notify *SIRINotifyGeneralMessage) BuildXML() (string, error) {
	var buffer bytes.Buffer
	var notifyDelivery = template.Must(template.New("generalMessageNotifyTemplate").Parse(generalMessageNotifyTemplate))
	if err := notifyDelivery.Execute(&buffer, notify); err != nil {
		return "", err
	}
	return buffer.String(), nil
}
