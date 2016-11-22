package core

import (
	"testing"

	"github.com/af83/edwig/model"
)

func Test_CollectManager_Partners(t *testing.T) {
	partners := NewPartnerManager(model.NewMemoryModel())
	collectManager := NewCollectManager(partners)

	if collectManager.Partners() != partners {
		t.Errorf("CollectManager Partners() should return correct value")
	}
}
