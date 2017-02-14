package model

import (
	"strconv"
	"time"
)

type TransactionalModel struct {
	parent Model

	stopAreas       *TransactionalStopAreas
	stopVisits      *TransactionalStopVisits
	vehicleJourneys *TransactionalVehicleJourneys
	lines           *TransactionalLines
	date            time.Time
}

func NewTransactionalModel(parent Model) *TransactionalModel {
	model := &TransactionalModel{parent: parent}
	model.stopAreas = NewTransactionalStopAreas(parent)
	model.stopVisits = NewTransactionalStopVisits(parent)
	model.vehicleJourneys = NewTransactionalVehicleJourneys(parent)
	model.lines = NewTransactionalLines(parent)
	model.date = parent.GetDate()
	return model
}

func (model *TransactionalModel) SetDate(reloadHour string) time.Time {
	hour, minute := 4, 0
	if len(reloadHour) == 5 {
		hour, _ = strconv.Atoi(reloadHour[0:2])
		minute, _ = strconv.Atoi(reloadHour[3:5])
	}
	loc_cet, _ := time.LoadLocation("CET")
	now := time.Now().In(loc_cet)
	model.date = time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, loc_cet)

	return model.date
}

func (model *TransactionalModel) GetDate() time.Time {
	return model.date
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

func (model *TransactionalModel) Lines() Lines {
	return model.lines
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
	return nil
}
