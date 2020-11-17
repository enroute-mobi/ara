package model

import "bitbucket.org/enroute-mobi/ara/uuid"

type TransactionalLines struct {
	uuid.UUIDConsumer

	model   Model
	saved   map[LineId]*Line
	deleted map[LineId]*Line
}

func NewTransactionalLines(model Model) *TransactionalLines {
	lines := TransactionalLines{model: model}
	lines.resetCaches()
	return &lines
}

func (manager *TransactionalLines) resetCaches() {
	manager.saved = make(map[LineId]*Line)
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
	for _, line := range manager.saved {
		lineObjectId, _ := line.ObjectID(objectid.Kind())
		if lineObjectId.Value() == objectid.Value() {
			return *line, true
		}
	}
	return manager.model.Lines().FindByObjectId(objectid)
}

func (manager *TransactionalLines) FindAll() []Line {
	lines := []Line{}
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
	return lines
}

func (manager *TransactionalLines) Save(line *Line) bool {
	if line.Id() == "" {
		line.id = LineId(manager.NewUUID())
	}
	manager.saved[line.Id()] = line
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
