package gql

import (
	"bitbucket.org/enroute-mobi/ara/core"
	"bitbucket.org/enroute-mobi/ara/model"
	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/errors"
)

type Resolver struct {
	Partner *core.Partner
}

func (r *Resolver) Vehicles() (l []*vehicleResolver) {
	for _, v := range r.Partner.Model().Vehicles().FindAll() {
		l = append(l, r.resolverFromVehicle(v))
	}
	return l
}

func (r *Resolver) Vehicle(args struct{ Code string }) (*vehicleResolver, error) {
	c := model.NewCode(r.Partner.RemoteCodeSpace(), args.Code)
	v, ok := r.Partner.Model().Vehicles().FindByCode(c)
	if !ok {
		return nil, errors.Errorf("Can't find Vehicle with code %s", args.Code)
	}
	return r.resolverFromVehicle(v), nil
}

func (r *Resolver) UpdateVehicle(args struct {
	Code  string
	Input vehicleInput
}) (*vehicleResolver, error) {
	c := model.NewCode(r.Partner.RemoteCodeSpace(), args.Code)
	v, ok := r.Partner.Model().Vehicles().FindByCode(c)
	if !ok {
		return nil, errors.Errorf("Can't find Vehicle with code %s", args.Code)
	}
	if args.Input.OccupancyStatus != nil && r.Partner.IsMutable(OccupancyStatus) {
		v.Occupancy = *args.Input.OccupancyStatus
	}
	if args.Input.OccupancyRate != nil && r.Partner.IsMutable(OccupancyRate) {
		v.Percentage = *args.Input.OccupancyRate
	}

	v.Save()

	return r.resolverFromVehicle(v), nil
}

func (r *Resolver) resolverFromVehicle(v *model.Vehicle) (res *vehicleResolver) {
	code, _ := v.Code(r.Partner.RemoteCodeSpace())
	res = &vehicleResolver{
		v: &vehicle{
			ID:              graphql.ID(v.Id()),
			StopArea:        graphql.ID(v.StopAreaId),
			Line:            graphql.ID(v.LineId),
			VehicleJourney:  graphql.ID(v.VehicleJourneyId),
			NextStopVisit:   graphql.ID(v.NextStopVisitId),
			Code:            code.Value(),
			OccupancyStatus: v.Occupancy,
			DriverRef:       v.DriverRef,
			LinkDistance:    v.LinkDistance,
			OccupancyRate:   v.Percentage,
			Longitude:       v.Longitude,
			Latitude:        v.Latitude,
			Bearing:         v.Bearing,
			RecordedAt:      graphql.Time{Time: v.RecordedAtTime},
			ValidUntil:      graphql.Time{Time: v.ValidUntilTime},
		}}
	return
}
