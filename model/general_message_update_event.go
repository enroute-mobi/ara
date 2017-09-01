package model

import "time"

type SituationUpdateRequestId string

type SituationUpdateEvent struct {
	id                  SituationUpdateRequestId
	CreatedAt           time.Time
	RecordedAt          time.Time
	SituationObjectID   ObjectID
	Version             int
	ProducerRef         string
	SituationAttributes SituationAttributes
}

type SituationAttributes struct {
	Format       string
	Channel      string
	References   []*Reference
	LineSections []*References
	Messages     []*Message
	ValidUntil   time.Time
}

func (event *SituationUpdateEvent) Id() SituationUpdateRequestId {
	return event.id
}

func (event *SituationUpdateEvent) SetId(id SituationUpdateRequestId) {
	event.id = id
}
