package audit

import (
	"context"
	"time"

	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
	"cloud.google.com/go/bigquery"
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

type BigQuery interface {
	model.Startable
	model.Stopable

	WriteMessage(message *BigQueryMessage) error
}

/**** Null struct to disable BQ by default ****/
type NullBigQuery struct{}

func (bq *NullBigQuery) WriteMessage(_ *BigQueryMessage) error {
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
	messages []*BigQueryMessage
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

func (bq *FakeBigQuery) Messages() []*BigQueryMessage {
	return bq.messages
}

/**** Real BQ ****/
type BigQueryClient struct {
	model.UUIDConsumer

	projectID string
	dataset   string
	table     string
	ctx       context.Context
	cancel    context.CancelFunc
	client    *bigquery.Client
	inserter  *bigquery.Inserter
	messages  chan *BigQueryMessage
	stop      chan struct{}
}

func NewBigQueryClient(projectID, dataset, table string) *BigQueryClient {
	return &BigQueryClient{
		projectID: projectID,
		dataset:   dataset,
		table:     table,
		messages:  make(chan *BigQueryMessage, 5),
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

func (bq *BigQueryClient) run() {
	bq.connect()

	for {
		select {
		case <-bq.stop:
			bq.client.Close()
			bq.cancel()
			return
		case message := <-bq.messages:
			if bq.inserter == nil {
				continue
			}
			ss := bigquery.StructSaver{Struct: message, InsertID: bq.NewUUID()}
			ctx, cancel := context.WithTimeout(bq.ctx, 5*time.Second)
			defer cancel()
			if err := bq.inserter.Put(ctx, &ss); err != nil {
				logger.Log.Debugf("BigQuery inserter error: %v", err)
			}
		}
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

	bq.inserter = bq.client.Dataset(bq.dataset).Table(bq.table).Inserter()
}
