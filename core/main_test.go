package core

import (
	"os"
	"testing"

	"bitbucket.org/enroute-mobi/ara/config"
)

func TestMain(m *testing.M) {
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
func newTestReferential(t *testing.T, testCollectManager ...bool) (*MemoryReferentials, *Referential) {
	if len(config.Config.CodeSpaces) == 0 {
		config.Config.CodeSpaces = []string{"internal", "external"}
	}
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	if len(testCollectManager) != 0 {
		referential.collectManager = NewTestCollectManager()
	}
	referential.Save()

	if config.Config.RedisAddr != "" {
		t.Cleanup(referential.RedisClient().FlushAll)
	}

	return referentials, referential
}

func newTestPartnerManager(t *testing.T) *PartnerManager {
	_, r := newTestReferential(t)
	return r.Partners().(*PartnerManager)
}
