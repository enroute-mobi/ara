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
	Version             int
}

type SituationAttributes struct {
	ValidUntil   time.Time
	Format       string
	Channel      string
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
