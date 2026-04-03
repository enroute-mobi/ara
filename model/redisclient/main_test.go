package redisclient

import (
	"os"
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/config"
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
	if len(config.Config.CodeSpaces) == 0 {
		config.Config.CodeSpaces = []string{"internal", "external"}
	}

	if config.Config.RedisAddr == "" {
		return
	}

	TestClient, err = New("test", config.Config.CodeSpaces)
	if err != nil {
		panic(err)
	}
	err = TestClient.Start(time.Now())
	if err != nil {
		panic(err)
	}

	c := m.Run()

	if config.Config.RedisAddr != "" {
		TestClient.FlushAll()
		TestClient.Stop()
	}

	os.Exit(c)
}
