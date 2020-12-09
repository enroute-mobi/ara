package audit

import (
	"time"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/civil"
)

const (
	BQ_MESSAGE       = "message"
	BQ_VEHICLE_EVENT = "vehicle"
	BQ_PARTNER_EVENT = "partner"
)

type BigQueryEvent interface {
	EventType() string
	SetTimeStamp(time.Time)
}

type BigQueryMessage struct {
	Timestamp               time.Time `bigquery:"timestamp"`
	IPAddress               string    `bigquery:"ip_address"`
	Protocol                string    `bigquery:"protocol"`  // "siri", "siri-lite", "gtfs", "push"
	Type                    string    `bigquery:"type"`      // "siri-checkstatus", "gtfs-trip-update", …
	Direction               string    `bigquery:"direction"` // "sent" (by Ara), "received" (by Ara)
	Partner                 string    `bigquery:"partner"`   // partner slug
	Status                  string    `bigquery:"status"`    // "OK", "Error"
	ErrorDetails            string    `bigquery:"error_details"`
	RequestRawMessage       string    `bigquery:"request_raw_message"`  // XML or JSON for GTFS-RT
	ResponseRawMessage      string    `bigquery:"response_raw_message"` // XML or JSON for GTFS-RT
	RequestIdentifier       string    `bigquery:"request_identifier"`
	ResponseIdentifier      string    `bigquery:"response_identifier"`
	RequestSize             int       `bigquery:"request_size"`
	ResponseSize            int       `bigquery:"response_size"`
	ProcessingTime          float64   `bigquery:"processing_time"`          // in seconds
	SubscriptionIdentifiers []string  `bigquery:"subscription_identifiers"` // array of ids
	StopAreas               []string  `bigquery:"stop_areas"`               // array of objectid values
	Lines                   []string  `bigquery:"lines"`                    // array of objectid values
	Vehicles                []string  `bigquery:"vehicles"`                 // array of objectid values
}

func (bq *BigQueryMessage) EventType() string        { return BQ_MESSAGE }
func (bq *BigQueryMessage) SetTimeStamp(t time.Time) { bq.Timestamp = t }

var bqMessageSchema = bigquery.Schema{
	{Name: "timestamp", Required: false, Type: bigquery.TimestampFieldType},
	{Name: "ip_address", Required: false, Type: bigquery.StringFieldType},
	{Name: "protocol", Required: false, Type: bigquery.StringFieldType},
	{Name: "type", Required: false, Type: bigquery.StringFieldType},
	{Name: "direction", Required: false, Type: bigquery.StringFieldType},
	{Name: "partner", Required: false, Type: bigquery.StringFieldType},
	{Name: "status", Required: false, Type: bigquery.StringFieldType},
	{Name: "error_details", Required: false, Type: bigquery.StringFieldType},
	{Name: "request_raw_message", Required: false, Type: bigquery.StringFieldType},
	{Name: "response_raw_message", Required: false, Type: bigquery.StringFieldType},
	{Name: "request_identifier", Required: false, Type: bigquery.StringFieldType},
	{Name: "response_identifier", Required: false, Type: bigquery.StringFieldType},
	{Name: "request_size", Required: false, Type: bigquery.IntegerFieldType},
	{Name: "response_size", Required: false, Type: bigquery.IntegerFieldType},
	{Name: "processing_time", Required: false, Type: bigquery.FloatFieldType},
	{Name: "subscription_identifiers", Repeated: true, Type: bigquery.StringFieldType},
	{Name: "lines", Repeated: true, Type: bigquery.StringFieldType},
	{Name: "stop_areas", Repeated: true, Type: bigquery.StringFieldType},
	{Name: "vehicles", Repeated: true, Type: bigquery.StringFieldType},
}

type BigQueryPartnerEvent struct {
	Timestamp                time.Time      `bigquery:"timestamp"`
	Slug                     string         `bigquery:"slug"`
	PreviousStatus           string         `bigquery:"previous_status"`
	PreviousServiceStartedAt civil.DateTime `bigquery:"previous_service_started_at"`
	NewStatus                string         `bigquery:"new_status"`
	NewServiceStartedAt      civil.DateTime `bigquery:"new_service_started_at"`
}

func (bq *BigQueryPartnerEvent) EventType() string        { return BQ_PARTNER_EVENT }
func (bq *BigQueryPartnerEvent) SetTimeStamp(t time.Time) { bq.Timestamp = t }

var bqPartnerSchema = bigquery.Schema{
	{Name: "timestamp", Required: false, Type: bigquery.TimestampFieldType},
	{Name: "slug", Required: false, Type: bigquery.StringFieldType},
	{Name: "previous_status", Required: false, Type: bigquery.StringFieldType},
	{Name: "previous_service_started_at", Required: false, Type: bigquery.DateTimeFieldType},
	{Name: "new_status", Required: false, Type: bigquery.StringFieldType},
	{Name: "new_service_started_at", Required: false, Type: bigquery.DateTimeFieldType},
}

type BigQueryVehicleEvent struct {
	Timestamp      time.Time      `bigquery:"timestamp"`
	ID             string         `bigquery:"id"`
	ObjectIDs      []string       `bigquery:"objectids"`
	Longitude      float64        `bigquery:"longitude"`
	Latitude       float64        `bigquery:"latitude"`
	Bearing        float64        `bigquery:"bearing"`
	RecordedAtTime civil.DateTime `bigquery:"recorded_at_time"`
}

func (bq *BigQueryVehicleEvent) EventType() string        { return BQ_VEHICLE_EVENT }
func (bq *BigQueryVehicleEvent) SetTimeStamp(t time.Time) { bq.Timestamp = t }

var bqVehicleSchema = bigquery.Schema{
	{Name: "timestamp", Required: false, Type: bigquery.TimestampFieldType},
	{Name: "id", Required: false, Type: bigquery.StringFieldType},
	{Name: "longitude", Required: false, Type: bigquery.FloatFieldType},
	{Name: "latitude", Required: false, Type: bigquery.FloatFieldType},
	{Name: "bearing", Required: false, Type: bigquery.FloatFieldType},
	{Name: "recorded_at_time", Required: false, Type: bigquery.DateTimeFieldType},
	{Name: "objectids", Repeated: true, Type: bigquery.StringFieldType},
}
