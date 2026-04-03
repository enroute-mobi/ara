package api

import (
	"os"
	"testing"

	"bitbucket.org/enroute-mobi/ara/config"
	"bitbucket.org/enroute-mobi/ara/core"
	"bitbucket.org/enroute-mobi/ara/logger"
)

func TestMain(m *testing.M) {
	logger.Log.Print("Init api test")
	config.SetEnvironment("test")
	// Load configuration
	err := config.LoadConfig("")
	if err != nil {
		panic(err)
	}
	config.Config.ApiKey = ""

	c := m.Run()

	os.Exit(c)
}

// Default will create with 2 codespace values: internal and external
func newTestReferential(t *testing.T, tokens ...string) (*core.MemoryReferentials, *core.Referential) {
	if len(config.Config.CodeSpaces) == 0 {
		config.Config.CodeSpaces = []string{"internal", "external"}
	}
	referentials := core.NewMemoryReferentials()
	referential := referentials.New("referential")
	referential.Tokens = tokens
	referential.Save()

	if config.Config.RedisAddr != "" {
		t.Cleanup(referential.RedisClient().FlushAll)
	}

	return referentials, referential
}

func newTestServer(t *testing.T, tokens ...string) (*Server, *core.Referential) {
	if len(config.Config.CodeSpaces) == 0 {
		config.Config.CodeSpaces = []string{"internal", "external"}
	}
	referentials := core.NewMemoryReferentials()
	server := &Server{}
	server.SetReferentials(referentials)
	server.startedTime = server.Clock().Now()

	referential := server.CurrentReferentials().New("referential")
	referential.Tokens = tokens
	referential.Save()

	if config.Config.RedisAddr != "" {
		t.Cleanup(referential.RedisClient().FlushAll)
	}

	return server, referential
}
