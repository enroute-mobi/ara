package core

import (
	"testing"
	"time"

	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
)

func Test_PartnerGuardian_Run(t *testing.T) {
	logger.Log.Debug = true
	partners := NewPartnerManager()
	partner := Partner{
		ConnectorTypes: []string{"test-check-status-client"},
		connectors:     make(map[string]Connector),
	}
	partner.RefreshConnectors()
	partners.Save(&partner)

	logger.Log.Debugf("Partner: %v", partner)
	logger.Log.Debugf("Partner connectors: %v", partner.connectors)

	fakeClock := model.NewFakeClock()
	partners.Guardian().SetClock(fakeClock)

	partners.Start()
	defer partners.Stop()

	// Wait for the guardian to launch Run
	fakeClock.BlockUntil(1)
	// Advance time
	fakeClock.Advance(31 * time.Second)
	// Wait for the Test CheckStatusClient to finish Status()
	select {
	case <-partner.CheckStatusClient().(*TestCheckStatusClient).Done:
		if partner.OperationnalStatus() != OPERATIONNAL_STATUS_UP {
			t.Errorf("Partner OperationnalStatus should be UP when guardian is running, got: %v", partner.OperationnalStatus())
		}
	case <-time.After(5 * time.Second):
		t.Errorf("Guardian CheckPartnerStatus with TestCheckStatusClient timed out")
	}

	// Test a change in status
	partner.CheckStatusClient().(*TestCheckStatusClient).SetStatus(OPERATIONNAL_STATUS_DOWN)
	fakeClock.Advance(31 * time.Second)
	select {
	case <-partner.CheckStatusClient().(*TestCheckStatusClient).Done:
		if partner.OperationnalStatus() != OPERATIONNAL_STATUS_DOWN {
			t.Errorf("Partner OperationnalStatus should be DOWN when guardian is running, got: %v", partner.OperationnalStatus())
		}
	case <-time.After(5 * time.Second):
		t.Errorf("Guardian CheckPartnerStatus with TestCheckStatusClient timed out")
	}
}
