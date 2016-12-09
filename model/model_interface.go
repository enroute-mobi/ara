package model

type Model interface {
	StopAreas() StopAreas
	StopVisits() StopVisits
	// ...
}

type MemoryModel struct {
	stopAreas  StopAreas
	stopVisits StopVisits
}

func NewMemoryModel() *MemoryModel {
	model := &MemoryModel{}

	stopAreas := NewMemoryStopAreas()
	stopAreas.model = model
	model.stopAreas = stopAreas

	stopVisits := NewMemoryStopVisits()
	stopVisits.model = model
	model.stopVisits = stopVisits

	return model
}

func (model *MemoryModel) StopAreas() StopAreas {
	return model.stopAreas
}

func (model *MemoryModel) StopVisits() StopVisits {
	return model.stopVisits
}

func (model *MemoryModel) NewTransaction() *Transaction {
	return NewTransaction(model)
}
