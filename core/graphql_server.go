package core

type GraphqlServerFactory struct{}

func (factory *GraphqlServerFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfLocalCredentials()
}

func (factory *GraphqlServerFactory) CreateConnector(partner *Partner) Connector {
	return NewGtfsRequestCollector(partner)
}

// Empty shell for now
type GraphqlServer struct {
	connector
}

func NewGraphqlServer(partner *Partner) *GraphqlServer {
	graphqlServer := &GraphqlServer{}
	graphqlServer.partner = partner
	return graphqlServer
}
