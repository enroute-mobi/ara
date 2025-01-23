package core

import (
	"slices"

	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/core/idgen"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/siri"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

type BroadcastGeneralMessageBuilder struct {
	clock.ClockConsumer
	uuid.UUIDConsumer

	partner         *Partner
	remoteCodeSpace string
	lineRef         map[string]struct{}
	stopPointRef    map[string]struct{}

	InfoChannelRef []string
}

func NewBroadcastGeneralMessageBuilder(partner *Partner, connector string) *BroadcastGeneralMessageBuilder {
	return &BroadcastGeneralMessageBuilder{
		partner:         partner,
		remoteCodeSpace: partner.RemoteCodeSpace(connector),
		lineRef:         make(map[string]struct{}),
		stopPointRef:    make(map[string]struct{}),
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

func (builder *BroadcastGeneralMessageBuilder) BuildGeneralMessageCancellation(situation model.Situation) *siri.SIRIGeneralMessageCancellation {
	if !builder.canBroadcast(situation) {
		return nil
	}

	var infoMessageIdentifier string
	code, present := situation.Code(builder.remoteCodeSpace)
	if present {
		infoMessageIdentifier = code.Value()
	} else {
		code, present = situation.Code("_default")
		if !present {
			return nil
		}
		infoMessageIdentifier = builder.partner.NewIdentifier(idgen.IdentifierAttributes{Type: "InfoMessage", Id: code.Value()})
	}

	siriGeneralMessageCancellation := &siri.SIRIGeneralMessageCancellation{
		ItemIdentifier:        builder.partner.NewIdentifier(idgen.IdentifierAttributes{Type: "Item", Id: builder.NewUUID()}),
		InfoMessageIdentifier: infoMessageIdentifier,
		RecordedAtTime:        situation.RecordedAt,
	}
	return siriGeneralMessageCancellation
}

func (builder *BroadcastGeneralMessageBuilder) BuildGeneralMessage(situation model.Situation) *siri.SIRIGeneralMessage {
	if !builder.canBroadcast(situation) {
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
	code, present := situation.Code(builder.remoteCodeSpace)
	if present {
		infoMessageIdentifier = code.Value()
	} else {
		code, present = situation.Code(model.Default)
		if !present {
			return nil
		}
		infoMessageIdentifier = builder.partner.NewIdentifier(idgen.IdentifierAttributes{Type: "InfoMessage", Id: code.Value()})
	}

	siriGeneralMessage := &siri.SIRIGeneralMessage{
		ItemIdentifier:        builder.partner.NewIdentifier(idgen.IdentifierAttributes{Type: "Item", Id: builder.NewUUID()}),
		InfoMessageIdentifier: infoMessageIdentifier,
		InfoChannelRef:        channel,
		InfoMessageVersion:    situation.Version,
		ValidUntilTime:        situation.GMValidUntil(),
		RecordedAtTime:        situation.RecordedAt,
		FormatRef:             "STIF-IDF",
	}

	// Build from Affects
	for _, affect := range situation.Affects {
		switch affect.GetType() {
		case model.SituationTypeStopArea:
			builder.buildAffectedStopArea(siriGeneralMessage, affect)
		case model.SituationTypeLine:
			builder.buildAffectedLine(siriGeneralMessage, affect)
		}
	}

	if !builder.checkAffectFilter(siriGeneralMessage.AffectedRefs) {
		return nil
	}

	if len(siriGeneralMessage.AffectedRefs) == 0 && len(siriGeneralMessage.LineSections) == 0 {
		return nil
	}

	if situation.Summary != nil {
		sts := siri.SIRITranslatedString{
			Tag:                       "MessageText",
			TranslatedString: *situation.Summary,
		}

		siriMessage := &siri.SIRIMessage{
			Type:                 "shortMessage",
			SIRITranslatedString: sts,
		}

		siriGeneralMessage.Messages = append(siriGeneralMessage.Messages, siriMessage)
	}

	if situation.Description != nil {
		sts := siri.SIRITranslatedString{
			Tag:                       "MessageText",
			TranslatedString: *situation.Description,
		}
		siriMessage := &siri.SIRIMessage{
			Type:                 "textOnly",
			SIRITranslatedString: sts,
		}

		siriGeneralMessage.Messages = append(siriGeneralMessage.Messages, siriMessage)
	}

	return siriGeneralMessage
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

func (builder *BroadcastGeneralMessageBuilder) buildAffectedStopArea(message *siri.SIRIGeneralMessage, affect model.Affect) {
	affectedStopAreaId, ok := builder.resolveStopAreaRef(model.StopAreaId(affect.GetId()))
	if !ok {
		logger.Log.Debugf("Unknown StopArea %s", affect.GetId())
		return
	}

	affectedStopAreaRef := &siri.SIRIAffectedRef{
		Kind: "StopPointRef",
		Id:   affectedStopAreaId,
	}
	message.AffectedRefs = append(message.AffectedRefs, affectedStopAreaRef)
}

func (builder *BroadcastGeneralMessageBuilder) buildAffectedLine(message *siri.SIRIGeneralMessage, affect model.Affect) {
	affectedLineId, ok := builder.resolveAffectedLineRef(affect)
	if !ok {
		logger.Log.Debugf("Unknown Line %s", affect.GetId())
		return
	}
	affectedLineRef := &siri.SIRIAffectedRef{
		Kind: "LineRef",
		Id:   affectedLineId,
	}
	message.AffectedRefs = append(message.AffectedRefs, affectedLineRef)
	for _, affectedDestination := range affect.(*model.AffectedLine).AffectedDestinations {
		affectedDestinationId, ok := builder.resolveStopAreaRef(model.StopAreaId(affectedDestination.StopAreaId))
		if !ok {
			continue
		}
		affectedDestinationRef := &siri.SIRIAffectedRef{
			Kind: "DestinationRef",
			Id:   affectedDestinationId,
		}
		message.AffectedRefs = append(message.AffectedRefs, affectedDestinationRef)
	}
	for _, affectedSection := range affect.(*model.AffectedLine).AffectedSections {
		firstStopId, ok := builder.resolveStopAreaRef(model.StopAreaId(affectedSection.FirstStop))
		if !ok {
			continue
		}
		lastStopId, ok := builder.resolveStopAreaRef(model.StopAreaId(affectedSection.LastStop))
		if !ok {
			continue
		}
		affectedSectionRef := &siri.SIRILineSection{
			FirstStop: firstStopId,
			LastStop:  lastStopId,
			LineRef:   affectedLineId,
		}
		message.LineSections = append(message.LineSections, affectedSectionRef)
	}
	for _, affectedRoute := range affect.(*model.AffectedLine).AffectedRoutes {
		affectedRouteRef := &siri.SIRIAffectedRef{
			Kind: "RouteRef",
			Id:   affectedRoute.RouteRef,
		}
		message.AffectedRefs = append(message.AffectedRefs, affectedRouteRef)
	}
}

func (builder *BroadcastGeneralMessageBuilder) resolveAffectedLineRef(affect model.Affect) (string, bool) {
	line, ok := builder.partner.Model().Lines().Find(model.LineId(affect.GetId()))
	if !ok {
		return "", false
	}
	lineCode, ok := line.Code(builder.remoteCodeSpace)
	if !ok {
		return "", false
	}
	return lineCode.Value(), true
}

func (builder *BroadcastGeneralMessageBuilder) resolveStopAreaRef(stopAreaId model.StopAreaId) (string, bool) {
	stopArea, ok := builder.partner.Model().StopAreas().Find(stopAreaId)
	if !ok {
		return "", false
	}
	stopAreaCode, ok := stopArea.ReferentOrSelfCode(builder.remoteCodeSpace)
	if !ok {
		return "", false
	}
	return stopAreaCode.Value(), true
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

func (builder *BroadcastGeneralMessageBuilder) canBroadcast(situation model.Situation) bool {
	if situation.Origin == string(builder.partner.Slug()) {
		return false
	}

	requestPeriod := &model.TimeRange{
		StartTime: builder.Clock().Now().Add(-builder.partner.SituationsTTL()),
	}

	if !situation.BroadcastPeriod().Overlaps(requestPeriod) {
		return false
	}

	tagsToBroadcast := builder.partner.BroadcastSituationsInternalTags()
	if len(tagsToBroadcast) != 0 {
		for _, tag := range situation.InternalTags {
			if slices.Contains(tagsToBroadcast, tag) {
				return true
			}
		}
		return false
	}

	return true
}
