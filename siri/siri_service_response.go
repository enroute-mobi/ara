package siri

import (
	"bytes"
	"text/template"
	"time"
)

const siriServiceResponseTemplate = `<ns1:GetSiriServiceResponse xmlns:ns1="http://wsdl.siri.org.uk">
	<Answer xmlns:ns3="http://www.siri.org.uk/siri"
					xmlns:ns4="http://www.ifopt.org.uk/acsb"
					xmlns:ns5="http://www.ifopt.org.uk/ifopt"
					xmlns:ns6="http://datex2.eu/schema/2_0RC1/2_0"
					xmlns:ns7="http://scma/siri"
					xmlns:ns8="http://wsdl.siri.org.uk"
					xmlns:ns9="http://wsdl.siri.org.uk/siri">
		<ns3:ResponseTimestamp>{{ .ResponseTimestamp.Format "2006-01-02T15:04:05.000Z07:00" }}</ns3:ResponseTimestamp>
		<ns3:ProducerRef>{{ .ProducerRef }}</ns3:ProducerRef>
		<ns3:ResponseMessageIdentifier>{{ .ResponseMessageIdentifier }}</ns3:ResponseMessageIdentifier>
		<ns3:RequestMessageRef>{{ .RequestMessageRef }}</ns3:RequestMessageRef>
		<ns3:Status>{{ .Status }}</ns3:Status>{{ range .Deliveries }}
		{{ .BuildStopMonitoringDeliveryXML }}{{end}}
	</Answer>
</ns1:GetSiriServiceResponse>`

type SIRIServiceResponse struct {
	ProducerRef               string
	ResponseMessageIdentifier string
	RequestMessageRef         string
	Status                    bool
	// ErrorType                 string
	// ErrorNumber               int
	// ErrorText                 string

	ResponseTimestamp time.Time

	Deliveries []SIRIStopMonitoringDelivery
}

func (response *SIRIServiceResponse) BuildXML() (string, error) {
	var buffer bytes.Buffer
	var siriResponse = template.Must(template.New("siriResponse").Parse(siriServiceResponseTemplate))
	if err := siriResponse.Execute(&buffer, response); err != nil {
		return "", err
	}
	return buffer.String(), nil
}
