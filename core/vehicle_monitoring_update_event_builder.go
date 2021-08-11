package core

import (
	"strconv"

	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri"
	"bitbucket.org/enroute-mobi/ara/uuid"
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
		remoteObjectidKind:            partner.Setting(REMOTE_OBJECTID_KIND),
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

func (builder *VehicleMonitoringUpdateEventBuilder) buildUpdateEvents(xmlVehicleActivity *siri.XMLVehicleActivity) {
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
			DestinationName: xmlVehicleActivity.DestinationName(),
			Monitored:       xmlVehicleActivity.Monitored(),

			ObjectidKind: builder.remoteObjectidKind,
			SiriXML:      &xmlVehicleActivity.XMLMonitoredVehicleJourney,
		}

		builder.vehicleMonitoringUpdateEvents.VehicleJourneys[xmlVehicleActivity.DatedVehicleJourneyRef()] = vjEvent
	}

	// Vehicles
	_, ok = builder.vehicleMonitoringUpdateEvents.Vehicles[xmlVehicleActivity.VehicleRef()]
	if !ok {
		vObjectId := model.NewObjectID(builder.remoteObjectidKind, xmlVehicleActivity.VehicleRef())
		longitude, _ := strconv.ParseFloat(xmlVehicleActivity.Longitude(), 64)
		latitude, _ := strconv.ParseFloat(xmlVehicleActivity.Latitude(), 64)
		bearing, _ := strconv.ParseFloat(xmlVehicleActivity.Bearing(), 64)

		vEvent := &model.VehicleUpdateEvent{
			Origin:                 origin,
			ObjectId:               vObjectId,
			VehicleJourneyObjectId: vjObjectId,
			SRSName:                xmlVehicleActivity.SRSName(),
			Coordinates:            xmlVehicleActivity.Coordinates(),
			DriverRef:              xmlVehicleActivity.DriverRef(),
			Longitude:              longitude,
			Latitude:               latitude,
			Bearing:                bearing,
			LinkDistance:           xmlVehicleActivity.LinkDistance(),
			Percentage:             xmlVehicleActivity.Percentage(),
			ValidUntilTime:         xmlVehicleActivity.ValidUntilTime(),
		}

		builder.vehicleMonitoringUpdateEvents.Vehicles[xmlVehicleActivity.DatedVehicleJourneyRef()] = vEvent
		builder.vehicleMonitoringUpdateEvents.VehicleRefs[xmlVehicleActivity.StopPointRef()] = struct{}{}
	}
}

func (builder *VehicleMonitoringUpdateEventBuilder) SetUpdateEvents(activities []*siri.XMLVehicleActivity) {
	for _, xmlVehicleActivity := range activities {
		builder.buildUpdateEvents(xmlVehicleActivity)
	}
}

func (builder *VehicleMonitoringUpdateEventBuilder) UpdateEvents() VehicleMonitoringUpdateEvents {
	return *builder.vehicleMonitoringUpdateEvents
}
