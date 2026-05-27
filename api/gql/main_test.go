package gql

import (
	"os"
	"testing"

	"bitbucket.org/enroute-mobi/ara/config"
	"bitbucket.org/enroute-mobi/ara/core"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

func TestMain(m *testing.M) {
	config.SetEnvironment("test")
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

func newTestReferential(t *testing.T) (*core.MemoryReferentials, *core.Referential) {
	referentials := core.NewMemoryReferentials()
	referentials.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	referential := referentials.New("referential")
	referential.Save()
	referential.StartRedisClient()

	if config.Config.RedisAddr != "" {
		t.Cleanup(referential.RedisClient().FlushAll)
	}

	return referentials, referential
}
