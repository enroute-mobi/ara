package model

type StopVisitBroadcastEvent struct {
	Id             StopVisitId
	StopAreaId     StopAreaId
	SubscriptionId string
	Schedules      StopVisitSchedules
}
