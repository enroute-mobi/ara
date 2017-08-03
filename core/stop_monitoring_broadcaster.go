package core

import (
	"fmt"
	"time"

	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type SIRIStopMonitoringBroadcaster interface {
	model.Stopable

	Run()
}

type SMBroadcaster struct {
	model.ClockConsumer

	connector *SIRIStopMonitoringSubscriptionBroadcaster
}

type StopMonitoringBroadcaster struct {
	SMBroadcaster

	stop chan struct{}
}

type FakeStopMonitoringBroadcaster struct {
	SMBroadcaster

	model.ClockConsumer
}

func NewFakeStopMonitoringBroadcaster(connector *SIRIStopMonitoringSubscriptionBroadcaster) SIRIStopMonitoringBroadcaster {
	broadcaster := &FakeStopMonitoringBroadcaster{}
	broadcaster.connector = connector
	return broadcaster
}

func (broadcaster *FakeStopMonitoringBroadcaster) Run() {
	broadcaster.prepareSIRIStopMonitoringNotify()
}

func (broadcaster *FakeStopMonitoringBroadcaster) Stop() {}

func NewSIRIStopMonitoringBroadcaster(connector *SIRIStopMonitoringSubscriptionBroadcaster) SIRIStopMonitoringBroadcaster {
	broadcaster := &StopMonitoringBroadcaster{}
	broadcaster.connector = connector

	return broadcaster
}

func (smb *StopMonitoringBroadcaster) Run() {
	if smb.stop != nil {
		return
	}

	logger.Log.Debugf("Start StopMonitoringBroadcaster")

	smb.stop = make(chan struct{})
	go smb.run()
}

func (smb *StopMonitoringBroadcaster) run() {
	c := smb.Clock().After(5 * time.Second)

	for {
		select {
		case <-smb.stop:
			return
		case <-c:
			logger.Log.Debugf("SIRIStopMonitoringBroadcaster visit")

			smb.prepareSIRIStopMonitoringNotify()

			c = smb.Clock().After(5 * time.Second)
		}
	}
}

func (smb *StopMonitoringBroadcaster) Stop() {
	if smb.stop != nil {
		close(smb.stop)
	}
}

func (smb *SMBroadcaster) RemoteObjectIDKind() string {
	if smb.connector.partner.Setting("siri-stop-monitoring-request-broadcaster.remote_objectid_kind") != "" {
		return smb.connector.partner.Setting("siri-stop-monitoring-request-broadcaster.remote_objectid_kind")
	}
	return smb.connector.partner.Setting("remote_objectid_kind")
}

func (smb *SMBroadcaster) prepareSIRIStopMonitoringNotify() {
	smb.connector.mutex.Lock()

	events := smb.connector.events
	smb.connector.events = make(map[SubscriptionId][]*model.StopVisitBroadcastEvent)

	smb.connector.mutex.Unlock()

	tx := smb.connector.Partner().Referential().NewTransaction()
	defer tx.Close()

	for key, sub := range events {

		//Voir pour le RequestMessageRef

		delivery := siri.SIRINotifyStopMonitoring{
			ProducerRef:               smb.connector.Partner().Setting("remote_credential"),
			ResponseMessageIdentifier: smb.connector.NewUUID(),
			SubscriberRef:             smb.connector.SIRIPartner().RequestorRef(),
			SubscriptionIdentifier:    fmt.Sprintf("Edwig:Subscription::%v:LOC", key),
			ResponseTimestamp:         smb.connector.Clock().Now(),

			Status: true,
		}

		svIds := make(map[string]bool) //Making sure not to send 2 times the same SV

		for _, event := range sub {
			for _, sm := range smb.handleEvent(event, tx) {
				if _, ok := svIds[sm.ItemIdentifier]; ok {
					continue
				}
				delivery.MonitoredStopVisits = append(delivery.MonitoredStopVisits, sm)
			}
		}
		smb.connector.SIRIPartner().SOAPClient().NotifyStopMonitoring(&delivery)
	}
}

func (smb *SMBroadcaster) handleEvent(event *model.StopVisitBroadcastEvent, tx *model.Transaction) []*siri.SIRIMonitoredStopVisit {
	switch event.ModelType {
	case "StopVisit":
		sv := smb.getMonitoredStopVisit(model.StopVisitId(event.ModelId), tx)
		//Update Subscription status
		return []*siri.SIRIMonitoredStopVisit{sv}
	default:
		return nil
	}
}

func (smb *SMBroadcaster) getMonitoredStopVisit(sv model.StopVisitId, tx *model.Transaction) *siri.SIRIMonitoredStopVisit {
	stopVisit, ok := smb.connector.Partner().Model().StopVisits().Find(model.StopVisitId(sv))
	if !ok {
		return nil
	}

	stopArea, ok := smb.connector.Partner().Model().StopAreas().Find(stopVisit.StopAreaId)
	if !ok {
		return nil
	}
	objectidKind := smb.connector.Partner().Setting("remote_objectid_kind")

	objectid, ok := stopArea.ObjectID(objectidKind)
	if !ok {
		return nil
	}

	var itemIdentifier string
	stopVisitId, ok := stopVisit.ObjectID(objectidKind)
	if ok {
		itemIdentifier = stopVisitId.Value()
	} else {
		defaultObjectID, ok := stopVisit.ObjectID("_default")
		if !ok {
			return nil
		}
		itemIdentifier = fmt.Sprintf("RATPDev:Item::%s:LOC", defaultObjectID.HashValue())
	}

	schedules := stopVisit.Schedules

	vehicleJourney := stopVisit.VehicleJourney()
	if vehicleJourney == nil {
		logger.Log.Printf("Ignore StopVisit %s without Vehiclejourney", stopVisit.Id())
		return nil
	}
	line := vehicleJourney.Line()
	if line == nil {
		logger.Log.Printf("Ignore StopVisit %s without Line", stopVisit.Id())
		return nil
	}
	if _, ok := line.ObjectID(objectidKind); !ok {
		logger.Log.Printf("Ignore StopVisit %s without Line without correct ObjectID", stopVisit.Id())
		return nil
	}

	vehicleJourneyId, ok := vehicleJourney.ObjectID(objectidKind)
	var dataVehicleJourneyRef string
	if ok {
		dataVehicleJourneyRef = vehicleJourneyId.Value()
	} else {
		defaultObjectID, ok := vehicleJourney.ObjectID("_default")
		if !ok {
			return nil
		}
		dataVehicleJourneyRef = fmt.Sprintf("RATPDev:VehicleJourney::%s:LOC", defaultObjectID.HashValue())
	}

	modelDate := tx.Model().Date()

	lineObjectId, _ := line.ObjectID(objectidKind)

	stopPointRef, _ := tx.Model().StopAreas().Find(stopVisit.StopAreaId)
	stopPointRefObjectId, _ := stopPointRef.ObjectID(objectidKind)

	monitoredStopVisit := &siri.SIRIMonitoredStopVisit{
		ItemIdentifier: itemIdentifier,
		MonitoringRef:  objectid.Value(),
		StopPointRef:   stopPointRefObjectId.Value(),
		StopPointName:  stopPointRef.Name,

		VehicleJourneyName:     vehicleJourney.Name,
		LineRef:                lineObjectId.Value(),
		DatedVehicleJourneyRef: dataVehicleJourneyRef,
		DataFrameRef:           fmt.Sprintf("RATPDev:DataFrame::%s:LOC", modelDate.String()),
		RecordedAt:             stopVisit.RecordedAt,
		PublishedLineName:      line.Name,
		DepartureStatus:        string(stopVisit.DepartureStatus),
		ArrivalStatus:          string(stopVisit.ArrivalStatus),
		Order:                  stopVisit.PassageOrder,
		VehicleAtStop:          stopVisit.VehicleAtStop,
		Attributes:             make(map[string]map[string]string),
		References:             make(map[string]map[string]model.Reference),
	}

	monitoredStopVisit.AimedArrivalTime = schedules.Schedule(model.STOP_VISIT_SCHEDULE_AIMED).ArrivalTime()
	monitoredStopVisit.ExpectedArrivalTime = schedules.Schedule(model.STOP_VISIT_SCHEDULE_EXPECTED).ArrivalTime()
	monitoredStopVisit.ActualArrivalTime = schedules.Schedule(model.STOP_VISIT_SCHEDULE_ACTUAL).ArrivalTime()

	monitoredStopVisit.AimedDepartureTime = schedules.Schedule(model.STOP_VISIT_SCHEDULE_AIMED).DepartureTime()
	monitoredStopVisit.ExpectedDepartureTime = schedules.Schedule(model.STOP_VISIT_SCHEDULE_EXPECTED).DepartureTime()
	monitoredStopVisit.ActualDepartureTime = schedules.Schedule(model.STOP_VISIT_SCHEDULE_ACTUAL).DepartureTime()

	vehicleJourneyRefCopy := vehicleJourney.References.Copy()
	stopVisitRefCopy := stopVisit.References.Copy()

	smb.resolveVehiculeJourneyReferences(vehicleJourneyRefCopy, tx.Model().StopAreas())
	smb.resolveOperator(stopVisitRefCopy)

	smb.reformatReferences(vehicleJourney.ToFormat(), vehicleJourneyRefCopy)

	monitoredStopVisit.Attributes["StopVisitAttributes"] = stopVisit.Attributes
	monitoredStopVisit.References["StopVisitReferences"] = stopVisitRefCopy

	monitoredStopVisit.Attributes["VehicleJourneyAttributes"] = vehicleJourney.Attributes
	monitoredStopVisit.References["VehicleJourney"] = vehicleJourneyRefCopy

	return monitoredStopVisit
}

func (connector *SMBroadcaster) resolveVehiculeJourneyReferences(references model.References, manager model.StopAreas) {
	toResolve := []string{"PlaceRef", "OriginRef", "DestinationRef"}

	for _, ref := range toResolve {
		if references[ref] == (model.Reference{}) {
			continue
		}
		if foundStopArea, ok := manager.Find(model.StopAreaId(references[ref].Id)); ok {
			obj, ok := foundStopArea.ObjectID(connector.RemoteObjectIDKind())
			if ok {
				tmp := references[ref]
				tmp.ObjectId = &obj
				references[ref] = tmp
			}
		} else {
			tmp := references[ref]
			tmp.ObjectId.SetValue(tmp.Getformat(ref, tmp.GetSha1()))
		}
	}
}

func (smb *SMBroadcaster) resolveOperator(references model.References) {
	operatorRef, _ := references["OperatorRef"]
	operator, ok := smb.connector.Partner().Model().Operators().Find(model.OperatorId(operatorRef.Id))
	if !ok {
		return
	}

	obj, ok := operator.ObjectID(smb.connector.Partner().Setting("remote_objectid_kind"))
	if !ok {
		return
	}
	references["OperatorRef"].ObjectId.SetValue(obj.Value())
}

func (connector *SMBroadcaster) reformatReferences(toReformat []string, references model.References) {
	for _, ref := range toReformat {
		if references[ref] != (model.Reference{}) {
			tmp := references[ref]
			tmp.ObjectId.SetValue(tmp.Getformat(ref, tmp.GetSha1()))
		}
	}
}
