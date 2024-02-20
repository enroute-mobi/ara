package gql

import (
	graphql "github.com/graph-gophers/graphql-go"
)

const ( // Const for mutable attributes
	OccupancyStatus = "vehicle.occupancyStatus"
	OccupancyRate   = "vehicle.occupancyRate"
)

type vehicle struct {
	ID              graphql.ID
	StopArea        graphql.ID // Temp, should be a model when implemented
	Line            graphql.ID // Temp, should be a model when implemented
	VehicleJourney  graphql.ID // Temp, should be a model when implemented
	NextStopVisit   graphql.ID // Temp, should be a model when implemented
	Code            string
	OccupancyStatus string
	DriverRef       string
	LinkDistance    float64
	OccupancyRate   float64
	Longitude       float64
	Latitude        float64
	Bearing         float64
	RecordedAtTime  graphql.Time
	ValidUntilTime  graphql.Time
}

type vehicleInput struct {
	OccupancyStatus *string
	OccupancyRate   *float64
}

type vehicleResolver struct {
	v *vehicle
}

func (v *vehicleResolver) ID() graphql.ID {
	return v.v.ID
}
func (v *vehicleResolver) StopArea() graphql.ID {
	return v.v.StopArea
}
func (v *vehicleResolver) Line() graphql.ID {
	return v.v.Line
}
func (v *vehicleResolver) VehicleJourney() graphql.ID {
	return v.v.VehicleJourney
}
func (v *vehicleResolver) NextStopVisit() graphql.ID {
	return v.v.NextStopVisit
}
func (v *vehicleResolver) Code() string {
	return v.v.Code
}
func (v *vehicleResolver) OccupancyStatus() string {
	return v.v.OccupancyStatus
}
func (v *vehicleResolver) DriverRef() string {
	return v.v.DriverRef
}
func (v *vehicleResolver) LinkDistance() float64 {
	return v.v.LinkDistance
}
func (v *vehicleResolver) OccupancyRate() float64 {
	return v.v.OccupancyRate
}
func (v *vehicleResolver) Longitude() float64 {
	return v.v.Longitude
}
func (v *vehicleResolver) Latitude() float64 {
	return v.v.Latitude
}
func (v *vehicleResolver) Bearing() float64 {
	return v.v.Bearing
}
func (v *vehicleResolver) RecordedAtTime() graphql.Time {
	return v.v.RecordedAtTime
}
func (v *vehicleResolver) ValidUntilTime() graphql.Time {
	return v.v.ValidUntilTime
}
