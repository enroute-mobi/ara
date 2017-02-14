package model

import (
	"strconv"
	"time"
)

type Model interface {
	StopAreas() StopAreas
	StopVisits() StopVisits
	VehicleJourneys() VehicleJourneys
	Lines() Lines
	GetDate() time.Time
	SetDate(string) time.Time
	// ...
}

type MemoryModel struct {
	stopAreas       StopAreas
	stopVisits      StopVisits
	vehicleJourneys VehicleJourneys
	lines           Lines
	date            time.Time
}

func NewMemoryModel() *MemoryModel {
	model := &MemoryModel{}

	stopAreas := NewMemoryStopAreas()
	stopAreas.model = model
	model.stopAreas = stopAreas

	stopVisits := NewMemoryStopVisits()
	stopVisits.model = model
	model.stopVisits = stopVisits

	vehicleJourneys := NewMemoryVehicleJourneys()
	vehicleJourneys.model = model
	model.vehicleJourneys = vehicleJourneys

	lines := NewMemoryLines()
	lines.model = model
	model.lines = lines

	return model
}

func (model *MemoryModel) SetDate(reloadHour string) time.Time {
	hour, minute := 4, 0
	if len(reloadHour) == 5 {
		hour, _ = strconv.Atoi(reloadHour[0:2])
		minute, _ = strconv.Atoi(reloadHour[3:5])
	}
	loc_cet, _ := time.LoadLocation("CET")
	now := time.Now().In(loc_cet)
	model.date = time.Date(now.Year(), now.Month(), now.Day()+1, hour, minute, 0, 0, loc_cet)
	return model.date
}

func (model *MemoryModel) GetDate() time.Time {
	return model.date
}

func (model *MemoryModel) StopAreas() StopAreas {
	return model.stopAreas
}

func (model *MemoryModel) StopVisits() StopVisits {
	return model.stopVisits
}

func (model *MemoryModel) VehicleJourneys() VehicleJourneys {
	return model.vehicleJourneys
}

func (model *MemoryModel) Lines() Lines {
	return model.lines
}

func (model *MemoryModel) NewTransaction() *Transaction {
	return NewTransaction(model)
}
