package model

type TransactionalModel struct {
	parent Model

	lines           *TransactionalLines
	situations      *TransactionalSituations
	stopAreas       *TransactionalStopAreas
	stopVisits      *TransactionalStopVisits
	vehicleJourneys *TransactionalVehicleJourneys
	operators       *TransactionalOperators
}

func NewTransactionalModel(parent Model) *TransactionalModel {
	model := &TransactionalModel{parent: parent}
	model.lines = NewTransactionalLines(parent)
	model.situations = NewTransactionalSituations(parent)
	model.stopAreas = NewTransactionalStopAreas(parent)
	model.stopVisits = NewTransactionalStopVisits(parent)
	model.vehicleJourneys = NewTransactionalVehicleJourneys(parent)
	model.operators = NewTransactionalOperators(parent)
	return model
}

func (model *TransactionalModel) Date() Date {
	return model.parent.Date()
}

func (model *TransactionalModel) Lines() Lines {
	return model.lines
}

func (model *TransactionalModel) Situations() Situations {
	return model.situations
}

func (model *TransactionalModel) StopAreas() StopAreas {
	return model.stopAreas
}

func (model *TransactionalModel) StopVisits() StopVisits {
	return model.stopVisits
}

func (model *TransactionalModel) VehicleJourneys() VehicleJourneys {
	return model.vehicleJourneys
}

func (model *TransactionalModel) Operators() Operators {
	return model.operators
}

func (model *TransactionalModel) NewTransaction() *Transaction {
	return NewTransaction(model)
}

func (model *TransactionalModel) Commit() error {
	if err := model.stopAreas.Commit(); err != nil {
		return err
	}
	if err := model.stopVisits.Commit(); err != nil {
		return err
	}
	if err := model.vehicleJourneys.Commit(); err != nil {
		return err
	}
	if err := model.lines.Commit(); err != nil {
		return err
	}
	if err := model.situations.Commit(); err != nil {
		return err
	}
	return nil
}

func (model *TransactionalModel) Rollback() error {
	if err := model.stopAreas.Rollback(); err != nil {
		return err
	}
	if err := model.stopVisits.Rollback(); err != nil {
		return err
	}
	if err := model.vehicleJourneys.Rollback(); err != nil {
		return err
	}
	if err := model.lines.Rollback(); err != nil {
		return err
	}
	if err := model.situations.Rollback(); err != nil {
		return err
	}
	return nil
}
