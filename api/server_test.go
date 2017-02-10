package api

import "github.com/af83/edwig/core"

func NewTestServer() *Server {
	server := Server{}
	referentials := core.NewMemoryReferentials()
	server.SetReferentials(referentials)
	server.startedTime = server.Clock().Now()
	return &server
}
