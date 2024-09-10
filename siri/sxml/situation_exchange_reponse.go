package sxml

import (
	"fmt"
	"strings"
	"time"

	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLSituationExchangeResponse struct {
	ResponseXMLStructureWithStatus

	deliveries []*XMLSituationExchangeDelivery
}

type XMLSituationExchangeDelivery struct {
	SubscriptionDeliveryXMLStructure

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
	reality        string
	participantRef string
	summaries      map[string]string
	descriptions   map[string]string

	affects               []*XMLAffect
	consequences          []*XMLConsequence
	publishToWebAction    *XMLPublishToWebAction
	publishToMobileAction *XMLPublishToMobileAction
}

type XMLActionData struct {
	XMLStructure

	name       string
	actionType string
	value      string
	prompt     map[string]string
}

type XMLCommonPublishingAction struct {
	XMLActionData

	actionStatus       string
	descriptions       map[string]string
	publicationWindows []*XMLPeriod
}
type XMLPublishToWebAction struct {
	XMLCommonPublishingAction

	incident       *bool
	homepage       *bool
	ticker         *bool
	socialNetworks []string
}

type XMLPublishToMobileAction struct {
	XMLCommonPublishingAction

	publicationWindows []*XMLPeriod
	incident           *bool
	homepage           *bool
}

func NewXMLActionData(node XMLNode) *XMLActionData {
	xmlActionData := &XMLActionData{}
	xmlActionData.node = node
	return xmlActionData
}

func NewXMLPublishToWebAction(node XMLNode) *XMLPublishToWebAction {
	xmlPublishToWebAction := &XMLPublishToWebAction{}
	xmlPublishToWebAction.node = node
	return xmlPublishToWebAction
}

func NewXMLPublishToMobileAction(node XMLNode) *XMLPublishToMobileAction {
	xmlPublishToMobileAction := &XMLPublishToMobileAction{}
	xmlPublishToMobileAction.node = node
	return xmlPublishToMobileAction
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
	condition      string
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
		nodes := response.findNodes(siri_attributes.SituationExchangeDelivery)
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
		nodes := delivery.findNodes(siri_attributes.PtSituationElement)
		for _, node := range nodes {
			situations = append(situations, NewXMLPtSituationElement(node))
		}
		delivery.situations = situations
	}
	return delivery.situations
}

func FindTranslations(nodes []XMLNode) map[string]string {
	translations := make(map[string]string)
	for _, node := range nodes {
		translations[node.NativeNode().Attr("lang")] = node.NativeNode().Content()
	}
	return translations
}

func (s *XMLPtSituationElement) Summaries() map[string]string {
	if s.summaries == nil {
		translations := FindTranslations(s.findNodes(siri_attributes.Summary))
		if translations != nil {
			s.summaries = translations
		}
	}
	return s.summaries
}

func (s *XMLPtSituationElement) Descriptions() map[string]string {
	if s.descriptions == nil {
		translations := FindTranslations(s.findDirectChildrenNodes((siri_attributes.Description)))
		if translations != nil {
			s.descriptions = translations
		}
	}
	return s.descriptions
}

func (s *XMLPtSituationElement) AlertCause() string {
	if s.alertCause == "" {
		s.alertCause = s.findStringChildContent(siri_attributes.AlertCause)
	}
	return s.alertCause
}

func (response *XMLSituationExchangeResponse) ErrorString() string {
	return fmt.Sprintf("%v: %v", response.errorType(), response.ErrorText())
}

func (response *XMLSituationExchangeResponse) errorType() string {
	if response.ErrorType() == siri_attributes.OtherError {
		return fmt.Sprintf("%v %v", response.ErrorType(), response.ErrorNumber())
	}
	return response.ErrorType()
}

func (visit *XMLPtSituationElement) RecordedAtTime() time.Time {
	if visit.recordedAtTime.IsZero() {
		visit.recordedAtTime = visit.VersionedAtTime()
		if visit.recordedAtTime.IsZero() {
			visit.recordedAtTime = visit.findTimeChildContent(siri_attributes.CreationTime)
		}
	}
	return visit.recordedAtTime
}

func (visit *XMLPtSituationElement) VersionedAtTime() time.Time {
	if visit.versionedAtTime.IsZero() {
		visit.versionedAtTime = visit.findTimeChildContent(siri_attributes.VersionedAtTime)
	}
	return visit.versionedAtTime
}

func (visit *XMLPtSituationElement) SituationNumber() string {
	if visit.situationNumber == "" {
		visit.situationNumber = visit.findStringChildContent(siri_attributes.SituationNumber)
	}
	return visit.situationNumber
}

func (visit *XMLPtSituationElement) Version() int {
	if !visit.version.Defined {
		visit.version.SetValueWithDefault(visit.findIntChildContent(siri_attributes.Version), 1)
	}
	return visit.version.Value
}

func (visit *XMLPtSituationElement) Keywords() []string {
	if len(visit.keywords) == 0 {
		keywords := strings.Split(visit.findStringChildContent(siri_attributes.Keywords), " ")
		visit.keywords = keywords
	}

	return visit.keywords
}

func (visit *XMLPtSituationElement) ReportType() string {
	if visit.reportType == "" {
		visit.reportType = visit.findStringChildContent(siri_attributes.ReportType)
	}
	return visit.reportType
}

func (visit *XMLPtSituationElement) PublicationWindows() []*XMLPeriod {
	if visit.publicationWindows == nil {
		publicationWindows := []*XMLPeriod{}
		nodes := visit.findNodes(siri_attributes.PublicationWindow)
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
		nodes := visit.findNodes(siri_attributes.ValidityPeriod)
		for _, node := range nodes {
			validityPeriods = append(validityPeriods, NewXMLPeriod(node))
		}
		visit.validityPeriods = validityPeriods
	}
	return visit.validityPeriods
}

func (v *XMLPeriod) StartTime() time.Time {
	if v.startTime.IsZero() {
		v.startTime = v.findTimeChildContent(siri_attributes.StartTime)
	}
	return v.startTime
}

func (v *XMLPeriod) EndTime() time.Time {
	if v.endTime.IsZero() {
		v.endTime = v.findTimeChildContent(siri_attributes.EndTime)
	}
	return v.endTime
}

func (visit *XMLPtSituationElement) Severity() string {
	if visit.severity == "" {
		visit.severity = visit.findStringChildContent(siri_attributes.Severity)
	}
	return visit.severity
}

func (visit *XMLPtSituationElement) Reality() string {
	if visit.reality == "" {
		visit.reality = visit.findStringChildContent(siri_attributes.Reality)
	}
	return visit.reality
}

func (visit *XMLPtSituationElement) Progress() string {
	if visit.progress == "" {
		visit.progress = visit.findStringChildContent(siri_attributes.Progress)
	}
	return visit.progress
}

func (visit *XMLPtSituationElement) ParticipantRef() string {
	if visit.participantRef == "" {
		visit.participantRef = visit.findStringChildContent(siri_attributes.ParticipantRef)
	}
	return visit.participantRef
}

func (visit *XMLPtSituationElement) Consequences() []*XMLConsequence {
	if visit.consequences == nil {
		consequences := []*XMLConsequence{}
		nodes := visit.findNodes(siri_attributes.Consequences)
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
		nodes := consequence.findNodes(siri_attributes.Period)
		for _, node := range nodes {
			periods = append(periods, NewXMLPeriod(node))
		}
		consequence.periods = periods
	}
	return consequence.periods
}

func (consequence *XMLConsequence) Condition() string {
	if consequence.condition == "" {
		consequence.condition = consequence.findStringChildContent(siri_attributes.Condition)
	}
	return consequence.condition
}

func (consequence *XMLConsequence) Severity() string {
	if consequence.severity == "" {
		consequence.severity = consequence.findStringChildContent(siri_attributes.Severity)
	}
	return consequence.severity
}

func (c *XMLConsequence) Affects() []*XMLAffect {
	if c.affects == nil {
		affects := []*XMLAffect{}
		nodes := c.findNodes(siri_attributes.Affects)
		for _, node := range nodes {
			affects = append(affects, NewXMLAffect(node))
		}
		c.affects = affects
	}
	return c.affects
}

func (c *XMLConsequence) HasBlocking() bool {
	node := c.findNode(siri_attributes.Blocking)
	if node != nil {
		c.hasBlocking = true
	}
	return c.hasBlocking
}

func (c *XMLConsequence) JourneyPlanner() bool {
	if !c.journeyPlanner.Defined {
		node := c.findNode(siri_attributes.Blocking)
		if node != nil {
			c.journeyPlanner.SetValue(c.findBoolChildContent(siri_attributes.JourneyPlanner))
		}
	}
	return c.journeyPlanner.Value
}

func (c *XMLConsequence) RealTime() bool {
	if !c.realTime.Defined {
		node := c.findNode(siri_attributes.Blocking)
		if node != nil {
			c.journeyPlanner.SetValue(c.findBoolChildContent(siri_attributes.JourneyPlanner))
		}
		c.realTime.SetValue(c.findBoolChildContent(siri_attributes.RealTime))
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
		nodes := a.findNodes(siri_attributes.AffectedNetwork)
		for _, affectedNetwork := range nodes {
			a.affectedNetworks = append(a.affectedNetworks, NewXMLAffectedNetwork(affectedNetwork))
		}
	}
	return a.affectedNetworks
}

func (an *XMLAffectedNetwork) LineRefs() []string {
	if len(an.lineRefs) == 0 {
		nodes := an.findNodes(siri_attributes.LineRef)
		for _, lineRef := range nodes {
			an.lineRefs = append(an.lineRefs, strings.TrimSpace(lineRef.NativeNode().Content()))
		}
	}
	return an.lineRefs
}

func (an *XMLAffectedNetwork) AffectedRoutes() []*XMLAffectedRoute {
	if len(an.affectedRoutes) == 0 {
		nodes := an.findNodes(siri_attributes.AffectedRoute)
		for _, affectedRoute := range nodes {
			an.affectedRoutes = append(an.affectedRoutes, NewXMLAffectedRoute(affectedRoute))
		}
	}
	return an.affectedRoutes
}

func (ar *XMLAffectedRoute) RouteRef() string {
	if ar.routeRef == "" {
		ar.routeRef = ar.findStringChildContent(siri_attributes.RouteRef)
	}
	return ar.routeRef
}
func (ar *XMLAffectedRoute) AffectedStopPoints() []*XMLAffectedStopPoint {
	if len(ar.affectedStopPoints) == 0 {
		stopPointsNodes := ar.findDirectChildrenNodes("StopPoints")
		if stopPointsNodes != nil {
			xmlStopPoints := NewXMLAffectedStopPoint(stopPointsNodes[0])
			nodes := xmlStopPoints.findNodes(siri_attributes.AffectedStopPoint)
			for _, affectedStopPoint := range nodes {
				ar.affectedStopPoints = append(ar.affectedStopPoints, NewXMLAffectedStopPoint(affectedStopPoint))
			}
		}

	}
	return ar.affectedStopPoints
}

func (an *XMLAffectedNetwork) AffectedSections() []*XMLAffectedSection {
	if len(an.affectedSections) == 0 {
		nodes := an.findNodes(siri_attributes.AffectedSection)
		for _, section := range nodes {
			an.affectedSections = append(an.affectedSections, NewXMLAffectedSection(section))
		}
	}
	return an.affectedSections
}

func (s *XMLAffectedSection) FirstStop() string {
	if s.firstStop == "" {
		s.firstStop = s.findStringChildContent(siri_attributes.FirstStopPointRef)
	}
	return s.firstStop
}

func (s *XMLAffectedSection) LastStop() string {
	if s.lastStop == "" {
		s.lastStop = s.findStringChildContent(siri_attributes.LastStopPointRef)
	}
	return s.lastStop
}

func (an *XMLAffectedNetwork) AffectedDestinations() []string {
	if len(an.affectedDestinations) == 0 {
		nodes := an.findNodes(siri_attributes.StopPlaceRef)
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
			nodes := xmlStopPoints.findNodes(siri_attributes.AffectedStopPoint)
			for _, affectedStopPoint := range nodes {
				a.affectedStopPoints = append(a.affectedStopPoints, NewXMLAffectedStopPoint(affectedStopPoint))
			}
		}

	}
	return a.affectedStopPoints
}

func (asp *XMLAffectedStopPoint) StopPointRef() string {
	if asp.stopPointRef == "" {
		asp.stopPointRef = asp.findStringChildContent(siri_attributes.StopPointRef)
	}
	return asp.stopPointRef
}

func (asp *XMLAffectedStopPoint) LineRefs() []string {
	if len(asp.lineRefs) == 0 {
		nodes := asp.findNodes(siri_attributes.LineRef)
		for _, lineRef := range nodes {
			asp.lineRefs = append(asp.lineRefs, strings.TrimSpace(lineRef.NativeNode().Content()))
		}
	}
	return asp.lineRefs
}

func (visit *XMLPtSituationElement) PublishToWebAction() *XMLPublishToWebAction {
	if visit.publishToWebAction == nil {
		nodes := visit.findNodes("PublishToWebAction")
		if nodes != nil {
			wa := NewXMLPublishToWebAction(nodes[0])
			visit.publishToWebAction = wa
		}
	}
	return visit.publishToWebAction
}

func (c *XMLActionData) Name() string {
	if c.name == "" {
		name := c.findStringChildContent("Name")
		c.name = name
	}
	return c.name
}

func (c *XMLActionData) Type() string {
	if c.actionType == "" {
		actionType := c.findStringChildContent("Type")
		c.actionType = actionType
	}
	return c.actionType
}

func (c *XMLActionData) Value() string {
	if c.value == "" {
		value := c.findStringChildContent("Value")
		c.value = value
	}
	return c.value
}

func (c *XMLActionData) Prompt() map[string]string {
	if c.prompt == nil {
		translations := FindTranslations(c.findNodes("Prompt"))
		if translations != nil {
			c.prompt = translations
		}
	}
	return c.prompt
}

func (wa *XMLPublishToWebAction) ActionData() *XMLActionData {
	if nodes := wa.findNodes("ActionData"); nodes != nil {
		actionData := NewXMLActionData(nodes[0])
		if actionData != nil {
			return actionData
		}
	}
	return nil
}

func (c *XMLPublishToWebAction) ActionStatus() string {
	if c.actionStatus == "" {
		actionStatus := c.findStringChildContent("ActionStatus")
		c.actionStatus = actionStatus
	}
	return c.actionStatus
}

func (c *XMLPublishToWebAction) Descriptions() map[string]string {
	if c.descriptions == nil {
		translations := FindTranslations(c.findDirectChildrenNodes(siri_attributes.Description))
		if translations != nil {
			c.descriptions = translations
		}
	}
	return c.descriptions
}

func (c *XMLPublishToWebAction) PublicationWindows() []*XMLPeriod {
	if c.publicationWindows == nil {
		publicationWindows := []*XMLPeriod{}
		nodes := c.findNodes(siri_attributes.PublicationWindow)
		for _, node := range nodes {
			publicationWindows = append(publicationWindows, NewXMLPeriod(node))
		}
		c.publicationWindows = publicationWindows
	}
	return c.publicationWindows
}

func (wa *XMLPublishToWebAction) Incident() *bool {
	if wa.incident == nil {
		if wa.findNode("Incident") != nil {
			incident := wa.findBoolChildContent("Incident")
			wa.incident = &incident
		}
	}
	return wa.incident
}

func (wa *XMLPublishToWebAction) HomePage() *bool {
	if wa.homepage == nil {
		if wa.findNode("HomePage") != nil {
			homepage := wa.findBoolChildContent("HomePage")
			wa.homepage = &homepage
		}
	}
	return wa.homepage
}

func (wa *XMLPublishToWebAction) Ticker() *bool {
	if wa.ticker == nil {
		if wa.findNode("Ticker") != nil {
			ticker := wa.findBoolChildContent("Ticker")
			wa.ticker = &ticker
		}
	}
	return wa.ticker
}

func (wa *XMLPublishToWebAction) SocialNetworks() []string {
	if len(wa.socialNetworks) == 0 {
		socialNetworks := wa.findNodes("SocialNetwork")
		for _, network := range socialNetworks {
			wa.socialNetworks = append(wa.socialNetworks, strings.TrimSpace(network.NativeNode().Content()))
		}
	}
	return wa.socialNetworks
}

func (visit *XMLPtSituationElement) PublishToMobileAction() *XMLPublishToMobileAction {
	if visit.publishToMobileAction == nil {
		nodes := visit.findNodes("PublishToMobileAction")
		if nodes != nil {
			ma := NewXMLPublishToMobileAction(nodes[0])
			visit.publishToMobileAction = ma
		}
	}
	return visit.publishToMobileAction
}

func (c *XMLPublishToMobileAction) ActionStatus() string {
	if c.actionStatus == "" {
		actionStatus := c.findStringChildContent("ActionStatus")
		c.actionStatus = actionStatus
	}
	return c.actionStatus
}

func (ma *XMLPublishToMobileAction) ActionData() *XMLActionData {
	if nodes := ma.findNodes("ActionData"); nodes != nil {
		actionData := NewXMLActionData(nodes[0])
		if actionData != nil {
			return actionData
		}
	}
	return nil
}

func (c *XMLPublishToMobileAction) Descriptions() map[string]string {
	if c.descriptions == nil {
		translations := FindTranslations(c.findDirectChildrenNodes(siri_attributes.Description))
		if translations != nil {
			c.descriptions = translations
		}
	}
	return c.descriptions
}

func (c *XMLPublishToMobileAction) PublicationWindows() []*XMLPeriod {
	if c.publicationWindows == nil {
		publicationWindows := []*XMLPeriod{}
		nodes := c.findNodes(siri_attributes.PublicationWindow)
		for _, node := range nodes {
			publicationWindows = append(publicationWindows, NewXMLPeriod(node))
		}
		c.publicationWindows = publicationWindows
	}
	return c.publicationWindows
}

func (wa *XMLPublishToMobileAction) Incident() *bool {
	if wa.incident == nil {
		if wa.findNode("Incident") != nil {
			incident := wa.findBoolChildContent("Incident")
			wa.incident = &incident
		}
	}
	return wa.incident
}

func (wa *XMLPublishToMobileAction) HomePage() *bool {
	if wa.homepage == nil {
		if wa.findNode("HomePage") != nil {
			homepage := wa.findBoolChildContent("HomePage")
			wa.homepage = &homepage
		}
	}
	return wa.homepage
}
