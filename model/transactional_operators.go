package model

import "bitbucket.org/enroute-mobi/ara/uuid"

type TransactionalOperators struct {
	uuid.UUIDConsumer

	model   Model
	saved   map[OperatorId]*Operator
	deleted map[OperatorId]*Operator
}

func NewTransactionalOperators(model Model) *TransactionalOperators {
	operators := TransactionalOperators{model: model}
	operators.resetCaches()
	return &operators
}

func (manager *TransactionalOperators) resetCaches() {
	manager.saved = make(map[OperatorId]*Operator)
	manager.deleted = make(map[OperatorId]*Operator)
}

func (manager *TransactionalOperators) New() Operator {
	return *NewOperator(manager.model)
}

func (manager *TransactionalOperators) Find(id OperatorId) (Operator, bool) {
	operator, ok := manager.saved[id]
	if ok {
		return *operator, ok
	}

	return manager.model.Operators().Find(id)
}

func (manager *TransactionalOperators) FindByObjectId(objectid ObjectID) (Operator, bool) {
	for _, operator := range manager.saved {
		operatorObjectId, _ := operator.ObjectID(objectid.Kind())
		if operatorObjectId.Value() == objectid.Value() {
			return *operator, true
		}
	}
	return manager.model.Operators().FindByObjectId(objectid)
}

func (manager *TransactionalOperators) FindAll() []Operator {
	operators := []Operator{}
	for _, operator := range manager.saved {
		operators = append(operators, *operator)
	}
	savedLines := manager.model.Operators().FindAll()
	for _, operator := range savedLines {
		_, ok := manager.saved[operator.Id()]
		if !ok {
			operators = append(operators, operator)
		}
	}
	return operators
}

func (manager *TransactionalOperators) Save(operator *Operator) bool {
	if operator.Id() == "" {
		operator.id = OperatorId(manager.NewUUID())
	}
	manager.saved[operator.Id()] = operator
	return true
}

func (manager *TransactionalOperators) Delete(operator *Operator) bool {
	manager.deleted[operator.Id()] = operator
	return true
}

func (manager *TransactionalOperators) Commit() error {
	for _, operator := range manager.deleted {
		manager.model.Operators().Delete(operator)
	}
	for _, operator := range manager.saved {
		manager.model.Operators().Save(operator)
	}
	return nil
}

func (manager *TransactionalOperators) Rollback() error {
	manager.resetCaches()
	return nil
}
