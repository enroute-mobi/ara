package audit

import (
	"context"
	"sync"
	"time"

	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/state"
	"bitbucket.org/enroute-mobi/ara/uuid"
	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

const (
	EXCHANGE_TABLE = "exchanges"
	PARTNER_TABLE  = "partners"
	VEHICLE_TABLE  = "vehicles"
)

type BigQuery interface {
	state.Startable
	state.Stopable

	WriteMessage(message *BigQueryMessage) error
	WritePartnerEvent(partnerEvent *BigQueryPartnerEvent) error
	WriteVehicleEvent(vehicleEvent *BigQueryVehicleEvent) error
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

func (bq *NullBigQuery) WriteMessage(_ *BigQueryMessage) error {
	return nil
}

func (bq *NullBigQuery) WritePartnerEvent(_ *BigQueryPartnerEvent) error {
	return nil
}

func (bq *NullBigQuery) WriteVehicleEvent(_ *BigQueryVehicleEvent) error {
	return nil
}

func (bq *NullBigQuery) Start() {}
func (bq *NullBigQuery) Stop()  {}

func NewNullBigQuery() BigQuery {
	return &NullBigQuery{}
}

/**** Test Structure ****/
type FakeBigQuery struct {
	messages      []*BigQueryMessage
	partnerEvents []*BigQueryPartnerEvent
	vehicleEvents []*BigQueryVehicleEvent
}

func NewFakeBigQuery() *FakeBigQuery {
	return &FakeBigQuery{}
}

func (bq *FakeBigQuery) Start() {}
func (bq *FakeBigQuery) Stop()  {}

func (bq *FakeBigQuery) WriteMessage(message *BigQueryMessage) error {
	bq.messages = append(bq.messages, message)
	return nil
}

func (bq *FakeBigQuery) WritePartnerEvent(partnerEvent *BigQueryPartnerEvent) error {
	bq.partnerEvents = append(bq.partnerEvents, partnerEvent)
	return nil
}

func (bq *FakeBigQuery) WriteVehicleEvent(vehicleEvent *BigQueryVehicleEvent) error {
	bq.vehicleEvents = append(bq.vehicleEvents, vehicleEvent)
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

/**** Real BQ ****/
type BigQueryClient struct {
	uuid.UUIDConsumer

	projectID       string
	dataset         string
	ctx             context.Context
	client          *bigquery.Client
	inserter        *bigquery.Inserter
	vehicleInserter *bigquery.Inserter
	partnerInserter *bigquery.Inserter
	messages        chan *BigQueryMessage
	partnerEvents   chan *BigQueryPartnerEvent
	vehicleEvents   chan *BigQueryVehicleEvent
	stop            chan struct{}
}

func NewBigQueryClient(projectID, dataset string) *BigQueryClient {
	return &BigQueryClient{
		projectID:     projectID,
		dataset:       dataset,
		messages:      make(chan *BigQueryMessage, 500),
		partnerEvents: make(chan *BigQueryPartnerEvent, 500),
		vehicleEvents: make(chan *BigQueryVehicleEvent, 500),
	}
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

func (bq *BigQueryClient) WriteMessage(message *BigQueryMessage) error {
	select {
	case bq.messages <- message:
	default:
		logger.Log.Debugf("BigQuery queue is full")
	}
	return nil
}

func (bq *BigQueryClient) WritePartnerEvent(partnerEvent *BigQueryPartnerEvent) error {
	select {
	case bq.partnerEvents <- partnerEvent:
	default:
		logger.Log.Debugf("BigQuery partner queue is full")
	}
	return nil
}

func (bq *BigQueryClient) WriteVehicleEvent(vehicleEvent *BigQueryVehicleEvent) error {
	select {
	case bq.vehicleEvents <- vehicleEvent:
	default:
		logger.Log.Debugf("BigQuery vehicle queue is full")
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

	return dataset, nil
}
