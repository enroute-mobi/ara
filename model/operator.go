package model

import (
	"encoding/json"
	"fmt"
	"sync"

	"bitbucket.org/enroute-mobi/ara/uuid"
)

type OperatorId string

type Operator struct {
	ObjectIDConsumer

	model Model

	id   OperatorId
	Name string `json:",omitempty"`
}

func NewOperator(model Model) *Operator {
	operator := &Operator{
		model: model,
	}

	operator.objectids = make(ObjectIDs)
	return operator
}

func (operator *Operator) modelId() ModelId {
	return ModelId(operator.id)
}

func (operator *Operator) copy() *Operator {
	o := *operator
	return &o
}

func (operator *Operator) Id() OperatorId {
	return operator.id
}

func (operator *Operator) Save() bool {
	return operator.model.Operators().Save(operator)
}

func (operator *Operator) MarshalJSON() ([]byte, error) {
	type Alias Operator

	aux := struct {
		Id        OperatorId
		ObjectIDs ObjectIDs `json:",omitempty"`
		*Alias
	}{
		Id:    operator.id,
		Alias: (*Alias)(operator),
	}

	if !operator.ObjectIDs().Empty() {
		aux.ObjectIDs = operator.ObjectIDs()
	}

	return json.Marshal(&aux)
}

func (operator *Operator) UnmarshalJSON(data []byte) error {
	type Alias Operator
	aux := &struct {
		ObjectIDs map[string]string
		*Alias
	}{
		Alias: (*Alias)(operator),
	}
	err := json.Unmarshal(data, aux)
	if err != nil {
		return err
	}

	if aux.ObjectIDs != nil {
		operator.ObjectIDConsumer.objectids = NewObjectIDsFromMap(aux.ObjectIDs)
	}

	return nil
}

type MemoryOperators struct {
	uuid.UUIDConsumer

	model *MemoryModel
	mutex *sync.RWMutex

	byIdentifier map[OperatorId]*Operator
	byObjectId   *ObjectIdIndex
}

type Operators interface {
	uuid.UUIDInterface

	New() *Operator
	Find(OperatorId) (*Operator, bool)
	FindByObjectId(ObjectID) (*Operator, bool)
	FindAll() []*Operator
	Save(*Operator) bool
	Delete(*Operator) bool
}

func NewMemoryOperators() *MemoryOperators {
	return &MemoryOperators{
		mutex:        &sync.RWMutex{},
		byIdentifier: make(map[OperatorId]*Operator),
		byObjectId:   NewObjectIdIndex(),
	}
}

func (manager *MemoryOperators) New() *Operator {
	return NewOperator(manager.model)
}

func (manager *MemoryOperators) Find(id OperatorId) (*Operator, bool) {
	manager.mutex.RLock()
	operator, ok := manager.byIdentifier[id]
	manager.mutex.RUnlock()

	if ok {
		return operator.copy(), true
	}
	return &Operator{}, false
}

func (manager *MemoryOperators) FindAll() (operators []*Operator) {
	manager.mutex.RLock()

	for _, operator := range manager.byIdentifier {
		operators = append(operators, operator.copy())
	}

	manager.mutex.RUnlock()
	return
}

func (manager *MemoryOperators) FindByObjectId(objectid ObjectID) (*Operator, bool) {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	id, ok := manager.byObjectId.Find(objectid)
	if ok {
		return manager.byIdentifier[OperatorId(id)].copy(), true
	}

	return &Operator{}, false
}

func (manager *MemoryOperators) Save(operator *Operator) bool {
	manager.mutex.Lock()

	if operator.Id() == "" {
		operator.id = OperatorId(manager.NewUUID())
	}
	operator.model = manager.model
	manager.byIdentifier[operator.Id()] = operator
	manager.byObjectId.Index(operator)

	manager.mutex.Unlock()
	return true
}

func (manager *MemoryOperators) Delete(operator *Operator) bool {
	manager.mutex.Lock()

	delete(manager.byIdentifier, operator.Id())
	manager.byObjectId.Delete(ModelId(operator.id))

	manager.mutex.Unlock()
	return true
}

func (manager *MemoryOperators) Load(referentialSlug string) error {
	var selectOperators []SelectOperator
	modelName := manager.model.Date()

	sqlQuery := fmt.Sprintf("select * from operators where referential_slug = '%s' and model_name = '%s'", referentialSlug, modelName.String())

	_, err := Database.Select(&selectOperators, sqlQuery)
	if err != nil {
		return err
	}

	for _, so := range selectOperators {
		operator := manager.New()
		operator.id = OperatorId(so.Id)
		if so.Name.Valid {
			operator.Name = so.Name.String
		}

		if so.ObjectIDs.Valid && len(so.ObjectIDs.String) > 0 {
			objectIdMap := make(map[string]string)
			if err = json.Unmarshal([]byte(so.ObjectIDs.String), &objectIdMap); err != nil {
				return err
			}

			operator.objectids = NewObjectIDsFromMap(objectIdMap)
		}
		manager.Save(operator)
	}
	return nil
}
