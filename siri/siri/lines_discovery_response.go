package siri

import (
	"bytes"
	"fmt"
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

func (response *SIRILinesDiscoveryResponse) BuildXML(envelopeType ...string) (string, error) {
	var buffer bytes.Buffer
	var envType string
	var templateName string

	if len(envelopeType) != 0 && envelopeType[0] != "soap" && envelopeType[0] != "" {
		envType = "_" + envelopeType[0]
	}

	templateName = fmt.Sprintf("lines_discovery_response%s.template", envType)

	if err := templates.ExecuteTemplate(&buffer, templateName, response); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}
