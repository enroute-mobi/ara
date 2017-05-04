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

const stopDiscoveryResponseTemplate = `
<ns8:StopPointsDiscoveryResponse xmlns:ns8="http://wsdl.siri.org.uk" xmlns:ns3="http://www.siri.org.uk/siri" xmlns:ns4="http://www.ifopt.org.uk/acsb" xmlns:ns5="http://www.ifopt.org.uk/ifopt" xmlns:ns6="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns7="http://scma/siri" xmlns:ns9="http://wsdl.siri.org.uk/siri">
   <Answer version="2.0">
      <ns3:ResponseTimestamp>{{ .ResponseTimestamp.Format "2006-01-02T15:04:05.000Z07:00" }}</ns3:ResponseTimestamp>
      <ns3:Address>{{.Address}}</ns3:Address>
      <ns3:ProducerRef>{{.ProducerRef}}</ns3:ProducerRef>
      <ns3:RequestMessageRef>{{.RequestMessageRef}}</ns3:RequestMessageRef>
      <ns3:ResponseMessageIdentifier>{{.ResponseMessageIdentifier}}</ns3:ResponseMessageIdentifier>
      <ns3:Status>{{.Status}}</ns3:Status>{{range .AnnotatedStopPoints}}
      <ns3:AnnotatedStopPointRef>
         <ns3:StopPointRef>{{.StopPointRef}}</ns3:StopPointRef>
         <ns3:StopName>{{.StopName}}</ns3:StopName>
      </ns3:AnnotatedStopPointRef>{{ end }}
   </Answer>
   <AnswerExtension />
</ns8:StopPointsDiscoveryResponse>`

func (response *SIRIStopPointsDiscoveryResponse) BuildXML() (string, error) {
	var buffer bytes.Buffer
	var siriResponse = template.Must(template.New("siriResponse").Parse(stopDiscoveryResponseTemplate))
	if err := siriResponse.Execute(&buffer, response); err != nil {
		return "", err
	}
	return buffer.String(), nil
}
