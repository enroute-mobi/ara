package model

import "time"

type SituationUpdateRequestId string

type SituationUpdateEvent struct {
	CreatedAt           time.Time
	RecordedAt          time.Time
	SituationObjectID   ObjectID
	id                  SituationUpdateRequestId
	Origin              string
	ProducerRef         string
	SituationAttributes SituationAttributes
	ValidityPeriods     []*TimeRange
	Keywords            []string
	ReportType          ReportType
	Version             int
}

type SituationAttributes struct {
	Format       string
	References   []*Reference
	LineSections []*References
	Messages     []*Message
}

func (event *SituationUpdateEvent) Id() SituationUpdateRequestId {
	return event.id
}

func (event *SituationUpdateEvent) SetId(id SituationUpdateRequestId) {
	event.id = id
}
