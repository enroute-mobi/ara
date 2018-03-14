package siri

import (
	"bytes"
	"strings"
	"text/template"
	"time"
)

type SIRILinesDiscoveryResponse struct {
	Status            bool
	ResponseTimestamp time.Time

	AnnotatedLines []*SIRIAnnotatedLine
}

type SIRIAnnotatedLine struct {
	LineRef   string
	LineName  string
	Monitored bool
}

type SIRIAnnotatedLineByLineRef []*SIRIAnnotatedLine

func (a SIRIAnnotatedLineByLineRef) Len() int      { return len(a) }
func (a SIRIAnnotatedLineByLineRef) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a SIRIAnnotatedLineByLineRef) Less(i, j int) bool {
	return strings.Compare(a[i].LineRef, a[j].LineRef) < 0
}

const linesDiscoveryResponseTemplate = `<sw:LinesDiscoveryResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
	<Answer version="2.0">
		<siri:ResponseTimestamp>{{ .ResponseTimestamp.Format "2006-01-02T15:04:05.000Z07:00" }}</siri:ResponseTimestamp>
		<siri:Status>{{ .Status }}</siri:Status>{{ range .AnnotatedLines }}
		<siri:AnnotatedLineRef>
			<siri:LineRef>{{ .LineRef }}</siri:LineRef>
			<siri:LineName>{{ .LineName }}</siri:LineName>
			<siri:Monitored>{{ .Monitored }}</siri:Monitored>
		</siri:AnnotatedLineRef>{{ end }}
	</Answer>
	<AnswerExtension/>
</sw:LinesDiscoveryResponse>`

func (response *SIRILinesDiscoveryResponse) BuildXML() (string, error) {
	var buffer bytes.Buffer
	var siriResponse = template.Must(template.New("siriResponse").Parse(linesDiscoveryResponseTemplate))
	if err := siriResponse.Execute(&buffer, response); err != nil {
		return "", err
	}
	return buffer.String(), nil
}
