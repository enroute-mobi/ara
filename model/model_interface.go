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

func (memoryModel *MemoryModel) StopAreas() StopAreas {
	return memoryModel.stopAreas
}
