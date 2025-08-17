package core

import (
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

type FacilityMonitoringUpdateEventBuilder struct {
	clock.ClockConsumer
	uuid.UUIDConsumer

	partner              *Partner
	remoteCodeSpace      string
	facilityUpdateEvents *CollectUpdateEvents
}

func NewFacilityMonitoringUpdateEventBuilder(partner *Partner) FacilityMonitoringUpdateEventBuilder {
	return FacilityMonitoringUpdateEventBuilder{
		partner:              partner,
		remoteCodeSpace:      partner.RemoteCodeSpace(),
		facilityUpdateEvents: NewCollectUpdateEvents(),
	}
}

func (builder *FacilityMonitoringUpdateEventBuilder) buildUpdateEvents(xmlFacilityEvent *sxml.XMLFacilityCondition) {
	origin := string(builder.partner.Slug())

	// Facilities
	facilityCode := model.NewCode(builder.remoteCodeSpace, xmlFacilityEvent.FacilityRef())

	_, ok := builder.facilityUpdateEvents.Facilities[xmlFacilityEvent.FacilityRef()]
	if !ok {
		event := &model.FacilityUpdateEvent{
			Origin: origin,
			Code:   facilityCode,
			Status: xmlFacilityEvent.FacilityStatus(),
		}

		builder.facilityUpdateEvents.Facilities[xmlFacilityEvent.FacilityRef()] = event
		builder.facilityUpdateEvents.FacilityRefs[xmlFacilityEvent.FacilityRef()] = struct{}{}
	}
}

func (builder *FacilityMonitoringUpdateEventBuilder) SetUpdateEvents(xmlFacilities []*sxml.XMLFacilityCondition) {
	for _, xmlFacilityCondition := range xmlFacilities {
		builder.buildUpdateEvents(xmlFacilityCondition)
	}
}

func (builder *FacilityMonitoringUpdateEventBuilder) UpdateEvents() CollectUpdateEvents {
	return *builder.facilityUpdateEvents
}
