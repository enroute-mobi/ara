package core

import (
	"time"

	"github.com/af83/edwig/model"
)

type StopMonitoringSubscriptionCollector interface {
	RequestStopAreaUpdate(request *StopAreaUpdateRequest)
}

type SIRIStopMonitoringSubscriptionCollector struct {
	model.ClockConsumer
	model.UUIDConsumer

	siriConnector
}

type SIRIStopMonitoringSubscriptionCollectorFactory struct{}

func NewSIRIStopMonitoringSubscriptionCollector(partner *Partner) *SIRIStopMonitoringSubscriptionCollector {
	siriStopMonitoringSubscriptionCollector := &SIRIStopMonitoringSubscriptionCollector{}
	siriStopMonitoringSubscriptionCollector.partner = partner

	return siriStopMonitoringSubscriptionCollector
}

func (factory *SIRIStopMonitoringSubscriptionCollectorFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRIStopMonitoringSubscriptionCollector(partner)
}

func (factory *SIRIStopMonitoringSubscriptionCollectorFactory) Validate(apiPartner *APIPartner) bool {
	ok := apiPartner.ValidatePresenceOfSetting("remote_objectid_kind")
	ok = ok && apiPartner.ValidatePresenceOfSetting("local_credential")
	return ok
}

func (connector *SIRIStopMonitoringSubscriptionCollector) RequestStopAreaUpdate(request *StopAreaUpdateRequest) {
	subscription := connector.partner.Subscriptions().FindOrCreateByKind("StopMonitoring")

	for _, sr := range subscription.resourcesByObjectID {
		if sr.Reference.ObjectId.Value() == string(request.StopAreaId()) {
			sr.SubscribedUntil = sr.SubscribedUntil.Add(1 * time.Minute)
			return
		}
	}

	objId := model.NewObjectID(connector.partner.Setting("remote_objectid_kind"), string(request.StopAreaId()))
	ref := model.Reference{
		ObjectId: &objId,
		Id:       string(request.StopAreaId()),
		Type:     "StopArea",
	}

	subscription.CreateAddNewResource(ref)
}
