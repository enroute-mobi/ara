package audit

import (
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatasetFormat(t *testing.T) {
	bq := NewBigQuery("dataset-with_hyphen").(*BigQueryClient)

	if bq.dataset != "dataset_with_hyphen" {
		t.Errorf("Wrong format for BQ Client: expected 'dataset_with_hyphen' got '%s'", bq.dataset)
	}
}

func TestCurrentBigQuery_ReturnsNullWhenNotSet(t *testing.T) {
	bq := CurrentBigQuery("slug-never-registered-" + t.Name())
	_, ok := bq.(*NullBigQuery)
	assert.True(t, ok)
}

func TestCurrentBigQuery_ReturnsRegisteredInstance(t *testing.T) {
	slug := "test-slug-" + t.Name()
	fake := NewFakeBigQuery()
	SetCurrentBigQuery(slug, fake)

	got := CurrentBigQuery(slug)
	assert.Equal(t, fake, got)
}

func TestSetCurrentBigQuery_Overwrites(t *testing.T) {
	slug := "test-slug-overwrite-" + t.Name()
	first := NewFakeBigQuery()
	second := NewFakeBigQuery()

	SetCurrentBigQuery(slug, first)
	SetCurrentBigQuery(slug, second)

	assert.Equal(t, second, CurrentBigQuery(slug))
}

func TestNullBigQuery_WriteEvent_ReturnsNil(t *testing.T) {
	bq := NewNullBigQuery()
	err := bq.WriteEvent(&BigQueryMessage{})
	assert.NoError(t, err)
}

func TestFakeBigQuery_WriteEvent_RoutesAllTypes(t *testing.T) {
	bq := NewFakeBigQuery()

	require.NoError(t, bq.WriteEvent(&BigQueryMessage{UUID: "m1"}))
	require.NoError(t, bq.WriteEvent(&BigQueryMessage{UUID: "m2"}))
	require.NoError(t, bq.WriteEvent(&BigQueryPartnerEvent{UUID: "p1"}))
	require.NoError(t, bq.WriteEvent(&BigQueryVehicleEvent{UUID: "v1"}))
	require.NoError(t, bq.WriteEvent(&BigQueryLongTermStopVisitEvent{UUID: "ltsv1"}))
	require.NoError(t, bq.WriteEvent(&BigQueryControlEvent{UUID: "c1"}))

	require.Len(t, bq.Messages(), 2)
	assert.Equal(t, "m1", bq.Messages()[0].UUID)
	assert.Equal(t, "m2", bq.Messages()[1].UUID)

	require.Len(t, bq.PartnerEvents(), 1)
	assert.Equal(t, "p1", bq.PartnerEvents()[0].UUID)

	require.Len(t, bq.VehicleEvents(), 1)
	assert.Equal(t, "v1", bq.VehicleEvents()[0].UUID)

	require.Len(t, bq.LongTermStopVisitEvents(), 1)
	assert.Equal(t, "ltsv1", bq.LongTermStopVisitEvents()[0].UUID)

	require.Len(t, bq.ControlEvents(), 1)
	assert.Equal(t, "c1", bq.ControlEvents()[0].UUID)
}

func TestBigQueryClient_writeMessage_EnqueuesOnChannel(t *testing.T) {
	bq := &BigQueryClient{
		messages: make(chan *BigQueryMessage, 1),
	}
	msg := &BigQueryMessage{UUID: "queued"}

	require.NoError(t, bq.writeMessage(msg))

	select {
	case got := <-bq.messages:
		assert.Equal(t, msg, got)
	default:
		t.Fatal("expected message in channel, got none")
	}
}

func TestBigQueryClient_writeMessage_DropsWhenFull(t *testing.T) {
	bq := &BigQueryClient{
		messages: make(chan *BigQueryMessage, 1),
	}
	bq.messages <- &BigQueryMessage{UUID: "existing"}

	require.NoError(t, bq.writeMessage(&BigQueryMessage{UUID: "dropped"}))

	assert.Equal(t, int64(1), bq.lostMessagesCount.Load())
	assert.Len(t, bq.messages, 1, "channel should still hold the original message")
}

func TestBigQueryClient_writeMessage_ResetsLostCountOnRecovery(t *testing.T) {
	bq := &BigQueryClient{
		messages: make(chan *BigQueryMessage, 1),
	}
	bq.messages <- &BigQueryMessage{}
	bq.writeMessage(&BigQueryMessage{}) // drops, count = 1
	bq.writeMessage(&BigQueryMessage{}) // drops, count = 2
	assert.Equal(t, int64(2), bq.lostMessagesCount.Load())

	<-bq.messages                       // drain to make room
	bq.writeMessage(&BigQueryMessage{}) // succeeds, should reset count to 0

	assert.Equal(t, int64(0), bq.lostMessagesCount.Load())
}

func TestBigQueryClient_writePartnerEvent_DropsWhenFull(t *testing.T) {
	bq := &BigQueryClient{
		partnerEvents: make(chan *BigQueryPartnerEvent, 1),
	}
	bq.partnerEvents <- &BigQueryPartnerEvent{}

	require.NoError(t, bq.writePartnerEvent(&BigQueryPartnerEvent{}))

	assert.Equal(t, int64(1), bq.lostPartnerEventsCount.Load())
}

func TestBigQueryClient_writeVehicleEvent_DropsWhenFull(t *testing.T) {
	bq := &BigQueryClient{
		vehicleEvents: make(chan *BigQueryVehicleEvent, 1),
	}
	bq.vehicleEvents <- &BigQueryVehicleEvent{}

	require.NoError(t, bq.writeVehicleEvent(&BigQueryVehicleEvent{}))

	assert.Equal(t, int64(1), bq.lostVehicleEventsCount.Load())
}

func TestBigQueryClient_writeLongTermStopVisitEvent_DropsWhenFull(t *testing.T) {
	bq := &BigQueryClient{
		longTermStopVisitEvents: make(chan *BigQueryLongTermStopVisitEvent, 1),
	}
	bq.longTermStopVisitEvents <- &BigQueryLongTermStopVisitEvent{}

	require.NoError(t, bq.writeLongTermStopVisitEvent(&BigQueryLongTermStopVisitEvent{}))

	assert.Equal(t, int64(1), bq.lostLongTermStopVisitsEventsCount.Load())
}

func TestBigQueryClient_writeControlEvent_DropsWhenFull(t *testing.T) {
	bq := &BigQueryClient{
		controlEvents: make(chan *BigQueryControlEvent, 1),
	}
	bq.controlEvents <- &BigQueryControlEvent{}

	require.NoError(t, bq.writeControlEvent(&BigQueryControlEvent{}))

	assert.Equal(t, int64(1), bq.lostControlEventsCount.Load())
}

func TestBigQueryClient_WriteEvent_SetsUUIDAndTimestamp(t *testing.T) {
	fixedTime := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	fakeClock := clock.NewFakeClockAt(fixedTime)
	fakeUUID := uuid.NewFakeUUIDGenerator()

	bq := &BigQueryClient{
		messages: make(chan *BigQueryMessage, 1),
	}
	bq.SetClock(fakeClock)
	bq.SetUUIDGenerator(fakeUUID)

	msg := &BigQueryMessage{}
	require.NoError(t, bq.WriteEvent(msg))

	select {
	case got := <-bq.messages:
		assert.Equal(t, fixedTime, got.Timestamp)
		assert.NotEmpty(t, got.UUID)
	default:
		t.Fatal("expected message in channel")
	}
}

func TestBigQueryClient_WriteEvent_RoutesAllTypes(t *testing.T) {
	bq := &BigQueryClient{
		messages:                make(chan *BigQueryMessage, 1),
		partnerEvents:           make(chan *BigQueryPartnerEvent, 1),
		vehicleEvents:           make(chan *BigQueryVehicleEvent, 1),
		longTermStopVisitEvents: make(chan *BigQueryLongTermStopVisitEvent, 1),
		controlEvents:           make(chan *BigQueryControlEvent, 1),
	}
	bq.SetClock(clock.NewFakeClockAt(time.Now()))
	bq.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())

	require.NoError(t, bq.WriteEvent(&BigQueryMessage{}))
	require.NoError(t, bq.WriteEvent(&BigQueryPartnerEvent{}))
	require.NoError(t, bq.WriteEvent(&BigQueryVehicleEvent{}))
	require.NoError(t, bq.WriteEvent(&BigQueryLongTermStopVisitEvent{}))
	require.NoError(t, bq.WriteEvent(&BigQueryControlEvent{}))

	assert.Len(t, bq.messages, 1)
	assert.Len(t, bq.partnerEvents, 1)
	assert.Len(t, bq.vehicleEvents, 1)
	assert.Len(t, bq.longTermStopVisitEvents, 1)
	assert.Len(t, bq.controlEvents, 1)
}
