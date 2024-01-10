package sxml

import (
	"fmt"
	"strings"
	"time"

	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLSituationExchangeResponse struct {
	ResponseXMLStructureWithStatus

	deliveries []*XMLSituationExchangeDelivery
}

type XMLSituationExchangeDelivery struct {
	XMLStructure

	situations []*XMLPtSituationElement
}

type XMLPtSituationElement struct {
	XMLStructure

	situationNumber string

	version Int

	keywords       []string
	reportType     string
	recordedAtTime time.Time

	validityPeriods []*XMLValidityPeriod

	summary     string
	description string

	affects []*XMLAffect
}

type XMLValidityPeriod struct {
	XMLStructure

	startTime time.Time
	endTime   time.Time
}

type XMLAffect struct {
	XMLStructure

	lineRefs             []string
	affectedRoutes       []string
	affectedSections     []*XMLAffectedSection
	affectedDestinations []string

	stopPoints []string
}

type XMLAffectedSection struct {
	XMLStructure

	firstStop string
	lastStop  string
}

func NewXMLAffect(node XMLNode) *XMLAffect {
	xmlAffect := &XMLAffect{}
	xmlAffect.node = node
	return xmlAffect
}

func NewXMLValidityPeriod(node XMLNode) *XMLValidityPeriod {
	xmlValidityPeriod := &XMLValidityPeriod{}
	xmlValidityPeriod.node = node
	return xmlValidityPeriod
}

func NewXMLAffectedSection(node XMLNode) *XMLAffectedSection {
	xmlAffectedSection := &XMLAffectedSection{}
	xmlAffectedSection.node = node
	return xmlAffectedSection
}

func NewXMLSituationExchangeResponse(node xml.Node) *XMLSituationExchangeResponse {
	xmlSituationExchangeResponse := &XMLSituationExchangeResponse{}
	xmlSituationExchangeResponse.node = NewSubXMLNode(node)
	return xmlSituationExchangeResponse
}

func NewXMLSituationExchangeResponseFromContent(content []byte) (*XMLSituationExchangeResponse, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	response := NewXMLSituationExchangeResponse(doc.Root().XmlNode)
	return response, nil
}

func NewXMLSituationExchangeDelivery(node XMLNode) *XMLSituationExchangeDelivery {
	delivery := &XMLSituationExchangeDelivery{}
	delivery.node = node
	return delivery
}

func NewXMLPtSituationElement(node XMLNode) *XMLPtSituationElement {
	situation := &XMLPtSituationElement{}
	situation.node = node
	return situation
}

func (response *XMLSituationExchangeResponse) SituationExchangeDeliveries() []*XMLSituationExchangeDelivery {
	if response.deliveries == nil {
		deliveries := []*XMLSituationExchangeDelivery{}
		nodes := response.findNodes("SituationExchangeDelivery")
		for _, node := range nodes {
			deliveries = append(deliveries, NewXMLSituationExchangeDelivery(node))
		}
		response.deliveries = deliveries
	}
	return response.deliveries
}

func (delivery *XMLSituationExchangeDelivery) Situations() []*XMLPtSituationElement {
	if delivery.situations == nil {
		situations := []*XMLPtSituationElement{}
		nodes := delivery.findNodes("PtSituationElement")
		for _, node := range nodes {
			situations = append(situations, NewXMLPtSituationElement(node))
		}
		delivery.situations = situations
	}
	return delivery.situations
}

func (s *XMLPtSituationElement) Summary() string {
	if s.summary == "" {
		s.summary = s.findStringChildContent("Summary")
	}
	return s.summary
}

func (s *XMLPtSituationElement) Description() string {
	if s.description == "" {
		s.description = s.findStringChildContent("Description")
	}
	return s.description
}

func (response *XMLSituationExchangeResponse) ErrorString() string {
	return fmt.Sprintf("%v: %v", response.errorType(), response.ErrorText())
}

func (response *XMLSituationExchangeResponse) errorType() string {
	if response.ErrorType() == "OtherError" {
		return fmt.Sprintf("%v %v", response.ErrorType(), response.ErrorNumber())
	}
	return response.ErrorType()
}

func (visit *XMLPtSituationElement) RecordedAtTime() time.Time {
	if visit.recordedAtTime.IsZero() {
		visit.recordedAtTime = visit.findTimeChildContent("CreationTime")
	}
	return visit.recordedAtTime
}

func (visit *XMLPtSituationElement) SituationNumber() string {
	if visit.situationNumber == "" {
		visit.situationNumber = visit.findStringChildContent("SituationNumber")
	}
	return visit.situationNumber
}

func (visit *XMLPtSituationElement) Version() int {
	if !visit.version.Defined {
		visit.version.SetValueWithDefault(visit.findIntChildContent("Version"), 1)
	}
	return visit.version.Value
}

func (visit *XMLPtSituationElement) Keywords() []string {
	if len(visit.keywords) == 0 {
		keywords := strings.Split(visit.findStringChildContent("Keywords"), " ")
		visit.keywords = keywords
	}

	return visit.keywords
}

func (visit *XMLPtSituationElement) ReportType() string {
	if visit.reportType == "" {
		visit.reportType = visit.findStringChildContent("ReportType")
	}
	return visit.reportType
}

func (visit *XMLPtSituationElement) ValidityPeriods() []*XMLValidityPeriod {
	if visit.validityPeriods == nil {
		validityPeriods := []*XMLValidityPeriod{}
		nodes := visit.findNodes("ValidityPeriod")
		for _, node := range nodes {
			validityPeriods = append(validityPeriods, NewXMLValidityPeriod(node))
		}
		visit.validityPeriods = validityPeriods
	}
	return visit.validityPeriods
}

func (v *XMLValidityPeriod) StartTime() time.Time {
	if v.startTime.IsZero() {
		v.startTime = v.findTimeChildContent("StartTime")
	}
	return v.startTime
}

func (v *XMLValidityPeriod) EndTime() time.Time {
	if v.endTime.IsZero() {
		v.endTime = v.findTimeChildContent("EndTime")
	}
	return v.endTime
}

func (visit *XMLPtSituationElement) Affects() []*XMLAffect {
	if visit.affects == nil {
		affects := []*XMLAffect{}
		nodes := visit.findNodes("Affects")
		for _, node := range nodes {
			affects = append(affects, NewXMLAffect(node))
		}
		visit.affects = affects
	}
	return visit.affects
}

func (a *XMLAffect) LineRefs() []string {
	if len(a.lineRefs) == 0 {
		nodes := a.findNodes("LineRef")
		for _, lineRef := range nodes {
			a.lineRefs = append(a.lineRefs, strings.TrimSpace(lineRef.NativeNode().Content()))
		}
	}
	return a.lineRefs
}

func (a *XMLAffect) AffectedRoutes() []string {
	if len(a.affectedRoutes) == 0 {
		nodes := a.findNodes("RouteRef")
		for _, routeRef := range nodes {
			a.affectedRoutes = append(a.affectedRoutes, strings.TrimSpace(routeRef.NativeNode().Content()))
		}
	}
	return a.affectedRoutes
}

func (a *XMLAffect) AffectedSections() []*XMLAffectedSection {
	if len(a.affectedSections) == 0 {
		nodes := a.findNodes("AffectedSection")
		for _, section := range nodes {
			a.affectedSections = append(a.affectedSections, NewXMLAffectedSection(section))
		}
	}
	return a.affectedSections
}

func (s *XMLAffectedSection) FirstStop() string {
	if s.firstStop == "" {
		s.firstStop = s.findStringChildContent("FirstStopPointRef")
	}
	return s.firstStop
}

func (s *XMLAffectedSection) LastStop() string {
	if s.lastStop == "" {
		s.lastStop = s.findStringChildContent("LastStopPointRef")
	}
	return s.lastStop
}

func (a *XMLAffect) AffectedDestinations() []string {
	if len(a.affectedDestinations) == 0 {
		nodes := a.findNodes("StopPlaceRef")
		for _, routeRef := range nodes {
			a.affectedDestinations = append(a.affectedDestinations, strings.TrimSpace(routeRef.NativeNode().Content()))
		}
	}
	return a.affectedDestinations
}

func (a *XMLAffect) StopPoints() []string {
	if len(a.stopPoints) == 0 {
		nodes := a.findNodes("StopPointRef")
		for _, stopPointRef := range nodes {
			a.stopPoints = append(a.stopPoints, strings.TrimSpace(stopPointRef.NativeNode().Content()))
		}
	}
	return a.stopPoints
}
