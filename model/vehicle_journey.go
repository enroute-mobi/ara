package model

import (
	"encoding/json"
	"fmt"
)

type VehicleJourneyId string

type VehicleJourneyAttributes struct {
	ObjectId     ObjectID
	LineObjectId ObjectID
	Attributes   Attributes
	References   References
}

type VehicleJourney struct {
	ObjectIDConsumer
	model Model

	id VehicleJourneyId

	LineId LineId `json:",omitempty"`
	Name   string `json:",omitempty"`

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

func (vehicleJourney *VehicleJourney) Id() VehicleJourneyId {
	return vehicleJourney.id
}

func (vehicleJourney *VehicleJourney) Line() *Line {
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
		ObjectIDs  ObjectIDs     `json:",omitempty"`
		StopVisits []StopVisitId `json:",omitempty"`
		Attributes Attributes    `json:",omitempty"`
		References References    `json:",omitempty"`
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
		aux.References = vehicleJourney.References
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
	value, present := vehicleJourney.References[key]
	return value, present
}

func (vehicleJourney *VehicleJourney) UnmarshalJSON(data []byte) error {
	type Alias VehicleJourney
	aux := &struct {
		ObjectIDs map[string]string
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
	return nil
}

func (vehicleJourney *VehicleJourney) Save() (ok bool) {
	ok = vehicleJourney.model.VehicleJourneys().Save(vehicleJourney)
	return
}

type MemoryVehicleJourneys struct {
	UUIDConsumer

	model Model

	byIdentifier map[VehicleJourneyId]*VehicleJourney
}

type VehicleJourneys interface {
	UUIDInterface

	New() VehicleJourney
	Find(id VehicleJourneyId) (VehicleJourney, bool)
	FindByObjectId(objectid ObjectID) (VehicleJourney, bool)
	FindByLineId(id LineId) []VehicleJourney
	FindAll() []VehicleJourney
	Save(vehicleJourney *VehicleJourney) bool
	Delete(vehicleJourney *VehicleJourney) bool
}

func NewMemoryVehicleJourneys() *MemoryVehicleJourneys {
	return &MemoryVehicleJourneys{
		byIdentifier: make(map[VehicleJourneyId]*VehicleJourney),
	}
}

func (manager *MemoryVehicleJourneys) New() VehicleJourney {
	vehicleJourney := NewVehicleJourney(manager.model)
	return *vehicleJourney
}

func (manager *MemoryVehicleJourneys) Find(id VehicleJourneyId) (VehicleJourney, bool) {
	vehicleJourney, ok := manager.byIdentifier[id]
	if ok {
		return *vehicleJourney, true
	} else {
		return VehicleJourney{}, false
	}
}

func (manager *MemoryVehicleJourneys) FindByObjectId(objectid ObjectID) (VehicleJourney, bool) {
	for _, vehicleJourney := range manager.byIdentifier {
		vehicleJourneyObjectId, _ := vehicleJourney.ObjectID(objectid.Kind())
		if vehicleJourneyObjectId.Value() == objectid.Value() {
			return *vehicleJourney, true
		}
	}
	return VehicleJourney{}, false
}

func (manager *MemoryVehicleJourneys) FindByLineId(id LineId) (vehicleJourneys []VehicleJourney) {
	for _, vehicleJourney := range manager.byIdentifier {
		if vehicleJourney.LineId == id {
			vehicleJourneys = append(vehicleJourneys, *vehicleJourney)
		}
	}
	return
}

func (manager *MemoryVehicleJourneys) FindAll() (vehicleJourneys []VehicleJourney) {
	if len(manager.byIdentifier) == 0 {
		return []VehicleJourney{}
	}
	for _, vehicleJourney := range manager.byIdentifier {
		vehicleJourneys = append(vehicleJourneys, *vehicleJourney)
	}
	return
}

func (manager *MemoryVehicleJourneys) Save(vehicleJourney *VehicleJourney) bool {
	if vehicleJourney.Id() == "" {
		vehicleJourney.id = VehicleJourneyId(manager.NewUUID())
	}
	vehicleJourney.model = manager.model
	manager.byIdentifier[vehicleJourney.Id()] = vehicleJourney
	return true
}

func (manager *MemoryVehicleJourneys) Delete(vehicleJourney *VehicleJourney) bool {
	delete(manager.byIdentifier, vehicleJourney.Id())
	return true
}

func (manager *MemoryVehicleJourneys) Load(referentialId string) error {
	var selectVehicleJourneys []SelectVehicleJourney
	modelName := manager.model.Date()
	sqlQuery := fmt.Sprintf("select * from vehicle_journeys where referential_id = '%s' and model_name = '%s'", referentialId, modelName.String())
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

		if vj.Attributes.Valid && len(vj.Attributes.String) > 0 {
			if err = json.Unmarshal([]byte(vj.Attributes.String), &vehicleJourney.Attributes); err != nil {
				return err
			}
		}

		if vj.References.Valid && len(vj.References.String) > 0 {
			if err = json.Unmarshal([]byte(vj.References.String), &vehicleJourney.References); err != nil {
				return err
			}
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
