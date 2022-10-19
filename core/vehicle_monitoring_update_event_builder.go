package core

import (
	"fmt"
	"strconv"
	"strings"

	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
	"bitbucket.org/enroute-mobi/ara/uuid"
	"github.com/wroge/wgs84"
)

type VehicleMonitoringUpdateEventBuilder struct {
	clock.ClockConsumer
	uuid.UUIDConsumer

	partner            *Partner
	remoteObjectidKind string

	vehicleMonitoringUpdateEvents *VehicleMonitoringUpdateEvents
}

type VehicleMonitoringUpdateEvents struct {
	StopAreas       map[string]*model.StopAreaUpdateEvent
	Lines           map[string]*model.LineUpdateEvent
	VehicleJourneys map[string]*model.VehicleJourneyUpdateEvent
	Vehicles        map[string]*model.VehicleUpdateEvent
	VehicleRefs     map[string]struct{}
}

func NewVehicleMonitoringUpdateEventBuilder(partner *Partner) VehicleMonitoringUpdateEventBuilder {
	return VehicleMonitoringUpdateEventBuilder{
		partner:                       partner,
		remoteObjectidKind:            partner.RemoteObjectIDKind(),
		vehicleMonitoringUpdateEvents: newVehicleMonitoringUpdateEvents(),
	}
}

func newVehicleMonitoringUpdateEvents() *VehicleMonitoringUpdateEvents {
	return &VehicleMonitoringUpdateEvents{
		StopAreas:       make(map[string]*model.StopAreaUpdateEvent),
		Lines:           make(map[string]*model.LineUpdateEvent),
		VehicleJourneys: make(map[string]*model.VehicleJourneyUpdateEvent),
		Vehicles:        make(map[string]*model.VehicleUpdateEvent),
		VehicleRefs:     make(map[string]struct{}),
	}
}

func (builder *VehicleMonitoringUpdateEventBuilder) buildUpdateEvents(xmlVehicleActivity *sxml.XMLVehicleActivity) {
	origin := string(builder.partner.Slug())

	// StopAreas
	stopAreaObjectId := model.NewObjectID(builder.remoteObjectidKind, xmlVehicleActivity.StopPointRef())

	_, ok := builder.vehicleMonitoringUpdateEvents.StopAreas[xmlVehicleActivity.StopPointRef()]
	if !ok {
		// CollectedAlways is false by default
		event := &model.StopAreaUpdateEvent{
			Origin:   origin,
			ObjectId: stopAreaObjectId,
			Name:     xmlVehicleActivity.StopPointName(),
		}

		builder.vehicleMonitoringUpdateEvents.StopAreas[xmlVehicleActivity.StopPointRef()] = event
	}

	// Lines
	lineObjectId := model.NewObjectID(builder.remoteObjectidKind, xmlVehicleActivity.LineRef())

	_, ok = builder.vehicleMonitoringUpdateEvents.Lines[xmlVehicleActivity.LineRef()]
	if !ok {
		// CollectedAlways is false by default
		lineEvent := &model.LineUpdateEvent{
			Origin:   origin,
			ObjectId: lineObjectId,
			Name:     xmlVehicleActivity.PublishedLineName(),
		}

		builder.vehicleMonitoringUpdateEvents.Lines[xmlVehicleActivity.LineRef()] = lineEvent
	}

	// VehicleJourneys
	vjObjectId := model.NewObjectID(builder.remoteObjectidKind, xmlVehicleActivity.DatedVehicleJourneyRef())

	_, ok = builder.vehicleMonitoringUpdateEvents.VehicleJourneys[xmlVehicleActivity.DatedVehicleJourneyRef()]
	if !ok {
		vjEvent := &model.VehicleJourneyUpdateEvent{
			Origin:          origin,
			ObjectId:        vjObjectId,
			LineObjectId:    lineObjectId,
			OriginRef:       xmlVehicleActivity.OriginRef(),
			OriginName:      xmlVehicleActivity.OriginName(),
			DestinationRef:  xmlVehicleActivity.DestinationRef(),
			DirectionType:   xmlVehicleActivity.DirectionRef(),
			DestinationName: xmlVehicleActivity.DestinationName(),
			Monitored:       xmlVehicleActivity.Monitored(),
			Occupancy:       model.NormalizedOccupancyName(xmlVehicleActivity.Occupancy()),

			ObjectidKind: builder.remoteObjectidKind,
			SiriXML:      &xmlVehicleActivity.XMLMonitoredVehicleJourney,
		}

		builder.vehicleMonitoringUpdateEvents.VehicleJourneys[xmlVehicleActivity.DatedVehicleJourneyRef()] = vjEvent
	}

	// Vehicles
	_, ok = builder.vehicleMonitoringUpdateEvents.Vehicles[xmlVehicleActivity.VehicleRef()]
	if !ok {
		vObjectId := model.NewObjectID(builder.remoteObjectidKind, xmlVehicleActivity.VehicleRef())
		bearing, _ := strconv.ParseFloat(xmlVehicleActivity.Bearing(), 64)
		linkDistance, _ := strconv.ParseFloat(xmlVehicleActivity.LinkDistance(), 64)
		percentage, _ := strconv.ParseFloat(xmlVehicleActivity.Percentage(), 64)

		vEvent := &model.VehicleUpdateEvent{
			Origin:                 origin,
			ObjectId:               vObjectId,
			StopAreaObjectId:       stopAreaObjectId,
			VehicleJourneyObjectId: vjObjectId,
			DriverRef:              xmlVehicleActivity.DriverRef(),
			Bearing:                bearing,
			LinkDistance:           linkDistance,
			Percentage:             percentage,
			ValidUntilTime:         xmlVehicleActivity.ValidUntilTime(),
			RecordedAt:             xmlVehicleActivity.RecordedAtTime(),
			Occupancy:              model.NormalizedOccupancyName(xmlVehicleActivity.Occupancy()),
		}

		longitude, latitude, err := builder.handleCoordinates(xmlVehicleActivity)
		if err == nil {
			vEvent.Longitude = longitude
			vEvent.Latitude = latitude
		}

		builder.vehicleMonitoringUpdateEvents.Vehicles[xmlVehicleActivity.DatedVehicleJourneyRef()] = vEvent
		builder.vehicleMonitoringUpdateEvents.VehicleRefs[xmlVehicleActivity.StopPointRef()] = struct{}{}
	}
}

func (builder *VehicleMonitoringUpdateEventBuilder) handleCoordinates(xmlVehicleActivity *sxml.XMLVehicleActivity) (lon, lat float64, e error) {
	longitude, _ := strconv.ParseFloat(xmlVehicleActivity.Longitude(), 64)
	latitude, _ := strconv.ParseFloat(xmlVehicleActivity.Latitude(), 64)

	if latitude != 0 || longitude != 0 {
		return longitude, latitude, nil
	}

	if xmlVehicleActivity.Coordinates() == "" {
		e = fmt.Errorf("no coordinates")
		return
	}

	cs := strings.Split(xmlVehicleActivity.Coordinates(), " ")
	var x, y float64
	x, e = strconv.ParseFloat(cs[0], 64)
	if e != nil {
		return
	}
	y, e = strconv.ParseFloat(cs[1], 64)
	if e != nil {
		return
	}

	var formatSRSName int
	formatSRSName, e = builder.formatSRSNameWithDefaut(xmlVehicleActivity.SRSName())
	if e != nil {
		return lon, lat, e
	}

	epsg := wgs84.EPSG()
	lon, lat, _, e = epsg.SafeTransform(formatSRSName, 4326)(x, y, 0)

	return lon, lat, e
}

func (builder *VehicleMonitoringUpdateEventBuilder) formatSRSNameWithDefaut(srs string) (int, error) {
	if srs == "" {
		return convertSRSNameToValue(strings.TrimPrefix(builder.partner.DefaultSRSName(), "EPSG:"))
	}

	if strings.HasPrefix(srs, "EPSG:") {
		return convertSRSNameToValue(strings.TrimPrefix(srs, "EPSG:"))
	}

	return convertSRSNameToValue(srs)
}

func convertSRSNameToValue(srs string) (int, error) {
	srsValue, err := strconv.Atoi(srs)
	if err != nil {
		return 0, err
	}

	return srsValue, nil
}

func (builder *VehicleMonitoringUpdateEventBuilder) SetUpdateEvents(activities []*sxml.XMLVehicleActivity) {
	for _, xmlVehicleActivity := range activities {
		builder.buildUpdateEvents(xmlVehicleActivity)
	}
}

func (builder *VehicleMonitoringUpdateEventBuilder) UpdateEvents() VehicleMonitoringUpdateEvents {
	return *builder.vehicleMonitoringUpdateEvents
}
