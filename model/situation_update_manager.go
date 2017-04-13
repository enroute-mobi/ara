package model

type SituationUpdateManager struct {
	ClockConsumer

	transactionProvider TransactionProvider
}

type SituationUpdater struct {
	ClockConsumer

	tx    *Transaction
	event []*SituationUpdateEvent
}

func NewSituationUpdateManager(transactionProvider TransactionProvider) func([]*SituationUpdateEvent) {
	manager := newSituationUpdateManager(transactionProvider)
	return manager.UpdateSituation
}

func newSituationUpdateManager(transactionProvider TransactionProvider) *SituationUpdateManager {
	return &SituationUpdateManager{transactionProvider: transactionProvider}
}

func (manager *SituationUpdateManager) UpdateSituation(event []*SituationUpdateEvent) {
	tx := manager.transactionProvider.NewTransaction()
	defer tx.Close()

	NewSituationUpdater(tx, event).Update()

	tx.Commit()
}

func NewSituationUpdater(tx *Transaction, event []*SituationUpdateEvent) *SituationUpdater {
	return &SituationUpdater{tx: tx, event: event}
}

func (updater *SituationUpdater) CreateSituationFromEvent(event *SituationUpdateEvent) *Situation {

	situation := updater.tx.Model().Situations().New()
	situation.SetObjectID(event.SituationObjectID)
	situation.References = event.SituationAttributes.References
	situation.Messages = event.SituationAttributes.Messages
	situation.ValidUntil = event.SituationAttributes.ValidUntil
	situation.Channel = event.SituationAttributes.Channel
	situation.Format = event.SituationAttributes.Format

	situation.Save()
	return &situation
}

func (updater *SituationUpdater) Update() {
	for _, event := range updater.event {
		existingSituation, ok := updater.tx.Model().Situations().FindByObjectId(event.SituationObjectID)
		if ok && existingSituation.Version != event.Version {
			existingSituation.RecordedAt = event.RecordedAt
			existingSituation.Version = event.Version

			existingSituation.References = event.SituationAttributes.References
			existingSituation.Messages = event.SituationAttributes.Messages
			existingSituation.ValidUntil = event.SituationAttributes.ValidUntil
			existingSituation.Channel = event.SituationAttributes.Channel
			existingSituation.Format = event.SituationAttributes.Format
			existingSituation.Save()
			return
		}
		updater.CreateSituationFromEvent(event)
	}
}
