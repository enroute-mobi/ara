package model

import (
	"os"
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/config"
	"bitbucket.org/enroute-mobi/ara/model/redisclient"
)

func TestMain(m *testing.M) {
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
func newTestModel(t testing.TB) Model {
	if config.Config.RedisAddr != "" {
		c, err := redisclient.New("referential", config.Config.CodeSpaces)
		if err != nil {
			panic(err)
		}
		err = c.Start(time.Now())
		if err != nil {
			panic(err)
		}

		t.Cleanup(c.FlushAll)
		return NewHybridModel("referential", c)
	}
	return NewTestMemoryModel()
}
