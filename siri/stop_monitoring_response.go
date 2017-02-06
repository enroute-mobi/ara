package siri

import (
	"bytes"
	"text/template"
	"time"

	"github.com/af83/edwig/logger"
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLStopMonitoringResponse struct {
	ResponseXMLStructure

	monitoredStopVisits []*XMLMonitoredStopVisit
}

type XMLMonitoredStopVisit struct {
	XMLStructure

	// parent *XMLStopMonitoringResponse

	itemIdentifier         string
	stopPointRef           string
	stopPointName          string
	datedVehicleJourneyRef string
	lineRef                string
	publishedLineName      string
	departureStatus        string
	arrivalStatus          string

	order int

	aimedArrivalTime    time.Time
	expectedArrivalTime time.Time
	actualArrivalTime   time.Time

	aimedDepartureTime    time.Time
	expectedDepartureTime time.Time
	actualDepartureTime   time.Time
}

type SIRIStopMonitoringResponse struct {
	Address                   string
	ProducerRef               string
	RequestMessageRef         string
	ResponseMessageIdentifier string
	Status                    bool
	// ErrorType                 string
	// ErrorNumber               int
	// ErrorText                 string
	ResponseTimestamp time.Time

	MonitoredStopVisits []*SIRIMonitoredStopVisit
}

type SIRIMonitoredStopVisit struct{}

const stopMonitoringResponseTemplate = `<ns8:GetStopMonitoringResponse xmlns:ns3="http://www.siri.org.uk/siri"
                               xmlns:ns4="http://www.ifopt.org.uk/acsb"
                               xmlns:ns5="http://www.ifopt.org.uk/ifopt"
                               xmlns:ns6="http://datex2.eu/schema/2_0RC1/2_0"
                               xmlns:ns7="http://scma/siri"
                               xmlns:ns8="http://wsdl.siri.org.uk"
                               xmlns:ns9="http://wsdl.siri.org.uk/siri">
  <ServiceDeliveryInfo>
    <ns3:ResponseTimestamp>{{ .ResponseTimestamp.Format "2006-01-02T15:04:05.000Z07:00" }}</ns3:ResponseTimestamp>
    <ns3:ProducerRef>{{ .ProducerRef }}</ns3:ProducerRef>
    <ns3:Address>{{ .Address }}</ns3:Address>
    <ns3:ResponseMessageIdentifier>{{ .ResponseMessageIdentifier }}</ns3:ResponseMessageIdentifier>
    <ns3:RequestMessageRef>{{ .RequestMessageRef }}</ns3:RequestMessageRef>
  </ServiceDeliveryInfo>
  <Answer>
    <ns3:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
      <ns3:ResponseTimestamp>{{ .ResponseTimestamp.Format "2006-01-02T15:04:05.000Z07:00" }}</ns3:ResponseTimestamp>
      <ns3:RequestMessageRef>{{ .RequestMessageRef }}</ns3:RequestMessageRef>
      <ns3:Status>{{ .Status }}</ns3:Status>{{ range .MonitoredStopVisits }}
      <ns3:MonitoredStopVisit>
      </ns3:MonitoredStopVisit>{{ end }}
    </ns3:StopMonitoringDelivery>
  </Answer>
  <AnswerExtension />
</ns8:GetStopMonitoringResponse>`

func NewXMLStopMonitoringResponse(node xml.Node) *XMLStopMonitoringResponse {
	xmlStopMonitoringResponse := &XMLStopMonitoringResponse{}
	xmlStopMonitoringResponse.node = NewXMLNode(node)
	return xmlStopMonitoringResponse
}

func NewXMLStopMonitoringResponseFromContent(content []byte) (*XMLStopMonitoringResponse, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	response := NewXMLStopMonitoringResponse(doc.Root().XmlNode)
	return response, nil
}

func NewSIRIStopMonitoringResponse(
	address string,
	producerRef string,
	requestMessageRef string,
	responseMessageIdentifier string,
	status bool,
	// errorType string,
	// errorNumber int,
	// errorText string,
	responseTimestamp time.Time) *SIRIStopMonitoringResponse {
	return &SIRIStopMonitoringResponse{
		Address:                   address,
		ProducerRef:               producerRef,
		RequestMessageRef:         requestMessageRef,
		ResponseMessageIdentifier: responseMessageIdentifier,
		Status: status,
		// ErrorType:          errorType,
		// ErrorNumber:        errorNumber,
		// ErrorText:          errorText,
		ResponseTimestamp: responseTimestamp,
	}
}

func (response *XMLStopMonitoringResponse) XMLMonitoredStopVisits() []*XMLMonitoredStopVisit {
	if len(response.monitoredStopVisits) == 0 {
		nodes := response.findNodes("MonitoredStopVisit")
		if nodes == nil {
			return response.monitoredStopVisits
		}
		for _, stopVisitNode := range nodes {
			response.monitoredStopVisits = append(response.monitoredStopVisits, NewXMLMonitoredStopVisit(stopVisitNode))
		}
	}
	return response.monitoredStopVisits
}

func NewXMLMonitoredStopVisit(node XMLNode) *XMLMonitoredStopVisit {
	stopVisit := &XMLMonitoredStopVisit{}
	stopVisit.node = node
	return stopVisit
}

func (visit *XMLMonitoredStopVisit) ItemIdentifier() string {
	if visit.itemIdentifier == "" {
		visit.itemIdentifier = visit.findStringChildContent("ItemIdentifier")
	}
	return visit.itemIdentifier
}

func (visit *XMLMonitoredStopVisit) StopPointRef() string {
	if visit.stopPointRef == "" {
		visit.stopPointRef = visit.findStringChildContent("StopPointRef")
	}
	return visit.stopPointRef
}

func (visit *XMLMonitoredStopVisit) StopPointName() string {
	if visit.stopPointName == "" {
		visit.stopPointName = visit.findStringChildContent("StopPointName")
	}
	return visit.stopPointName
}

func (visit *XMLMonitoredStopVisit) DatedVehicleJourneyRef() string {
	if visit.datedVehicleJourneyRef == "" {
		visit.datedVehicleJourneyRef = visit.findStringChildContent("DatedVehicleJourneyRef")
	}
	return visit.datedVehicleJourneyRef
}

func (visit *XMLMonitoredStopVisit) LineRef() string {
	if visit.lineRef == "" {
		visit.lineRef = visit.findStringChildContent("LineRef")
	}
	return visit.lineRef
}

func (visit *XMLMonitoredStopVisit) PublishedLineName() string {
	if visit.publishedLineName == "" {
		visit.publishedLineName = visit.findStringChildContent("PublishedLineName")
	}
	return visit.publishedLineName
}

func (visit *XMLMonitoredStopVisit) DepartureStatus() string {
	if visit.departureStatus == "" {
		visit.departureStatus = visit.findStringChildContent("DepartureStatus")
	}
	return visit.departureStatus
}

func (visit *XMLMonitoredStopVisit) ArrivalStatus() string {
	if visit.arrivalStatus == "" {
		visit.arrivalStatus = visit.findStringChildContent("ArrivalStatus")
	}
	return visit.arrivalStatus
}

func (visit *XMLMonitoredStopVisit) Order() int {
	if visit.order == 0 {
		visit.order = visit.findIntChildContent("Order")
	}
	return visit.order
}

func (visit *XMLMonitoredStopVisit) AimedArrivalTime() time.Time {
	if visit.aimedArrivalTime.IsZero() {
		visit.aimedArrivalTime = visit.findTimeChildContent("AimedArrivalTime")
	}
	return visit.aimedArrivalTime
}

func (visit *XMLMonitoredStopVisit) ExpectedArrivalTime() time.Time {
	if visit.expectedArrivalTime.IsZero() {
		visit.expectedArrivalTime = visit.findTimeChildContent("ExpectedArrivalTime")
	}
	return visit.expectedArrivalTime
}

func (visit *XMLMonitoredStopVisit) ActualArrivalTime() time.Time {
	if visit.actualArrivalTime.IsZero() {
		visit.actualArrivalTime = visit.findTimeChildContent("ActualArrivalTime")
	}
	return visit.actualArrivalTime
}

func (visit *XMLMonitoredStopVisit) AimedDepartureTime() time.Time {
	if visit.aimedDepartureTime.IsZero() {
		visit.aimedDepartureTime = visit.findTimeChildContent("AimedDepartureTime")
	}
	return visit.aimedDepartureTime
}

func (visit *XMLMonitoredStopVisit) ExpectedDepartureTime() time.Time {
	if visit.expectedDepartureTime.IsZero() {
		visit.expectedDepartureTime = visit.findTimeChildContent("ExpectedDepartureTime")
	}
	return visit.expectedDepartureTime
}

func (visit *XMLMonitoredStopVisit) ActualDepartureTime() time.Time {
	if visit.actualDepartureTime.IsZero() {
		visit.actualDepartureTime = visit.findTimeChildContent("ActualDepartureTime")
	}
	return visit.actualDepartureTime
}

func (response *SIRIStopMonitoringResponse) BuildXML() string {
	var buffer bytes.Buffer
	var siriResponse = template.Must(template.New("siriResponse").Parse(stopMonitoringResponseTemplate))
	if err := siriResponse.Execute(&buffer, response); err != nil {
		logger.Log.Panicf("Error while using response template: %v", err)
	}
	return buffer.String()
}
