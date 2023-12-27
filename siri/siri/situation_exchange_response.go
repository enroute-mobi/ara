package siri

import (
	"bytes"
	"fmt"
	"time"

	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
)

type SIRISituationExchangeResponse struct {
	SIRISituationExchangeDelivery

	Address                   string
	ProducerRef               string
	ResponseMessageIdentifier string
}

type SIRISituationExchangeDelivery struct {
	RequestMessageRef string

	ResponseTimestamp time.Time

	Status      bool
	ErrorType   string
	ErrorNumber int
	ErrorText   string

	Situations []*SIRIPtSituationElement

	LineRefs       map[string]struct{} `json:"-"`
	MonitoringRefs map[string]struct{} `json:"-"`
}

type SIRIPtSituationElement struct {
	CreationTime    time.Time
	SituationNumber string
	Version         int
	VersionedAtTime time.Time
	ValidityPeriods []*model.TimeRange
	Keywords        string
	ReportType      model.ReportType
	ParticipantRef  string
	Summary         string
	Description     string

	AffectedLines      []*AffectedLine
	AffectedStopPoints []*AffectedStopPoint
}

type AffectedStopPoint struct {
	StopPointRef  string
	StopPointName string
}

type AffectedLine struct {
	LineRef      string
	Destinations []SIRIAffectedDestination
	Sections     []SIRIAffectedSection
	Routes       []SIRIAffectedRoute
}

type SIRIAffectedDestination struct {
	StopPlaceRef  string
	StopPlaceName string
}

type SIRIAffectedSection struct {
	FirstStopPointRef string
	LastStopPointRef  string
}

type SIRIAffectedRoute struct {
	RouteRef string
}

func (response *SIRISituationExchangeResponse) BuildXML(envelopeType ...string) (string, error) {
	var buffer bytes.Buffer
	var envType string
	var templateName string

	if len(envelopeType) != 0 && envelopeType[0] != "soap" && envelopeType[0] != "" {
		envType = "_" + envelopeType[0]
	}

	templateName = fmt.Sprintf("situation_exchange_response%s.template", envType)

	if err := templates.ExecuteTemplate(&buffer, templateName, response); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}

func (delivery *SIRISituationExchangeDelivery) BuildSituationExchangeDeliveryXML() (string, error) {
	var buffer bytes.Buffer
	if err := templates.ExecuteTemplate(&buffer, "situation_exchange_delivery.template", delivery); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}

func (message *SIRIPtSituationElement) BuildSituationExchangeXML() (string, error) {
	var buffer bytes.Buffer
	if err := templates.ExecuteTemplate(&buffer, "situation_exchange.template", message); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}
