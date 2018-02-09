package model

type SituationUpdateManager struct {
	ClockConsumer

	transactionProvider TransactionProvider
}

type SituationUpdater struct {
	ClockConsumer

	tx     *Transaction
	events []*SituationUpdateEvent
}

func NewSituationUpdateManager(transactionProvider TransactionProvider) func([]*SituationUpdateEvent) {
	manager := newSituationUpdateManager(transactionProvider)
	return manager.UpdateSituation
}

func newSituationUpdateManager(transactionProvider TransactionProvider) *SituationUpdateManager {
	return &SituationUpdateManager{transactionProvider: transactionProvider}
}

func (manager *SituationUpdateManager) UpdateSituation(events []*SituationUpdateEvent) {
	tx := manager.transactionProvider.NewTransaction()
	defer tx.Close()

	NewSituationUpdater(tx, events).Update()

	tx.Commit()
}

func NewSituationUpdater(tx *Transaction, events []*SituationUpdateEvent) *SituationUpdater {
	return &SituationUpdater{tx: tx, events: events}
}

func (updater *SituationUpdater) Update() {
	for _, event := range updater.events {
		situation, ok := updater.tx.Model().Situations().FindByObjectId(event.SituationObjectID)
		if ok && situation.RecordedAt == event.RecordedAt {
			return
		}

		if !ok {
			situation = updater.tx.Model().Situations().New()
			situation.Origin = event.Origin
			situation.SetObjectID(event.SituationObjectID)
			situation.SetObjectID(NewObjectID("_default", event.SituationObjectID.HashValue()))
		}

		situation.RecordedAt = event.RecordedAt
		situation.Version = event.Version
		situation.ProducerRef = event.ProducerRef

		situation.References = event.SituationAttributes.References
		situation.LineSections = event.SituationAttributes.LineSections
		situation.Messages = event.SituationAttributes.Messages
		situation.ValidUntil = event.SituationAttributes.ValidUntil
		situation.Channel = event.SituationAttributes.Channel
		situation.Format = event.SituationAttributes.Format

		situation.Save()
	}
}
