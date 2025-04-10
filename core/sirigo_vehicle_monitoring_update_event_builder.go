package core

import (
	"fmt"
	"strconv"
	"strings"

	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/uuid"
	"bitbucket.org/enroute-mobi/sirigo/siristructs"
	"github.com/wroge/wgs84"
)

type SirigoVehicleMonitoringUpdateEventBuilder struct {
	clock.ClockConsumer
	uuid.UUIDConsumer

	partner         *Partner
	remoteCodeSpace string

	vehicleMonitoringUpdateEvents *VehicleMonitoringUpdateEvents
}

// type VehicleMonitoringUpdateEvents struct {
// 	StopAreas          map[string]*model.StopAreaUpdateEvent
// 	Lines              map[string]*model.LineUpdateEvent
// 	VehicleJourneys    map[string]*model.VehicleJourneyUpdateEvent
// 	Vehicles           map[string]*model.VehicleUpdateEvent
// 	VehicleRefs        map[string]struct{}
// 	LineRefs           map[string]struct{}
// 	VehicleJourneyRefs map[string]struct{}
// 	MonitoringRefs     map[string]struct{}
// }

func NewSirigoVehicleMonitoringUpdateEventBuilder(partner *Partner) SirigoVehicleMonitoringUpdateEventBuilder {
	return SirigoVehicleMonitoringUpdateEventBuilder{
		partner:                       partner,
		remoteCodeSpace:               partner.RemoteCodeSpace(),
		vehicleMonitoringUpdateEvents: newVehicleMonitoringUpdateEvents(),
	}
}

// func newVehicleMonitoringUpdateEvents() *VehicleMonitoringUpdateEvents {
// 	return &VehicleMonitoringUpdateEvents{
// 		StopAreas:          make(map[string]*model.StopAreaUpdateEvent),
// 		Lines:              make(map[string]*model.LineUpdateEvent),
// 		VehicleJourneys:    make(map[string]*model.VehicleJourneyUpdateEvent),
// 		Vehicles:           make(map[string]*model.VehicleUpdateEvent),
// 		VehicleRefs:        make(map[string]struct{}),
// 		LineRefs:           make(map[string]struct{}),
// 		VehicleJourneyRefs: make(map[string]struct{}),
// 		MonitoringRefs:     make(map[string]struct{}),
// 	}
// }

func (builder *SirigoVehicleMonitoringUpdateEventBuilder) buildUpdateEvents(xmlVehicleActivity *siristructs.VehicleActivity) {
	origin := string(builder.partner.Slug())

	// StopAreas
	stopAreaCode := model.NewCode(builder.remoteCodeSpace, xmlVehicleActivity.StopPointRef.String)

	_, ok := builder.vehicleMonitoringUpdateEvents.StopAreas[xmlVehicleActivity.StopPointRef.String]
	if !ok {
		// CollectedAlways is false by default
		event := &model.StopAreaUpdateEvent{
			Origin: origin,
			Code:   stopAreaCode,
			Name:   xmlVehicleActivity.StopPointName.String,
		}

		builder.vehicleMonitoringUpdateEvents.StopAreas[xmlVehicleActivity.StopPointRef.String] = event
		builder.vehicleMonitoringUpdateEvents.MonitoringRefs[xmlVehicleActivity.StopPointRef.String] = struct{}{}
	}

	// Lines
	lineCode := model.NewCode(builder.remoteCodeSpace, xmlVehicleActivity.LineRef.String)

	_, ok = builder.vehicleMonitoringUpdateEvents.Lines[xmlVehicleActivity.LineRef.String]
	if !ok {
		// CollectedAlways is false by default
		lineEvent := &model.LineUpdateEvent{
			Origin: origin,
			Code:   lineCode,
			Name:   xmlVehicleActivity.PublishedLineName.String,
		}

		builder.vehicleMonitoringUpdateEvents.Lines[xmlVehicleActivity.LineRef.String] = lineEvent
		builder.vehicleMonitoringUpdateEvents.LineRefs[xmlVehicleActivity.LineRef.String] = struct{}{}
	}

	// VehicleJourneys
	vjCode := model.NewCode(builder.remoteCodeSpace, xmlVehicleActivity.DatedVehicleJourneyRef.String)

	_, ok = builder.vehicleMonitoringUpdateEvents.VehicleJourneys[xmlVehicleActivity.DatedVehicleJourneyRef.String]
	if !ok {
		vjEvent := &model.VehicleJourneyUpdateEvent{
			Origin:          origin,
			Code:            vjCode,
			LineCode:        lineCode,
			OriginRef:       xmlVehicleActivity.OriginRef.String,
			OriginName:      xmlVehicleActivity.OriginName.String,
			DestinationRef:  xmlVehicleActivity.DestinationRef.String,
			DirectionType:   builder.directionRef(xmlVehicleActivity.DirectionRef.String),
			DestinationName: xmlVehicleActivity.DestinationName.String,
			Monitored:       xmlVehicleActivity.Monitored.Bool,
			Occupancy:       model.NormalizedOccupancyName(xmlVehicleActivity.Occupancy.String),

			CodeSpace: builder.remoteCodeSpace,
		}

		builder.vehicleMonitoringUpdateEvents.VehicleJourneys[xmlVehicleActivity.DatedVehicleJourneyRef.String] = vjEvent
		builder.vehicleMonitoringUpdateEvents.VehicleJourneyRefs[xmlVehicleActivity.DatedVehicleJourneyRef.String] = struct{}{}
	}

	// Vehicles
	_, ok = builder.vehicleMonitoringUpdateEvents.Vehicles[xmlVehicleActivity.VehicleRef.String]
	if !ok {
		vCode := model.NewCode(builder.remoteCodeSpace, xmlVehicleActivity.VehicleRef.String)
		bearing := xmlVehicleActivity.Bearing.Float64
		linkDistance := xmlVehicleActivity.LinkDistance.Float64
		percentage := xmlVehicleActivity.Percentage.Float64

		vEvent := &model.VehicleUpdateEvent{
			Origin:             origin,
			Code:               vCode,
			StopAreaCode:       stopAreaCode,
			VehicleJourneyCode: vjCode,
			DriverRef:          xmlVehicleActivity.DriverRef.String,
			Bearing:            bearing,
			LinkDistance:       linkDistance,
			Percentage:         percentage,
			ValidUntilTime:     xmlVehicleActivity.ValidUntilTime.Time,
			RecordedAt:         xmlVehicleActivity.RecordedAtTime.Time,
			Occupancy:          model.NormalizedOccupancyName(xmlVehicleActivity.Occupancy.String),
			NextStopPointOrder: xmlVehicleActivity.Order.Int,
		}

		longitude, latitude, err := builder.handleCoordinates(xmlVehicleActivity)
		if err == nil {
			vEvent.Longitude = longitude
			vEvent.Latitude = latitude
		}

		builder.vehicleMonitoringUpdateEvents.Vehicles[xmlVehicleActivity.DatedVehicleJourneyRef.String] = vEvent
		builder.vehicleMonitoringUpdateEvents.VehicleRefs[xmlVehicleActivity.VehicleRef.String] = struct{}{}
	}
}

func (builder *SirigoVehicleMonitoringUpdateEventBuilder) directionRef(direction string) (dir string) {
	in, out, err := builder.partner.PartnerSettings.SIRIDirectionType()
	if err {
		return direction
	}

	switch direction {
	case in:
		dir = model.VEHICLE_DIRECTION_INBOUND
	case out:
		dir = model.VEHICLE_DIRECTION_OUTBOUND
	default:
		dir = direction
	}

	return dir
}

func (builder *SirigoVehicleMonitoringUpdateEventBuilder) handleCoordinates(xmlVehicleActivity *siristructs.VehicleActivity) (lon, lat float64, e error) {
	if xmlVehicleActivity.Latitude.Valid && xmlVehicleActivity.Longitude.Valid {
		return xmlVehicleActivity.Longitude.Float64, xmlVehicleActivity.Latitude.Float64, nil
	}

	if !xmlVehicleActivity.Coordinates.Valid {
		e = fmt.Errorf("no coordinates")
		return
	}

	cs := strings.Split(xmlVehicleActivity.Coordinates.String, " ")
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
	formatSRSName, e = builder.formatSRSNameWithDefaut(xmlVehicleActivity.SRSName.String)
	if e != nil {
		return lon, lat, e
	}

	epsg := wgs84.EPSG()
	lon, lat, _, e = epsg.SafeTransform(formatSRSName, 4326)(x, y, 0)

	return lon, lat, e
}

func (builder *SirigoVehicleMonitoringUpdateEventBuilder) formatSRSNameWithDefaut(srs string) (int, error) {
	if srs == "" {
		return convertSRSNameToValue(strings.TrimPrefix(builder.partner.DefaultSRSName(), "EPSG:"))
	}

	if strings.HasPrefix(srs, "EPSG:") {
		return convertSRSNameToValue(strings.TrimPrefix(srs, "EPSG:"))
	}

	return convertSRSNameToValue(srs)
}

// func convertSRSNameToValue(srs string) (int, error) {
// 	srsValue, err := strconv.Atoi(srs)
// 	if err != nil {
// 		return 0, err
// 	}

// 	return srsValue, nil
// }

func (builder *SirigoVehicleMonitoringUpdateEventBuilder) SetUpdateEvents(activities []*siristructs.VehicleActivity) {
	for _, xmlVehicleActivity := range activities {
		builder.buildUpdateEvents(xmlVehicleActivity)
	}
}

func (builder *SirigoVehicleMonitoringUpdateEventBuilder) UpdateEvents() VehicleMonitoringUpdateEvents {
	return *builder.vehicleMonitoringUpdateEvents
}
