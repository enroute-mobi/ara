package model

type UpdateManager struct {
	ClockConsumer
	UUIDConsumer

	transactionProvider TransactionProvider
}

func NewUpdateManager(transactionProvider TransactionProvider) func(UpdateEvent) {
	manager := &UpdateManager{transactionProvider: transactionProvider}
	return manager.Update
}

func (manager *UpdateManager) Update(event UpdateEvent) {
	switch event.EventKind() {
	case STOP_AREA_EVENT:
		manager.updateStopArea(event.(*StopAreaUpdateEvent))
	case LINE_EVENT:
		manager.updateLine(event.(*LineUpdateEvent))
	}
}

func (manager *UpdateManager) updateStopArea(event *StopAreaUpdateEvent) {
	tx := manager.transactionProvider.NewTransaction()

	stopArea, found := tx.Model().StopAreas().FindByObjectId(event.ObjectId)
	if !found {
		stopArea = tx.Model().StopAreas().New()

		stopArea.SetObjectID(event.ObjectId)
		stopArea.CollectGeneralMessages = true
	}

	stopArea.Name = event.Name
	stopArea.CollectedAlways = event.CollectedAlways
	stopArea.Longitude = event.Longitude
	stopArea.Latitude = event.Latitude

	if stopArea.ParentId == "" && event.ParentObjectId.Value() != "" {
		parentSA, _ := tx.Model().StopAreas().FindByObjectId(event.ParentObjectId)
		stopArea.ParentId = parentSA.Id()
	}

	stopArea.Updated(manager.Clock().Now())

	tx.Model().StopAreas().Save(&stopArea)
	tx.Commit()
	tx.Close()
}

func (manager *UpdateManager) updateLine(event *LineUpdateEvent) {
	tx := manager.transactionProvider.NewTransaction()

	line, found := tx.Model().Lines().FindByObjectId(event.ObjectId)
	if !found {
		line = tx.Model().Lines().New()

		line.SetObjectID(event.ObjectId)
		line.SetObjectID(NewObjectID("_default", event.ObjectId.HashValue()))
	}

	line.Name = event.Name
	line.CollectGeneralMessages = true
	line.SetOrigin(event.Origin)

	line.Updated(manager.Clock().Now())

	tx.Model().Lines().Save(&line)
	tx.Commit()
	tx.Close()
}
