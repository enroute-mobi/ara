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

	if config.Config.RedisAddr != "" {
		var err error
		redisclient.TestClient, err = redisclient.New("test", []string{})
		if err != nil {
			panic(err)
		}
		err = redisclient.TestClient.Start(time.Now())
		if err != nil {
			panic(err)
		}
	}

	c := m.Run()

	if config.Config.RedisAddr != "" {
		redisclient.TestClient.Stop()
	}

	os.Exit(c)
}

// Default will create with 2 codespace values: internal and external
func newTestModel(t *testing.T) Model {
	if len(config.Config.CodeSpaces) == 0 {
		config.Config.CodeSpaces = []string{"internal", "external"}
	}
	if config.Config.RedisAddr != "" {
		t.Cleanup(redisclient.TestClient.FlushAll)
		return NewTestHybridModel()
	}
	return NewTestMemoryModel()
}
