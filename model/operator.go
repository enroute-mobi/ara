package model

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type OperatorId string

type Operator struct {
	ObjectIDConsumer

	model Model

	id       OperatorId `json:",omitempty"`
	Name     string     `json:",omitempty"`
	Objectid *ObjectID  `json:",omitempty"`
}

func NewOperator(model Model) *Operator {
	operator := &Operator{
		model: model,
	}

	operator.objectids = make(ObjectIDs)
	return operator
}

func (operator *Operator) Id() OperatorId {
	return operator.id
}

func (operator *Operator) Save() (ok bool) {
	ok = operator.model.Operators().Save(operator)
	return
}

func (operator *Operator) MarshalJSON() ([]byte, error) {
	type Alias Operator

	aux := struct {
		Id OperatorId
		*Alias
	}{
		Id:    operator.id,
		Alias: (*Alias)(operator),
	}

	return json.Marshal(&aux)
}

type MemoryOperators struct {
	UUIDConsumer

	model *MemoryModel

	byIdentifier map[OperatorId]*Operator
}

type Operators interface {
	UUIDInterface

	New() Operator
	Find(id OperatorId) (Operator, bool)
	FindByObjectId(objectid ObjectID) (Operator, bool)
	FindAll() []Operator
	Save(operator *Operator) bool
	Delete(operator *Operator) bool
}

func NewMemoryOperators() *MemoryOperators {
	return &MemoryOperators{
		byIdentifier: make(map[OperatorId]*Operator),
	}
}

func (manager *MemoryOperators) New() Operator {
	operator := NewOperator(manager.model)
	return *operator
}

func (manager *MemoryOperators) Find(id OperatorId) (Operator, bool) {
	operator, ok := manager.byIdentifier[id]
	if ok {
		return *operator, true
	} else {
		return Operator{}, false
	}
}

func (manager *MemoryOperators) FindAll() (operators []Operator) {
	if len(manager.byIdentifier) == 0 {
		return []Operator{}
	}
	for _, operator := range manager.byIdentifier {
		operators = append(operators, *operator)
	}
	return
}

func (manager *MemoryOperators) FindByObjectId(objectid ObjectID) (Operator, bool) {
	for _, operator := range manager.byIdentifier {
		operatorObjectId, _ := operator.ObjectID(objectid.Kind())
		if operatorObjectId.Value() == objectid.Value() {
			return *operator, true
		}
	}
	return Operator{}, false
}

func (manager *MemoryOperators) Save(operator *Operator) bool {
	if operator.Id() == "" {
		operator.id = OperatorId(manager.NewUUID())
	}
	operator.model = manager.model
	manager.byIdentifier[operator.Id()] = operator
	return true
}

func (manager *MemoryOperators) Delete(operator *Operator) bool {
	delete(manager.byIdentifier, operator.Id())
	return true
}

func (manager *MemoryOperators) Load(referentialId string) error {
	var selectOperators []struct {
		Id            string
		ReferentialId string `db:"referential_id"`
		Name          sql.NullString
		ObjectID      sql.NullString `db:"object_id"`
		ObjectIDs     sql.NullString `db:"object_ids"`
	}
	sqlQuery := fmt.Sprintf("select * from operators where referential_id = '%s'", referentialId)

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
			if so.ObjectID.Valid && len(so.ObjectID.String) > 0 {
				objectid := &ObjectID{}
				if err = json.Unmarshal([]byte(so.ObjectID.String), objectid); err != nil {
					return err
				}
				operator.Objectid = objectid
			}
		}

		manager.Save(&operator)
	}
	return nil
}
