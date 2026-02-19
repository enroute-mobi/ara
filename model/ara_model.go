package model

import (
	"bitbucket.org/enroute-mobi/ara/uuid"
)

type ModelId string

type ModelInstance interface {
	CodeConsumerInterface

	ModelId() ModelId
}

type ModelManager[Id, O any] interface {
	uuid.UUIDInterface
	New() O
	Find(Id) (O, bool)
	FindAll() []O
	Save(O) bool
	Delete(O) bool
	SetModel(Model)
}

type CodeValidator interface {
	CodeExists(Code) bool
}

type CodeHandler[O any] interface {
	CodeValidator
	FindByCode(Code) (O, bool)
}

type Loadable interface {
	Load(string) error
}

type Broadcaster[O any] interface {
	SetBroadcaster(func(O), ...string)
}

type Model interface {
	Date() Date
	SetDate(Date)
	Referential() string
	SetReferential(string)

	Lines() Lines
	LineGroups() LineGroups
	StopAreaGroups() StopAreaGroups
	Situations() Situations
	StopAreas() StopAreas
	StopVisits() StopVisits
	ScheduledStopVisits() StopVisits
	VehicleJourneys() VehicleJourneys
	Operators() Operators
	Vehicles() Vehicles
	Macros() Macros
	Controls() Controls
	Facilities() Facilities

	Load() error
	Reload() Model
	SetBroadcastSMChan(chan StopMonitoringBroadcastEvent)
	SetBroadcastGMChan(chan SituationBroadcastEvent)
	SetBroadcastSXChan(chan SituationBroadcastEvent)
	SetBroadcastVeChan(chan VehicleBroadcastEvent)
	SetBroadcastFMChan(chan FacilityBroadcastEvent)
}
