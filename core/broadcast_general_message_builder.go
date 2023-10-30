package core

import (
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/core/idgen"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/siri"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

type BroadcastGeneralMessageBuilder struct {
	clock.ClockConsumer
	uuid.UUIDConsumer

	partner            *Partner
	referenceGenerator *idgen.IdentifierGenerator
	remoteObjectidKind string
	lineRef            map[string]struct{}
	stopPointRef       map[string]struct{}

	InfoChannelRef []string
}

func NewBroadcastGeneralMessageBuilder(partner *Partner, connector string) *BroadcastGeneralMessageBuilder {
	return &BroadcastGeneralMessageBuilder{
		partner:            partner,
		referenceGenerator: partner.ReferenceIdentifierGenerator(),
		remoteObjectidKind: partner.RemoteObjectIDKind(connector),
		lineRef:            make(map[string]struct{}),
		stopPointRef:       make(map[string]struct{}),
	}
}

func (builder *BroadcastGeneralMessageBuilder) SetLineRef(lineRef []string) {
	if (len(lineRef) == 0) || (len(lineRef) == 1 && lineRef[0] == "") {
		return
	}

	for i := range lineRef {
		builder.lineRef[lineRef[i]] = struct{}{}
	}
}

func (builder *BroadcastGeneralMessageBuilder) SetStopPointRef(stopPointRef []string) {
	if (len(stopPointRef) == 0) || (len(stopPointRef) == 1 && stopPointRef[0] == "") {
		return
	}

	for i := range stopPointRef {
		builder.stopPointRef[stopPointRef[i]] = struct{}{}
	}
}

func (builder *BroadcastGeneralMessageBuilder) BuildGeneralMessage(situation model.Situation) *siri.SIRIGeneralMessage {
	if situation.Origin == string(builder.partner.Slug()) || situation.GMValidUntil().Before(builder.Clock().Now()) {
		return nil
	}

	// Filter by expected channel for GM
	channel, ok := situation.GetGMChannel()
	if !ok {
		if situation.ReportType == model.SituationReportTypeIncident {
			channel = "Perturbation"
		} else {
			channel = "Information"
		}
	}

	// InfoChannelRef filter
	ok = builder.checkInfoChannelRef(builder.InfoChannelRef, channel)
	if !ok {
		return nil
	}

	var infoMessageIdentifier string
	objectid, present := situation.ObjectID(builder.remoteObjectidKind)
	if present {
		infoMessageIdentifier = objectid.Value()
	} else {
		objectid, present = situation.ObjectID("_default")
		if !present {
			return nil
		}
		infoMessageIdentifier = builder.referenceGenerator.NewIdentifier(idgen.IdentifierAttributes{Type: "InfoMessage", Id: objectid.Value()})
	}

	siriGeneralMessage := &siri.SIRIGeneralMessage{
		ItemIdentifier:        builder.referenceGenerator.NewIdentifier(idgen.IdentifierAttributes{Type: "Item", Id: builder.NewUUID()}),
		InfoMessageIdentifier: infoMessageIdentifier,
		InfoChannelRef:        channel,
		InfoMessageVersion:    situation.Version,
		ValidUntilTime:        situation.GMValidUntil(),
		RecordedAtTime:        situation.RecordedAt,
		FormatRef:             "STIF-IDF",
	}

	for _, affect := range situation.Affects {
		switch Type := affect.GetType(); Type {
		case "StopArea":
			affectedStopAreaId, ok := builder.resolveAffectedStopAreaRef(affect)
			if !ok {
				continue
			}
			affectedStopAreaRef := &siri.SIRIAffectedRef{
				Kind: builder.setReferenceKind(affect),
				Id:   affectedStopAreaId,
			}
			siriGeneralMessage.AffectedRefs = append(siriGeneralMessage.AffectedRefs, affectedStopAreaRef)
		case "Line":
			affectedLineId, ok := builder.resolveAffectedLineRef(affect)
			if !ok {
				continue
			}
			affectedLineRef := &siri.SIRIAffectedRef{
				Kind: builder.setReferenceKind(affect),
				Id:   affectedLineId,
			}
			siriGeneralMessage.AffectedRefs = append(siriGeneralMessage.AffectedRefs, affectedLineRef)
			for _, affectedDestination := range affect.(*model.AffectedLine).AffectedDestinations {
				affectedDestinationId, ok := builder.resolveAffectedDestinationRef(model.StopAreaId(affectedDestination.StopAreaId))
				if !ok {
					continue
				}
				affectedDestinationRef := &siri.SIRIAffectedRef{
					Kind: "DestinationRef",
					Id:   affectedDestinationId,
				}
				siriGeneralMessage.AffectedRefs = append(siriGeneralMessage.AffectedRefs, affectedDestinationRef)
			}
			for _, affectedSection := range affect.(*model.AffectedLine).AffectedSections {
				// PLEASE CHANGE ME, make a standard resolveStopAreaRef !!!
				firstStopId, ok := builder.resolveAffectedDestinationRef(model.StopAreaId(affectedSection.FirstStop))
				if !ok {
					continue
				}
				lastStopId, ok := builder.resolveAffectedDestinationRef(model.StopAreaId(affectedSection.LastStop))
				if !ok {
					continue
				}
				affectedSectionRef := &siri.SIRILineSection{
					FirstStop: firstStopId,
					LastStop:  lastStopId,
					LineRef:   affectedLineId,
				}
				siriGeneralMessage.LineSections = append(siriGeneralMessage.LineSections, affectedSectionRef)
			}
			for _, affectedRoute := range affect.(*model.AffectedLine).AffectedRoutes {
				affectedRouteRef := &siri.SIRIAffectedRef{
					Kind: "RouteRef",
					Id:   affectedRoute.RouteRef,
				}
				siriGeneralMessage.AffectedRefs = append(siriGeneralMessage.AffectedRefs, affectedRouteRef)
			}
		}
	}

	for _, reference := range situation.References {
		id, ok := builder.resolveReference(reference)
		if !ok {
			continue
		}
		siriGeneralMessage.References = append(siriGeneralMessage.References, &siri.SIRIReference{Kind: reference.Type, Id: id})
	}

	if !builder.checkAffectFilter(siriGeneralMessage.AffectedRefs) {
		return nil
	}

	if len(siriGeneralMessage.References) == 0 && len(siriGeneralMessage.AffectedRefs) == 0 && len(siriGeneralMessage.LineSections) == 0 {
		return nil
	}

	var siriMessage siri.SIRIMessage
	if situation.Summary != nil {
		siriMessage.Content = situation.Summary.DefaultValue
		siriMessage.Type = "shortMessage"
		siriGeneralMessage.Messages = append(siriGeneralMessage.Messages, &siriMessage)
	}
	if situation.Description != nil {
		siriMessage.Content = situation.Description.DefaultValue
		siriMessage.Type = "longMessage"
		siriGeneralMessage.Messages = append(siriGeneralMessage.Messages, &siriMessage)
	}
	return siriGeneralMessage
}

func (builder *BroadcastGeneralMessageBuilder) setReferenceKind(affect model.Affect) string {
	switch affect.GetType() {
	case "Line":
		return "LineRef"
	case "StopArea":
		return "StopPointRef"
	}
	return ""
}

func (builder *BroadcastGeneralMessageBuilder) checkInfoChannelRef(requestChannels []string, channel string) bool {
	if (len(requestChannels) == 1 && requestChannels[0] == "") || len(requestChannels) == 0 {
		return true
	}

	for i := range requestChannels {
		if requestChannels[i] == channel {
			return true
		}
	}

	return false
}

func (builder *BroadcastGeneralMessageBuilder) resolveReference(reference *model.Reference) (string, bool) {
	switch reference.Type {
	case "DestinationRef", "FirstStop", "LastStop":
		return builder.resolveStopAreaRef(reference)
	default:
		kind := reference.Type
		return builder.referenceGenerator.NewIdentifier(idgen.IdentifierAttributes{Type: kind[:len(kind)-3], Id: reference.GetSha1()}), true
	}
}

func (builder *BroadcastGeneralMessageBuilder) resolveAffectedLineRef(affect model.Affect) (string, bool) {
	line, ok := builder.partner.Model().Lines().Find(model.LineId(affect.GetId()))
	if !ok {
		return "", false
	}
	lineObjectId, ok := line.ObjectID(builder.remoteObjectidKind)
	if !ok {
		return "", false
	}
	return lineObjectId.Value(), true
}

func (builder *BroadcastGeneralMessageBuilder) resolveAffectedStopAreaRef(affect model.Affect) (string, bool) {
	stopArea, ok := builder.partner.Model().StopAreas().Find(model.StopAreaId(affect.GetId()))
	if !ok {
		return "", false
	}
	stopAreaObjectId, ok := stopArea.ReferentOrSelfObjectId(builder.remoteObjectidKind)
	if !ok {
		return "", false
	}
	return stopAreaObjectId.Value(), true
}

func (builder *BroadcastGeneralMessageBuilder) resolveAffectedDestinationRef(stopAreaId model.StopAreaId) (string, bool) {
	stopArea, ok := builder.partner.Model().StopAreas().Find(stopAreaId)
	if !ok {
		return "", false
	}
	stopAreaObjectId, ok := stopArea.ReferentOrSelfObjectId(builder.remoteObjectidKind)
	if !ok {
		return "", false
	}
	return stopAreaObjectId.Value(), true
}

func (builder *BroadcastGeneralMessageBuilder) resolveStopAreaRef(reference *model.Reference) (string, bool) {
	stopArea, ok := builder.partner.Model().StopAreas().FindByObjectId(*reference.ObjectId)
	if !ok {
		return "", false
	}
	stopAreaObjectId, ok := stopArea.ReferentOrSelfObjectId(builder.remoteObjectidKind)
	if !ok {
		return "", false
	}
	return stopAreaObjectId.Value(), true
}

func (builder *BroadcastGeneralMessageBuilder) checkAffectFilter(affectedRefs []*siri.SIRIAffectedRef) bool {
	if len(builder.lineRef) == 0 && len(builder.stopPointRef) == 0 {
		return true
	}

	for _, affected := range affectedRefs {
		switch affected.Kind {
		case "LineRef":
			if _, ok := builder.lineRef[affected.Id]; ok {
				return true
			}
		case "StopAreaRef":
			if _, ok := builder.stopPointRef[affected.Id]; ok {
				return true
			}
		}
	}
	return false
}
