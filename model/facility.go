package model

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"bitbucket.org/enroute-mobi/ara/uuid"
)

type FacilityId ModelId

type FacilityStatus string

const (
	FacilityStatusUnknown            FacilityStatus = "unknown"
	FacilityStatusAvailable          FacilityStatus = "available"
	FacilityStatusNotAvailable       FacilityStatus = "notAvailable"
	FacilityStatusPartiallyAvailable FacilityStatus = "partiallyAvailable"
	FacilityStatusRemoved            FacilityStatus = "removed"
)

type Facility struct {
	Collectable
	model Model
	id    FacilityId
	CodeConsumer
	Status FacilityStatus `json:",omitempty"`
	Origin string
}

func NewFacility(model Model) *Facility {
	facility := &Facility{
		model: model,
	}
	facility.codes = make(Codes)
	return facility
}

func (facility *Facility) ModelId() ModelId {
	return ModelId(facility.id)
}

func (facility *Facility) copy() *Facility {
	return &Facility{
		Collectable: Collectable{
			nextCollectAt: facility.nextCollectAt,
			collectedAt:   facility.collectedAt,
		},
		CodeConsumer: facility.CodeConsumer.Copy(),
		model:        facility.model,
		id:           facility.id,
		Status:       facility.Status,
	}
}

func (facility *Facility) Id() FacilityId {
	return facility.id
}

type MemoryFacilities struct {
	uuid.UUIDConsumer

	model        *MemoryModel
	mutex        *sync.RWMutex
	byIdentifier map[FacilityId]*Facility
	byCode       *CodeIndex
}

type Facilities interface {
	uuid.UUIDInterface

	New() *Facility
	FindAll() []*Facility
	Find(FacilityId) (*Facility, bool)
	FindByCode(Code) (*Facility, bool)
	Save(*Facility) bool
	Delete(*Facility) bool
}

func (facility *Facility) Save() bool {
	return facility.model.Facilities().Save(facility)
}

func NewMemoryFacilities() *MemoryFacilities {
	return &MemoryFacilities{
		mutex:        &sync.RWMutex{},
		byIdentifier: make(map[FacilityId]*Facility),
		byCode:       NewCodeIndex(),
	}
}

func (manager *MemoryFacilities) New() *Facility {
	return NewFacility(manager.model)
}

func (manager *MemoryFacilities) Find(id FacilityId) (*Facility, bool) {
	manager.mutex.RLock()
	facility, ok := manager.byIdentifier[id]
	manager.mutex.RUnlock()

	if ok {
		return facility.copy(), true
	}
	return &Facility{}, false
}

func (manager *MemoryFacilities) FindByCode(code Code) (*Facility, bool) {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	id, ok := manager.byCode.Find(code)
	if ok {
		return manager.byIdentifier[FacilityId(id)].copy(), true
	}

	return &Facility{}, false
}

func (manager *MemoryFacilities) FindAll() (facilitys []*Facility) {
	manager.mutex.RLock()

	for _, facility := range manager.byIdentifier {
		facilitys = append(facilitys, facility.copy())
	}

	manager.mutex.RUnlock()
	return
}

func (manager *MemoryFacilities) Save(facility *Facility) bool {
	if facility.Id() == "" {
		facility.id = FacilityId(manager.NewUUID())
	}

	manager.mutex.Lock()

	facility.model = manager.model
	manager.byIdentifier[facility.Id()] = facility
	manager.byCode.Index(facility)

	manager.mutex.Unlock()

	return true
}

func (manager *MemoryFacilities) Delete(facility *Facility) bool {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	delete(manager.byIdentifier, facility.Id())
	manager.byCode.Delete(ModelId(facility.id))

	return true
}

func (facility *Facility) MarshalJSON() ([]byte, error) {
	type Alias Facility

	aux := struct {
		Codes         Codes      `json:",omitempty"`
		NextCollectAt *time.Time `json:",omitempty"`
		CollectedAt   *time.Time `json:",omitempty"`
		Id            FacilityId `json:",omitempty"`
		*Alias
	}{
		Id:    facility.id,
		Alias: (*Alias)(facility),
	}

	if !facility.Codes().Empty() {
		aux.Codes = facility.Codes()
	}

	if !facility.nextCollectAt.IsZero() {
		aux.NextCollectAt = &facility.nextCollectAt
	}
	if !facility.collectedAt.IsZero() {
		aux.CollectedAt = &facility.collectedAt
	}
	return json.Marshal(&aux)
}

func (facility *Facility) UnmarshalJSON(data []byte) error {
	type Alias Facility
	aux := &struct {
		Codes map[string]string
		*Alias
	}{
		Alias: (*Alias)(facility),
	}
	err := json.Unmarshal(data, aux)
	if err != nil {
		return err
	}

	if aux.Codes != nil {
		facility.CodeConsumer.codes = NewCodesFromMap(aux.Codes)
	}

	return nil
}

func (manager *MemoryFacilities) Load(referentialSlug string) error {
	var selectFacilities []SelectFacility
	modelName := manager.model.Date()

	sqlQuery := fmt.Sprintf("select * from facilities where referential_slug = '%s' and model_name = '%s'", referentialSlug, modelName.String())

	_, err := Database.Select(&selectFacilities, sqlQuery)
	if err != nil {
		return err
	}

	for _, so := range selectFacilities {
		facility := manager.New()
		facility.id = FacilityId(so.Id)

		if so.Codes.Valid && len(so.Codes.String) > 0 {
			codeMap := make(map[string]string)
			if err = json.Unmarshal([]byte(so.Codes.String), &codeMap); err != nil {
				return err
			}

			facility.codes = NewCodesFromMap(codeMap)
		}
		manager.Save(facility)
	}
	return nil
}
