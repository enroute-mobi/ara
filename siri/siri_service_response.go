package siri

import (
	"bytes"
	"text/template"
	"time"
)

const siriServiceResponseTemplate = `<sw:GetSiriServiceResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
	<Answer>
		<siri:ResponseTimestamp>{{ .ResponseTimestamp.Format "2006-01-02T15:04:05.000Z07:00" }}</siri:ResponseTimestamp>
		<siri:ProducerRef>{{ .ProducerRef }}</siri:ProducerRef>
		<siri:ResponseMessageIdentifier>{{ .ResponseMessageIdentifier }}</siri:ResponseMessageIdentifier>
		<siri:RequestMessageRef>{{ .RequestMessageRef }}</siri:RequestMessageRef>
		<siri:Status>{{ .Status }}</siri:Status>{{ range .StopMonitoringDeliveries }}
		{{ .BuildStopMonitoringDeliveryXML }}{{ end }}{{ range .GeneralMessageDeliveries }}
		{{ .BuildGeneralMessageDeliveryXML }}{{ end }}{{ range .EstimatedTimetableDeliveries }}
		{{ .BuildEstimatedTimetableDeliveryXML }}{{ end }}
	</Answer>
	<AnswerExtension />
</sw:GetSiriServiceResponse>`

type SIRIServiceResponse struct {
	ProducerRef               string
	ResponseMessageIdentifier string
	RequestMessageRef         string
	Status                    bool
	// ErrorType                 string
	// ErrorNumber               int
	// ErrorText                 string

	ResponseTimestamp time.Time

	StopMonitoringDeliveries     []*SIRIStopMonitoringDelivery
	GeneralMessageDeliveries     []*SIRIGeneralMessageDelivery
	EstimatedTimetableDeliveries []*SIRIEstimatedTimetableDelivery
}

func (response *SIRIServiceResponse) BuildXML() (string, error) {
	var buffer bytes.Buffer
	var siriResponse = template.Must(template.New("siriResponse").Parse(siriServiceResponseTemplate))
	if err := siriResponse.Execute(&buffer, response); err != nil {
		return "", err
	}
	return buffer.String(), nil
}
