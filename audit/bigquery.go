package audit

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"bitbucket.org/enroute-mobi/ara/audit/controlpb"
	"bitbucket.org/enroute-mobi/ara/audit/exchangepb"
	"bitbucket.org/enroute-mobi/ara/audit/longtermsvpb"
	"bitbucket.org/enroute-mobi/ara/audit/partnerpb"
	"bitbucket.org/enroute-mobi/ara/audit/vehiclepb"
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/config"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/state"
	"bitbucket.org/enroute-mobi/ara/uuid"
	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
	"google.golang.org/protobuf/proto"
)

const (
	EXCHANGE_TABLE             = "exchanges"
	PARTNER_TABLE              = "partners"
	VEHICLE_TABLE              = "vehicles"
	LONG_TERM_STOP_VISIT_TABLE = "long_term_stop_visits"
	CONTROL_TABLE              = "control_messages"
)

var tablesSchemas = map[string]bigquery.Schema{
	EXCHANGE_TABLE:             bqMessageSchema,
	PARTNER_TABLE:              bqPartnerSchema,
	VEHICLE_TABLE:              bqVehicleSchema,
	LONG_TERM_STOP_VISIT_TABLE: bqLongTermStopVisitsSchema,
	CONTROL_TABLE:              bqControlSchema,
}

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

	logger.Log.Printf("BigQuery WriteEvent failed: %v", err)

	return err
}

/**** Real BQ ****/
type BigQueryClient struct {
	uuid.UUIDConsumer
	clock.ClockConsumer

	projectID                         string
	dataset                           string
	ctx                               context.Context
	client                            *bigquery.Client
	exchangeWriter                    *storageWriter
	vehicleWriter                     *storageWriter
	longTermStopVisitWriter           *storageWriter
	controlWriter                     *storageWriter
	partnerWriter                     *storageWriter
	messages                          chan *BigQueryMessage
	partnerEvents                     chan *BigQueryPartnerEvent
	vehicleEvents                     chan *BigQueryVehicleEvent
	longTermStopVisitEvents           chan *BigQueryLongTermStopVisitEvent
	controlEvents                     chan *BigQueryControlEvent
	stop                              chan struct{}
	lostMessagesCount                 atomic.Int64
	lostPartnerEventsCount            atomic.Int64
	lostVehicleEventsCount            atomic.Int64
	lostLongTermStopVisitsEventsCount atomic.Int64
	lostControlEventsCount            atomic.Int64
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
		if n := bq.lostMessagesCount.Swap(0); n > 0 {
			logger.Log.Printf("BigQuery message queue was full: %d lost messages", n)
		}
	default:
		bq.lostMessagesCount.Add(1)
	}
	return nil
}

func (bq *BigQueryClient) writePartnerEvent(partnerEvent *BigQueryPartnerEvent) error {
	select {
	case bq.partnerEvents <- partnerEvent:
		if n := bq.lostPartnerEventsCount.Swap(0); n > 0 {
			logger.Log.Printf("BigQuery partnerEvent queue was full: %d lost messages", n)
		}
	default:
		bq.lostPartnerEventsCount.Add(1)
	}
	return nil
}

func (bq *BigQueryClient) writeVehicleEvent(vehicleEvent *BigQueryVehicleEvent) error {
	select {
	case bq.vehicleEvents <- vehicleEvent:
		if n := bq.lostVehicleEventsCount.Swap(0); n > 0 {
			logger.Log.Printf("BigQuery vehicleEvent queue was full: %d lost messages", n)
		}
	default:
		bq.lostVehicleEventsCount.Add(1)
	}
	return nil
}

func (bq *BigQueryClient) writeLongTermStopVisitEvent(longTermStopVisitEvent *BigQueryLongTermStopVisitEvent) error {
	select {
	case bq.longTermStopVisitEvents <- longTermStopVisitEvent:
		if n := bq.lostLongTermStopVisitsEventsCount.Swap(0); n > 0 {
			logger.Log.Printf("BigQuery longTermStopVisitEvent queue was full: %d lost messages", n)
		}
	default:
		bq.lostLongTermStopVisitsEventsCount.Add(1)
	}
	return nil
}

func (bq *BigQueryClient) writeControlEvent(controlEvent *BigQueryControlEvent) error {
	select {
	case bq.controlEvents <- controlEvent:
		if n := bq.lostControlEventsCount.Swap(0); n > 0 {
			logger.Log.Printf("BigQuery controleEvent queue was full: %d lost messages", n)
		}
	default:
		bq.lostControlEventsCount.Add(1)
	}
	return nil
}

func (bq *BigQueryClient) run() {
	bq.connect()
	for {
		select {
		case <-bq.stop:
			bq.closeWriters()
			bq.client.Close()
			return
		case message := <-bq.messages:
			bq.sendWithWriter(EXCHANGE_TABLE, bq.exchangeWriter, func() ([]byte, error) { return encodeExchange(message) })
		case partnerMessage := <-bq.partnerEvents:
			bq.sendWithWriter(PARTNER_TABLE, bq.partnerWriter, func() ([]byte, error) { return encodePartner(partnerMessage) })
		case vehicleMessage := <-bq.vehicleEvents:
			bq.sendWithWriter(VEHICLE_TABLE, bq.vehicleWriter, func() ([]byte, error) { return encodeVehicle(vehicleMessage) })
		case longTermStopVisitMessage := <-bq.longTermStopVisitEvents:
			if os.Getenv("ENABLE_BIGQUERY_LTS") != "false" {
				bq.sendWithWriter(LONG_TERM_STOP_VISIT_TABLE, bq.longTermStopVisitWriter, func() ([]byte, error) {
					return encodeLongTermStopVisit(longTermStopVisitMessage)
				})
			}
		case controlMessage := <-bq.controlEvents:
			bq.sendWithWriter(CONTROL_TABLE, bq.controlWriter, func() ([]byte, error) { return encodeControl(controlMessage) })
		}
	}
}

func (bq *BigQueryClient) closeWriters() {
	for _, w := range []*storageWriter{bq.exchangeWriter, bq.partnerWriter, bq.vehicleWriter, bq.longTermStopVisitWriter, bq.controlWriter} {
		if w != nil {
			w.close()
		}
	}
}

func (bq *BigQueryClient) sendWithWriter(table string, w *storageWriter, encode func() ([]byte, error)) {
	if w == nil {
		return
	}
	data, err := encode()
	if err != nil {
		logger.Log.Printf("BigQuery encode error (%s): %v", table, err)
		return
	}
	if err := w.send(bq.ctx, data); err != nil {
		logger.Log.Printf("BigQuery storage writer error (%s): %v", table, err)
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

	if err := bq.findOrCreateDatasetAndTables(); err != nil {
		logger.Log.Printf("error while finding or creating the dataset: %v", err)
		return
	}
	writers := []struct {
		field **storageWriter
		table string
		proto proto.Message
	}{
		{&bq.exchangeWriter, EXCHANGE_TABLE, &exchangepb.BigQueryMessage{}},
		{&bq.partnerWriter, PARTNER_TABLE, &partnerpb.BigQueryPartnerEvent{}},
		{&bq.vehicleWriter, VEHICLE_TABLE, &vehiclepb.BigQueryVehicleEvent{}},
		{&bq.longTermStopVisitWriter, LONG_TERM_STOP_VISIT_TABLE, &longtermsvpb.BigQueryLongTermStopVisitEvent{}},
		{&bq.controlWriter, CONTROL_TABLE, &controlpb.BigQueryControlEvent{}},
	}
	for _, w := range writers {
		writer, err := newStorageWriter(bq.ctx, bq.projectID, bq.dataset, w.table, w.proto)
		if err != nil {
			logger.Log.Printf("error creating storage writer for %s: %v", w.table, err)
			continue
		}
		*w.field = writer
	}
}

func (bq *BigQueryClient) findOrCreateDatasetAndTables() error {
	dataset, err := bq.findOrCreateDataset()
	if err != nil {
		return err
	}

	for t, s := range tablesSchemas {
		if err = bq.findOrCreateTable(dataset, t, s); err != nil {
			return err
		}
	}

	return nil
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

	logger.Log.Printf("Creating New Dataset")
	dataset := bq.client.Dataset(bq.dataset)
	if err := dataset.Create(bq.ctx, &bigquery.DatasetMetadata{Location: "EU"}); err != nil {
		return nil, err
	}

	return dataset, nil
}

func (bq *BigQueryClient) findOrCreateTable(dataset *bigquery.Dataset, tableName string, schema bigquery.Schema) error {
	t := dataset.Table(tableName)
	if _, err := t.Metadata(bq.ctx); err == nil {
		return nil
	}

	logger.Log.Printf("Creating New Table %s for Dataset %s", tableName, bq.dataset)

	p := &bigquery.TimePartitioning{
		Field:      "timestamp",
		Expiration: 30 * 24 * time.Hour,
	}

	if err := t.Create(bq.ctx, &bigquery.TableMetadata{TimePartitioning: p, Schema: schema}); err != nil {
		return err
	}

	return nil

}
