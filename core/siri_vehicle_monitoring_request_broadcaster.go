package core

import (
	"fmt"
	"sort"
	"strings"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/core/idgen"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/siri"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
	"bitbucket.org/enroute-mobi/ara/state"
)

type VehicleMonitoringRequestBroadcaster interface {
	RequestVehicles(*sxml.XMLGetVehicleMonitoring, *audit.BigQueryMessage) *siri.SIRIVehicleMonitoringResponse
}

type SIRIVehicleMonitoringRequestBroadcaster struct {
	state.Startable

	connector

	vjRemoteObjectidKinds      []string
	vehicleRemoteObjectidKinds []string
}

type SIRIVehicleMonitoringRequestBroadcasterFactory struct{}

func NewSIRIVehicleMonitoringRequestBroadcaster(partner *Partner) *SIRIVehicleMonitoringRequestBroadcaster {
	connector := &SIRIVehicleMonitoringRequestBroadcaster{}

	connector.partner = partner
	return connector
}

func (connector *SIRIVehicleMonitoringRequestBroadcaster) Start() {
	connector.vjRemoteObjectidKinds = connector.partner.VehicleJourneyRemoteObjectIDKindWithFallback(SIRI_VEHICLE_MONITORING_REQUEST_BROADCASTER)
	connector.vehicleRemoteObjectidKinds = connector.partner.VehicleRemoteObjectIDKindWithFallback(SIRI_VEHICLE_MONITORING_REQUEST_BROADCASTER)
	connector.remoteObjectidKind = connector.partner.RemoteObjectIDKind(SIRI_VEHICLE_MONITORING_REQUEST_BROADCASTER)
}

func (connector *SIRIVehicleMonitoringRequestBroadcaster) RequestVehicles(request *sxml.XMLGetVehicleMonitoring, message *audit.BigQueryMessage) (siriResponse *siri.SIRIVehicleMonitoringResponse) {
	lineRef := request.LineRef()

	messageIdentifier := request.MessageIdentifier()

	message.RequestIdentifier = messageIdentifier
	message.Lines = []string{lineRef}

	siriResponse = &siri.SIRIVehicleMonitoringResponse{
		ResponseTimestamp:         connector.Clock().Now(),
		ProducerRef:               connector.Partner().ProducerRef(),
		ResponseMessageIdentifier: connector.Partner().ResponseMessageIdentifierGenerator().NewMessageIdentifier(),
		RequestMessageRef:         messageIdentifier,
	}

	response := siri.SIRIVehicleMonitoringDelivery{
		Version:           "2.0:FR-IDF-2.4",
		ResponseTimestamp: connector.Clock().Now(),
		RequestMessageRef: messageIdentifier,
	}

	objectid := model.NewObjectID(connector.remoteObjectidKind, lineRef)
	line, ok := connector.partner.Model().Lines().FindByObjectId(objectid)
	if !ok {
		response.ErrorCondition = &siri.ErrorCondition{
			ErrorType: "InvalidDataReferencesError",
			ErrorText: fmt.Sprintf("Line %v not found", objectid.Value()),
		}
		message.Status = "Error"
		message.ErrorDetails = response.ErrorCondition.ErrorText
		siriResponse.SIRIVehicleMonitoringDelivery = response

		return
	}

	response.Status = true
	siriResponse.SIRIVehicleMonitoringDelivery = response

	var vehicleIds []string

	vs := connector.partner.Model().Vehicles().FindByLineId(line.Id())
	for i := range vs {
		vehicleId, ok := vs[i].ObjectIDWithFallback(connector.vehicleRemoteObjectidKinds)
		if !ok {
			continue
		}

		vj := vs[i].VehicleJourney()
		if vj == nil {
			continue
		}
		dvj, ok := connector.datedVehicleJourneyRef(vj)
		if !ok {
			continue
		}

		refs := vj.References.Copy()

		activity := &siri.SIRIVehicleActivity{
			RecordedAtTime:       vs[i].RecordedAtTime,
			ValidUntilTime:       vs[i].ValidUntilTime,
			VehicleMonitoringRef: vehicleId.Value(),
			ProgressBetweenStops: connector.handleProgressBetweenStops(vs[i]),
		}

		monitoredVehicleJourney := &siri.SIRIMonitoredVehicleJourney{
			LineRef:            lineRef,
			PublishedLineName:  line.Name,
			DirectionName:      vj.Attributes["DirectionName"],
			DirectionType:      vj.DirectionType,
			OriginName:         vj.OriginName,
			DestinationName:    vj.DestinationName,
			Monitored:          vj.Monitored,
			Bearing:            vs[i].Bearing,
			DriverRef:          vs[i].DriverRef,
			Occupancy:          vj.Occupancy,
			OriginRef:          connector.handleRef("OriginRef", vj.Origin, refs),
			DestinationRef:     connector.handleRef("DestinationRef", vj.Origin, refs),
			JourneyPatternRef:  connector.handleJourneyPatternRef(refs),
			JourneyPatternName: connector.handleJourneyPatternName(refs),
			VehicleLocation:    connector.handleVehicleLocation(vs[i]),
		}

		framedVehicleJourneyRef := &siri.SIRIFramedVehicleJourneyRef{}
		modelDate := connector.partner.Model().Date()
		framedVehicleJourneyRef.DataFrameRef =
			connector.Partner().DataFrameIdentifierGenerator().NewIdentifier(idgen.IdentifierAttributes{Id: modelDate.String()})
		framedVehicleJourneyRef.DatedVehicleJourneyRef = dvj

		monitoredVehicleJourney.FramedVehicleJourneyRef = framedVehicleJourneyRef
		activity.MonitoredVehicleJourney = monitoredVehicleJourney
		response.VehicleActivity = append(response.VehicleActivity, activity)

		vehicleIds = append(vehicleIds, vehicleId.Value())
	}

	if connector.partner.PartnerSettings.SortPaylodForTest() {
		sort.Sort(siri.SortByVehicleMonitoringRef{VehicleActivities: response.VehicleActivity})
	}

	message.Vehicles = vehicleIds

	siriResponse.SIRIVehicleMonitoringDelivery = response

	return siriResponse
}

func (connector *SIRIVehicleMonitoringRequestBroadcaster) datedVehicleJourneyRef(vehicleJourney *model.VehicleJourney) (string, bool) {
	vehicleJourneyId, ok := vehicleJourney.ObjectIDWithFallback(connector.vjRemoteObjectidKinds)

	var dataVehicleJourneyRef string
	if ok {
		dataVehicleJourneyRef = vehicleJourneyId.Value()
	} else {
		defaultObjectID, ok := vehicleJourney.ObjectID("_default")
		if !ok {
			return "", false
		}
		dataVehicleJourneyRef =
			connector.Partner().ReferenceIdentifierGenerator().NewIdentifier(idgen.IdentifierAttributes{Type: "VehicleJourney", Id: defaultObjectID.Value()})
	}
	return dataVehicleJourneyRef, true
}

func (connector *SIRIVehicleMonitoringRequestBroadcaster) handleVehicleLocation(v *model.Vehicle) *siri.SIRIVehicleLocation {
	var lat = v.Latitude
	var lon = v.Longitude
	if lat != 0. || lon != 0. {
		vehicleLocation := &siri.SIRIVehicleLocation{
			Longitude: lon,
			Latitude:  lat,
		}
		return vehicleLocation
	}
	return nil
}

func (connector *SIRIVehicleMonitoringRequestBroadcaster) handleProgressBetweenStops(v *model.Vehicle) *siri.SIRIProgressBetweenStops {
	var dist = v.LinkDistance
	var percent = v.Percentage
	if dist != 0. || percent != 0. {
		progressBetweenStops := &siri.SIRIProgressBetweenStops{
			LinkDistance: dist,
			Percentage:   percent,
		}
		return progressBetweenStops
	}
	return nil
}

func (connector *SIRIVehicleMonitoringRequestBroadcaster) handleJourneyPatternRef(refs model.References) string {
	journeyPatternRef, ok := refs.Get("JourneyPatternRef")
	if ok {
		if connector.remoteObjectidKind == journeyPatternRef.ObjectId.Kind() {
			return journeyPatternRef.ObjectId.Value()
		}
	}

	return ""
}

func (connector *SIRIVehicleMonitoringRequestBroadcaster) handleJourneyPatternName(refs model.References) string {
	journeyPatternName, ok := refs.Get("JourneyPatternName")
	if ok {
		if connector.remoteObjectidKind == journeyPatternName.ObjectId.Kind() {
			return journeyPatternName.ObjectId.Value()
		}
	}

	return ""
}

func (connector *SIRIVehicleMonitoringRequestBroadcaster) handleRef(refType, origin string, references model.References) string {
	reference, ok := references.Get(refType)
	if !ok || reference.ObjectId == nil || (refType == "DestinationRef" && connector.noDestinationRefRewritingFrom(origin)) {
		return ""
	}
	return connector.resolveStopAreaRef(reference)
}

func (connector *SIRIVehicleMonitoringRequestBroadcaster) noDestinationRefRewritingFrom(origin string) bool {
	ndrrf := connector.Partner().NoDestinationRefRewritingFrom()
	for _, o := range ndrrf {
		if origin == strings.TrimSpace(o) {
			return true
		}
	}
	return false
}

func (connector *SIRIVehicleMonitoringRequestBroadcaster) resolveStopAreaRef(reference model.Reference) string {
	stopArea, ok := connector.partner.Model().StopAreas().FindByObjectId(*reference.ObjectId)
	if ok {
		obj, ok := stopArea.ReferentOrSelfObjectId(connector.remoteObjectidKind)
		if ok {
			return obj.Value()
		}
	}
	return connector.Partner().ReferenceStopAreaIdentifierGenerator().NewIdentifier(idgen.IdentifierAttributes{Id: reference.GetSha1()})
}

func (factory *SIRIVehicleMonitoringRequestBroadcasterFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteObjectIdKind()
	apiPartner.ValidatePresenceOfLocalCredentials()
}

func (factory *SIRIVehicleMonitoringRequestBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRIVehicleMonitoringRequestBroadcaster(partner)
}
