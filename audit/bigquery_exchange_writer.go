package audit

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"sync"

	"bitbucket.org/enroute-mobi/ara/logger"

	"bitbucket.org/enroute-mobi/ara/audit/controlpb"
	"bitbucket.org/enroute-mobi/ara/audit/exchangepb"
	"bitbucket.org/enroute-mobi/ara/audit/longtermsvpb"
	"bitbucket.org/enroute-mobi/ara/audit/partnerpb"
	"bitbucket.org/enroute-mobi/ara/audit/vehiclepb"
	"cloud.google.com/go/bigquery/storage/managedwriter"
	"cloud.google.com/go/bigquery/storage/managedwriter/adapt"
	"cloud.google.com/go/civil"
	"google.golang.org/protobuf/proto"
)

const maxConcurrentWrites = 256

type appendResult interface {
	GetResult(ctx context.Context) error
}

type writeStream interface {
	AppendRows(ctx context.Context, data [][]byte) (appendResult, error)
	Close() error
}

type managedStreamAdapter struct {
	s *managedwriter.ManagedStream
}

func (a *managedStreamAdapter) AppendRows(ctx context.Context, data [][]byte) (appendResult, error) {
	r, err := a.s.AppendRows(ctx, data)
	if err != nil {
		return nil, err
	}
	return &managedAppendResult{r: r}, nil
}

func (a *managedStreamAdapter) Close() error { return a.s.Close() }

type managedAppendResult struct {
	r *managedwriter.AppendResult
}

func (a *managedAppendResult) GetResult(ctx context.Context) error {
	_, err := a.r.GetResult(ctx)
	return err
}

type storageWriter struct {
	client      *managedwriter.Client
	stream      writeStream
	tableRef    string
	protoMsg    proto.Message
	newStreamFn func(ctx context.Context, client *managedwriter.Client, tableRef string, protoMsg proto.Message) (writeStream, error)
	wg          sync.WaitGroup
	done        chan struct{}
	sem         chan struct{}
}

func newStorageWriter(ctx context.Context, projectID, dataset, table string, protoMsg proto.Message) (*storageWriter, error) {
	client, err := managedwriter.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("managedwriter.NewClient: %w", err)
	}

	tableRef := fmt.Sprintf("projects/%s/datasets/%s/tables/%s", projectID, dataset, table)
	stream, err := newManagedStream(ctx, client, tableRef, protoMsg)
	if err != nil {
		client.Close()
		return nil, err
	}

	return &storageWriter{
		client:      client,
		stream:      stream,
		tableRef:    tableRef,
		protoMsg:    protoMsg,
		newStreamFn: newManagedStream,
		done:        make(chan struct{}),
		sem:         make(chan struct{}, maxConcurrentWrites),
	}, nil
}

func newManagedStream(ctx context.Context, client *managedwriter.Client, tableRef string, protoMsg proto.Message) (writeStream, error) {
	dp, err := adapt.NormalizeDescriptor(protoMsg.ProtoReflect().Descriptor())
	if err != nil {
		return nil, fmt.Errorf("NormalizeDescriptor: %w", err)
	}

	stream, err := client.NewManagedStream(ctx,
		managedwriter.WithDestinationTable(tableRef),
		managedwriter.WithType(managedwriter.DefaultStream),
		managedwriter.WithSchemaDescriptor(dp),
	)
	if err != nil {
		return nil, fmt.Errorf("NewManagedStream: %w", err)
	}

	return &managedStreamAdapter{s: stream}, nil
}

func (w *storageWriter) reconnect(ctx context.Context) error {
	w.stream.Close()
	stream, err := w.newStreamFn(ctx, w.client, w.tableRef, w.protoMsg)
	if err != nil {
		return err
	}
	w.stream = stream
	return nil
}

func (w *storageWriter) close() {
	close(w.done)
	w.stream.Close()
	w.wg.Wait()
	w.client.Close()
}

func (w *storageWriter) send(ctx context.Context, data []byte, label string) error {
	result, err := w.stream.AppendRows(ctx, [][]byte{data})
	if err != nil {
		if errors.Is(err, io.EOF) {
			logger.Log.Printf("BigQuery storage stream closed (%s), reconnecting", w.tableRef)
			if reconnErr := w.reconnect(ctx); reconnErr != nil {
				return fmt.Errorf("AppendRows: %w; reconnect failed: %v", err, reconnErr)
			}
			result, err = w.stream.AppendRows(ctx, [][]byte{data})
			if err != nil {
				return fmt.Errorf("AppendRows after reconnect: %w", err)
			}
		} else {
			return fmt.Errorf("AppendRows: %w", err)
		}
	}

	// Acquire a semaphore slot before spawning the goroutine, providing back-pressure
	// when too many writes are already in-flight.
	select {
	case w.sem <- struct{}{}:
	case <-ctx.Done():
		return fmt.Errorf("send: %w", ctx.Err())
	case <-w.done:
		return fmt.Errorf("send: writer is closing")
	}

	// GetResult blocks until the row is confirmed by BigQuery. Running it in a
	// goroutine keeps send() non-blocking so the caller's critical path is not
	// stalled while waiting for the write acknowledgement.
	rowSize := len(data)
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		defer func() { <-w.sem }()
		if err := result.GetResult(context.Background()); err != nil {
			select {
			case <-w.done:
			default:
				logger.Log.Printf("BigQuery storage write error (table: %s, %d bytes%s): %v", w.tableRef, rowSize, label, err)
			}
		}
	}()
	return nil
}

func encodeExchange(msg *BigQueryMessage) ([]byte, error) {
	pbMsg := &exchangepb.BigQueryMessage{
		Uuid:                    proto.String(msg.UUID),
		Timestamp:               proto.Int64(msg.Timestamp.UnixMicro()),
		IpAddress:               proto.String(msg.IPAddress),
		Protocol:                proto.String(msg.Protocol),
		Type:                    proto.String(string(msg.Type)),
		Direction:               proto.String(msg.Direction),
		Partner:                 proto.String(msg.Partner),
		Status:                  proto.String(msg.Status),
		ErrorDetails:            proto.String(msg.ErrorDetails),
		RequestRawMessage:       proto.String(msg.RequestRawMessage),
		ResponseRawMessage:      proto.String(msg.ResponseRawMessage),
		RequestIdentifier:       proto.String(msg.RequestIdentifier),
		ResponseIdentifier:      proto.String(msg.ResponseIdentifier),
		RequestSize:             proto.Int64(msg.RequestSize),
		ResponseSize:            proto.Int64(msg.ResponseSize),
		ProcessingTime:          proto.Float64(msg.ProcessingTime),
		SubscriptionIdentifiers: msg.SubscriptionIdentifiers,
		Lines:                   msg.Lines,
		StopAreas:               msg.StopAreas,
		Vehicles:                msg.Vehicles,
		VehicleJourneys:         msg.VehicleJourneys,
		Facilities:              msg.Facilities,
	}
	return proto.Marshal(pbMsg)
}

func civilDateTimeString(dt civil.DateTime) string {
	return fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d",
		dt.Date.Year, int(dt.Date.Month), dt.Date.Day,
		dt.Time.Hour, dt.Time.Minute, dt.Time.Second,
	)
}

func encodeVehicle(msg *BigQueryVehicleEvent) ([]byte, error) {
	pbMsg := &vehiclepb.BigQueryVehicleEvent{
		Uuid:           proto.String(msg.UUID),
		Timestamp:      proto.Int64(msg.Timestamp.UnixMicro()),
		Id:             proto.String(msg.ID),
		Codes:          msg.Codes,
		Longitude:      proto.Float64(msg.Longitude),
		Latitude:       proto.Float64(msg.Latitude),
		Bearing:        proto.Float64(msg.Bearing),
		RecordedAtTime: proto.String(civilDateTimeString(msg.RecordedAtTime)),
		Occupancy:      proto.String(msg.Occupancy),
	}
	return proto.Marshal(pbMsg)
}

func encodePartner(msg *BigQueryPartnerEvent) ([]byte, error) {
	pbMsg := &partnerpb.BigQueryPartnerEvent{
		Uuid:                     proto.String(msg.UUID),
		Timestamp:                proto.Int64(msg.Timestamp.UnixMicro()),
		Slug:                     proto.String(msg.Slug),
		PartnerUuid:              proto.String(msg.PartnerUUID),
		PreviousStatus:           proto.String(msg.PreviousStatus),
		PreviousServiceStartedAt: proto.String(civilDateTimeString(msg.PreviousServiceStartedAt)),
		NewStatus:                proto.String(msg.NewStatus),
		NewServiceStartedAt:      proto.String(civilDateTimeString(msg.NewServiceStartedAt)),
	}
	return proto.Marshal(pbMsg)
}

func encodeControl(msg *BigQueryControlEvent) ([]byte, error) {
	pbMsg := &controlpb.BigQueryControlEvent{
		Uuid:                             proto.String(msg.UUID),
		Timestamp:                        proto.Int64(msg.Timestamp.UnixMicro()),
		Criticity:                        proto.String(msg.Criticity),
		ControlType:                      proto.String(msg.ControlType),
		InternalCode:                     proto.String(msg.InternalCode),
		TargetModelClass:                 proto.String(msg.TargetModelClass),
		TargetModelUuid:                  proto.String(msg.TargetModelUUID),
		TranslationInfoMessageKey:        proto.String(msg.TranslationInfoMessageKey),
		TranslationInfoMessageAttributes: proto.String(msg.TranslationInfoMessageAttributes),
	}
	return proto.Marshal(pbMsg)
}

func encodeLongTermStopVisit(msg *BigQueryLongTermStopVisitEvent) ([]byte, error) {
	pbMsg := &longtermsvpb.BigQueryLongTermStopVisitEvent{
		Uuid:                          proto.String(msg.UUID),
		Timestamp:                     proto.Int64(msg.Timestamp.UnixMicro()),
		StopVisitUuid:                 proto.String(msg.StopVisitUUID),
		PassageOrder:                  proto.String(strconv.Itoa(msg.PassageOrder)),
		DepartureStatus:               proto.String(msg.DepartureStatus),
		ArrivalStatus:                 proto.String(msg.ArrivalStatus),
		StopAreaName:                  proto.String(msg.StopAreaName),
		StopAreaCoordinates:           proto.String(msg.StopAreaCoordinates),
		LineName:                      proto.String(msg.LineName),
		LineNumber:                    proto.String(msg.LineNumber),
		TransportMode:                 proto.String(msg.TransportMode),
		VehicleJourneyDirectionType:   proto.String(msg.VehicleJourneyDirectionType),
		VehicleJourneyOriginName:      proto.String(msg.VehicleJourneyOriginName),
		VehicleJourneyDestinationName: proto.String(msg.VehicleJourneyDestinationName),
		VehicleDriverRef:              proto.String(msg.VehicleDriverRef),
		VehicleOccupancy:              proto.String(msg.VehicleOccupancy),
	}

	if msg.AimedDepartureTime.Valid {
		pbMsg.AimedDepartureTime = proto.Int64(msg.AimedDepartureTime.Timestamp.UnixMicro())
	}
	if msg.AimedArrivalTime.Valid {
		pbMsg.AimedArrivalTime = proto.Int64(msg.AimedArrivalTime.Timestamp.UnixMicro())
	}
	if msg.ExpectedDepartureTime.Valid {
		pbMsg.ExpectedDepartureTime = proto.Int64(msg.ExpectedDepartureTime.Timestamp.UnixMicro())
	}
	if msg.ExpectedArrivalTime.Valid {
		pbMsg.ExpectedArrivalTime = proto.Int64(msg.ExpectedArrivalTime.Timestamp.UnixMicro())
	}
	if msg.ActualDepartureTime.Valid {
		pbMsg.ActualDepartureTime = proto.Int64(msg.ActualDepartureTime.Timestamp.UnixMicro())
	}
	if msg.ActualArrivalTime.Valid {
		pbMsg.ActualArrivalTime = proto.Int64(msg.ActualArrivalTime.Timestamp.UnixMicro())
	}

	for _, c := range msg.StopAreaCodes {
		pbMsg.StopAreaCodes = append(pbMsg.StopAreaCodes, &longtermsvpb.Code{
			CodeSpace: proto.String(c.CodeSpace),
			Value:     proto.String(c.Value),
		})
	}
	for _, c := range msg.LineCodes {
		pbMsg.LineCodes = append(pbMsg.LineCodes, &longtermsvpb.Code{
			CodeSpace: proto.String(c.CodeSpace),
			Value:     proto.String(c.Value),
		})
	}
	for _, c := range msg.VehicleJourneyCodes {
		pbMsg.VehicleJourneyCodes = append(pbMsg.VehicleJourneyCodes, &longtermsvpb.Code{
			CodeSpace: proto.String(c.CodeSpace),
			Value:     proto.String(c.Value),
		})
	}

	return proto.Marshal(pbMsg)
}
