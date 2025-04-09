package audit

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/config"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/state"
	"bitbucket.org/enroute-mobi/ara/uuid"
	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

const (
	EXCHANGE_TABLE             = "exchanges"
	PARTNER_TABLE              = "partners"
	VEHICLE_TABLE              = "vehicles"
	LONG_TERM_STOP_VISIT_TABLE = "long_term_stop_visits"
	CONTROL_TABLE              = "control_messages"
)

type BigQuery interface {
	state.Startable
	state.Stopable

	WriteEvent(event BigQueryEvent) error
}

/**** Manager ****/

type BigQueryManager struct {
	mutex *sync.RWMutex
	bq    map[string]BigQuery
}

var manager = BigQueryManager{
	mutex: &sync.RWMutex{},
	bq:    make(map[string]BigQuery),
}

func CurrentBigQuery(slug string) BigQuery {
	manager.mutex.Lock()
	bq, ok := manager.bq[slug]
	if !ok {
		bq = NewNullBigQuery()
		manager.bq[slug] = bq
	}
	manager.mutex.Unlock()
	return bq
}

func SetCurrentBigQuery(slug string, bq BigQuery) {
	manager.mutex.Lock()
	manager.bq[slug] = bq
	manager.mutex.Unlock()
}

/**** Null struct to disable BQ by default ****/
type NullBigQuery struct{}

func (bq *NullBigQuery) WriteEvent(_ BigQueryEvent) error { return nil }

func (bq *NullBigQuery) Start() {}
func (bq *NullBigQuery) Stop()  {}

func NewNullBigQuery() BigQuery {
	return &NullBigQuery{}
}

/**** Test Memory Structure ****/
type FakeBigQuery struct {
	messages                []*BigQueryMessage
	partnerEvents           []*BigQueryPartnerEvent
	vehicleEvents           []*BigQueryVehicleEvent
	longTermStopVisitEvents []*BigQueryLongTermStopVisitEvent
	controlEvents           []*BigQueryControlEvent
}

func NewFakeBigQuery() *FakeBigQuery {
	return &FakeBigQuery{}
}

func (bq *FakeBigQuery) Start() {}
func (bq *FakeBigQuery) Stop()  {}

func (bq *FakeBigQuery) WriteEvent(e BigQueryEvent) error {
	switch e.EventType() {
	case BQ_MESSAGE:
		bq.messages = append(bq.messages, e.(*BigQueryMessage))
	case BQ_PARTNER_EVENT:
		bq.partnerEvents = append(bq.partnerEvents, e.(*BigQueryPartnerEvent))
	case BQ_VEHICLE_EVENT:
		bq.vehicleEvents = append(bq.vehicleEvents, e.(*BigQueryVehicleEvent))
	case BQ_LONG_TERM_STOP_VISIT_EVENT:
		bq.longTermStopVisitEvents = append(bq.longTermStopVisitEvents, e.(*BigQueryLongTermStopVisitEvent))
	case BQ_CONTROL_EVENT:
		bq.controlEvents = append(bq.controlEvents, e.(*BigQueryControlEvent))
	}
	return nil
}

func (bq *FakeBigQuery) Messages() []*BigQueryMessage {
	return bq.messages
}

func (bq *FakeBigQuery) PartnerEvents() []*BigQueryPartnerEvent {
	return bq.partnerEvents
}

func (bq *FakeBigQuery) VehicleEvents() []*BigQueryVehicleEvent {
	return bq.vehicleEvents
}

func (bq *FakeBigQuery) LongTermStopVisitEvents() []*BigQueryLongTermStopVisitEvent {
	return bq.longTermStopVisitEvents
}

func (bq *FakeBigQuery) ControlEvents() []*BigQueryControlEvent {
	return bq.controlEvents
}

/**** Test External Structure ****/

type TestBigQuery struct {
	clock.ClockConsumer

	target  string
	dataset string
}

func NewTestBigQuery(dataset string) *TestBigQuery {
	return &TestBigQuery{
		dataset: dataset,
		target:  config.Config.BigQueryTest,
	}
}

func (bq *TestBigQuery) Start() {}
func (bq *TestBigQuery) Stop()  {}

type TestBigQueryMessage struct {
	*BigQueryMessage
	Dataset string
}

func (bq *TestBigQuery) WriteEvent(e BigQueryEvent) error {
	e.SetTimeStamp(bq.Clock().Now())
	logger.Log.Debugf("WriteEvent %v", e)

	switch e.EventType() {
	case BQ_MESSAGE:
		e = &TestBigQueryMessage{
			BigQueryMessage: e.(*BigQueryMessage),
			Dataset:         bq.dataset,
		}
		// case BQ_PARTNER_EVENT:
		// case BQ_VEHICLE_EVENT:
	}

	json, _ := json.Marshal(e)

	logger.Log.Debugf("Send JSON %v", string(json))

	_, err := http.Post(
		bq.target,
		"application/json",
		bytes.NewBuffer(json),
	)

	logger.Log.Debugf("WriteEvent err %v", err)

	return err
}

/**** Real BQ ****/
type BigQueryClient struct {
	uuid.UUIDConsumer
	clock.ClockConsumer

	projectID                 string
	dataset                   string
	ctx                       context.Context
	client                    *bigquery.Client
	inserter                  *bigquery.Inserter
	vehicleInserter           *bigquery.Inserter
	partnerInserter           *bigquery.Inserter
	longTermStopVisitInserter *bigquery.Inserter
	controlInserter           *bigquery.Inserter
	messages                  chan *BigQueryMessage
	partnerEvents             chan *BigQueryPartnerEvent
	vehicleEvents             chan *BigQueryVehicleEvent
	longTermStopVisitEvents   chan *BigQueryLongTermStopVisitEvent
	controlEvents             chan *BigQueryControlEvent
	stop                      chan struct{}
}

func NewBigQuery(dataset string) BigQuery {
	formattedDataset := formatDatasetName(dataset)
	if config.Config.BigQueryTestMode() {
		return NewTestBigQuery(formattedDataset)
	} else {
		return NewBigQueryClient(formattedDataset)
	}
}

func NewBigQueryClient(dataset string) *BigQueryClient {
	return &BigQueryClient{
		dataset:                 dataset,
		projectID:               config.Config.BigQueryProjectID,
		messages:                make(chan *BigQueryMessage, 500),
		partnerEvents:           make(chan *BigQueryPartnerEvent, 500),
		vehicleEvents:           make(chan *BigQueryVehicleEvent, 500),
		longTermStopVisitEvents: make(chan *BigQueryLongTermStopVisitEvent, 500),
		controlEvents:           make(chan *BigQueryControlEvent, 500),
	}
}

func formatDatasetName(dataset string) string {
	return strings.ReplaceAll(dataset, "-", "_")
}

func (bq *BigQueryClient) Start() {
	bq.stop = make(chan struct{})
	go bq.run()
}

func (bq *BigQueryClient) Stop() {
	if bq.stop != nil {
		close(bq.stop)
	}
}

func (bq *BigQueryClient) WriteEvent(e BigQueryEvent) error {
	e.SetUUID(bq.NewUUID())
	e.SetTimeStamp(bq.Clock().Now())
	switch e.EventType() {
	case BQ_MESSAGE:
		return bq.writeMessage(e.(*BigQueryMessage))
	case BQ_PARTNER_EVENT:
		return bq.writePartnerEvent(e.(*BigQueryPartnerEvent))
	case BQ_VEHICLE_EVENT:
		return bq.writeVehicleEvent(e.(*BigQueryVehicleEvent))
	case BQ_LONG_TERM_STOP_VISIT_EVENT:
		return bq.writeLongTermStopVisitEvent(e.(*BigQueryLongTermStopVisitEvent))
	case BQ_CONTROL_EVENT:
		return bq.writeControlEvent(e.(*BigQueryControlEvent))
	}
	logger.Log.Debugf("Unknown BigQueryMessage type")
	return nil
}

func (bq *BigQueryClient) writeMessage(message *BigQueryMessage) error {
	select {
	case bq.messages <- message:
	default:
		logger.Log.Debugf("BigQuery queue is full")
	}
	return nil
}

func (bq *BigQueryClient) writePartnerEvent(partnerEvent *BigQueryPartnerEvent) error {
	select {
	case bq.partnerEvents <- partnerEvent:
	default:
		logger.Log.Debugf("BigQuery partner queue is full")
	}
	return nil
}

func (bq *BigQueryClient) writeVehicleEvent(vehicleEvent *BigQueryVehicleEvent) error {
	select {
	case bq.vehicleEvents <- vehicleEvent:
	default:
		logger.Log.Debugf("BigQuery vehicle queue is full")
	}
	return nil
}

func (bq *BigQueryClient) writeLongTermStopVisitEvent(longTermStopVisitEvent *BigQueryLongTermStopVisitEvent) error {
	select {
	case bq.longTermStopVisitEvents <- longTermStopVisitEvent:
	default:
		logger.Log.Debugf("BigQuery longTermStopVisit queue is full")
	}
	return nil
}

func (bq *BigQueryClient) writeControlEvent(controlEvent *BigQueryControlEvent) error {
	select {
	case bq.controlEvents <- controlEvent:
	default:
		logger.Log.Debugf("BigQuery control queue is full")
	}
	return nil
}

func (bq *BigQueryClient) run() {
	bq.connect()

	for {
		select {
		case <-bq.stop:
			bq.client.Close()
			return
		case message := <-bq.messages:
			bq.send(message, bq.inserter)
		case partnerMessage := <-bq.partnerEvents:
			bq.send(partnerMessage, bq.partnerInserter)
		case vehicleMessage := <-bq.vehicleEvents:
			bq.send(vehicleMessage, bq.vehicleInserter)
		case longTermStopVisitMessage := <-bq.longTermStopVisitEvents:
			if os.Getenv("ENABLE_BIGQUERY_LTS") != "false" {
				bq.send(longTermStopVisitMessage, bq.longTermStopVisitInserter)
			}
		case controlMessage := <-bq.controlEvents:
			bq.send(controlMessage, bq.controlInserter)
		}
	}
}

func (bq *BigQueryClient) send(message interface{}, inserter *bigquery.Inserter) {
	if inserter == nil {
		return
	}
	ss := bigquery.StructSaver{Struct: message, InsertID: bq.NewUUID()}
	ctx, cancel := context.WithTimeout(bq.ctx, 5*time.Second)
	defer cancel()
	if err := inserter.Put(ctx, &ss); err != nil {
		logger.Log.Debugf("BigQuery inserter error: %v", err)
	}
}

func (bq *BigQueryClient) connect() {
	bq.ctx = context.Background()

	var err error
	bq.client, err = bigquery.NewClient(bq.ctx, bq.projectID)
	if err != nil {
		logger.Log.Printf("can't connect to BigQuery: %v", err)
		return
	}

	dataset, err := bq.findOrCreateDataset()
	if err != nil {
		logger.Log.Printf("error while finding or creating the dataset: %v", err)
		return
	}
	bq.inserter = dataset.Table(EXCHANGE_TABLE).Inserter()
	bq.partnerInserter = dataset.Table(PARTNER_TABLE).Inserter()
	bq.vehicleInserter = dataset.Table(VEHICLE_TABLE).Inserter()
	bq.longTermStopVisitInserter = dataset.Table(LONG_TERM_STOP_VISIT_TABLE).Inserter()
	bq.controlInserter = dataset.Table(CONTROL_TABLE).Inserter()
}

func (bq *BigQueryClient) findOrCreateDataset() (*bigquery.Dataset, error) {
	it := bq.client.Datasets(bq.ctx)
	for {
		dataset, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		if dataset.DatasetID == bq.dataset {
			logger.Log.Printf("Found dataset %v", bq.dataset)
			return dataset, nil
		}
	}

	logger.Log.Printf("Creating New Dataset and tables")
	dataset := bq.client.Dataset(bq.dataset)
	if err := dataset.Create(bq.ctx, &bigquery.DatasetMetadata{Location: "EU"}); err != nil {
		return nil, err
	}

	p := &bigquery.TimePartitioning{
		Field:      "timestamp",
		Expiration: 30 * 24 * time.Hour,
	}

	if err := dataset.Table(EXCHANGE_TABLE).Create(bq.ctx, &bigquery.TableMetadata{TimePartitioning: p, Schema: bqMessageSchema}); err != nil {
		return nil, err
	}

	if err := dataset.Table(PARTNER_TABLE).Create(bq.ctx, &bigquery.TableMetadata{TimePartitioning: p, Schema: bqPartnerSchema}); err != nil {
		return nil, err
	}

	if err := dataset.Table(VEHICLE_TABLE).Create(bq.ctx, &bigquery.TableMetadata{TimePartitioning: p, Schema: bqVehicleSchema}); err != nil {
		return nil, err
	}

	if err := dataset.Table(LONG_TERM_STOP_VISIT_TABLE).Create(bq.ctx, &bigquery.TableMetadata{TimePartitioning: p, Schema: bqLongTermStopVisitsSchema}); err != nil {
		return nil, err
	}

	if err := dataset.Table(CONTROL_TABLE).Create(bq.ctx, &bigquery.TableMetadata{TimePartitioning: p, Schema: bqControlSchema}); err != nil {
		return nil, err
	}

	return dataset, nil
}
