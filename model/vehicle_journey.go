package model

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"bitbucket.org/enroute-mobi/ara/model/schedules"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

const (
	VEHICLE_DIRECTION_INBOUND  = "inbound"
	VEHICLE_DIRECTION_OUTBOUND = "outbound"
)

type VehicleJourneyId ModelId

type VehicleJourney struct {
	References References
	model      Model
	Attributes Attributes
	CodeConsumer
	LineId                  LineId `json:",omitempty"`
	Name                    string `json:",omitempty"`
	id                      VehicleJourneyId
	DestinationName         string `json:",omitempty"`
	Occupancy               string `json:",omitempty"`
	DirectionType           string `json:",omitempty"`
	Origin                  string `json:",omitempty"`
	OriginName              string `json:",omitempty"`
	Monitored               bool
	HasCompleteStopSequence bool
	Cancellation            bool
	DetailedStopVisits      []DetailedStopVisit `json:"DetailedStopVisits,omitempty"`
	AimedStopVisitCount     int                 `json:",omitempty"`
}

type DetailedStopVisit struct {
	Order           int                           `json:"Order"`
	StopAreaName    string                        `json:"StopAreaName"`
	StopAreaId      StopAreaId                    `json:"StopAreaId"`
	Schedules       []schedules.StopVisitSchedule `json:",omitempty"`
	ArrivalStatus   StopVisitArrivalStatus        `json:",omitempty"`
	DepartureStatus StopVisitDepartureStatus      `json:",omitempty"`
	CollectedAt     time.Time                     `json:",omitempty"`
}

func NewVehicleJourney(model Model) *VehicleJourney {
	vehicleJourney := &VehicleJourney{
		model:      model,
		Attributes: NewAttributes(),
		References: NewReferences(),
	}
	vehicleJourney.codes = make(Codes)
	return vehicleJourney
}

func (vehicleJourney *VehicleJourney) ModelId() ModelId {
	return ModelId(vehicleJourney.id)
}

func (vehicleJourney *VehicleJourney) copy() *VehicleJourney {
	vj := *vehicleJourney
	vj.Attributes = vehicleJourney.Attributes.Copy()
	vj.References = vehicleJourney.References.Copy()
	return &vj
}

func (vehicleJourney *VehicleJourney) Id() VehicleJourneyId {
	return vehicleJourney.id
}

func (vehicleJourney *VehicleJourney) Line() *Line {
	if vehicleJourney.model == nil {
		return nil
	}
	line, ok := vehicleJourney.model.Lines().Find(vehicleJourney.LineId)
	if !ok {
		return nil
	}
	return line
}

func (vehicleJourney *VehicleJourney) MarshalJSON() ([]byte, error) {
	type Alias VehicleJourney

	aux := struct {
		Codes      Codes                `json:",omitempty"`
		Attributes Attributes           `json:",omitempty"`
		References map[string]Reference `json:",omitempty"`
		*Alias
		Id                 VehicleJourneyId
		StopVisits         []StopVisitId       `json:",omitempty"`
		DetailedStopVisits []DetailedStopVisit `json:",omitempty"`
	}{
		Id:    vehicleJourney.id,
		Alias: (*Alias)(vehicleJourney),
	}

	if !vehicleJourney.Codes().Empty() {
		aux.Codes = vehicleJourney.Codes()
	}
	if !vehicleJourney.Attributes.IsEmpty() {
		aux.Attributes = vehicleJourney.Attributes
	}
	if !vehicleJourney.References.IsEmpty() {
		aux.References = vehicleJourney.References.GetReferences()
	}

	if len(vehicleJourney.DetailedStopVisits) != 0 {
		aux.DetailedStopVisits = vehicleJourney.DetailedStopVisits
	}

	svs := vehicleJourney.model.StopVisits().FindByVehicleJourneyId(vehicleJourney.id)
	var stopVisits []StopVisitId
	for i := range svs {
		stopVisits = append(stopVisits, svs[i].Id())
	}
	if len(stopVisits) > 0 {
		aux.StopVisits = stopVisits
	}

	return json.Marshal(&aux)
}

func (vehicleJourney *VehicleJourney) ToFormat() []string {
	return []string{"RouteRef", "JourneyPatternRef", "DatedVehicleJourneyRef"}
}

func (vehicleJourney *VehicleJourney) Attribute(key string) (string, bool) {
	value, present := vehicleJourney.Attributes[key]
	return value, present
}

func (vehicleJourney *VehicleJourney) Reference(key string) (Reference, bool) {
	value, present := vehicleJourney.References.Get(key)
	return value, present
}

func (vehicleJourney *VehicleJourney) UnmarshalJSON(data []byte) error {
	type Alias VehicleJourney
	aux := &struct {
		Codes      map[string]string
		References map[string]Reference
		*Alias
	}{
		Alias: (*Alias)(vehicleJourney),
	}
	err := json.Unmarshal(data, aux)
	if err != nil {
		return err
	}

	if aux.Codes != nil {
		vehicleJourney.CodeConsumer.codes = NewCodesFromMap(aux.Codes)
	}

	if aux.References != nil {
		vehicleJourney.References.SetReferences(aux.References)
	}

	return nil
}

func (vehicleJourney *VehicleJourney) GtfsDirectionId() *uint32 {
	var directionId uint32
	switch vehicleJourney.DirectionType {
	case VEHICLE_DIRECTION_OUTBOUND:
		directionId = uint32(0)
	case VEHICLE_DIRECTION_INBOUND:
		directionId = uint32(1)
	default:
		return nil
	}
	return &directionId
}

func (vehicleJourney *VehicleJourney) Save() bool {
	return vehicleJourney.model.VehicleJourneys().Save(vehicleJourney)
}

type MemoryVehicleJourneys struct {
	uuid.UUIDConsumer

	model Model

	mutex             *sync.RWMutex
	byIdentifier      map[VehicleJourneyId]*VehicleJourney
	byCode            *CodeIndex
	byLine            *Index
	byBroadcastedFull map[string]map[VehicleJourneyId]struct{}
}

type VehicleJourneys interface {
	uuid.UUIDInterface

	New() *VehicleJourney
	Find(VehicleJourneyId) (*VehicleJourney, bool)
	FindByCode(code Code) (*VehicleJourney, bool)
	FindByLineId(LineId) []*VehicleJourney
	FullVehicleJourneyExistBySubscriptionId(string, VehicleJourneyId) bool
	FindAll() []*VehicleJourney
	Save(*VehicleJourney) bool
	SetFullVehicleJourneyBySubscriptionId(string, VehicleJourneyId)
	Delete(*VehicleJourney) bool
	DeleteById(VehicleJourneyId) bool
}

func NewMemoryVehicleJourneys() *MemoryVehicleJourneys {
	extractor := func(instance ModelInstance) ModelId { return ModelId((instance.(*VehicleJourney)).LineId) }

	return &MemoryVehicleJourneys{
		mutex:             &sync.RWMutex{},
		byIdentifier:      make(map[VehicleJourneyId]*VehicleJourney),
		byCode:            NewCodeIndex(),
		byLine:            NewIndex(extractor),
		byBroadcastedFull: make(map[string]map[VehicleJourneyId]struct{}),
	}
}

func (manager *MemoryVehicleJourneys) New() *VehicleJourney {
	return NewVehicleJourney(manager.model)
}

func (manager *MemoryVehicleJourneys) SetFullVehicleJourneyBySubscriptionId(id string, vehicleJourneyId VehicleJourneyId) {
	manager.mutex.Lock()
	vjIds, ok := manager.byBroadcastedFull[id]
	if !ok {
		vjIds = make(map[VehicleJourneyId]struct{})
		manager.byBroadcastedFull[id] = vjIds
	}
	vjIds[vehicleJourneyId] = struct{}{}
	manager.mutex.Unlock()
}

func (manager *MemoryVehicleJourneys) FullVehicleJourneyExistBySubscriptionId(id string, vehicleJourneyId VehicleJourneyId) bool {
	manager.mutex.RLock()
	_, ok := manager.byBroadcastedFull[id][vehicleJourneyId]
	manager.mutex.RUnlock()

	return ok
}

func (manager *MemoryVehicleJourneys) TestLenFullVehicleJourneyBySubscriptionId() int {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()
	return len(manager.byBroadcastedFull)
}

func (manager *MemoryVehicleJourneys) Find(id VehicleJourneyId) (*VehicleJourney, bool) {
	manager.mutex.RLock()
	vehicleJourney, ok := manager.byIdentifier[id]
	manager.mutex.RUnlock()

	if ok {
		return vehicleJourney.copy(), true
	}
	return &VehicleJourney{}, false
}

func (manager *MemoryVehicleJourneys) FindByCode(code Code) (*VehicleJourney, bool) {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	id, ok := manager.byCode.Find(code)
	if ok {
		return manager.byIdentifier[VehicleJourneyId(id)].copy(), true
	}

	return &VehicleJourney{}, false
}

func (manager *MemoryVehicleJourneys) CodeExists(code Code) bool {
	manager.mutex.RLock()
	_, ok := manager.byCode.Find(code)
	manager.mutex.RUnlock()

	return ok
}

func (manager *MemoryVehicleJourneys) FindByLineId(id LineId) (vehicleJourneys []*VehicleJourney) {
	manager.mutex.RLock()

	ids, _ := manager.byLine.Find(ModelId(id))

	for _, id := range ids {
		vj := manager.byIdentifier[VehicleJourneyId(id)]
		vehicleJourneys = append(vehicleJourneys, vj.copy())
	}

	manager.mutex.RUnlock()
	return
}

func (manager *MemoryVehicleJourneys) FindAll() (vehicleJourneys []*VehicleJourney) {
	manager.mutex.RLock()

	for _, vehicleJourney := range manager.byIdentifier {
		vehicleJourneys = append(vehicleJourneys, vehicleJourney.copy())
	}

	manager.mutex.RUnlock()
	return
}

func (manager *MemoryVehicleJourneys) Save(vehicleJourney *VehicleJourney) bool {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	if vehicleJourney.Id() == "" {
		vehicleJourney.id = VehicleJourneyId(manager.NewUUID())
	}

	vehicleJourney.model = manager.model
	manager.byIdentifier[vehicleJourney.Id()] = vehicleJourney
	manager.byCode.Index(vehicleJourney)
	manager.byLine.Index(vehicleJourney)

	return true
}

func (manager *MemoryVehicleJourneys) Delete(vehicleJourney *VehicleJourney) bool {
	return manager.DeleteById(vehicleJourney.id)
}

func (manager *MemoryVehicleJourneys) DeleteById(id VehicleJourneyId) bool {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	delete(manager.byIdentifier, id)
	manager.byCode.Delete(ModelId(id))
	manager.byLine.Delete(ModelId(id))
	for subscriptionId, vehicleJourneyIds := range manager.byBroadcastedFull {
		delete(vehicleJourneyIds, id)
		if len(vehicleJourneyIds) == 0 {
			delete(manager.byBroadcastedFull, subscriptionId)
		}
	}
	return true
}

func (manager *MemoryVehicleJourneys) Load(referentialSlug string) error {
	var selectVehicleJourneys []SelectVehicleJourney
	modelName := manager.model.Date()
	sqlQuery := fmt.Sprintf("select * from vehicle_journeys where referential_slug = '%s' and model_name = '%s'", referentialSlug, modelName.String())
	_, err := Database.Select(&selectVehicleJourneys, sqlQuery)
	if err != nil {
		return err
	}
	for _, vj := range selectVehicleJourneys {
		vehicleJourney := manager.New()
		vehicleJourney.id = VehicleJourneyId(vj.Id)
		if vj.Name.Valid {
			vehicleJourney.Name = vj.Name.String
		}
		if vj.LineId.Valid {
			vehicleJourney.LineId = LineId(vj.LineId.String)
		}
		if vj.OriginName.Valid {
			vehicleJourney.OriginName = vj.OriginName.String
		}
		if vj.DestinationName.Valid {
			vehicleJourney.DestinationName = vj.DestinationName.String
		}
		if vj.DirectionType.Valid {
			vehicleJourney.DirectionType = vj.DirectionType.String
		}

		if vj.AimedStopVisitCount.Valid {
			vehicleJourney.AimedStopVisitCount = int(vj.AimedStopVisitCount.Int64)
		}

		if vj.Attributes.Valid && len(vj.Attributes.String) > 0 {
			if err = json.Unmarshal([]byte(vj.Attributes.String), &vehicleJourney.Attributes); err != nil {
				return err
			}
		}

		if vj.References.Valid && len(vj.References.String) > 0 {
			references := make(map[string]Reference)
			if err = json.Unmarshal([]byte(vj.References.String), &references); err != nil {
				return err
			}
			vehicleJourney.References.SetReferences(references)
		}

		if vj.Codes.Valid && len(vj.Codes.String) > 0 {
			codeMap := make(map[string]string)
			if err = json.Unmarshal([]byte(vj.Codes.String), &codeMap); err != nil {
				return err
			}
			vehicleJourney.codes = NewCodesFromMap(codeMap)
		}

		manager.Save(vehicleJourney)
	}
	return nil
}
