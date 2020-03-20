package siri

import (
	"bytes"
	"strings"
	"time"

	"bitbucket.org/enroute-mobi/ara/logger"
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

func (response *SIRILinesDiscoveryResponse) BuildXML() (string, error) {
	var buffer bytes.Buffer
	if err := templates.ExecuteTemplate(&buffer, "lines_discovery_response.template", response); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}
