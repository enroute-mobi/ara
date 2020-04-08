package siri

import (
	"bytes"
	"strings"
	"time"

	"bitbucket.org/enroute-mobi/edwig/logger"
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

func (response *SIRIStopPointsDiscoveryResponse) BuildXML() (string, error) {
	var buffer bytes.Buffer
	if err := templates.ExecuteTemplate(&buffer, "stop_discovery_response.template", response); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}
