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

	keywords           []string
	reportType         string
	alertCause         string
	recordedAtTime     time.Time
	versionedAtTime    time.Time
	validityPeriods    []*XMLPeriod
	publicationWindows []*XMLPeriod

	progress       string
	severity       string
	participantRef string
	summary        string
	description    string

	affects      []*XMLAffect
	consequences []*XMLConsequence
}

type XMLPeriod struct {
	XMLStructure

	startTime time.Time
	endTime   time.Time
}

type XMLAffect struct {
	XMLStructure

	affectedNetworks   []*XMLAffectedNetwork
	affectedStopPoints []*XMLAffectedStopPoint
}

type XMLAffectedRoute struct {
	XMLStructure

	routeRef           string
	affectedStopPoints []*XMLAffectedStopPoint
}

func NewXMLAffectedRoute(node XMLNode) *XMLAffectedRoute {
	xmlAffectedRoute := &XMLAffectedRoute{}
	xmlAffectedRoute.node = node
	return xmlAffectedRoute
}

type XMLAffectedNetwork struct {
	XMLStructure

	lineRefs             []string
	affectedSections     []*XMLAffectedSection
	affectedDestinations []string
	affectedRoutes       []*XMLAffectedRoute
}

func NewXMLAffectedNetwork(node XMLNode) *XMLAffectedNetwork {
	xmlAffectedNetwork := &XMLAffectedNetwork{}
	xmlAffectedNetwork.node = node
	return xmlAffectedNetwork
}

type XMLAffectedStopPoint struct {
	XMLStructure

	stopPointRef string
	lineRefs     []string
}

func NewXMLAffectedStopPoint(node XMLNode) *XMLAffectedStopPoint {
	xmlAffectedStopPoint := &XMLAffectedStopPoint{}
	xmlAffectedStopPoint.node = node
	return xmlAffectedStopPoint
}

type XMLAffectedSection struct {
	XMLStructure

	firstStop string
	lastStop  string
}

type XMLConsequence struct {
	XMLStructure

	periods        []*XMLPeriod
	severity       string
	affects        []*XMLAffect
	hasBlocking    bool
	journeyPlanner Bool
	realTime       Bool
}

func NewXMLConsequence(node XMLNode) *XMLConsequence {
	xmlConsequence := &XMLConsequence{}
	xmlConsequence.node = node
	return xmlConsequence
}

func NewXMLAffect(node XMLNode) *XMLAffect {
	xmlAffect := &XMLAffect{}
	xmlAffect.node = node
	return xmlAffect
}

func NewXMLPeriod(node XMLNode) *XMLPeriod {
	xmlPeriod := &XMLPeriod{}
	xmlPeriod.node = node
	return xmlPeriod
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

func (s *XMLPtSituationElement) AlertCause() string {
	if s.alertCause == "" {
		s.alertCause = s.findStringChildContent("AlertCause")
	}
	return s.alertCause
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

func (visit *XMLPtSituationElement) VersionedAtime() time.Time {
	if visit.versionedAtTime.IsZero() {
		visit.versionedAtTime = visit.findTimeChildContent("VersionedAtTime")
	}
	return visit.versionedAtTime
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

func (visit *XMLPtSituationElement) PublicationWindows() []*XMLPeriod {
	if visit.publicationWindows == nil {
		publicationWindows := []*XMLPeriod{}
		nodes := visit.findNodes("PublicationWindow")
		for _, node := range nodes {
			publicationWindows = append(publicationWindows, NewXMLPeriod(node))
		}
		visit.publicationWindows = publicationWindows
	}
	return visit.publicationWindows
}

func (visit *XMLPtSituationElement) ValidityPeriods() []*XMLPeriod {
	if visit.validityPeriods == nil {
		validityPeriods := []*XMLPeriod{}
		nodes := visit.findNodes("ValidityPeriod")
		for _, node := range nodes {
			validityPeriods = append(validityPeriods, NewXMLPeriod(node))
		}
		visit.validityPeriods = validityPeriods
	}
	return visit.validityPeriods
}

func (v *XMLPeriod) StartTime() time.Time {
	if v.startTime.IsZero() {
		v.startTime = v.findTimeChildContent("StartTime")
	}
	return v.startTime
}

func (v *XMLPeriod) EndTime() time.Time {
	if v.endTime.IsZero() {
		v.endTime = v.findTimeChildContent("EndTime")
	}
	return v.endTime
}

func (visit *XMLPtSituationElement) Severity() string {
	if visit.severity == "" {
		visit.severity = visit.findStringChildContent("Severity")
	}
	return visit.severity
}

func (visit *XMLPtSituationElement) Progress() string {
	if visit.progress == "" {
		visit.progress = visit.findStringChildContent("Progress")
	}
	return visit.progress
}

func (visit *XMLPtSituationElement) ParticipantRef() string {
	if visit.participantRef == "" {
		visit.participantRef = visit.findStringChildContent("ParticipantRef")
	}
	return visit.participantRef
}

func (visit *XMLPtSituationElement) Consequences() []*XMLConsequence {
	if visit.consequences == nil {
		consequences := []*XMLConsequence{}
		nodes := visit.findNodes("Consequences")
		for _, node := range nodes {
			consequences = append(consequences, NewXMLConsequence(node))
		}
		visit.consequences = consequences
	}
	return visit.consequences
}

func (consequence *XMLConsequence) Periods() []*XMLPeriod {
	if consequence.periods == nil {
		periods := []*XMLPeriod{}
		nodes := consequence.findNodes("Period")
		for _, node := range nodes {
			periods = append(periods, NewXMLPeriod(node))
		}
		consequence.periods = periods
	}
	return consequence.periods
}

func (consequence *XMLConsequence) Severity() string {
	if consequence.severity == "" {
		consequence.severity = consequence.findStringChildContent("Severity")
	}
	return consequence.severity
}

func (c *XMLConsequence) Affects() []*XMLAffect {
	if c.affects == nil {
		affects := []*XMLAffect{}
		nodes := c.findNodes("Affects")
		for _, node := range nodes {
			affects = append(affects, NewXMLAffect(node))
		}
		c.affects = affects
	}
	return c.affects
}

func (c *XMLConsequence) HasBlocking() bool {
	node := c.findNode("Blocking")
	if node != nil {
		c.hasBlocking = true
	}
	return c.hasBlocking
}

func (c *XMLConsequence) JourneyPlanner() bool {
	if !c.journeyPlanner.Defined {
		node := c.findNode("Blocking")
		if node != nil {
			c.journeyPlanner.SetValue(c.findBoolChildContent("JourneyPlanner"))
		}
	}
	return c.journeyPlanner.Value
}

func (c *XMLConsequence) RealTime() bool {
	if !c.realTime.Defined {
		node := c.findNode("Blocking")
		if node != nil {
			c.journeyPlanner.SetValue(c.findBoolChildContent("JourneyPlanner"))
		}
		c.realTime.SetValue(c.findBoolChildContent("RealTime"))
	}
	return c.realTime.Value
}

func (visit *XMLPtSituationElement) Affects() []*XMLAffect {
	if visit.affects == nil {
		affects := []*XMLAffect{}
		nodes := visit.findDirectChildrenNodes("Affects")
		for _, node := range nodes {
			affects = append(affects, NewXMLAffect(node))
		}
		visit.affects = affects
	}
	return visit.affects
}

func (a *XMLAffect) AffectedNetworks() []*XMLAffectedNetwork {
	if len(a.affectedNetworks) == 0 {
		nodes := a.findNodes("AffectedNetwork")
		for _, affectedNetwork := range nodes {
			a.affectedNetworks = append(a.affectedNetworks, NewXMLAffectedNetwork(affectedNetwork))
		}
	}
	return a.affectedNetworks
}

func (an *XMLAffectedNetwork) LineRefs() []string {
	if len(an.lineRefs) == 0 {
		nodes := an.findNodes("LineRef")
		for _, lineRef := range nodes {
			an.lineRefs = append(an.lineRefs, strings.TrimSpace(lineRef.NativeNode().Content()))
		}
	}
	return an.lineRefs
}

func (an *XMLAffectedNetwork) AffectedRoutes() []*XMLAffectedRoute {
	if len(an.affectedRoutes) == 0 {
		nodes := an.findNodes("AffectedRoute")
		for _, affectedRoute := range nodes {
			an.affectedRoutes = append(an.affectedRoutes, NewXMLAffectedRoute(affectedRoute))
		}
	}
	return an.affectedRoutes
}

func (ar *XMLAffectedRoute) RouteRef() string {
	if ar.routeRef == "" {
		ar.routeRef = ar.findStringChildContent("RouteRef")
	}
	return ar.routeRef
}
func (ar *XMLAffectedRoute) AffectedStopPoints() []*XMLAffectedStopPoint {
	if len(ar.affectedStopPoints) == 0 {
		stopPointsNodes := ar.findDirectChildrenNodes("StopPoints")
		if stopPointsNodes != nil {
			xmlStopPoints := NewXMLAffectedStopPoint(stopPointsNodes[0])
			nodes := xmlStopPoints.findNodes("AffectedStopPoint")
			for _, affectedStopPoint := range nodes {
				ar.affectedStopPoints = append(ar.affectedStopPoints, NewXMLAffectedStopPoint(affectedStopPoint))
			}
		}

	}
	return ar.affectedStopPoints
}

func (an *XMLAffectedNetwork) AffectedSections() []*XMLAffectedSection {
	if len(an.affectedSections) == 0 {
		nodes := an.findNodes("AffectedSection")
		for _, section := range nodes {
			an.affectedSections = append(an.affectedSections, NewXMLAffectedSection(section))
		}
	}
	return an.affectedSections
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

func (an *XMLAffectedNetwork) AffectedDestinations() []string {
	if len(an.affectedDestinations) == 0 {
		nodes := an.findNodes("StopPlaceRef")
		for _, routeRef := range nodes {
			an.affectedDestinations = append(an.affectedDestinations, strings.TrimSpace(routeRef.NativeNode().Content()))
		}
	}
	return an.affectedDestinations
}

func (a *XMLAffect) AffectedStopPoints() []*XMLAffectedStopPoint {
	if len(a.affectedStopPoints) == 0 {
		stopPointsNodes := a.findDirectChildrenNodes("StopPoints")
		if stopPointsNodes != nil {
			xmlStopPoints := NewXMLAffectedStopPoint(stopPointsNodes[0])
			nodes := xmlStopPoints.findNodes("AffectedStopPoint")
			for _, affectedStopPoint := range nodes {
				a.affectedStopPoints = append(a.affectedStopPoints, NewXMLAffectedStopPoint(affectedStopPoint))
			}
		}

	}
	return a.affectedStopPoints
}

func (asp *XMLAffectedStopPoint) StopPointRef() string {
	if asp.stopPointRef == "" {
		asp.stopPointRef = asp.findStringChildContent("StopPointRef")
	}
	return asp.stopPointRef
}

func (asp *XMLAffectedStopPoint) LineRefs() []string {
	if len(asp.lineRefs) == 0 {
		nodes := asp.findNodes("LineRef")
		for _, lineRef := range nodes {
			asp.lineRefs = append(asp.lineRefs, strings.TrimSpace(lineRef.NativeNode().Content()))
		}
	}
	return asp.lineRefs
}
