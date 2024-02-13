package model

import (
	"bitbucket.org/enroute-mobi/ara/clock"
)

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
		situation, ok := manager.model.Situations().FindByCode(event.SituationCode)
		if ok &&
			situation.RecordedAt == event.RecordedAt &&
			situation.Version == event.Version {
			continue
		}

		if !ok {
			situation = manager.model.Situations().New()
			situation.Origin = event.Origin
			situation.SetCode(event.SituationCode)
			situation.SetCode(NewCode("_default", event.SituationCode.HashValue()))
		}

		situation.RecordedAt = event.RecordedAt
		situation.Version = event.Version
		situation.ProducerRef = event.ProducerRef
		situation.ParticipantRef = event.ParticipantRef

		situation.Summary = event.Summary
		situation.Description = event.Description

		situation.VersionedAt = event.VersionedAt
		situation.ValidityPeriods = event.ValidityPeriods
		situation.PublicationWindows = event.PublicationWindows
		situation.Keywords = event.Keywords
		situation.ReportType = event.ReportType
		situation.AlertCause = event.AlertCause
		situation.Severity = event.Severity
		situation.Progress = event.Progress
		situation.Format = event.Format
		situation.Affects = event.Affects
		situation.Consequences = event.Consequences

		situation.Save()
	}
}
