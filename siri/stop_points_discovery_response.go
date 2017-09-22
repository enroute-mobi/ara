package siri

import (
	"bytes"
	"strings"
	"text/template"
	"time"
)

type SIRIStopPointsDiscoveryResponse struct {
	Address                   string
	ProducerRef               string
	RequestMessageRef         string
	ResponseMessageIdentifier string
	Status                    bool
	ResponseTimestamp         time.Time

	AnnotatedStopPoints []*SIRIAnnotatedStopPoint
}

type SIRIAnnotatedStopPoint struct {
	StopPointRef string
	StopName     string
	Lines        []string
	Monitored    bool
	TimingPoint  bool
}

type SIRIAnnotatedStopPointByStopPointRef []*SIRIAnnotatedStopPoint

func (a SIRIAnnotatedStopPointByStopPointRef) Len() int      { return len(a) }
func (a SIRIAnnotatedStopPointByStopPointRef) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a SIRIAnnotatedStopPointByStopPointRef) Less(i, j int) bool {
	return strings.Compare(a[i].StopPointRef, a[j].StopPointRef) < 0
}

const stopDiscoveryResponseTemplate = `<sw:StopPointsDiscoveryResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
	<Answer version="2.0:FR-IDF-2.4">
		<siri:ResponseTimestamp>{{ .ResponseTimestamp.Format "2006-01-02T15:04:05.000Z07:00" }}</siri:ResponseTimestamp>{{ if .Address }}
		<siri:Address>{{ .Address }}</siri:Address>{{ end }}
		<siri:RequestMessageRef>{{ .RequestMessageRef }}</siri:RequestMessageRef>
		<siri:ResponseMessageIdentifier>{{ .ResponseMessageIdentifier }}</siri:ResponseMessageIdentifier>
		<siri:Status>{{ .Status }}</siri:Status>{{ range .AnnotatedStopPoints }}
		<siri:AnnotatedStopPointRef>
			<siri:StopPointRef>{{ .StopPointRef }}</siri:StopPointRef>
			<siri:Monitored>{{ .Monitored }}</siri:Monitored>
			<siri:StopName>{{ .StopName }}</siri:StopName>{{ if .Lines }}
			<siri:Lines>{{ range .Lines }}
				<siri:LineRef>{{ . }}</siri:LineRef>{{ end }}
			</siri:Lines>{{ end }}
		</siri:AnnotatedStopPointRef>{{ end }}
	</Answer>
	<AnswerExtension />
</sw:StopPointsDiscoveryResponse>`

func (response *SIRIStopPointsDiscoveryResponse) BuildXML() (string, error) {
	var buffer bytes.Buffer
	var siriResponse = template.Must(template.New("siriResponse").Parse(stopDiscoveryResponseTemplate))
	if err := siriResponse.Execute(&buffer, response); err != nil {
		return "", err
	}
	return buffer.String(), nil
}
