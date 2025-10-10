package partners

import "time"

type Partner interface {
	Id() Id
	FromTemplate() Id
}

type OperationnalStatus string

const (
	OperationnalStatusUnknown OperationnalStatus = "unknown"
	OperationnalStatusUp      OperationnalStatus = "up"
	OperationnalStatusDown    OperationnalStatus = "down"
)

type Id string
type Slug string

type SlugAndCredentialsHandler interface {
	UniqCredentials(Id, string, ...string) bool
	UniqSlug(Id, Slug) bool
}

type Status struct {
	OperationnalStatus OperationnalStatus
	RetryCount         int
	ServiceStartedAt   time.Time
}

type StatusCheck struct {
	LastCheck time.Time
	Status    OperationnalStatus
}
