package siri

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"bitbucket.org/enroute-mobi/ara/logger"
)

type SIRIStopPointsDiscoveryResponse struct {
	Status            bool
	ResponseTimestamp time.Time

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

func (response *SIRIStopPointsDiscoveryResponse) BuildXML(envelopeType ...string) (string, error) {
	var buffer bytes.Buffer
	var envType string
	var templateName string

	if len(envelopeType) != 0 && envelopeType[0] != "soap" {
		envType = "_" + envelopeType[0]
	}

	templateName = fmt.Sprintf("stop_points_discovery_response%s.template", envType)

	if err := templates.ExecuteTemplate(&buffer, templateName, response); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}
