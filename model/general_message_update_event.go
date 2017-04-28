package model

import "time"

type SituationUpdateRequestId string

type SituationAttributes struct {
	Format     string
	Channel    string
	References References
	Messages   []*Message
	ValidUntil time.Time
}

type SituationUpdateEvent struct {
	id                  SituationUpdateRequestId
	CreatedAt           time.Time
	RecordedAt          time.Time
	SituationObjectID   ObjectID
	Version             int64
	ProducerRef         string
	SituationAttributes SituationAttributes
}

func (event *SituationUpdateEvent) Id() SituationUpdateRequestId {
	return event.id
}

func (event *SituationUpdateEvent) SetId(id SituationUpdateRequestId) {
	event.id = id
}
