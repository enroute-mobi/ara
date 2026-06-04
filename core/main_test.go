package core

import (
	"os"
	"testing"

	"bitbucket.org/enroute-mobi/ara/config"
	"bitbucket.org/enroute-mobi/ara/uuid"
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
func newTestReferential(t testing.TB, testCollectManager ...bool) (*MemoryReferentials, *Referential) {
	referentials := NewMemoryReferentials()
	referentials.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	referential := referentials.New("referential")
	if len(testCollectManager) != 0 {
		referential.collectManager = NewTestCollectManager()
	}
	referential.Save()
	referential.StartRedisClient()

	if config.Config.RedisAddr != "" {
		t.Cleanup(referential.RedisClient().FlushAll)
	}

	return referentials, referential
}

func newTestPartnerManager(t testing.TB) *PartnerManager {
	_, r := newTestReferential(t)
	return r.Partners().(*PartnerManager)
}

// func newTestPartner(t *testing.T) (*Referential, *Partner) {
// 	_, r := newTestReferential(t)
// 	return r, r.Partners().New("partner")
// }
