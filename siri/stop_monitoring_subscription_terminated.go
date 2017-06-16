package siri

import "time"

type ErrorCondition struct {
	number      int
	description string
}

type XMLStopMonitoringSubscriptionTerminated struct {
	responseTimestamp time.Time
	producerRef       string
	subscriberRef     string
	subscriptionRef   string
}
