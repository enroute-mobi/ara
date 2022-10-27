package core

import (
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/model"
)

func Test_StopAreaUpdateRequest_StopAreaId(t *testing.T) {
	stopAreaUpdateRequest := NewStopAreaUpdateRequest("StopAreaId")

	if expected := model.StopAreaId("StopAreaId"); stopAreaUpdateRequest.StopAreaId() != expected {
		t.Errorf("StopAreaUpdateRequest.StopAreaId() returns wrong value, got: %s, required: %s", stopAreaUpdateRequest.StopAreaId(), expected)
	}
}

func Test_StopAreaUpdateRequest_CreatedAt(t *testing.T) {
	testTime := time.Date(1984, time.April, 4, 0, 0, 0, 0, time.UTC)
	clock.SetDefaultClock(clock.NewFakeClockAt(testTime))
	stopAreaUpdateRequest := NewStopAreaUpdateRequest("stopAreaId")

	if !stopAreaUpdateRequest.CreatedAt().Equal(testTime) {
		t.Errorf("StopAreaUpdateRequest.CreatedAt() returns wrong value, got: %s, required: %s", stopAreaUpdateRequest.CreatedAt(), testTime)
	}
}
