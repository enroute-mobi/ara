package model

import (
	"encoding/json"
	"fmt"
	"sync"

	"bitbucket.org/enroute-mobi/ara/uuid"
)

type VehicleJourneyId ModelId

type VehicleJourney struct {
	ObjectIDConsumer

	model  Model
	Origin string `json:",omitempty"`

	id VehicleJourneyId

	LineId          LineId `json:",omitempty"`
	Name            string `json:",omitempty"`
	OriginName      string `json:",omitempty"`
	DestinationName string `json:",omitempty"`

	Monitored bool

	Attributes Attributes
	References References
}

func NewVehicleJourney(model Model) *VehicleJourney {
	vehicleJourney := &VehicleJourney{
		model:      model,
		Attributes: NewAttributes(),
		References: NewReferences(),
	}
	vehicleJourney.objectids = make(ObjectIDs)
	return vehicleJourney
}

func (vehicleJourney *VehicleJourney) modelId() ModelId {
	return ModelId(vehicleJourney.id)
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
	return &line
}

func (vehicleJourney *VehicleJourney) MarshalJSON() ([]byte, error) {
	type Alias VehicleJourney
	aux := struct {
		Id         VehicleJourneyId
		ObjectIDs  ObjectIDs            `json:",omitempty"`
		StopVisits []StopVisitId        `json:",omitempty"`
		Attributes Attributes           `json:",omitempty"`
		References map[string]Reference `json:",omitempty"`
		*Alias
	}{
		Id:    vehicleJourney.id,
		Alias: (*Alias)(vehicleJourney),
	}

	if !vehicleJourney.ObjectIDs().Empty() {
		aux.ObjectIDs = vehicleJourney.ObjectIDs()
	}
	if !vehicleJourney.Attributes.IsEmpty() {
		aux.Attributes = vehicleJourney.Attributes
	}
	if !vehicleJourney.References.IsEmpty() {
		aux.References = vehicleJourney.References.GetReferences()
	}

	stopVisitIds := []StopVisitId{}
	for _, stopVisit := range vehicleJourney.model.StopVisits().FindByVehicleJourneyId(vehicleJourney.id) {
		stopVisitIds = append(stopVisitIds, stopVisit.Id())
	}
	if len(stopVisitIds) > 0 {
		aux.StopVisits = stopVisitIds
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
		ObjectIDs  map[string]string
		References map[string]Reference
		*Alias
	}{
		Alias: (*Alias)(vehicleJourney),
	}
	err := json.Unmarshal(data, aux)
	if err != nil {
		return err
	}

	if aux.ObjectIDs != nil {
		vehicleJourney.ObjectIDConsumer.objectids = NewObjectIDsFromMap(aux.ObjectIDs)
	}

	if aux.References != nil {
		vehicleJourney.References.SetReferences(aux.References)
	}

	return nil
}

func (vehicleJourney *VehicleJourney) Save() (ok bool) {
	ok = vehicleJourney.model.VehicleJourneys().Save(vehicleJourney)
	return
}

type MemoryVehicleJourneys struct {
	uuid.UUIDConsumer

	model Model

	mutex        *sync.RWMutex
	byIdentifier map[VehicleJourneyId]*VehicleJourney
	byObjectId   *ObjectIdIndex
	byLine       *Index
}

type VehicleJourneys interface {
	uuid.UUIDInterface

	New() VehicleJourney
	Find(id VehicleJourneyId) (VehicleJourney, bool)
	FindByObjectId(objectid ObjectID) (VehicleJourney, bool)
	FindByLineId(id LineId) []VehicleJourney
	FindAll() []VehicleJourney
	Save(vehicleJourney *VehicleJourney) bool
	Delete(vehicleJourney *VehicleJourney) bool
}

func NewMemoryVehicleJourneys() *MemoryVehicleJourneys {
	extractor := func(instance ModelInstance) ModelId { return ModelId((instance.(*VehicleJourney)).LineId) }

	return &MemoryVehicleJourneys{
		mutex:        &sync.RWMutex{},
		byIdentifier: make(map[VehicleJourneyId]*VehicleJourney),
		byObjectId:   NewObjectIdIndex(),
		byLine:       NewIndex(extractor),
	}
}

func (manager *MemoryVehicleJourneys) New() VehicleJourney {
	vehicleJourney := NewVehicleJourney(manager.model)
	return *vehicleJourney
}

func (manager *MemoryVehicleJourneys) Find(id VehicleJourneyId) (VehicleJourney, bool) {
	manager.mutex.RLock()

	vehicleJourney, ok := manager.byIdentifier[id]
	if ok {
		manager.mutex.RUnlock()
		return *vehicleJourney, true
	} else {
		manager.mutex.RUnlock()
		return VehicleJourney{}, false
	}
}

func (manager *MemoryVehicleJourneys) FindByObjectId(objectid ObjectID) (VehicleJourney, bool) {
	manager.mutex.RLock()

	id, ok := manager.byObjectId.Find(objectid)
	if ok {
		manager.mutex.RUnlock()
		return *manager.byIdentifier[VehicleJourneyId(id)], true
	}

	manager.mutex.RUnlock()
	return VehicleJourney{}, false
}

func (manager *MemoryVehicleJourneys) FindByLineId(id LineId) (vehicleJourneys []VehicleJourney) {
	manager.mutex.RLock()

	ids, _ := manager.byLine.Find(ModelId(id))

	for _, id := range ids {
		vj, _ := manager.byIdentifier[VehicleJourneyId(id)]
		vehicleJourneys = append(vehicleJourneys, *vj)
	}

	manager.mutex.RUnlock()
	return
}

func (manager *MemoryVehicleJourneys) FindAll() (vehicleJourneys []VehicleJourney) {
	manager.mutex.RLock()

	for _, vehicleJourney := range manager.byIdentifier {
		vehicleJourneys = append(vehicleJourneys, *vehicleJourney)
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
	manager.byObjectId.Index(vehicleJourney)
	manager.byLine.Index(vehicleJourney)

	return true
}

func (manager *MemoryVehicleJourneys) Delete(vehicleJourney *VehicleJourney) bool {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	delete(manager.byIdentifier, vehicleJourney.id)
	manager.byObjectId.Delete(ModelId(vehicleJourney.id))
	manager.byLine.Delete(ModelId(vehicleJourney.id))

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

		if vj.ObjectIDs.Valid && len(vj.ObjectIDs.String) > 0 {
			objectIdMap := make(map[string]string)
			if err = json.Unmarshal([]byte(vj.ObjectIDs.String), &objectIdMap); err != nil {
				return err
			}
			vehicleJourney.objectids = NewObjectIDsFromMap(objectIdMap)
		}

		manager.Save(&vehicleJourney)
	}
	return nil
}
