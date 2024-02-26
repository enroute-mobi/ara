package gql

import (
	"context"
	"encoding/json"
	"testing"

	"bitbucket.org/enroute-mobi/ara/core"
	"bitbucket.org/enroute-mobi/ara/core/settings"
	"bitbucket.org/enroute-mobi/ara/model"
	graphql "github.com/graph-gophers/graphql-go"
	"github.com/stretchr/testify/assert"
)

type vehicleQueryResult struct {
	Vehicle v `json:"vehicle"`
}

type vehiclesQueryResult struct {
	Vehicles []v `json:"vehicles"`
}

type mutationResult struct {
	UpdateVehicle v `json:"updateVehicle"`
}

type v struct {
	Id              string  `json:"id"`
	Code            string  `json:"code"`
	DriverRef       string  `json:"driverRef"`
	OccupancyStatus string  `json:"occupancyStatus"`
	OccupancyRate   float64 `json:"occupancyRate"`
}

func TestResolverQuery(t *testing.T) {
	assert := assert.New(t)

	referentials := core.NewMemoryReferentials()
	referential := referentials.New(core.ReferentialSlug("referential"))
	referentials.Save(referential)

	partner := referential.Partners().New("slug")
	s := map[string]string{
		"remote_code_space": "internal",
	}
	partner.PartnerSettings = settings.NewPartnerSettings(partner.UUIDGenerator, s)
	partner.ConnectorTypes = []string{"graphql-server"}
	partner.Save()

	v := partner.Model().Vehicles().New()
	v.SetCode(model.NewCode("internal", "1234"))
	v.Occupancy = "seatsAvailable"
	v.Percentage = 25.6
	v.DriverRef = "Michel"
	v.Save()

	v2 := partner.Model().Vehicles().New()
	v2.SetCode(model.NewCode("internal", "5678"))
	v2.Occupancy = "noSeatsAvailable"
	v2.Percentage = 100
	v2.DriverRef = "Bob"
	v2.Save()

	schema := graphql.MustParseSchema(Schema, &Resolver{Partner: partner})

	query := `
		query {
			vehicles {
				id
				code
				occupancyStatus
				occupancyRate
				driverRef
			}
		}
	`

	res := schema.Exec(context.Background(), query, "", nil)

	/*
		// Test display for debug
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		err := enc.Encode(res)
		if err != nil {
			panic(err)
		}
	*/

	var vResults vehiclesQueryResult
	err := json.Unmarshal(res.Data, &vResults)
	if err != nil {
		panic(err)
	}
	assert.Len(vResults.Vehicles, 2)

	query = `
		query {
			vehicle(code: "1234") {
				id
				code
				occupancyStatus
				occupancyRate
				driverRef
			}
		}
	`

	res = schema.Exec(context.Background(), query, "", nil)

	var vResult vehicleQueryResult
	err = json.Unmarshal(res.Data, &vResult)
	if err != nil {
		panic(err)
	}
	assert.Equal("1234", vResult.Vehicle.Code)
	assert.Equal("seatsAvailable", vResult.Vehicle.OccupancyStatus)
	assert.Equal(25.6, vResult.Vehicle.OccupancyRate)
	assert.Equal("Michel", vResult.Vehicle.DriverRef)

	query = `
		query {
			vehicle(code: "5678") {
				id
				code
				occupancyStatus
				occupancyRate
				driverRef
			}
		}
	`

	res = schema.Exec(context.Background(), query, "", nil)

	err = json.Unmarshal(res.Data, &vResult)
	if err != nil {
		panic(err)
	}
	assert.Equal("5678", vResult.Vehicle.Code)
	assert.Equal("noSeatsAvailable", vResult.Vehicle.OccupancyStatus)
	assert.Equal(100.0, vResult.Vehicle.OccupancyRate)
	assert.Equal("Bob", vResult.Vehicle.DriverRef)

}

func TestResolverMutation(t *testing.T) {
	assert := assert.New(t)

	referentials := core.NewMemoryReferentials()
	referential := referentials.New(core.ReferentialSlug("referential"))
	referentials.Save(referential)

	partner := referential.Partners().New("slug")
	s := map[string]string{
		"remote_code_space":          "internal",
		"graphql.mutable_attributes": "vehicle.occupancyStatus,vehicle.occupancyRate",
	}
	partner.PartnerSettings = settings.NewPartnerSettings(partner.UUIDGenerator, s)
	partner.ConnectorTypes = []string{"graphql-server"}
	partner.Save()

	v := partner.Model().Vehicles().New()
	v.SetCode(model.NewCode("internal", "1234"))
	v.Occupancy = "noSeatsAvailable"
	v.Percentage = 25.6
	v.DriverRef = "Michel"
	v.Save()

	schema := graphql.MustParseSchema(Schema, &Resolver{Partner: partner})

	mutation := `mutation {
  updateVehicle(code: "1234", input: { occupancyStatus: "seatsAvailable", occupancyRate: 0.65 }) {
    code
    occupancyStatus
    occupancyRate
    driverRef
  }
}
	`

	res := schema.Exec(context.Background(), mutation, "", nil)

	var result mutationResult
	err := json.Unmarshal(res.Data, &result)
	if err != nil {
		panic(err)
	}
	assert.Equal("1234", result.UpdateVehicle.Code)
	assert.Equal("seatsAvailable", result.UpdateVehicle.OccupancyStatus)
	assert.Equal(0.65, result.UpdateVehicle.OccupancyRate)
	assert.Equal("Michel", result.UpdateVehicle.DriverRef)

	mutation = `mutation {
  updateVehicle(code: "1234", input: { occupancyStatus: "notALotOfSeatsAvailable" }) {
    code
    occupancyStatus
    occupancyRate
    driverRef
  }
}
	`

	res = schema.Exec(context.Background(), mutation, "", nil)

	err = json.Unmarshal(res.Data, &result)
	if err != nil {
		panic(err)
	}
	assert.Equal("1234", result.UpdateVehicle.Code)
	assert.Equal("notALotOfSeatsAvailable", result.UpdateVehicle.OccupancyStatus)
	assert.Equal(0.65, result.UpdateVehicle.OccupancyRate)
	assert.Equal("Michel", result.UpdateVehicle.DriverRef)

	mutation = `mutation {
  updateVehicle(code: "1234", input: { occupancyRate: 34.6 }) {
    code
    occupancyStatus
    occupancyRate
    driverRef
  }
}
	`

	res = schema.Exec(context.Background(), mutation, "", nil)

	err = json.Unmarshal(res.Data, &result)
	if err != nil {
		panic(err)
	}
	assert.Equal("1234", result.UpdateVehicle.Code)
	assert.Equal("notALotOfSeatsAvailable", result.UpdateVehicle.OccupancyStatus)
	assert.Equal(34.6, result.UpdateVehicle.OccupancyRate)
	assert.Equal("Michel", result.UpdateVehicle.DriverRef)

}

func TestResolverMutationWithoutAuthorization(t *testing.T) {
	assert := assert.New(t)

	referentials := core.NewMemoryReferentials()
	referential := referentials.New(core.ReferentialSlug("referential"))
	referentials.Save(referential)

	partner := referential.Partners().New("slug")
	s := map[string]string{
		"remote_code_space":          "internal",
		"graphql.mutable_attributes": "vehicle.occupancyStatus",
	}
	partner.PartnerSettings = settings.NewPartnerSettings(partner.UUIDGenerator, s)
	partner.ConnectorTypes = []string{"graphql-server"}
	partner.Save()

	v := partner.Model().Vehicles().New()
	v.SetCode(model.NewCode("internal", "1234"))
	v.Occupancy = "noSeatsAvailable"
	v.Percentage = 25.6
	v.DriverRef = "Michel"
	v.Save()

	schema := graphql.MustParseSchema(Schema, &Resolver{Partner: partner})

	mutation := `mutation {
  updateVehicle(code: "1234", input: { occupancyStatus: "seatsAvailable", occupancyRate: 0.65 }) {
    code
    occupancyStatus
    occupancyRate
    driverRef
  }
}
	`

	res := schema.Exec(context.Background(), mutation, "", nil)

	var result mutationResult
	err := json.Unmarshal(res.Data, &result)
	if err != nil {
		panic(err)
	}
	assert.Equal("1234", result.UpdateVehicle.Code)
	assert.Equal("seatsAvailable", result.UpdateVehicle.OccupancyStatus)
	assert.Equal(25.6, result.UpdateVehicle.OccupancyRate)
	assert.Equal("Michel", result.UpdateVehicle.DriverRef)
}
