package gql

var Schema = `
	scalar Time

	schema {
		query: Query
		mutation: Mutation
	}

	# The query type, represents all of the entry points into our object graph
	type Query {
		vehicles: [Vehicle]!
		vehicle(code: String!): Vehicle!
	}

	# The mutation type, represents all updates we can make to our data
	type Mutation {
		updateVehicle(code: String!, input: VehicleInput!): Vehicle!
	}

	type Vehicle {
		id:              ID!
		stopArea:        ID!
		line:            ID!
		vehicleJourney:  ID!
		nextStopVisit:   ID!
		code:            String!
		occupancyStatus: String!
		driverRef:       String!
		linkDistance:    Float!
		occupancyRate:   Float!
		longitude:       Float!
		latitude:        Float!
		bearing:         Float!
		recordedAt:      Time!
		validUntil:      Time!
	}

	input VehicleInput {
		occupancyStatus:	String
		occupancyRate:		Float
	}
`
