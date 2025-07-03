package core

import (
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/model/schedules"
	"bitbucket.org/enroute-mobi/ara/siri/siri"
	"github.com/stretchr/testify/assert"
)

func Test_EstimatedTimetableBroadcaster_BuildDepartureSchedues(t *testing.T) {
	assert := assert.New(t)

	fakeClock := clock.NewFakeClock()
	clock.SetDefaultClock(fakeClock)
	// uuidGenerator := uuid.NewFakeUUIDGenerator()

	referentials := NewMemoryReferentials()
	referential := referentials.New("Un Referential Plutot Cool")
	referential.SetClock(fakeClock)
	referential.Save()

	partner := referential.Partners().New("test")
	partner.Save()

	timeA, _ := time.Parse(time.RFC3339, "2016-09-22T07:58:34+02:00")
	timeB, _ := time.Parse(time.RFC3339, "2020-01-01T07:58:34+02:00")
	timeC, _ := time.Parse(time.RFC3339, "2024-01-31T07:58:34+02:00")

	type schedule struct {
		kind          schedules.StopVisitScheduleType
		departureTime time.Time
	}

	type schedules []schedule
	var TestCases = []struct {
		schedules             schedules
		expectedDepartureTime time.Time
		aimedDepartureTime    time.Time
		actualDepartureTime   time.Time
		message               string
	}{
		{
			schedules: []schedule{
				schedule{
					kind:          "expected",
					departureTime: timeA,
				},
			},
			expectedDepartureTime: timeA,
			message: `Should set only ExpectedDepartueTime
if only Expected time is provided`,
		},
		{
			schedules: []schedule{
				schedule{
					kind:          "aimed",
					departureTime: timeA,
				},
			},
			aimedDepartureTime: timeA,
			message: `Should set only AimedDepartueTime
if only Aimed time is provided`,
		},
		{
			schedules: []schedule{
				schedule{
					kind:          "actual",
					departureTime: timeA,
				},
			},
			actualDepartureTime: timeA,
			message: `Should set only ActualDepartueTime
if only Actual time is provided`,
		},
		{
			schedules: []schedule{
				schedule{
					kind:          "expected",
					departureTime: timeA,
				},
				schedule{
					kind:          "aimed",
					departureTime: timeB,
				},
			},
			expectedDepartureTime: timeA,
			aimedDepartureTime:    timeB,
			message: `Should set ExpectedDepartueTime
and AimedDepartureTime if Expected and Aimed time are provided`,
		},
		{
			schedules: []schedule{
				schedule{
					kind:          "expected",
					departureTime: timeA,
				},
				schedule{
					kind:          "aimed",
					departureTime: timeB,
				},
				schedule{
					kind:          "actual",
					departureTime: timeC,
				},
			},
			expectedDepartureTime: timeA,
			aimedDepartureTime:    timeB,
			message: `Should set ExpectedDepartueTime
 and AimedDepartureTime if Expected, Aimed and Actual time are provided`,
		},
		{
			schedules: []schedule{
				schedule{
					kind:          "actual",
					departureTime: timeC,
				},
				schedule{
					kind:          "aimed",
					departureTime: timeB,
				},
			},
			actualDepartureTime: timeC,
			aimedDepartureTime:  timeB,
			message: `Should set ActualDepartueTime
 and AimedDepartureTime if Actual and Aimed time are provided`,
		},
	}

	for _, tt := range TestCases {
		sv := referential.Model().StopVisits().New()

		for _, schedule := range tt.schedules {
			sv.Schedules.SetDepartureTime(schedule.kind, schedule.departureTime)
		}

		builder := NewBroadcastStopMonitoringBuilder(partner, "")

		siriStopVisit := &siri.SIRIMonitoredStopVisit{}
		builder.BuildDepartureSchedules(sv, siriStopVisit)

		assert.Equal(siriStopVisit.ExpectedDepartureTime, tt.expectedDepartureTime, tt.message)
		assert.Equal(siriStopVisit.AimedDepartureTime, tt.aimedDepartureTime, tt.message)
		assert.Equal(siriStopVisit.ActualDepartureTime, tt.actualDepartureTime, tt.message)
	}

}

func Test_EstimatedTimetableBroadcaster_BuildArrivalSchedues(t *testing.T) {
	assert := assert.New(t)

	fakeClock := clock.NewFakeClock()
	clock.SetDefaultClock(fakeClock)
	// uuidGenerator := uuid.NewFakeUUIDGenerator()

	referentials := NewMemoryReferentials()
	referential := referentials.New("Un Referential Plutot Cool")
	referential.SetClock(fakeClock)
	referential.Save()

	partner := referential.Partners().New("test")
	partner.Save()

	timeA, _ := time.Parse(time.RFC3339, "2016-09-22T07:58:34+02:00")
	timeB, _ := time.Parse(time.RFC3339, "2020-01-01T07:58:34+02:00")
	timeC, _ := time.Parse(time.RFC3339, "2024-01-31T07:58:34+02:00")

	type schedule struct {
		kind          schedules.StopVisitScheduleType
		departureTime time.Time
	}

	type schedules []schedule
	var TestCases = []struct {
		schedules           schedules
		expectedArrivalTime time.Time
		aimedArrivalTime    time.Time
		actualArrivalTime   time.Time
		message             string
	}{
		{
			schedules: []schedule{
				schedule{
					kind:          "expected",
					departureTime: timeA,
				},
			},
			expectedArrivalTime: timeA,
			message: `Should set only ExpectedDepartueTime
if only Expected time is provided`,
		},
		{
			schedules: []schedule{
				schedule{
					kind:          "aimed",
					departureTime: timeA,
				},
			},
			aimedArrivalTime: timeA,
			message: `Should set only AimedDepartueTime
if only Aimed time is provided`,
		},
		{
			schedules: []schedule{
				schedule{
					kind:          "actual",
					departureTime: timeA,
				},
			},
			actualArrivalTime: timeA,
			message: `Should set only ActualDepartueTime
if only Actual time is provided`,
		},
		{
			schedules: []schedule{
				schedule{
					kind:          "expected",
					departureTime: timeA,
				},
				schedule{
					kind:          "aimed",
					departureTime: timeB,
				},
			},
			expectedArrivalTime: timeA,
			aimedArrivalTime:    timeB,
			message: `Should set ExpectedDepartueTime
and AimedArrivalTime if Expected and Aimed time are provided`,
		},
		{
			schedules: []schedule{
				schedule{
					kind:          "expected",
					departureTime: timeA,
				},
				schedule{
					kind:          "aimed",
					departureTime: timeB,
				},
				schedule{
					kind:          "actual",
					departureTime: timeC,
				},
			},
			expectedArrivalTime: timeA,
			aimedArrivalTime:    timeB,
			message: `Should set ExpectedDepartueTime
 and AimedArrivalTime if Expected, Aimed and Actual time are provided`,
		},
		{
			schedules: []schedule{
				schedule{
					kind:          "actual",
					departureTime: timeC,
				},
				schedule{
					kind:          "aimed",
					departureTime: timeB,
				},
			},
			actualArrivalTime: timeC,
			aimedArrivalTime:  timeB,
			message: `Should set ActualDepartueTime
 and AimedArrivalTime if Actual and Aimed time are provided`,
		},
	}

	for _, tt := range TestCases {
		sv := referential.Model().StopVisits().New()

		for _, schedule := range tt.schedules {
			sv.Schedules.SetArrivalTime(schedule.kind, schedule.departureTime)
		}

		builder := NewBroadcastStopMonitoringBuilder(partner, "")

		siriStopVisit := &siri.SIRIMonitoredStopVisit{}
		builder.BuildArrivalSchedules(sv, siriStopVisit)

		assert.Equal(siriStopVisit.ExpectedArrivalTime, tt.expectedArrivalTime, tt.message)
		assert.Equal(siriStopVisit.AimedArrivalTime, tt.aimedArrivalTime, tt.message)
		assert.Equal(siriStopVisit.ActualArrivalTime, tt.actualArrivalTime, tt.message)
	}
}
