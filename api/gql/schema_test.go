package gql

import (
	"testing"

	graphql "github.com/graph-gophers/graphql-go"
)

func TestSchema(t *testing.T) {
	graphql.MustParseSchema(Schema, &Resolver{})
}
