package model

import "time"

type SituationUpdateRequestId string

type SituationUpdateEvent struct {
	CreatedAt          time.Time
	RecordedAt         time.Time
	SituationCode      Code
	id                 SituationUpdateRequestId
	Origin             string
	Format             string
	ProducerRef        string
	ParticipantRef     string
	VersionedAt        time.Time
	ValidityPeriods    []*TimeRange
	PublicationWindows []*TimeRange
	Keywords           []string
	AlertCause         SituationAlertCause
	Progress           SituationProgress
	ReportType         ReportType
	Version            int
	Summary            *SituationTranslatedString
	Description        *SituationTranslatedString
	Affects            []Affect
}

func (event *SituationUpdateEvent) Id() SituationUpdateRequestId {
	return event.id
}

func (event *SituationUpdateEvent) SetId(id SituationUpdateRequestId) {
	event.id = id
}

func (event *SituationUpdateEvent) TestFindAffectByLineId(lineId LineId) (bool, *AffectedLine) {
	for _, affect := range event.Affects {
		if affect.GetType() == SituationTypeLine &&
			affect.GetId() == ModelId(lineId) {
			return true, affect.(*AffectedLine)
		}
	}
	return false, nil
}

func (event *SituationUpdateEvent) TestFindAffectByStopAreaId(stopAreaId StopAreaId) (bool, *AffectedStopArea) {
	for _, affect := range event.Affects {
		if affect.GetType() == SituationTypeStopArea &&
			affect.GetId() == ModelId(stopAreaId) {
			return true, affect.(*AffectedStopArea)
		}
	}
	return false, nil
}
