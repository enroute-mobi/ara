package model

type Model interface {
	StopAreas() StopAreas
	StopVisits() StopVisits
	VehicleJourneys() VehicleJourneys
	// ...
}

type MemoryModel struct {
	stopAreas       StopAreas
	stopVisits      StopVisits
	vehicleJourneys VehicleJourneys
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

	return model
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

func (model *MemoryModel) NewTransaction() *Transaction {
	return NewTransaction(model)
}
