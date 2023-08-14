package model

import "bitbucket.org/enroute-mobi/ara/clock"

type SituationUpdateManager struct {
	clock.ClockConsumer

	model Model
}

func NewSituationUpdateManager(model Model) func([]*SituationUpdateEvent) {
	manager := newSituationUpdateManager(model)
	return manager.Update
}

func newSituationUpdateManager(model Model) *SituationUpdateManager {
	return &SituationUpdateManager{model: model}
}

func (manager *SituationUpdateManager) Update(events []*SituationUpdateEvent) {
	for _, event := range events {
		situation, ok := manager.model.Situations().FindByObjectId(event.SituationObjectID)
		if ok && situation.RecordedAt == event.RecordedAt {
			continue
		}

		if !ok {
			situation = manager.model.Situations().New()
			situation.Origin = event.Origin
			situation.SetObjectID(event.SituationObjectID)
			situation.SetObjectID(NewObjectID("_default", event.SituationObjectID.HashValue()))
		}

		situation.RecordedAt = event.RecordedAt
		situation.Version = event.Version
		situation.ProducerRef = event.ProducerRef

		situation.References = event.SituationAttributes.References
		situation.LineSections = event.SituationAttributes.LineSections
		situation.Messages = event.SituationAttributes.Messages
		situation.ValidityPeriods = event.ValidityPeriods
		situation.Channel = event.SituationAttributes.Channel
		situation.Format = event.SituationAttributes.Format

		situation.Save()
	}
}
