package model

type TransactionalLines struct {
	UUIDConsumer

	model           Model
	saved           map[LineId]*Line
	savedByObjectId map[string]map[string]LineId
	deleted         map[LineId]*Line
}

func NewTransactionalLines(model Model) *TransactionalLines {
	lines := TransactionalLines{model: model}
	lines.resetCaches()
	return &lines
}

func (manager *TransactionalLines) resetCaches() {
	manager.saved = make(map[LineId]*Line)
	manager.savedByObjectId = make(map[string]map[string]LineId)
	manager.deleted = make(map[LineId]*Line)
}

func (manager *TransactionalLines) New() Line {
	return *NewLine(manager.model)
}

func (manager *TransactionalLines) Find(id LineId) (Line, bool) {
	line, ok := manager.saved[id]
	if ok {
		return *line, ok
	}

	return manager.model.Lines().Find(id)
}

func (manager *TransactionalLines) FindByObjectId(objectid ObjectID) (Line, bool) {
	valueMap, ok := manager.savedByObjectId[objectid.Kind()]
	if !ok {
		return manager.model.Lines().FindByObjectId(objectid)
	}
	id, ok := valueMap[objectid.Value()]
	if !ok {
		return manager.model.Lines().FindByObjectId(objectid)
	}
	return *manager.saved[id], true
}

func (manager *TransactionalLines) FindAll() (lines []Line) {
	for _, line := range manager.saved {
		lines = append(lines, *line)
	}
	savedLines := manager.model.Lines().FindAll()
	for _, line := range savedLines {
		_, ok := manager.saved[line.Id()]
		if !ok {
			lines = append(lines, line)
		}
	}
	return
}

func (manager *TransactionalLines) Save(line *Line) bool {
	if line.Id() == "" {
		line.id = LineId(manager.NewUUID())
	}
	manager.saved[line.Id()] = line
	for _, objectid := range line.ObjectIDs() {
		_, ok := manager.savedByObjectId[objectid.Kind()]
		if !ok {
			manager.savedByObjectId[objectid.Kind()] = make(map[string]LineId)
		}
		manager.savedByObjectId[objectid.Kind()][objectid.Value()] = line.Id()
	}
	return true
}

func (manager *TransactionalLines) Delete(line *Line) bool {
	manager.deleted[line.Id()] = line
	return true
}

func (manager *TransactionalLines) Commit() error {
	for _, line := range manager.deleted {
		manager.model.Lines().Delete(line)
	}
	for _, line := range manager.saved {
		manager.model.Lines().Save(line)
	}
	return nil
}

func (manager *TransactionalLines) Rollback() error {
	manager.resetCaches()
	return nil
}
