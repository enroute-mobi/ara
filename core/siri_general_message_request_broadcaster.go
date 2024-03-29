package core

import (
	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/siri/siri"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
)

type GeneralMessageRequestBroadcaster interface {
	Situations(*sxml.XMLGetGeneralMessage, *audit.BigQueryMessage) (*siri.SIRIGeneralMessageResponse, error)
}

type SIRIGeneralMessageRequestBroadcaster struct {
	connector
}

type SIRIGeneralMessageRequestBroadcasterFactory struct{}

func NewSIRIGeneralMessageRequestBroadcaster(partner *Partner) *SIRIGeneralMessageRequestBroadcaster {
	siriGeneralMessageRequestBroadcaster := &SIRIGeneralMessageRequestBroadcaster{}
	siriGeneralMessageRequestBroadcaster.partner = partner
	return siriGeneralMessageRequestBroadcaster
}

func (connector *SIRIGeneralMessageRequestBroadcaster) Situations(request *sxml.XMLGetGeneralMessage, message *audit.BigQueryMessage) (*siri.SIRIGeneralMessageResponse, error) {
	lineRefs := make(map[string]struct{})
	monitoringRefs := make(map[string]struct{})

	response := &siri.SIRIGeneralMessageResponse{
		Address:                   connector.Partner().Address(),
		ProducerRef:               connector.Partner().ProducerRef(),
		ResponseMessageIdentifier: connector.Partner().NewResponseMessageIdentifier(),
	}

	response.SIRIGeneralMessageDelivery = connector.getGeneralMessageDelivery(&request.XMLGeneralMessageRequest)

	if !response.SIRIGeneralMessageDelivery.Status {
		message.Status = "Error"
		message.ErrorDetails = response.SIRIGeneralMessageDelivery.ErrorString()
	}

	for _, message := range response.GeneralMessages {
		for _, affectedRef := range message.AffectedRefs {
			switch affectedRef.Kind {
			case "LineRef":
				lineRefs[affectedRef.Id] = struct{}{}
			case "StopPointRef", "DestinationRef":
				monitoringRefs[affectedRef.Id] = struct{}{}
			}
		}
		for _, affectedLineSection := range message.LineSections {
			lineRefs[affectedLineSection.LineRef] = struct{}{}
			monitoringRefs[affectedLineSection.FirstStop] = struct{}{}
			monitoringRefs[affectedLineSection.LastStop] = struct{}{}
		}
	}

	message.Lines = GetModelReferenceSlice(lineRefs)
	message.StopAreas = GetModelReferenceSlice(monitoringRefs)
	message.RequestIdentifier = request.MessageIdentifier()
	message.ResponseIdentifier = response.ResponseMessageIdentifier

	return response, nil
}

func (connector *SIRIGeneralMessageRequestBroadcaster) getGeneralMessageDelivery(request *sxml.XMLGeneralMessageRequest) siri.SIRIGeneralMessageDelivery {
	delivery := siri.SIRIGeneralMessageDelivery{
		RequestMessageRef: request.MessageIdentifier(),
		Status:            true,
		ResponseTimestamp: connector.Clock().Now(),
	}

	builder := NewBroadcastGeneralMessageBuilder(connector.Partner(), SIRI_GENERAL_MESSAGE_REQUEST_BROADCASTER)
	builder.InfoChannelRef = request.InfoChannelRef()
	builder.SetLineRef(request.LineRef())
	builder.SetStopPointRef(request.StopPointRef())

	situations := connector.partner.Model().Situations().FindAll()
	for i := range situations {
		siriGeneralMessage := builder.BuildGeneralMessage(situations[i])
		if siriGeneralMessage == nil {
			continue
		}
		delivery.GeneralMessages = append(delivery.GeneralMessages, siriGeneralMessage)
	}

	return delivery
}

func (factory *SIRIGeneralMessageRequestBroadcasterFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteCodeSpace()
	apiPartner.ValidatePresenceOfLocalCredentials()
}

func (factory *SIRIGeneralMessageRequestBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRIGeneralMessageRequestBroadcaster(partner)
}
