package model

type Model interface {
	StopAreas() StopAreas
	// ...
}

type MemoryModel struct {
	stopAreas StopAreas
}

func NewMemoryModel() *MemoryModel {
	model := &MemoryModel{}

	stopAreas := NewMemoryStopAreas()
	stopAreas.model = model
	model.stopAreas = stopAreas

	return model
}

func (model *MemoryModel) StopAreas() StopAreas {
	return model.stopAreas
}

func (model *MemoryModel) NewTransaction() *Transaction {
	return NewTransaction(model)
}
