package core

import (
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

type BroadcastGeneralMessageBuilder struct {
	clock.ClockConsumer
	uuid.UUIDConsumer

	tx                 *model.Transaction
	partner            *Partner
	referenceGenerator *IdentifierGenerator
	remoteObjectidKind string
	lineRef            map[string]struct{}
	stopPointRef       map[string]struct{}

	InfoChannelRef []string
}

func NewBroadcastGeneralMessageBuilder(tx *model.Transaction, partner *Partner, connector string) *BroadcastGeneralMessageBuilder {
	return &BroadcastGeneralMessageBuilder{
		tx:                 tx,
		partner:            partner,
		referenceGenerator: partner.IdentifierGenerator(REFERENCE_IDENTIFIER),
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
	if situation.Origin == string(builder.partner.Slug()) || situation.Channel == "Commercial" || situation.ValidUntil.Before(builder.Clock().Now()) {
		return nil
	}

	// InfoChannelRef filter

	if !builder.checkInfoChannelRef(builder.InfoChannelRef, situation.Channel) {
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
		infoMessageIdentifier = builder.referenceGenerator.NewIdentifier(IdentifierAttributes{Type: "InfoMessage", Id: objectid.Value()})
	}

	siriGeneralMessage := &siri.SIRIGeneralMessage{
		ItemIdentifier:        builder.referenceGenerator.NewIdentifier(IdentifierAttributes{Type: "Item", Id: builder.NewUUID()}),
		InfoMessageIdentifier: infoMessageIdentifier,
		InfoChannelRef:        situation.Channel,
		InfoMessageVersion:    situation.Version,
		ValidUntilTime:        situation.ValidUntil,
		RecordedAtTime:        situation.RecordedAt,
		FormatRef:             "STIF-IDF",
	}
	for _, reference := range situation.References {
		id, ok := builder.resolveReference(reference)
		if !ok {
			continue
		}
		siriGeneralMessage.References = append(siriGeneralMessage.References, &siri.SIRIReference{Kind: reference.Type, Id: id})
	}
	if !builder.checkFilter(siriGeneralMessage.References) {
		return nil
	}
	for _, lineSection := range situation.LineSections {
		siriLineSection, ok := builder.handleLineSection(*lineSection)
		if !ok {
			continue
		}
		siriGeneralMessage.LineSections = append(siriGeneralMessage.LineSections, siriLineSection)
	}
	if len(siriGeneralMessage.References) == 0 && len(siriGeneralMessage.LineSections) == 0 {
		return nil
	}

	for _, message := range situation.Messages {
		siriMessage := &siri.SIRIMessage{
			Content:             message.Content,
			Type:                message.Type,
			NumberOfLines:       message.NumberOfLines,
			NumberOfCharPerLine: message.NumberOfCharPerLine,
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

func (builder *BroadcastGeneralMessageBuilder) handleLineSection(lineSection model.References) (*siri.SIRILineSection, bool) {
	siriLineSection := &siri.SIRILineSection{}
	lineSectionMap := make(map[string]string)

	for kind, reference := range lineSection.GetReferences() {
		ref, ok := builder.resolveReference(&reference)
		if !ok {
			return nil, false
		}
		lineSectionMap[kind] = ref
	}

	siriLineSection.FirstStop = lineSectionMap["FirstStop"]
	siriLineSection.LastStop = lineSectionMap["LastStop"]
	siriLineSection.LineRef = lineSectionMap["LineRef"]

	return siriLineSection, true
}

func (builder *BroadcastGeneralMessageBuilder) resolveReference(reference *model.Reference) (string, bool) {
	switch reference.Type {
	case "LineRef":
		return builder.resolveLineRef(reference)
	case "StopPointRef", "DestinationRef", "FirstStop", "LastStop":
		return builder.resolveStopAreaRef(reference)
	default:
		kind := reference.Type
		return builder.referenceGenerator.NewIdentifier(IdentifierAttributes{Type: kind[:len(kind)-3], Id: reference.GetSha1()}), true
	}
}

func (builder *BroadcastGeneralMessageBuilder) resolveLineRef(reference *model.Reference) (string, bool) {
	line, ok := builder.tx.Model().Lines().FindByObjectId(*reference.ObjectId)
	if !ok {
		return "", false
	}
	lineObjectId, ok := line.ObjectID(builder.remoteObjectidKind)
	if !ok {
		return "", false
	}
	return lineObjectId.Value(), true
}

func (builder *BroadcastGeneralMessageBuilder) resolveStopAreaRef(reference *model.Reference) (string, bool) {
	stopArea, ok := builder.tx.Model().StopAreas().FindByObjectId(*reference.ObjectId)
	if !ok {
		return "", false
	}
	stopAreaObjectId, ok := stopArea.ReferentOrSelfObjectId(builder.remoteObjectidKind)
	if !ok {
		return "", false
	}
	return stopAreaObjectId.Value(), true
}

func (builder *BroadcastGeneralMessageBuilder) checkFilter(references []*siri.SIRIReference) bool {
	if len(builder.lineRef) == 0 && len(builder.stopPointRef) == 0 {
		return true
	}

	for _, reference := range references {
		switch reference.Kind {
		case "LineRef":
			if _, ok := builder.lineRef[reference.Id]; ok {
				return true
			}
		case "StopPointRef":
			if _, ok := builder.stopPointRef[reference.Id]; ok {
				return true
			}
		}
	}

	return false
}
