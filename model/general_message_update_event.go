package model

import "time"

type SituationUpdateRequestId string

type SituationUpdateEvent struct {
	CreatedAt         time.Time
	RecordedAt        time.Time
	SituationObjectID ObjectID
	id                SituationUpdateRequestId
	Origin            string
	Format            string
	ProducerRef       string
	ValidityPeriods   []*TimeRange
	Keywords          []string
	ReportType        ReportType
	Version           int
	Summary           *SituationTranslatedString
	Description       *SituationTranslatedString
	Affects           []Affect
}

func (event *SituationUpdateEvent) Id() SituationUpdateRequestId {
	return event.id
}

func (event *SituationUpdateEvent) SetId(id SituationUpdateRequestId) {
	event.id = id
}
