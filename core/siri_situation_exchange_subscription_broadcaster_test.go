package core

import (
	"fmt"
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/clock"
	s "bitbucket.org/enroute-mobi/ara/core/settings"
	"bitbucket.org/enroute-mobi/ara/model"
	"github.com/stretchr/testify/assert"
)

func Test_SituationExchangeBroadcaster_Create_Events(t *testing.T) {
	assert := assert.New(t)
	clock.SetDefaultClock(clock.NewFakeClock())

	referentials := NewMemoryReferentials()
	referential := referentials.New("Un Referential Plutot Cool")
	referential.model = model.NewMemoryModel()

	referential.model.SetBroadcastSXChan(referential.broacasterManager.GetSituationExchangeBroadcastEventChan())
	referential.broacasterManager.Start()
	defer referential.broacasterManager.Stop()

	partner := referential.Partners().New("Un Partner tout autant cool")

	settings := map[string]string{
		"remote_code_space": "internal",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)

	partner.ConnectorTypes = []string{SIRI_SITUATION_EXCHANGE_SUBSCRIPTION_BROADCASTER}
	partner.RefreshConnectors()
	referential.Partners().Save(partner)

	// Create right Lines & StopArea
	code := model.NewCode("internal", "first")

	lineFirst := referential.Model().Lines().New()
	lineFirst.SetCode(code)
	lineFirst.Save()

	stopAreaFirst := referential.Model().StopAreas().New()
	stopAreaFirst.SetCode(code)
	stopAreaFirst.Save()

	// Create Dummy lines & StopArea
	code = model.NewCode("internal", "DUMMY")

	lineDummy := referential.Model().Lines().New()
	lineDummy.SetCode(code)
	lineDummy.Save()

	stopAreaDummy := referential.Model().StopAreas().New()
	stopAreaDummy.SetCode(code)
	stopAreaDummy.Save()

	code2 := model.NewCode("SituationResource", "Situation")
	reference := model.Reference{
		Code: &code2,
		Type: "Situation",
	}

	var TestCases = []struct {
		situationMultipleAffectedLines     bool
		situationAffectedLine              bool
		situationLineRef                   string
		situationMultipleAffectedStopPoint bool
		situationAffectedStopPoint         bool
		situationStopPointRef              string
		hasExternalResource                bool
		subscriptionLineValue              string
		subscriptionStopPointValue         string
		expectedBroadcastedEvent           int
		message                            string
	}{
		{
			situationAffectedLine:      false,
			situationAffectedStopPoint: false,
			hasExternalResource:        false,
			expectedBroadcastedEvent:   1,
			message: `Should broadcast the situation without
affected lines or affected StopPoint if no external resource is set`,
		},
		{
			situationAffectedLine:      true,
			situationLineRef:           "DUMMY",
			situationAffectedStopPoint: true,
			situationStopPointRef:      "first",
			hasExternalResource:        false,
			expectedBroadcastedEvent:   1,
			message: `Should broadcast the situation with
affected Line and with affected StopPoint if no external resource is set`,
		},
		{
			situationAffectedLine:      true,
			situationLineRef:           "DUMMY",
			situationAffectedStopPoint: false,
			hasExternalResource:        false,
			expectedBroadcastedEvent:   1,
			message: `Should broadcast the situation with
affected Line if no external resource is set`,
		},
		{
			situationAffectedLine:      true,
			situationLineRef:           "DUMMY",
			situationAffectedStopPoint: false,
			hasExternalResource:        true,
			subscriptionLineValue:      "first",
			expectedBroadcastedEvent:   0,
			message: `Should NOT broadcast the situation with
affected Line if the external resource does not match`,
		},
		{
			situationAffectedLine:      true,
			situationLineRef:           "first",
			situationAffectedStopPoint: false,
			hasExternalResource:        true,
			subscriptionLineValue:      "first",
			expectedBroadcastedEvent:   1,
			message: `Should broadcast the situation with
affected Line if the external resource match`,
		},
		{
			situationMultipleAffectedLines: true,
			hasExternalResource:            true,
			subscriptionLineValue:          "first",
			expectedBroadcastedEvent:       1,
			message: `Should broadcast the situation with
multiple affected Lines if at least one of the external resource match`,
		},
		{
			situationAffectedLine:      false,
			situationAffectedStopPoint: true,
			situationStopPointRef:      "DUMMY",
			hasExternalResource:        false,
			expectedBroadcastedEvent:   1,
			message: `Should broadcast the situation with
affected StopPoint if no external resource is set`,
		},
		{
			situationAffectedLine:      false,
			situationAffectedStopPoint: true,
			situationStopPointRef:      "DUMMY",
			hasExternalResource:        true,
			subscriptionLineValue:      "first",
			expectedBroadcastedEvent:   0,
			message: `Should NOT broadcast the situation with
affected stopPoint if the external resource does not match`,
		},
		{
			situationAffectedLine:      false,
			situationAffectedStopPoint: true,
			situationStopPointRef:      "first",
			hasExternalResource:        true,
			subscriptionStopPointValue: "first",
			expectedBroadcastedEvent:   1,
			message: `Should broadcast the situation with
affected stopPoint if the external resource match`,
		},
		{
			situationMultipleAffectedStopPoint: true,
			hasExternalResource:                true,
			subscriptionStopPointValue:         "first",
			expectedBroadcastedEvent:           1,
			message: `Should broadcast the situation with
multiple affected StopPoints if at least one of the external resource match`,
		},
	}

	for _, tt := range TestCases {
		// Setup Situation with the context
		situation := referential.Model().Situations().New()
		validityPeriod := &model.TimeRange{
			EndTime: partner.Referential().Clock().Now().Add(10 * time.Minute),
		}
		situation.ValidityPeriods = append(situation.ValidityPeriods, validityPeriod)
		if tt.situationMultipleAffectedLines {
			affect := model.NewAffectedLine()
			affect.LineId = lineDummy.Id()
			situation.Affects = append(situation.Affects, affect)
			affect = model.NewAffectedLine()
			affect.LineId = lineFirst.Id()
			situation.Affects = append(situation.Affects, affect)
		}

		if tt.situationMultipleAffectedStopPoint {
			affect := model.NewAffectedStopArea()
			affect.StopAreaId = stopAreaDummy.Id()
			situation.Affects = append(situation.Affects, affect)
			affect = model.NewAffectedStopArea()
			affect.StopAreaId = stopAreaFirst.Id()
			situation.Affects = append(situation.Affects, affect)
		}

		if tt.situationAffectedLine {
			affect := model.NewAffectedLine()
			if tt.situationLineRef == "DUMMY" {
				affect.LineId = lineDummy.Id()
			}
			if tt.situationLineRef == "first" {
				affect.LineId = lineFirst.Id()
			}
			situation.Affects = append(situation.Affects, affect)
		}

		if tt.situationAffectedStopPoint {
			affect := model.NewAffectedStopArea()
			if tt.situationStopPointRef == "DUMMY" {
				affect.StopAreaId = stopAreaDummy.Id()
			}
			if tt.situationStopPointRef == "first" {
				affect.StopAreaId = stopAreaFirst.Id()
			}
			situation.Affects = append(situation.Affects, affect)
		}
		situation.Save()

		// Setup Subscription with the context
		subs := partner.Subscriptions().New(SituationExchangeBroadcast)
		subs.CreateAndAddNewResource(reference)
		subs.SetExternalId("externalId")
		if tt.hasExternalResource {
			if tt.subscriptionLineValue != "" {
				subs.SetSubscriptionOption("LineRef", fmt.Sprintf("%s:%s", partner.RemoteCodeSpace(SIRI_SITUATION_EXCHANGE_SUBSCRIPTION_BROADCASTER), tt.subscriptionLineValue))
			}
			if tt.subscriptionStopPointValue != "" {
				subs.SetSubscriptionOption("StopPointRef", fmt.Sprintf("%s:%s", partner.RemoteCodeSpace(SIRI_SITUATION_EXCHANGE_SUBSCRIPTION_BROADCASTER), tt.subscriptionStopPointValue))
			}

		}

		subs.Save()

		connector, _ := partner.Connector(SIRI_SITUATION_EXCHANGE_SUBSCRIPTION_BROADCASTER)

		event := &model.SituationBroadcastEvent{
			SituationId: situation.Id(),
		}
		connector.(*SIRISituationExchangeSubscriptionBroadcaster).HandleSituationExchangeBroadcastEvent(event)

		// Test
		assert.Len(connector.(*SIRISituationExchangeSubscriptionBroadcaster).toBroadcast, tt.expectedBroadcastedEvent, tt.message)

		// Cleanup
		subs.Delete()
		connector.(*SIRISituationExchangeSubscriptionBroadcaster).Partner().Model().Situations().Delete(&situation)
		connector.(*SIRISituationExchangeSubscriptionBroadcaster).mutex.Lock()
		connector.(*SIRISituationExchangeSubscriptionBroadcaster).toBroadcast = make(map[SubscriptionId][]model.SituationId)
		connector.(*SIRISituationExchangeSubscriptionBroadcaster).mutex.Unlock()
	}
}
