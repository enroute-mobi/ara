package api

import (
	"os"
	"testing"

	"bitbucket.org/enroute-mobi/ara/config"
	"bitbucket.org/enroute-mobi/ara/core"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/uuid"
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
	if len(config.Config.CodeSpaces) == 0 {
		config.Config.CodeSpaces = []string{"internal", "external"}
	}

	c := m.Run()

	os.Exit(c)
}

// Default will create with 2 codespace values: internal and external
func newTestReferential(t testing.TB, tokens ...string) (*core.MemoryReferentials, *core.Referential) {
	referentials := core.NewMemoryReferentials()
	referentials.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	referential := referentials.New("referential")
	referential.Tokens = tokens
	referential.Save()
	referential.StartRedisClient()

	if config.Config.RedisAddr != "" {
		t.Cleanup(referential.RedisClient().FlushAll)
	}

	return referentials, referential
}

func newTestServer(t testing.TB, g ...uuid.UUIDGenerator) (*Server, *core.Referential) {
	referentials := core.NewMemoryReferentials()
	if len(g) == 1 {
		referentials.SetUUIDGenerator(g[0])
	} else {
		referentials.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	}

	server := &Server{}
	server.SetReferentials(referentials)
	server.startedTime = server.Clock().Now()

	referential := server.CurrentReferentials().New("referential")
	referential.Save()
	referential.StartRedisClient()

	if config.Config.RedisAddr != "" {
		t.Cleanup(referential.RedisClient().FlushAll)
	}

	return server, referential
}
