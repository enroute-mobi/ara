package api

import (
	"bytes"
	"compress/gzip"
	"net/http"
	"strconv"
	"time"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/core"
	"bitbucket.org/enroute-mobi/ara/logger"
	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"github.com/golang/protobuf/proto"
)

type GtfsHandler struct {
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

	logStashEvent := partner.NewLogStashEvent()
	logStashEvent["connector"] = "GtfsHandler"
	logStashEvent["resource"] = "resource"

	var gc []core.GtfsConnector
	var c core.Connector
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
		gc, ok = partner.GtfsConnectors()
	}
	if !ok {
		http.Error(response, "Partner doesn't have the required Gtfs connector", http.StatusNotImplemented)
		return
	}

	version := "2.0"
	timestamp := uint64(handler.referential.Clock().Now().Unix())
	incrementality := gtfs.FeedHeader_FULL_DATASET

	feed := &gtfs.FeedMessage{}
	feed.Header = &gtfs.FeedHeader{}
	feed.Header.GtfsRealtimeVersion = &version
	feed.Header.Incrementality = &incrementality
	feed.Header.Timestamp = &timestamp

	for i := range gc {
		gc[i].HandleGtfs(feed, logStashEvent)
	}

	data, err := proto.Marshal(feed)
	if err != nil {
		logger.Log.Debugf("Error while marshaling feed: %v", err)
		http.Error(response, "Internal error", http.StatusInternalServerError)
		return
	}

	// Prepare the http response
	var buffer bytes.Buffer
	if partner.GzipGtfs() {
		g := gzip.NewWriter(&buffer)
		if _, err = g.Write(data); err != nil {
			logger.Log.Debugf("Can't gzip feed: %v", err)
			http.Error(response, "Internal error", http.StatusInternalServerError)
			return
		}
		if err = g.Close(); err != nil {
			logger.Log.Debugf("Can't close gzip writer: %v", err)
			http.Error(response, "Internal error", http.StatusInternalServerError)
			return
		}
	} else {
		buffer.Write(data)
	}

	logStashEvent["protobuf_size"] = strconv.Itoa(buffer.Len())
	logStashEvent["response_time"] = time.Since(startTime).String()
	audit.CurrentLogStash().WriteEvent(logStashEvent)

	response.WriteHeader(http.StatusOK)
	request.Header.Set("Content-Type", "application/x-protobuf")
	if partner.GzipGtfs() {
		request.Header.Set("Content-Encoding", "gzip")
	}
	response.Write(buffer.Bytes())
}
