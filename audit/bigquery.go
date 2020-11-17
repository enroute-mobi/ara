package audit

import (
	"context"
	"time"

	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/state"
	"bitbucket.org/enroute-mobi/ara/uuid"
	"cloud.google.com/go/bigquery"
)

// FIXME: we need to see how we store the table names
const (
	EXCHANGE_TABLE = "exchange_events"
	PARTNER_TABLE  = "partner_events"
	VEHICLE_TABLE  = "vehicle_events"
)

type BigQueryMessage struct {
	Timestamp               time.Time `bigquery:"timestamp"`
	IPAddress               string    `bigquery:ip_address`
	Protocol                string    `bigquery:"protocol"`
	Type                    string    `bigquery:"type"`
	Direction               string    `bigquery:"direction"`
	Partner                 string    `bigquery:"partner"`
	Status                  string    `bigquery:"status"`
	ErrorDetails            string    `bigquery:"error_details"`
	RequestRawMessage       string    `bigquery:"request_raw_message"`
	ResponseRawMessage      string    `bigquery:"response_raw_message"`
	RequestIdentifier       string    `bigquery:"request_identifier"`
	ResponseIdentifier      string    `bigquery:"response_identifier"`
	RequestSize             int       `bigquery:"request_size"`
	ResponseSize            int       `bigquery:"response_size"`
	ProcessingTime          float64   `bigquery:"processing_time"`
	SubscriptionIdentifiers []string  `bigquery:"subscription_identifiers"`
	Lines                   []string  `bigquery:"lines"`
	StopAreas               []string  `bigquery:"stop_areas"`
	Vehicles                []string  `bigquery:"vehicles"`
}

type BigQueryPartnerEvent struct {
	Timestamp                time.Time `bigquery:"timestamp"`
	Slug                     string    `bigquery:"slug"`
	PreviousStatus           string    `bigquery:"previous_status"`
	PreviousServiceStartedAt time.Time `bigquery:"previous_service_started_at"`
	NewStatus                string    `bigquery:"new_status"`
	NewServiceStartedAt      time.Time `bigquery:"new_service_started_at"`
}

type BigQueryVehicleEvent struct {
	Timestamp      time.Time `bigquery:"timestamp"`
	ID             string    `bigquery:"id"`
	ObjectIDs      []string  `bigquery:"objectids"`
	Longitude      float64   `bigquery:"longitude"`
	Latitude       float64   `bigquery:"latitude"`
	Bearing        float64   `bigquery:"bearing"`
	RecordedAtTime time.Time `bigquery:"recorded_at_time"`
}

type BigQuery interface {
	state.Startable
	state.Stopable

	WriteMessage(message *BigQueryMessage) error
	WritePartnerEvent(partnerEvent *BigQueryPartnerEvent) error
	WriteVehicleEvent(vehicleEvent *BigQueryVehicleEvent) error
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

var currentBigQuery BigQuery = NewNullBigQuery()

func CurrentBigQuery() BigQuery {
	return currentBigQuery
}

func SetCurrentBigQuery(bq BigQuery) {
	currentBigQuery = bq
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

/**** Real BQ ****/
type BigQueryClient struct {
	uuid.UUIDConsumer

	projectID       string
	dataset         string
	ctx             context.Context
	cancel          context.CancelFunc
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
		messages:      make(chan *BigQueryMessage, 5),
		partnerEvents: make(chan *BigQueryPartnerEvent, 5),
		vehicleEvents: make(chan *BigQueryVehicleEvent, 5),
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
			bq.cancel()
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

	dataset := bq.client.Dataset(bq.dataset)
	bq.inserter = dataset.Table(EXCHANGE_TABLE).Inserter()
	bq.partnerInserter = dataset.Table(PARTNER_TABLE).Inserter()
	bq.vehicleInserter = dataset.Table(VEHICLE_TABLE).Inserter()
}
