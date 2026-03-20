package model

import (
	"os"
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/config"
	"bitbucket.org/enroute-mobi/ara/model/redisclient"
)

func TestMain(m *testing.M) {
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
		defer func() {
			redisclient.TestClient.FlushAll()
			redisclient.TestClient.Stop()
		}()
	}

	c := m.Run()

	os.Exit(c)
}
