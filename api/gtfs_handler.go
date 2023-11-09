package api

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"net/http"
	"time"

	"bitbucket.org/enroute-mobi/ara/api/rah"
	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/core"
	"bitbucket.org/enroute-mobi/ara/gtfs"
	"bitbucket.org/enroute-mobi/ara/logger"
	"google.golang.org/protobuf/proto"
)

type GtfsHandler struct {
	rah.RemoteAddressHandler

	referential *core.Referential
	token       string
}

func NewGtfsHandler(referential *core.Referential, token string) *GtfsHandler {
	return &GtfsHandler{
		referential: referential,
		token:       token,
	}
}

func (handler *GtfsHandler) serve(response http.ResponseWriter, request *http.Request, resource string) {
	// Find Partner by authorization Key
	partner, ok := handler.referential.Partners().FindByCredential(handler.token)
	if !ok {
		http.Error(response, "Invalid Authorization Token", http.StatusUnauthorized)
		return
	}

	startTime := handler.referential.Clock().Now()

	message := handler.newBQMessage(string(partner.Slug()), handler.HandleRemoteAddress(request))
	defer audit.CurrentBigQuery(string(handler.referential.Slug())).WriteEvent(message)

	var gc []core.GtfsConnector
	var c core.Connector
	messageType := resource

	if resource == "trip-updates" {
		c, ok = partner.Connector(core.GTFS_RT_TRIP_UPDATES_BROADCASTER)
		if ok {
			gc = []core.GtfsConnector{c.(core.GtfsConnector)}
		}
	} else if resource == "vehicle-positions" {
		c, ok = partner.Connector(core.GTFS_RT_VEHICLE_POSITIONS_BROADCASTER)
		if ok {
			gc = []core.GtfsConnector{c.(core.GtfsConnector)}
		}
	} else {
		messageType = string(audit.GTFS_TRIP_UPDATES_VEHICLE_POSITION)
		gc, ok = partner.GtfsConnectors()
	}

	message.Type = audit.BigQueryMessageType(messageType)

	if !ok {
		handler.logError(message, startTime, "Partner %v doesn't have the required Gtfs connector %v", partner.Slug(), resource)
		http.Error(response, "Partner doesn't have the required Gtfs connector", http.StatusNotImplemented)
		return
	}

	d, err := partner.GtfsCache().Fetch(messageType, func() (interface{}, error) { return handler.getFeed(gc) })

	if err != nil {
		handler.logError(message, startTime, "%v", err)
		http.Error(response, "Internal error", http.StatusInternalServerError)
		return
	}
	data := d.([]byte)

	// Prepare the http response
	var buffer bytes.Buffer
	if partner.GzipGtfs() {
		g := gzip.NewWriter(&buffer)
		if _, err := g.Write(data); err != nil {
			handler.logError(message, startTime, "Can't gzip feed: %v", err)
			http.Error(response, "Internal error", http.StatusInternalServerError)
			return
		}
		if err := g.Close(); err != nil {
			handler.logError(message, startTime, "Can't close gzip writer: %v", err)
			http.Error(response, "Internal error", http.StatusInternalServerError)
			return
		}
	} else {
		buffer.Write(data)
	}

	responseSize := buffer.Len()
	processingTime := handler.referential.Clock().Since(startTime)

	message.ResponseSize = int64(responseSize)
	message.ProcessingTime = processingTime.Seconds()

	response.WriteHeader(http.StatusOK)
	request.Header.Set("Content-Type", "application/x-protobuf")
	if partner.GzipGtfs() {
		request.Header.Set("Content-Encoding", "gzip")
	}
	response.Write(buffer.Bytes())
}

func (handler *GtfsHandler) getFeed(gc []core.GtfsConnector) ([]byte, error) {
	version := "2.0"
	timestamp := uint64(handler.referential.Clock().Now().Unix())
	incrementality := gtfs.FeedHeader_FULL_DATASET

	feed := &gtfs.FeedMessage{}
	feed.Header = &gtfs.FeedHeader{}
	feed.Header.GtfsRealtimeVersion = &version
	feed.Header.Incrementality = &incrementality
	feed.Header.Timestamp = &timestamp

	for i := range gc {
		gc[i].HandleGtfs(feed)
	}

	return proto.Marshal(feed)
}

func (handler *GtfsHandler) newBQMessage(slug, remoteAddress string) *audit.BigQueryMessage {
	return &audit.BigQueryMessage{
		Protocol:  "gtfs",
		Direction: "received",
		Status:    "OK",
		Partner:   slug,
		IPAddress: remoteAddress,
	}
}

func (handler *GtfsHandler) logError(m *audit.BigQueryMessage, startTime time.Time, format string, values ...interface{}) {
	m.ProcessingTime = handler.referential.Clock().Since(startTime).Seconds()
	m.Status = "Error"
	errorString := fmt.Sprintf(format, values...)

	m.ErrorDetails = errorString
	logger.Log.Debugf(errorString)
}
