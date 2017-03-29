package model

import "encoding/json"

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

	LineId     LineId
	Name       string
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

func (vehicleJourney *VehicleJourney) FillVehicleJourney(vehicleJourneyMap map[string]interface{}) {

	stopVisitIds := []StopVisitId{}
	for _, stopVisit := range vehicleJourney.model.StopVisits().FindByVehicleJourneyId(vehicleJourney.id) {
		stopVisitIds = append(stopVisitIds, stopVisit.Id())
	}

	if len(stopVisitIds) > 0 {
		vehicleJourneyMap["StopVisits"] = stopVisitIds
	}

	if vehicleJourney.id != "" {
		vehicleJourneyMap["Id"] = vehicleJourney.id
	}

	if vehicleJourney.LineId != "" {
		vehicleJourneyMap["LineId"] = vehicleJourney.LineId
	}

	if vehicleJourney.Name != "" {
		vehicleJourneyMap["Name"] = vehicleJourney.Name
	}

	if !vehicleJourney.Attributes.IsEmpty() {
		vehicleJourneyMap["Attributes"] = vehicleJourney.Attributes
	}

	if !vehicleJourney.References.IsEmpty() {
		vehicleJourneyMap["References"] = vehicleJourney.References
	}

	if !vehicleJourney.ObjectIDs().Empty() {
		vehicleJourneyMap["ObjectIDs"] = vehicleJourney.ObjectIDs()
	}
}

func (vehicleJourney *VehicleJourney) MarshalJSON() ([]byte, error) {
	stopVisitIds := []StopVisitId{}
	for _, stopVisit := range vehicleJourney.model.StopVisits().FindByVehicleJourneyId(vehicleJourney.id) {
		stopVisitIds = append(stopVisitIds, stopVisit.Id())
	}

	vehicleJourneyMap := make(map[string]interface{})
	vehicleJourney.FillVehicleJourney(vehicleJourneyMap)

	return json.Marshal(vehicleJourneyMap)
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
