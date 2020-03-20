package api

import "bitbucket.org/enroute-mobi/ara/core"

func NewTestServer() *Server {
	server := Server{}
	referentials := core.NewMemoryReferentials()
	server.SetReferentials(referentials)
	server.startedTime = server.Clock().Now()
	return &server
}
