package audit

import (
	"context"
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/audit/controlpb"
	"bitbucket.org/enroute-mobi/ara/audit/exchangepb"
	"bitbucket.org/enroute-mobi/ara/audit/longtermsvpb"
	"bitbucket.org/enroute-mobi/ara/audit/partnerpb"
	"bitbucket.org/enroute-mobi/ara/audit/vehiclepb"
	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/civil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

// fakeWriteStream is a writeStream whose AppendRows returns immediately and
// whose GetResult blocks until unblockCh is closed.
type fakeWriteStream struct {
	unblockCh chan struct{}
}

func (s *fakeWriteStream) AppendRows(_ context.Context, _ [][]byte) (appendResult, error) {
	return &fakeAppendResult{unblockCh: s.unblockCh}, nil
}

func (s *fakeWriteStream) Close() error { return nil }

type fakeAppendResult struct {
	unblockCh chan struct{}
}

func (r *fakeAppendResult) GetResult(ctx context.Context) error {
	select {
	case <-r.unblockCh:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func newTestStorageWriter(semSize int, stream writeStream) *storageWriter {
	return &storageWriter{
		stream: stream,
		done:   make(chan struct{}),
		sem:    make(chan struct{}, semSize),
	}
}

// TestStorageWriterSemaphoreBlocksAtCapacity verifies that send() blocks when
// all semaphore slots are taken and unblocks once a slot is freed.
func TestStorageWriterSemaphoreBlocksAtCapacity(t *testing.T) {
	unblock := make(chan struct{})
	w := newTestStorageWriter(1, &fakeWriteStream{unblockCh: unblock})

	require.NoError(t, w.send(context.Background(), []byte("msg1")))
	assert.Len(t, w.sem, 1, "semaphore slot should be taken")

	errCh := make(chan error, 1)
	go func() { errCh <- w.send(context.Background(), []byte("msg2")) }()

	select {
	case <-errCh:
		t.Fatal("second send should be blocked on the full semaphore")
	case <-time.After(20 * time.Millisecond):
	}

	close(unblock) // unblocks GetResult → slot freed → second send proceeds

	select {
	case err := <-errCh:
		assert.NoError(t, err)
	case <-time.After(time.Second):
		t.Fatal("second send did not unblock after slot was freed")
	}

	w.wg.Wait()
}

// TestStorageWriterSemaphoreContextCancellation verifies that send() returns
// context.Canceled when the context is cancelled while waiting for a slot.
func TestStorageWriterSemaphoreContextCancellation(t *testing.T) {
	w := newTestStorageWriter(1, &fakeWriteStream{unblockCh: make(chan struct{})})
	w.sem <- struct{}{} // fill the semaphore directly, no goroutine needed

	ctx, cancel := context.WithCancel(context.Background())

	errCh := make(chan error, 1)
	go func() { errCh <- w.send(ctx, []byte("msg")) }()

	time.Sleep(20 * time.Millisecond)
	cancel()

	select {
	case err := <-errCh:
		assert.ErrorIs(t, err, context.Canceled)
	case <-time.After(time.Second):
		t.Fatal("send did not return after context cancellation")
	}

	assert.Len(t, w.sem, 1, "semaphore slot should not have been consumed")
}

// TestStorageWriterSemaphoreWriterClosing verifies that send() returns an
// error when the writer is shut down while waiting for a semaphore slot.
func TestStorageWriterSemaphoreWriterClosing(t *testing.T) {
	w := newTestStorageWriter(1, &fakeWriteStream{unblockCh: make(chan struct{})})
	w.sem <- struct{}{} // fill the semaphore directly

	errCh := make(chan error, 1)
	go func() { errCh <- w.send(context.Background(), []byte("msg")) }()

	time.Sleep(20 * time.Millisecond)
	close(w.done)

	select {
	case err := <-errCh:
		assert.EqualError(t, err, "send: writer is closing")
	case <-time.After(time.Second):
		t.Fatal("send did not return after writer shutdown")
	}

	assert.Len(t, w.sem, 1, "semaphore slot should not have been consumed")
}

func unmarshalExchange(t *testing.T, msg *BigQueryMessage) *exchangepb.BigQueryMessage {
	t.Helper()
	b, err := encodeExchange(msg)
	require.NoError(t, err)
	decoded := &exchangepb.BigQueryMessage{}
	require.NoError(t, proto.Unmarshal(b, decoded))
	return decoded
}

func TestExchangeEncode(t *testing.T) {
	ts := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	msg := &BigQueryMessage{
		UUID:               "test-uuid",
		Timestamp:          ts,
		IPAddress:          "1.2.3.4",
		Protocol:           "siri",
		Type:               STOP_MONITORING_REQUEST,
		Direction:          "received",
		Partner:            "my-partner",
		Status:             "OK",
		ErrorDetails:       "none",
		RequestRawMessage:  "<xml>request</xml>",
		ResponseRawMessage: "<xml>response</xml>",
		RequestIdentifier:  "req-id",
		ResponseIdentifier: "resp-id",
		RequestSize:        1024,
		ResponseSize:       2048,
		ProcessingTime:     0.5,
		SubscriptionIdentifiers: []string{"sub-1", "sub-2"},
		Lines:           []string{"line-A"},
		StopAreas:       []string{"stop-1", "stop-2"},
		Vehicles:        []string{"veh-1"},
		VehicleJourneys: []string{"vj-1"},
		Facilities:      []string{"fac-1"},
	}

	decoded := unmarshalExchange(t, msg)

	assert.Equal(t, msg.UUID, decoded.GetUuid())
	assert.Equal(t, ts.UnixMicro(), decoded.GetTimestamp())
	assert.Equal(t, msg.IPAddress, decoded.GetIpAddress())
	assert.Equal(t, msg.Protocol, decoded.GetProtocol())
	assert.Equal(t, string(msg.Type), decoded.GetType())
	assert.Equal(t, msg.Direction, decoded.GetDirection())
	assert.Equal(t, msg.Partner, decoded.GetPartner())
	assert.Equal(t, msg.Status, decoded.GetStatus())
	assert.Equal(t, msg.ErrorDetails, decoded.GetErrorDetails())
	assert.Equal(t, msg.RequestRawMessage, decoded.GetRequestRawMessage())
	assert.Equal(t, msg.ResponseRawMessage, decoded.GetResponseRawMessage())
	assert.Equal(t, msg.RequestIdentifier, decoded.GetRequestIdentifier())
	assert.Equal(t, msg.ResponseIdentifier, decoded.GetResponseIdentifier())
	assert.Equal(t, msg.RequestSize, decoded.GetRequestSize())
	assert.Equal(t, msg.ResponseSize, decoded.GetResponseSize())
	assert.Equal(t, msg.ProcessingTime, decoded.GetProcessingTime())
	assert.Equal(t, msg.SubscriptionIdentifiers, decoded.GetSubscriptionIdentifiers())
	assert.Equal(t, msg.Lines, decoded.GetLines())
	assert.Equal(t, msg.StopAreas, decoded.GetStopAreas())
	assert.Equal(t, msg.Vehicles, decoded.GetVehicles())
	assert.Equal(t, msg.VehicleJourneys, decoded.GetVehicleJourneys())
	assert.Equal(t, msg.Facilities, decoded.GetFacilities())
}

func TestExchangeEncodeEmptyMessage(t *testing.T) {
	decoded := unmarshalExchange(t, &BigQueryMessage{})

	assert.Empty(t, decoded.GetUuid())
	assert.Empty(t, decoded.GetSubscriptionIdentifiers())
}

func TestExchangeEncodeLargeRawMessage(t *testing.T) {
	large := make([]byte, 11*1024*1024) // 11 MB
	for i := range large {
		large[i] = 'x'
	}
	msg := &BigQueryMessage{
		UUID:               "large-uuid",
		RequestRawMessage:  string(large),
		ResponseRawMessage: string(large),
	}

	decoded := unmarshalExchange(t, msg)

	assert.Len(t, decoded.GetRequestRawMessage(), len(large))
	assert.Len(t, decoded.GetResponseRawMessage(), len(large))
}

func TestVehicleEncode(t *testing.T) {
	ts := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	msg := &BigQueryVehicleEvent{
		UUID:      "veh-uuid",
		Timestamp: ts,
		ID:        "veh-1",
		Codes:     []string{"code-a", "code-b"},
		Longitude: 2.35,
		Latitude:  48.86,
		Bearing:   90.0,
		RecordedAtTime: civil.DateTime{
			Date: civil.Date{Year: 2024, Month: 6, Day: 1},
			Time: civil.Time{Hour: 11, Minute: 30, Second: 0},
		},
		Occupancy: "seatsAvailable",
	}

	b, err := encodeVehicle(msg)
	require.NoError(t, err)
	var decoded vehiclepb.BigQueryVehicleEvent
	require.NoError(t, proto.Unmarshal(b, &decoded))

	assert.Equal(t, msg.UUID, decoded.GetUuid())
	assert.Equal(t, ts.UnixMicro(), decoded.GetTimestamp())
	assert.Equal(t, msg.ID, decoded.GetId())
	assert.Equal(t, msg.Codes, decoded.GetCodes())
	assert.Equal(t, msg.Longitude, decoded.GetLongitude())
	assert.Equal(t, msg.Latitude, decoded.GetLatitude())
	assert.Equal(t, msg.Bearing, decoded.GetBearing())
	assert.Equal(t, "2024-06-01T11:30:00", decoded.GetRecordedAtTime())
	assert.Equal(t, msg.Occupancy, decoded.GetOccupancy())
}

func TestControlEncode(t *testing.T) {
	ts := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	msg := &BigQueryControlEvent{
		UUID:                             "ctrl-uuid",
		Timestamp:                        ts,
		Criticity:                        "warning",
		ControlType:                      "typeA",
		InternalCode:                     "internal-1",
		TargetModelClass:                 "ModelClass",
		TargetModelUUID:                  "model-uuid",
		TranslationInfoMessageKey:        "msg.key",
		TranslationInfoMessageAttributes: `{"key":"val"}`,
	}

	b, err := encodeControl(msg)
	require.NoError(t, err)
	var decoded controlpb.BigQueryControlEvent
	require.NoError(t, proto.Unmarshal(b, &decoded))

	assert.Equal(t, msg.UUID, decoded.GetUuid())
	assert.Equal(t, ts.UnixMicro(), decoded.GetTimestamp())
	assert.Equal(t, msg.Criticity, decoded.GetCriticity())
	assert.Equal(t, msg.ControlType, decoded.GetControlType())
	assert.Equal(t, msg.InternalCode, decoded.GetInternalCode())
	assert.Equal(t, msg.TargetModelClass, decoded.GetTargetModelClass())
	assert.Equal(t, msg.TargetModelUUID, decoded.GetTargetModelUuid())
	assert.Equal(t, msg.TranslationInfoMessageKey, decoded.GetTranslationInfoMessageKey())
	assert.Equal(t, msg.TranslationInfoMessageAttributes, decoded.GetTranslationInfoMessageAttributes())
}

func TestLongTermStopVisitEncode(t *testing.T) {
	ts := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	aimed := time.Date(2024, 6, 1, 14, 0, 0, 0, time.UTC)
	msg := &BigQueryLongTermStopVisitEvent{
		UUID:               "ltsv-uuid",
		Timestamp:          ts,
		StopVisitUUID:      "sv-uuid",
		PassageOrder:       3,
		AimedDepartureTime: bigquery.NullTimestamp{Timestamp: aimed, Valid: true},
		AimedArrivalTime:   bigquery.NullTimestamp{Valid: false},
		DepartureStatus:    "onTime",
		ArrivalStatus:      "onTime",
		StopAreaName:       "Gare du Nord",
		StopAreaCodes:      []Code{{CodeSpace: "STIF", Value: "stop-1"}},
		LineName:           "RER B",
		LineNumber:         "B",
		VehicleJourneyCodes: []Code{{CodeSpace: "STIF", Value: "vj-1"}},
	}

	b, err := encodeLongTermStopVisit(msg)
	require.NoError(t, err)
	var decoded longtermsvpb.BigQueryLongTermStopVisitEvent
	require.NoError(t, proto.Unmarshal(b, &decoded))

	assert.Equal(t, msg.UUID, decoded.GetUuid())
	assert.Equal(t, ts.UnixMicro(), decoded.GetTimestamp())
	assert.Equal(t, "3", decoded.GetPassageOrder())
	assert.Equal(t, aimed.UnixMicro(), decoded.GetAimedDepartureTime())
	assert.Nil(t, decoded.AimedArrivalTime, "AimedArrivalTime should be nil when invalid")
	require.Len(t, decoded.GetStopAreaCodes(), 1)
	assert.Equal(t, "stop-1", decoded.GetStopAreaCodes()[0].GetValue())
	assert.Equal(t, "STIF", decoded.GetStopAreaCodes()[0].GetCodeSpace())
	require.Len(t, decoded.GetVehicleJourneyCodes(), 1)
	assert.Equal(t, "vj-1", decoded.GetVehicleJourneyCodes()[0].GetValue())
}

func TestPartnerEncode(t *testing.T) {
	ts := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	msg := &BigQueryPartnerEvent{
		UUID:           "partner-uuid",
		Timestamp:      ts,
		Slug:           "my-partner",
		PartnerUUID:    "prt-uuid",
		PreviousStatus: "up",
		PreviousServiceStartedAt: civil.DateTime{
			Date: civil.Date{Year: 2024, Month: 5, Day: 31},
			Time: civil.Time{Hour: 8, Minute: 0, Second: 0},
		},
		NewStatus: "down",
		NewServiceStartedAt: civil.DateTime{
			Date: civil.Date{Year: 2024, Month: 6, Day: 1},
			Time: civil.Time{Hour: 12, Minute: 0, Second: 0},
		},
	}

	b, err := encodePartner(msg)
	require.NoError(t, err)
	var decoded partnerpb.BigQueryPartnerEvent
	require.NoError(t, proto.Unmarshal(b, &decoded))

	assert.Equal(t, msg.UUID, decoded.GetUuid())
	assert.Equal(t, ts.UnixMicro(), decoded.GetTimestamp())
	assert.Equal(t, msg.Slug, decoded.GetSlug())
	assert.Equal(t, msg.PartnerUUID, decoded.GetPartnerUuid())
	assert.Equal(t, msg.PreviousStatus, decoded.GetPreviousStatus())
	assert.Equal(t, "2024-05-31T08:00:00", decoded.GetPreviousServiceStartedAt())
	assert.Equal(t, msg.NewStatus, decoded.GetNewStatus())
	assert.Equal(t, "2024-06-01T12:00:00", decoded.GetNewServiceStartedAt())
}
