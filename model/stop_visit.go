package model

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync"
	"time"

	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

type StopVisitId ModelId

type StopVisit struct {
	RecordedAt  time.Time
	collectedAt time.Time
	model       Model
	References  References
	Attributes  Attributes
	Schedules   *StopVisitSchedules
	CodeConsumer
	VehicleJourneyId VehicleJourneyId         `json:",omitempty"`
	StopAreaId       StopAreaId               `json:",omitempty"`
	ArrivalStatus    StopVisitArrivalStatus   `json:",omitempty"`
	DepartureStatus  StopVisitDepartureStatus `json:",omitempty"`
	DataFrameRef     string                   `json:",omitempty"`
	id               StopVisitId
	Origin           string
	PassageOrder     int `json:",omitempty"`
	collected        bool
	VehicleAtStop    bool
}

func NewStopVisit(model Model) *StopVisit {
	stopVisit := &StopVisit{
		model:      model,
		Schedules:  NewStopVisitSchedules(),
		Attributes: NewAttributes(),
		References: NewReferences(),
	}
	stopVisit.codes = make(Codes)
	return stopVisit
}

func (stopVisit *StopVisit) modelId() ModelId {
	return ModelId(stopVisit.id)
}

func (stopVisit *StopVisit) copy() *StopVisit {
	return &StopVisit{
		CodeConsumer:     stopVisit.CodeConsumer.Copy(),
		model:            stopVisit.model,
		Origin:           stopVisit.Origin,
		id:               stopVisit.id,
		collected:        stopVisit.collected,
		collectedAt:      stopVisit.collectedAt,
		StopAreaId:       stopVisit.StopAreaId,
		VehicleJourneyId: stopVisit.VehicleJourneyId,
		Attributes:       stopVisit.Attributes.Copy(),
		References:       stopVisit.References.Copy(),
		ArrivalStatus:    stopVisit.ArrivalStatus,
		DepartureStatus:  stopVisit.DepartureStatus,
		DataFrameRef:     stopVisit.DataFrameRef,
		RecordedAt:       stopVisit.RecordedAt,
		Schedules:        stopVisit.Schedules.Copy(),
		VehicleAtStop:    stopVisit.VehicleAtStop,
		PassageOrder:     stopVisit.PassageOrder,
	}
}

func (stopVisit *StopVisit) IsCollected() bool {
	return stopVisit.collected
}

func (stopVisit *StopVisit) IsRecordable(now time.Time) bool {
	if stopVisit.DepartureStatus == STOP_VISIT_DEPARTURE_DEPARTED {
		return true
	}

	if stopVisit.ArrivalStatus == STOP_VISIT_ARRIVAL_ARRIVED &&
		stopVisit.DepartureStatus == STOP_VISIT_DEPARTURE_CANCELLED {
		return true
	}

	return stopVisit.ReferenceTime().Before(now)
}

func (stopVisit *StopVisit) IsArchivable() bool {
	return stopVisit.DepartureStatus == STOP_VISIT_DEPARTURE_DEPARTED ||
		(stopVisit.ArrivalStatus == STOP_VISIT_ARRIVAL_CANCELLED &&
			stopVisit.DepartureStatus == STOP_VISIT_DEPARTURE_CANCELLED)
}

func (stopVisit *StopVisit) NotCollected(notCollectedAt time.Time) {
	stopVisit.collected = false
	stopVisit.VehicleAtStop = false

	if !stopVisit.ArrivalStatus.Arrived() {
		stopVisit.ArrivalStatus = STOP_VISIT_ARRIVAL_ARRIVED
	}
	if !stopVisit.DepartureStatus.Departed() {
		stopVisit.DepartureStatus = STOP_VISIT_DEPARTURE_DEPARTED
	}

	stopVisit.Schedules.SetArrivalTimeIfNotDefined(STOP_VISIT_SCHEDULE_ACTUAL, notCollectedAt)
	stopVisit.Schedules.SetDepartureTimeIfNotDefined(STOP_VISIT_SCHEDULE_ACTUAL, notCollectedAt)
}

func (stopVisit *StopVisit) CollectedAt() time.Time {
	return stopVisit.collectedAt
}

func (stopVisit *StopVisit) Collected(t time.Time) {
	stopVisit.collected = true
	stopVisit.collectedAt = t
}

func (stopVisit *StopVisit) Id() StopVisitId {
	return stopVisit.id
}

func (stopVisit *StopVisit) StopArea() *StopArea {
	stopArea, _ := stopVisit.model.StopAreas().Find(stopVisit.StopAreaId)
	return stopArea
}

func (stopVisit *StopVisit) VehicleJourney() *VehicleJourney {
	vehicleJourney, ok := stopVisit.model.VehicleJourneys().Find(stopVisit.VehicleJourneyId)
	if !ok {
		return nil
	}
	return vehicleJourney
}

func (stopVisit *StopVisit) MarshalJSON() ([]byte, error) {
	type Alias StopVisit
	aux := struct {
		References map[string]Reference `json:",omitempty"`
		Codes      Codes                `json:",omitempty"`
		*Alias
		CollectedAt *time.Time `json:",omitempty"`
		RecordedAt  *time.Time `json:",omitempty"`
		Attributes  Attributes `json:",omitempty"`
		Id          StopVisitId
		Schedules   []StopVisitSchedule `json:",omitempty"`
		Collected   bool
	}{
		Id:        stopVisit.id,
		Collected: stopVisit.collected,
		Alias:     (*Alias)(stopVisit),
	}

	if !stopVisit.Codes().Empty() {
		aux.Codes = stopVisit.Codes()
	}
	if !stopVisit.Attributes.IsEmpty() {
		aux.Attributes = stopVisit.Attributes
	}
	if !stopVisit.References.IsEmpty() {
		aux.References = stopVisit.References.GetReferences()
	}
	if !stopVisit.RecordedAt.IsZero() {
		aux.RecordedAt = &stopVisit.RecordedAt
	}
	if !stopVisit.collectedAt.IsZero() {
		aux.CollectedAt = &stopVisit.collectedAt
	}

	scheduleSlice := stopVisit.Schedules.ToSlice()
	if len(scheduleSlice) != 0 {
		aux.Schedules = scheduleSlice
	}

	return json.Marshal(&aux)
}

func (stopVisit *StopVisit) UnmarshalJSON(data []byte) error {
	type Alias StopVisit
	aux := &struct {
		CollectedAt time.Time
		Codes       map[string]string
		References  map[string]Reference
		*Alias
		Schedules []StopVisitSchedule
	}{
		Alias: (*Alias)(stopVisit),
	}

	err := json.Unmarshal(data, aux)
	if err != nil {
		return err
	}

	if aux.Codes != nil {
		stopVisit.CodeConsumer.codes = NewCodesFromMap(aux.Codes)
	}

	if aux.References != nil {
		stopVisit.References.SetReferences(aux.References)
	}

	if aux.Schedules != nil {
		stopVisit.Schedules = NewStopVisitSchedules()
		for _, schedule := range aux.Schedules {
			stopVisit.Schedules.SetSchedule(schedule.Kind(), schedule.DepartureTime(), schedule.ArrivalTime())
		}
	}

	if !aux.CollectedAt.IsZero() {
		stopVisit.Collected(aux.CollectedAt)
	}
	return nil
}

func (stopVisit *StopVisit) Attribute(key string) (string, bool) {
	value, present := stopVisit.Attributes[key]
	return value, present
}

func (stopVisit *StopVisit) Save() bool {
	return stopVisit.model.StopVisits().Save(stopVisit)
}

func (stopVisit *StopVisit) Reference(key string) (Reference, bool) {
	value, present := stopVisit.References.Get(key)
	return value, present
}

func (stopVisit *StopVisit) ReferenceTime() time.Time {
	return stopVisit.Schedules.ReferenceTime()
}

func (stopVisit *StopVisit) ReferenceArrivalTime() time.Time {
	return stopVisit.Schedules.ReferenceArrivalTime()
}

func (stopVisit *StopVisit) ReferenceDepartureTime() time.Time {
	return stopVisit.Schedules.ReferenceDepartureTime()
}

type ByTime []*StopVisit

func (a ByTime) Len() int           { return len(a) }
func (a ByTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTime) Less(i, j int) bool { return !a[i].ReferenceTime().After(a[j].ReferenceTime()) }

type MemoryStopVisits struct {
	uuid.UUIDConsumer
	clock.ClockConsumer

	model Model

	mutex                             *sync.RWMutex
	byIdentifier                      map[StopVisitId]*StopVisit
	byCode                            *CodeIndex
	byStopArea                        *Index
	byVehicleJourney                  *Index
	byVehicleJourneyIdAndPassageOrder map[vehicleJourneyStopVisitOrder]StopVisitId

	broadcastEvent func(event StopMonitoringBroadcastEvent)
}

type vehicleJourneyStopVisitOrder struct {
	vehicleJourneyId      VehicleJourneyId
	stopVisitPassageOrder int
}

type StopVisits interface {
	uuid.UUIDInterface

	New() *StopVisit
	Find(StopVisitId) (*StopVisit, bool)
	FindByCode(Code) (*StopVisit, bool)
	FindByVehicleJourneyId(VehicleJourneyId) []*StopVisit
	FindByVehicleJourneyIdAndStopVisitOrder(VehicleJourneyId, int) *StopVisit
	VehicleJourneyHasStopVisits(VehicleJourneyId) bool
	FindByVehicleJourneyIdAfter(VehicleJourneyId, time.Time) []*StopVisit
	FindFollowingByVehicleJourneyId(VehicleJourneyId) []*StopVisit
	StopVisitsLenByVehicleJourney(VehicleJourneyId) int
	FindByStopAreaId(StopAreaId) []*StopVisit
	FindMonitoredByOriginByStopAreaId(StopAreaId, string) []*StopVisit
	FindFollowingByStopAreaId(StopAreaId) []*StopVisit
	FindFollowingByStopAreaIds([]StopAreaId) []*StopVisit
	FindAll() []*StopVisit
	UnsafeFindAll() []*StopVisit
	FindByVehicleJourneyIdAndStopAreaId(VehicleJourneyId, StopAreaId) []StopVisitId
	FindAllAfter(time.Time) []*StopVisit
	Save(*StopVisit) bool
	Delete(*StopVisit) bool
}

func NewMemoryStopVisits() *MemoryStopVisits {
	stopExtractor := func(instance ModelInstance) ModelId { return ModelId((instance.(*StopVisit)).StopAreaId) }
	vjExtractor := func(instance ModelInstance) ModelId { return ModelId((instance.(*StopVisit)).VehicleJourneyId) }

	return &MemoryStopVisits{
		mutex:                             &sync.RWMutex{},
		byIdentifier:                      make(map[StopVisitId]*StopVisit),
		byCode:                            NewCodeIndex(),
		byStopArea:                        NewIndex(stopExtractor),
		byVehicleJourney:                  NewIndex(vjExtractor),
		byVehicleJourneyIdAndPassageOrder: make(map[vehicleJourneyStopVisitOrder]StopVisitId),
	}
}

func (manager *MemoryStopVisits) New() *StopVisit {
	return NewStopVisit(manager.model)
}

func (manager *MemoryStopVisits) Find(id StopVisitId) (*StopVisit, bool) {
	manager.mutex.RLock()
	stopVisit, ok := manager.byIdentifier[id]
	manager.mutex.RUnlock()

	if ok {
		return stopVisit.copy(), true
	}
	return &StopVisit{}, false
}

func (manager *MemoryStopVisits) FindByCode(code Code) (*StopVisit, bool) {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	id, ok := manager.byCode.Find(code)
	if ok {
		return manager.byIdentifier[StopVisitId(id)].copy(), true
	}

	return &StopVisit{}, false
}

func (manager *MemoryStopVisits) FindByVehicleJourneyIdAndStopVisitOrder(vjId VehicleJourneyId, order int) *StopVisit {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()
	stopVisitId := manager.byVehicleJourneyIdAndPassageOrder[vehicleJourneyStopVisitOrder{
		vehicleJourneyId:      vjId,
		stopVisitPassageOrder: order}]
	if stopVisitId != "" {
		stopVisit, ok := manager.byIdentifier[stopVisitId]
		if ok {
			return stopVisit.copy()
		}
	}
	return &StopVisit{}
}

func (manager *MemoryStopVisits) FindByVehicleJourneyId(id VehicleJourneyId) (stopVisits []*StopVisit) {
	manager.mutex.RLock()

	ids, _ := manager.byVehicleJourney.Find(ModelId(id))

	for _, id := range ids {
		sv := manager.byIdentifier[StopVisitId(id)]
		stopVisits = append(stopVisits, sv.copy())
	}

	manager.mutex.RUnlock()
	return
}

func (manager *MemoryStopVisits) FindByVehicleJourneyIdAndStopAreaId(vjId VehicleJourneyId, saId StopAreaId) []StopVisitId {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	var stopVisitIds []StopVisitId

	ids, _ := manager.byVehicleJourney.Find(ModelId(vjId))
	for _, id := range ids {
		stopVisit := manager.byIdentifier[StopVisitId(id)]
		if stopVisit.StopAreaId == saId {
			stopVisitIds = append(stopVisitIds, StopVisitId(id))
		}
	}

	return stopVisitIds
}

func (manager *MemoryStopVisits) StopVisitsLenByVehicleJourney(id VehicleJourneyId) int {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()
	return manager.byVehicleJourney.IndexableLength(ModelId(id))
}

func (manager *MemoryStopVisits) VehicleJourneyHasStopVisits(id VehicleJourneyId) bool {
	return manager.StopVisitsLenByVehicleJourney(id) != 0
}

func (manager *MemoryStopVisits) FindFollowingByVehicleJourneyId(id VehicleJourneyId) (stopVisits []*StopVisit) {
	return manager.FindByVehicleJourneyIdAfter(id, manager.Clock().Now())
}

func (manager *MemoryStopVisits) FindByVehicleJourneyIdAfter(id VehicleJourneyId, t time.Time) (stopVisits []*StopVisit) {
	manager.mutex.RLock()

	ids, _ := manager.byVehicleJourney.Find(ModelId(id))

	for _, id := range ids {
		sv := manager.byIdentifier[StopVisitId(id)]
		if sv.ReferenceTime().After(t) {
			stopVisits = append(stopVisits, sv.copy())
		}
	}

	manager.mutex.RUnlock()
	sort.Sort(ByTime(stopVisits))
	return
}

func (manager *MemoryStopVisits) FindByStopAreaId(id StopAreaId) (stopVisits []*StopVisit) {
	manager.mutex.RLock()

	ids, _ := manager.byStopArea.Find(ModelId(id))

	for _, id := range ids {
		sv := manager.byIdentifier[StopVisitId(id)]
		stopVisits = append(stopVisits, sv.copy())
	}

	manager.mutex.RUnlock()
	return
}

func (manager *MemoryStopVisits) FindMonitoredByOriginByStopAreaId(id StopAreaId, origin string) (stopVisits []*StopVisit) {
	manager.mutex.RLock()

	for _, stopVisit := range manager.byIdentifier {
		if stopVisit.StopAreaId == id && stopVisit.collected && stopVisit.Origin == origin {
			stopVisits = append(stopVisits, stopVisit.copy())
		}
	}

	manager.mutex.RUnlock()
	return
}

func (manager *MemoryStopVisits) FindFollowingByStopAreaId(id StopAreaId) (stopVisits []*StopVisit) {
	manager.mutex.RLock()

	ids, _ := manager.byStopArea.Find(ModelId(id))

	for _, id := range ids {
		sv := manager.byIdentifier[StopVisitId(id)]
		if sv.ReferenceTime().After(manager.Clock().Now()) {
			stopVisits = append(stopVisits, sv.copy())
		}
	}

	manager.mutex.RUnlock()
	sort.Sort(ByTime(stopVisits))
	return
}

func (manager *MemoryStopVisits) FindFollowingByStopAreaIds(stopAreaIds []StopAreaId) (stopVisits []*StopVisit) {
	manager.mutex.RLock()

	var ids []ModelId
	for _, id := range stopAreaIds {
		saids, _ := manager.byStopArea.Find(ModelId(id))
		ids = append(ids, saids...)
	}

	for _, id := range ids {
		sv := manager.byIdentifier[StopVisitId(id)]
		if sv.ReferenceTime().After(manager.Clock().Now()) {
			stopVisits = append(stopVisits, sv.copy())
		}
	}

	manager.mutex.RUnlock()
	sort.Sort(ByTime(stopVisits))
	return
}

func (manager *MemoryStopVisits) FindAll() (stopVisits []*StopVisit) {
	manager.mutex.RLock()

	for _, stopVisit := range manager.byIdentifier {
		stopVisits = append(stopVisits, stopVisit.copy())
	}

	manager.mutex.RUnlock()
	return
}

// Warning, this method doesn't copy the StopVisits so we shouldn't access its maps
func (manager *MemoryStopVisits) UnsafeFindAll() (stopVisits []*StopVisit) {
	manager.mutex.RLock()

	for _, stopVisit := range manager.byIdentifier {
		stopVisits = append(stopVisits, stopVisit)
	}

	manager.mutex.RUnlock()
	return
}

func (manager *MemoryStopVisits) FindAllAfter(t time.Time) (stopVisits []*StopVisit) {
	manager.mutex.RLock()

	for _, stopVisit := range manager.byIdentifier {
		if stopVisit.ReferenceTime().After(t) {
			stopVisits = append(stopVisits, stopVisit.copy())
		}
	}

	manager.mutex.RUnlock()
	return
}

func (manager *MemoryStopVisits) Save(stopVisit *StopVisit) bool {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	if stopVisit.id == "" {
		stopVisit.id = StopVisitId(manager.NewUUID())
	}

	stopVisit.model = manager.model
	manager.byIdentifier[stopVisit.id] = stopVisit
	manager.byCode.Index(stopVisit)
	manager.byStopArea.Index(stopVisit)
	manager.byVehicleJourney.Index(stopVisit)
	manager.byVehicleJourneyIdAndPassageOrder[vehicleJourneyStopVisitOrder{
		vehicleJourneyId:      stopVisit.VehicleJourneyId,
		stopVisitPassageOrder: stopVisit.PassageOrder}] = stopVisit.id

	event := StopMonitoringBroadcastEvent{
		ModelId:   string(stopVisit.id),
		ModelType: "StopVisit",
	}

	if manager.broadcastEvent != nil {
		manager.broadcastEvent(event)
	}

	return true
}

func (manager *MemoryStopVisits) Delete(stopVisit *StopVisit) bool {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	delete(manager.byIdentifier, stopVisit.id)
	manager.byCode.Delete(ModelId(stopVisit.id))
	manager.byStopArea.Delete(ModelId(stopVisit.id))
	manager.byVehicleJourney.Delete(ModelId(stopVisit.id))
	delete(manager.byVehicleJourneyIdAndPassageOrder, vehicleJourneyStopVisitOrder{
		vehicleJourneyId:      stopVisit.VehicleJourneyId,
		stopVisitPassageOrder: stopVisit.PassageOrder,
	})

	return true
}

func (manager *MemoryStopVisits) Load(referentialSlug string) error {
	var selectStopVisits []SelectStopVisit
	modelName := manager.model.Date()

	sqlQuery := fmt.Sprintf("select * from stop_visits where referential_slug = '%s' and model_name = '%s'", referentialSlug, modelName.String())
	_, err := Database.Select(&selectStopVisits, sqlQuery)
	if err != nil {
		return err
	}
	for _, sv := range selectStopVisits {
		stopVisit := manager.New()
		stopVisit.id = StopVisitId(sv.Id)
		if sv.StopAreaId.Valid {
			stopVisit.StopAreaId = StopAreaId(sv.StopAreaId.String)
		}
		if sv.VehicleJourneyId.Valid {
			stopVisit.VehicleJourneyId = VehicleJourneyId(sv.VehicleJourneyId.String)
		}
		if sv.PassageOrder.Valid {
			stopVisit.PassageOrder = int(sv.PassageOrder.Int64)
		}

		if sv.Attributes.Valid && len(sv.Attributes.String) > 0 {
			if err = json.Unmarshal([]byte(sv.Attributes.String), &stopVisit.Attributes); err != nil {
				return err
			}
		}

		if sv.References.Valid && len(sv.References.String) > 0 {
			references := make(map[string]Reference)
			if err = json.Unmarshal([]byte(sv.References.String), &references); err != nil {
				return err
			}
			stopVisit.References.SetReferences(references)
		}

		if sv.Codes.Valid && len(sv.Codes.String) > 0 {
			codeMap := make(map[string]string)
			if err = json.Unmarshal([]byte(sv.Codes.String), &codeMap); err != nil {
				return err
			}
			stopVisit.codes = NewCodesFromMap(codeMap)
		}

		if sv.Schedules.Valid && len(sv.Schedules.String) > 0 {
			scheduleSlice := []StopVisitSchedule{}
			if err = json.Unmarshal([]byte(sv.Schedules.String), &scheduleSlice); err != nil {
				return err
			}
			stopVisit.Schedules = NewStopVisitSchedules()
			for _, schedule := range scheduleSlice {
				stopVisit.Schedules.SetSchedule(schedule.Kind(), schedule.DepartureTime(), schedule.ArrivalTime())
			}
		}

		manager.Save(stopVisit)
	}
	return nil
}
