package audit

import (
	"time"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/civil"
)

const (
	BQ_MESSAGE                    = "message"
	BQ_VEHICLE_EVENT              = "vehicle"
	BQ_PARTNER_EVENT              = "partner"
	BQ_CONTROL_EVENT              = "control"
	BQ_LONG_TERM_STOP_VISIT_EVENT = "long_term_stop_visit"

	CHECK_STATUS_REQUEST BigQueryMessageType = "CheckStatusRequest"

	NOTIFY_ESTIMATED_TIMETABLE  BigQueryMessageType = "NotifyEstimatedTimetable"
	NOTIFY_GENERAL_MESSAGE      BigQueryMessageType = "NotifyGeneralMessage"
	NOTIFY_PRODUCTION_TIMETABLE BigQueryMessageType = "NotifyProductionTimetable"
	NOTIFY_SITUATION_EXCHANGE   BigQueryMessageType = "NotifySituationExchange"
	NOTIFY_STOP_MONITORING      BigQueryMessageType = "NotifyStopMonitoring"
	NOTIFY_VEHICLE_MONITORING   BigQueryMessageType = "NotifyVehicleMonitoring"

	DELETE_SUBSCRIPTION_REQUEST    BigQueryMessageType = "DeleteSubscriptionRequest"
	NOTIFY_SUBSCRIPTION_TERMINATED BigQueryMessageType = "NotifySubscriptionTerminated"

	ESTIMATED_TIMETABLE_REQUEST   BigQueryMessageType = "EstimatedTimetableRequest"
	FACILITY_MONITORING_REQUEST   BigQueryMessageType = "FacilityMonitoringRequest"
	GENERAL_MESSAGE_REQUEST       BigQueryMessageType = "GeneralMessageRequest"
	LINES_DISCOVERY_REQUEST       BigQueryMessageType = "LinesDiscoveryRequest"
	SIRI_SERVICE_REQUEST          BigQueryMessageType = "SiriServiceRequest"
	SITUATION_EXCHANGE_REQUEST    BigQueryMessageType = "SituationExchangeRequest"
	STOP_MONITORING_REQUEST       BigQueryMessageType = "StopMonitoringRequest"
	STOP_POINTS_DISCOVERY_REQUEST BigQueryMessageType = "StopPointsDiscoveryRequest"
	VEHICLE_MONITORING_REQUEST    BigQueryMessageType = "VehicleMonitoringRequest"

	GTFS_TRIP_UPDATES_VEHICLE_POSITION_SERVICE_ALERTS BigQueryMessageType = "trip-updates,vehicle-position,service-alerts"
	GTFS_TRIP_UPDATES                                 BigQueryMessageType = "trip-updates"
	GTFS_VEHICLE_POSITION                             BigQueryMessageType = "vehicle-position"
	GTFS_REQUEST                                      BigQueryMessageType = "GtfsRequest"

	SUBSCRIPTION_TERMINATED_NOTIFICATION BigQueryMessageType = "SubscriptionTerminatedNotification"

	ESTIMATED_TIMETABLE_SUBSCRIPTION_REQUEST  BigQueryMessageType = "EstimatedTimetableSubscriptionRequest"
	GENERAL_MESSAGE_SUBSCRIPTION_REQUEST      BigQueryMessageType = "GeneralMessageSubscriptionRequest"
	PRODUCTION_TIMETABLE_SUBSCRIPTION_REQUEST BigQueryMessageType = "ProductionTimetableSubscriptionRequest"
	STOP_MONITORING_SUBSCRIPTION_REQUEST      BigQueryMessageType = "StopMonitoringSubscriptionRequest"
	VEHICLE_MONITORING_SUBSCRIPTION_REQUEST   BigQueryMessageType = "VehicleMonitoringSubscriptionRequest"
	SITUATION_EXCHANGE_SUBSCRIPTION_REQUEST   BigQueryMessageType = "SituationExchangeSubscriptionRequest"

	PUSH_NOTIFICATION BigQueryMessageType = "push-notification"

	GRAPHQL_REQUEST BigQueryMessageType = "GraphQLRequest"
)

var AraBigQuerySchemas = map[string]bigquery.Schema{
	"bqMessageSchema":            bqMessageSchema,
	"bqVehicleSchema":            bqVehicleSchema,
	"bqPartnerSchema":            bqPartnerSchema,
	"bqLongTermStopVisitsSchema": bqLongTermStopVisitsSchema,
	"bqControlSchema":            bqControlSchema,
}

type BigQueryMessageType string
type BigQueryEvent interface {
	EventType() string
	SetTimeStamp(time.Time)
	SetUUID(string)
}

type BigQueryMessage struct {
	UUID                    string              `bigquery:"uuid"`
	Timestamp               time.Time           `bigquery:"timestamp"`
	IPAddress               string              `bigquery:"ip_address"`
	Protocol                string              `bigquery:"protocol"`  // "siri", "siri-lite", "gtfs", "push"
	Type                    BigQueryMessageType `bigquery:"type"`      // "siri-checkstatus", "gtfs-trip-update", …
	Direction               string              `bigquery:"direction"` // "sent" (by Ara), "received" (by Ara)
	Partner                 string              `bigquery:"partner"`   // partner slug
	Status                  string              `bigquery:"status"`    // "OK", "Error"
	ErrorDetails            string              `bigquery:"error_details"`
	RequestRawMessage       string              `bigquery:"request_raw_message"`  // XML or JSON for GTFS-RT
	ResponseRawMessage      string              `bigquery:"response_raw_message"` // XML or JSON for GTFS-RT
	RequestIdentifier       string              `bigquery:"request_identifier"`
	ResponseIdentifier      string              `bigquery:"response_identifier"`
	RequestSize             int64               `bigquery:"request_size"`
	ResponseSize            int64               `bigquery:"response_size"`
	ProcessingTime          float64             `bigquery:"processing_time"`          // in seconds
	SubscriptionIdentifiers []string            `bigquery:"subscription_identifiers"` // array of ids
	StopAreas               []string            `bigquery:"stop_areas"`               // array of code values
	Lines                   []string            `bigquery:"lines"`                    // array of code values
	Vehicles                []string            `bigquery:"vehicles"`
	VehicleJourneys         []string            `bigquery:"vehicle_journeys"` // array of code values
}

func (bq *BigQueryMessage) EventType() string        { return BQ_MESSAGE }
func (bq *BigQueryMessage) SetTimeStamp(t time.Time) { bq.Timestamp = t }
func (bq *BigQueryMessage) SetUUID(u string)         { bq.UUID = u }

var bqMessageSchema = bigquery.Schema{
	{Name: "uuid", Required: false, Type: bigquery.StringFieldType},
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
	{Name: "vehicle_journeys", Repeated: true, Type: bigquery.StringFieldType},
}

type BigQueryPartnerEvent struct {
	UUID                     string         `bigquery:"uuid"`
	Timestamp                time.Time      `bigquery:"timestamp"`
	Slug                     string         `bigquery:"slug"`
	PartnerUUID              string         `bigquery:"partner_uuid"`
	PreviousStatus           string         `bigquery:"previous_status"`
	PreviousServiceStartedAt civil.DateTime `bigquery:"previous_service_started_at"`
	NewStatus                string         `bigquery:"new_status"`
	NewServiceStartedAt      civil.DateTime `bigquery:"new_service_started_at"`
}

func (bq *BigQueryPartnerEvent) EventType() string        { return BQ_PARTNER_EVENT }
func (bq *BigQueryPartnerEvent) SetTimeStamp(t time.Time) { bq.Timestamp = t }
func (bq *BigQueryPartnerEvent) SetUUID(u string)         { bq.UUID = u }

var bqPartnerSchema = bigquery.Schema{
	{Name: "uuid", Required: false, Type: bigquery.StringFieldType},
	{Name: "timestamp", Required: false, Type: bigquery.TimestampFieldType},
	{Name: "slug", Required: false, Type: bigquery.StringFieldType},
	{Name: "partner_uuid", Required: false, Type: bigquery.StringFieldType},
	{Name: "previous_status", Required: false, Type: bigquery.StringFieldType},
	{Name: "previous_service_started_at", Required: false, Type: bigquery.DateTimeFieldType},
	{Name: "new_status", Required: false, Type: bigquery.StringFieldType},
	{Name: "new_service_started_at", Required: false, Type: bigquery.DateTimeFieldType},
}

type BigQueryVehicleEvent struct {
	UUID           string         `bigquery:"uuid"`
	Timestamp      time.Time      `bigquery:"timestamp"`
	ID             string         `bigquery:"id"`
	Codes          []string       `bigquery:"codes"`
	Longitude      float64        `bigquery:"longitude"`
	Latitude       float64        `bigquery:"latitude"`
	Bearing        float64        `bigquery:"bearing"`
	RecordedAtTime civil.DateTime `bigquery:"recorded_at_time"`
	Occupancy      string         `bigquery:"occupancy"`
}

func (bq *BigQueryVehicleEvent) EventType() string        { return BQ_VEHICLE_EVENT }
func (bq *BigQueryVehicleEvent) SetTimeStamp(t time.Time) { bq.Timestamp = t }
func (bq *BigQueryVehicleEvent) SetUUID(u string)         { bq.UUID = u }

var bqVehicleSchema = bigquery.Schema{
	{Name: "uuid", Required: false, Type: bigquery.StringFieldType},
	{Name: "timestamp", Required: false, Type: bigquery.TimestampFieldType},
	{Name: "id", Required: false, Type: bigquery.StringFieldType},
	{Name: "longitude", Required: false, Type: bigquery.FloatFieldType},
	{Name: "latitude", Required: false, Type: bigquery.FloatFieldType},
	{Name: "bearing", Required: false, Type: bigquery.FloatFieldType},
	{Name: "recorded_at_time", Required: false, Type: bigquery.DateTimeFieldType},
	{Name: "codes", Repeated: true, Type: bigquery.StringFieldType},
	{Name: "occupancy", Required: false, Type: bigquery.StringFieldType},
}

type BigQueryLongTermStopVisitEvent struct {
	UUID      string    `bigquery:"uuid"`
	Timestamp time.Time `bigquery:"timestamp"`

	StopVisitUUID      string                 `bigquery:"stop_visit_uuid"`
	PassageOrder       int                    `bigquery:"passage_order"`
	AimedDepartureTime bigquery.NullTimestamp `bigquery:"aimed_departure_time"`
	AimedArrivalTime   bigquery.NullTimestamp `bigquery:"aimed_arrival_time"`

	ExpectedDepartureTime bigquery.NullTimestamp `bigquery:"expected_departure_time"`
	ExpectedArrivalTime   bigquery.NullTimestamp `bigquery:"expected_arrival_time"`

	ActualDepartureTime bigquery.NullTimestamp `bigquery:"actual_departure_time"`
	ActualArrivalTime   bigquery.NullTimestamp `bigquery:"actual_arrival_time"`

	DepartureStatus string `bigquery:"departure_status"`
	ArrivalStatus   string `bigquery:"arrival_status"`

	StopAreaName        string `bigquery:"stop_area_name"`
	StopAreaCodes       []Code `bigquery:"stop_area_codes"`
	StopAreaCoordinates string `bigquery:"stop_area_coordinates"`

	LineName      string `bigquery:"line_name"`
	LineNumber    string `bigquery:"line_number"`
	TransportMode string `bigquery:"transport_mode"`
	LineCodes     []Code `bigquery:"line_codes"`

	VehicleJourneyDirectionType   string `bigquery:"vehicle_journey_direction_type"`
	VehicleJourneyOriginName      string `bigquery:"vehicle_journey_origin_name"`
	VehicleJourneyDestinationName string `bigquery:"vehicle_journey_destination_name"`
	VehicleJourneyCodes           []Code `bigquery:"vehicle_journey_codes"`
	VehicleDriverRef              string `bigquery:"vehicle_driver_ref"`
	VehicleOccupancy              string `bigquery:"vehicle_occupancy"`
}

type Code struct {
	CodeSpace string `bigquery:"code_space"`
	Value     string `bigquery:"value"`
}

func (bq *BigQueryLongTermStopVisitEvent) EventType() string        { return BQ_LONG_TERM_STOP_VISIT_EVENT }
func (bq *BigQueryLongTermStopVisitEvent) SetTimeStamp(t time.Time) { bq.Timestamp = t }
func (bq *BigQueryLongTermStopVisitEvent) SetUUID(u string)         { bq.UUID = u }

var bqLongTermStopVisitsSchema = bigquery.Schema{
	{Name: "uuid", Required: true, Type: bigquery.StringFieldType},
	{Name: "timestamp", Required: true, Type: bigquery.TimestampFieldType},

	{Name: "stop_visit_uuid", Required: true, Type: bigquery.StringFieldType},
	{Name: "passage_order", Required: false, Type: bigquery.BigNumericFieldType},

	{Name: "aimed_departure_time", Required: false, Type: bigquery.TimestampFieldType},
	{Name: "aimed_arrival_time", Required: false, Type: bigquery.TimestampFieldType},

	{Name: "expected_departure_time", Required: false, Type: bigquery.TimestampFieldType},
	{Name: "expected_arrival_time", Required: false, Type: bigquery.TimestampFieldType},

	{Name: "actual_departure_time", Required: false, Type: bigquery.TimestampFieldType},
	{Name: "actual_arrival_time", Required: false, Type: bigquery.TimestampFieldType},

	{Name: "departure_status", Required: true, Type: bigquery.StringFieldType},
	{Name: "arrival_status", Required: true, Type: bigquery.StringFieldType},

	{Name: "stop_area_name", Required: false, Type: bigquery.StringFieldType},
	{Name: "stop_area_codes",
		Required: false,
		Repeated: true,
		Type:     bigquery.RecordFieldType,
		Schema: bigquery.Schema{
			{Name: "code_space", Type: bigquery.StringFieldType},
			{Name: "value", Type: bigquery.StringFieldType},
		},
	},

	{Name: "stop_area_coordinates", Required: false, Type: bigquery.GeographyFieldType},

	{Name: "line_name", Required: false, Type: bigquery.StringFieldType},
	{Name: "line_number", Required: false, Type: bigquery.StringFieldType},
	{Name: "transport_mode", Required: false, Type: bigquery.StringFieldType},
	{Name: "line_codes",
		Required: false,
		Repeated: true,
		Type:     bigquery.RecordFieldType,
		Schema: bigquery.Schema{
			{Name: "code_space", Type: bigquery.StringFieldType},
			{Name: "value", Type: bigquery.StringFieldType},
		},
	},

	{Name: "vehicle_journey_direction_type", Required: false, Type: bigquery.StringFieldType},
	{Name: "vehicle_journey_origin_name", Required: false, Type: bigquery.StringFieldType},
	{Name: "vehicle_journey_destination_name", Required: false, Type: bigquery.StringFieldType},

	{Name: "vehicle_journey_codes",
		Required: false,
		Repeated: true,
		Type:     bigquery.RecordFieldType,
		Schema: bigquery.Schema{
			{Name: "code_space", Type: bigquery.StringFieldType},
			{Name: "value", Type: bigquery.StringFieldType},
		},
	},

	{Name: "vehicle_driver_ref", Required: false, Type: bigquery.StringFieldType},
	{Name: "vehicle_occupancy", Required: false, Type: bigquery.StringFieldType},
}

type BigQueryControlEvent struct {
	UUID                             string    `bigquery:"uuid"`
	Timestamp                        time.Time `bigquery:"timestamp"`
	Criticity                        string    `bigquery:"criticity"`
	ControlType                      string    `bigquery:"control_type"`
	InternalCode                     string    `bigquery:"internal_code"`
	TargetModelClass                 string    `bigquery:"target_model_class"`
	TargetModelUUID                  string    `bigquery:"target_model_uuid"`
	TranslationInfoMessageKey        string    `bigquery:"translation_info_message_key"`
	TranslationInfoMessageAttributes string    `bigquery:"translation_info_message_attributes"`
}

func (bq *BigQueryControlEvent) EventType() string        { return BQ_CONTROL_EVENT }
func (bq *BigQueryControlEvent) SetTimeStamp(t time.Time) { bq.Timestamp = t }
func (bq *BigQueryControlEvent) SetUUID(u string)         { bq.UUID = u }

var bqControlSchema = bigquery.Schema{
	{Name: "uuid", Required: false, Type: bigquery.StringFieldType},
	{Name: "timestamp", Required: false, Type: bigquery.TimestampFieldType},
	{Name: "criticity", Required: false, Type: bigquery.StringFieldType},
	{Name: "control_type", Required: false, Type: bigquery.StringFieldType},
	{Name: "internal_code", Required: false, Type: bigquery.StringFieldType},
	{Name: "target_model_class", Required: false, Type: bigquery.StringFieldType},
	{Name: "target_model_uuid", Required: false, Type: bigquery.StringFieldType},
	{Name: "translation_info_message_key", Required: false, Type: bigquery.StringFieldType},
	{Name: "translation_info_message_attributes", Required: false, Type: bigquery.StringFieldType},
}
