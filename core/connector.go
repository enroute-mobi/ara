package core

import (
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/gtfs"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

const (
	SIRI_PARTNER = "siri-partner"

	// Connectors
	PUSH_COLLECTOR                                 = "push-collector"
	SIRI_STOP_POINTS_DISCOVERY_REQUEST_COLLECTOR   = "siri-stop-points-discovery-request-collector"
	SIRI_STOP_POINTS_DISCOVERY_REQUEST_BROADCASTER = "siri-stop-points-discovery-request-broadcaster"
	SIRI_LINES_DISCOVERY_REQUEST_COLLECTOR         = "siri-lines-discovery-request-collector"
	SIRI_LINES_DISCOVERY_REQUEST_BROADCASTER       = "siri-lines-discovery-request-broadcaster"
	SIRI_SERVICE_REQUEST_BROADCASTER               = "siri-service-request-broadcaster"
	SIRI_STOP_MONITORING_REQUEST_COLLECTOR         = "siri-stop-monitoring-request-collector"
	SIRI_LITE_STOP_MONITORING_REQUEST_COLLECTOR    = "siri-lite-stop-monitoring-request-collector"
	TEST_STOP_MONITORING_REQUEST_COLLECTOR         = "test-stop-monitoring-request-collector"
	SIRI_STOP_MONITORING_REQUEST_BROADCASTER       = "siri-stop-monitoring-request-broadcaster"
	SIRI_STOP_MONITORING_SUBSCRIPTION_COLLECTOR    = "siri-stop-monitoring-subscription-collector"
	SIRI_STOP_MONITORING_SUBSCRIPTION_BROADCASTER  = "siri-stop-monitoring-subscription-broadcaster"
	TEST_STOP_MONITORING_SUBSCRIPTION_BROADCASTER  = "siri-stop-monitoring-subscription-broadcaster-test"
	SIRI_GENERAL_MESSAGE_REQUEST_COLLECTOR         = "siri-general-message-request-collector"
	SIRI_GENERAL_MESSAGE_REQUEST_BROADCASTER       = "siri-general-message-request-broadcaster"
	SIRI_GENERAL_MESSAGE_SUBSCRIPTION_COLLECTOR    = "siri-general-message-subscription-collector"
	SIRI_GENERAL_MESSAGE_SUBSCRIPTION_BROADCASTER  = "siri-general-message-subscription-broadcaster"
	TEST_GENERAL_MESSAGE_SUBSCRIPTION_BROADCASTER  = "siri-general-message-subscription-broadcaster-test"
	// SIRI_ESTIMATED_TIMETABLE_REQUEST_COLLECTOR        = "siri-estimated-timetable-request-collector"
	SIRI_ESTIMATED_TIMETABLE_REQUEST_BROADCASTER      = "siri-estimated-timetable-request-broadcaster"
	SIRI_ESTIMATED_TIMETABLE_SUBSCRIPTION_COLLECTOR   = "siri-estimated-timetable-subscription-collector"
	SIRI_ESTIMATED_TIMETABLE_SUBSCRIPTION_BROADCASTER = "siri-estimated-timetable-subscription-broadcaster"
	TEST_ESTIMATED_TIMETABLE_SUBSCRIPTION_BROADCASTER = "siri-estimated-timetable-subscription-broadcaster-test"
	// SIRI_PRODUCTION_TIMETABLE_REQUEST_BROADCASTER      = "siri-production-timetable-request-broadcaster"
	SIRI_PRODUCTION_TIMETABLE_SUBSCRIPTION_BROADCASTER = "siri-production-timetable-subscription-broadcaster"
	SIRI_SITUATION_EXCHANGE_REQUEST_BROADCASTER        = "siri-situation-exchange-request-broadcaster"
	SIRI_SITUATION_EXCHANGE_REQUEST_COLLECTOR          = "siri-situation-exchange-request-collector"
	SIRI_SITUATION_EXCHANGE_SUBSCRIPTION_COLLECTOR     = "siri-situation-exchange-subscription-collector"
	SIRI_VEHICLE_MONITORING_REQUEST_COLLECTOR          = "siri-vehicle-monitoring-request-collector"
	SIRI_VEHICLE_MONITORING_SUBSCRIPTION_COLLECTOR     = "siri-vehicle-monitoring-subscription-collector"
	SIRI_SUBSCRIPTION_REQUEST_DISPATCHER               = "siri-subscription-request-dispatcher"
	SIRI_CHECK_STATUS_CLIENT_TYPE                      = "siri-check-status-client"
	TEST_CHECK_STATUS_CLIENT_TYPE                      = "test-check-status-client"
	SIRI_CHECK_STATUS_SERVER_TYPE                      = "siri-check-status-server"
	SIRI_LITE_VEHICLE_MONITORING_REQUEST_BROADCASTER   = "siri-lite-vehicle-monitoring-request-broadcaster"
	SIRI_VEHICLE_MONITORING_REQUEST_BROADCASTER        = "siri-vehicle-monitoring-request-broadcaster"
	SIRI_VEHICLE_MONITORING_SUBSCRIPTION_BROADCASTER   = "siri-vehicle-monitoring-subscription-broadcaster"
	TEST_VEHICLE_MONITORING_SUBSCRIPTION_BROADCASTER   = "siri-vehicle-monitoring-subscription-broadcaster-test"
	TEST_VALIDATION_CONNECTOR                          = "test-validation-connector"
	TEST_STARTABLE_CONNECTOR                           = "test-startable-connector-connector"
	GTFS_RT_REQUEST_COLLECTOR                          = "gtfs-rt-request-collector"
	GTFS_RT_TRIP_UPDATES_BROADCASTER                   = "gtfs-rt-trip-updates-broadcaster"
	GTFS_RT_VEHICLE_POSITIONS_BROADCASTER              = "gtfs-rt-vehicle-positions-broadcaster"
)

type Connector interface {
	clock.ClockInterface
	uuid.UUIDInterface

	Partner() *Partner
}

type GtfsConnector interface {
	HandleGtfs(*gtfs.FeedMessage)
}

type connector struct {
	uuid.UUIDConsumer
	clock.ClockConsumer

	partner         *Partner
	remoteCodeSpace string
}

func (c *connector) Partner() *Partner {
	return c.partner
}

type ConnectorFactory interface {
	Validate(*APIPartner) // Validate() fill the Errors of the APIPartner
	CreateConnector(*Partner) Connector
}

func NewConnectorFactory(connectorType string) ConnectorFactory {
	switch connectorType {
	case PUSH_COLLECTOR:
		return &PushCollectorFactory{}
	case SIRI_STOP_POINTS_DISCOVERY_REQUEST_COLLECTOR:
		return &SIRIStopPointsDiscoveryRequestCollectorFactory{}
	case SIRI_STOP_POINTS_DISCOVERY_REQUEST_BROADCASTER:
		return &SIRIStopPointsDiscoveryRequestBroadcasterFactory{}
	case SIRI_LINES_DISCOVERY_REQUEST_COLLECTOR:
		return &SIRILinesDiscoveryRequestCollectorFactory{}
	case SIRI_LINES_DISCOVERY_REQUEST_BROADCASTER:
		return &SIRILinesDiscoveryRequestBroadcasterFactory{}
	case SIRI_SERVICE_REQUEST_BROADCASTER:
		return &SIRIServiceRequestBroadcasterFactory{}
	case SIRI_STOP_MONITORING_REQUEST_COLLECTOR:
		return &SIRIStopMonitoringRequestCollectorFactory{}
	case SIRI_LITE_STOP_MONITORING_REQUEST_COLLECTOR:
		return &SIRILiteStopMonitoringRequestCollectorFactory{}
	case TEST_STOP_MONITORING_REQUEST_COLLECTOR:
		return &TestStopMonitoringRequestCollectorFactory{}
	case SIRI_STOP_MONITORING_REQUEST_BROADCASTER:
		return &SIRIStopMonitoringRequestBroadcasterFactory{}
	case SIRI_STOP_MONITORING_SUBSCRIPTION_COLLECTOR:
		return &SIRIStopMonitoringSubscriptionCollectorFactory{}
	case SIRI_STOP_MONITORING_SUBSCRIPTION_BROADCASTER:
		return &SIRIStopMonitoringSubscriptionBroadcasterFactory{}
	case TEST_STOP_MONITORING_SUBSCRIPTION_BROADCASTER:
		return &TestSIRIStopMonitoringSubscriptionBroadcasterFactory{}
	case SIRI_GENERAL_MESSAGE_REQUEST_COLLECTOR:
		return &SIRIGeneralMessageRequestCollectorFactory{}
	case SIRI_GENERAL_MESSAGE_REQUEST_BROADCASTER:
		return &SIRIGeneralMessageRequestBroadcasterFactory{}
	case SIRI_GENERAL_MESSAGE_SUBSCRIPTION_COLLECTOR:
		return &SIRIGeneralMessageSubscriptionCollectorFactory{}
	case SIRI_GENERAL_MESSAGE_SUBSCRIPTION_BROADCASTER:
		return &SIRIGeneralMessageSubscriptionBroadcasterFactory{}
	case TEST_GENERAL_MESSAGE_SUBSCRIPTION_BROADCASTER:
		return &TestSIRIGeneralMessageSubscriptionBroadcasterFactory{}
	// case SIRI_ESTIMATED_TIMETABLE_REQUEST_COLLECTOR:
	// 	return &SIRIEstimatedTimetableRequestCollectorFactory{}
	case SIRI_ESTIMATED_TIMETABLE_REQUEST_BROADCASTER:
		return &SIRIEstimatedTimetableRequestBroadcasterFactory{}
	case SIRI_ESTIMATED_TIMETABLE_SUBSCRIPTION_COLLECTOR:
		return &SIRIEstimatedTimetableSubscriptionCollectorFactory{}
	case SIRI_ESTIMATED_TIMETABLE_SUBSCRIPTION_BROADCASTER:
		return &SIRIEstimatedTimetableSubscriptionBroadcasterFactory{}
	case TEST_ESTIMATED_TIMETABLE_SUBSCRIPTION_BROADCASTER:
		return &TestSIRIETTSubscriptionBroadcasterFactory{}
	// case SIRI_PRODUCTION_TIMETABLE_REQUEST_BROADCASTER:
	// 	return &SIRIProductionTimetableSubscriptionBroadcasterFactory{}
	case SIRI_PRODUCTION_TIMETABLE_SUBSCRIPTION_BROADCASTER:
		return &SIRIProductionTimetableSubscriptionBroadcasterFactory{}
	case SIRI_SITUATION_EXCHANGE_REQUEST_BROADCASTER:
		return &SIRISituationExchangeRequestBroadcasterFactory{}
	case SIRI_SITUATION_EXCHANGE_SUBSCRIPTION_COLLECTOR:
		return &SIRISituationExchangeSubscriptionCollectorFactory{}
	case SIRI_SITUATION_EXCHANGE_REQUEST_COLLECTOR:
		return &SIRISituationExchangeRequestCollectorFactory{}
	case SIRI_VEHICLE_MONITORING_REQUEST_COLLECTOR:
		return &SIRIVehicleMonitoringRequestCollectorFactory{}
	case SIRI_VEHICLE_MONITORING_SUBSCRIPTION_COLLECTOR:
		return &SIRIVehicleMonitoringSubscriptionCollectorFactory{}
	case SIRI_CHECK_STATUS_CLIENT_TYPE:
		return &SIRICheckStatusClientFactory{}
	case SIRI_SUBSCRIPTION_REQUEST_DISPATCHER:
		return &SIRISubscriptionRequestDispatcherFactory{}
	case TEST_CHECK_STATUS_CLIENT_TYPE:
		return &TestCheckStatusClientFactory{}
	case SIRI_CHECK_STATUS_SERVER_TYPE:
		return &SIRICheckStatusServerFactory{}
	case SIRI_LITE_VEHICLE_MONITORING_REQUEST_BROADCASTER:
		return &SIRILiteVehicleMonitoringRequestBroadcasterFactory{}
	case SIRI_VEHICLE_MONITORING_REQUEST_BROADCASTER:
		return &SIRIVehicleMonitoringRequestBroadcasterFactory{}
	case SIRI_VEHICLE_MONITORING_SUBSCRIPTION_BROADCASTER:
		return &SIRIVehicleMonitoringSubscriptionBroadcasterFactory{}
	case TEST_VEHICLE_MONITORING_SUBSCRIPTION_BROADCASTER:
		return &TestSIRIVMSubscriptionBroadcasterFactory{}
	case GTFS_RT_REQUEST_COLLECTOR:
		return &GtfsRequestCollectorFactory{}
	case GTFS_RT_TRIP_UPDATES_BROADCASTER:
		return &TripUpdatesBroadcasterFactory{}
	case GTFS_RT_VEHICLE_POSITIONS_BROADCASTER:
		return &VehiclePositionBroadcasterFactory{}
	case TEST_VALIDATION_CONNECTOR:
		return &TestValidationFactory{}
	case TEST_STARTABLE_CONNECTOR:
		return &TestStartableFactory{}
	default:
		return nil
	}
}

type TestValidationFactory struct{}
type TestValidationConnector struct {
	connector
}

func (factory *TestValidationFactory) Validate(apiPartner *APIPartner) {
	if apiPartner.Slug == PartnerSlug("invalid_slug") {
		apiPartner.Errors.Add("slug", "Invalid format")
	}
}

func (factory *TestValidationFactory) CreateConnector(partner *Partner) Connector {
	return &TestValidationConnector{}
}

type TestStartableFactory struct{}
type TestStartableConnector struct {
	connector

	started bool
}

func (factory *TestStartableFactory) Validate(apiPartner *APIPartner) {} // Always valid

func (factory *TestStartableFactory) CreateConnector(partner *Partner) Connector {
	return &TestStartableConnector{}
}

func (connector *TestStartableConnector) Start() {
	connector.started = true
}

func (connector *TestStartableConnector) Stop() {
	connector.started = false
}
