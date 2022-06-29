package model

import (
	"sync"

	"bitbucket.org/enroute-mobi/ara/gtfs"
)

/* For Reference, these are the gtfs statuses:
0: "EMPTY",
1: "MANY_SEATS_AVAILABLE",
2: "FEW_SEATS_AVAILABLE",
3: "STANDING_ROOM_ONLY",
4: "CRUSHED_STANDING_ROOM_ONLY",
5: "FULL",
6: "NOT_ACCEPTING_PASSENGERS",
7: "NO_DATA_AVAILABLE",
8: "NOT_BOARDABLE",
*/

var (
	statusName = struct {
		sync.RWMutex
		m map[gtfs.VehiclePosition_OccupancyStatus]string
	}{
		m: map[gtfs.VehiclePosition_OccupancyStatus]string{
			0: "EMPTY",
			1: "MANY_SEATS_AVAILABLE",
			2: "FEW_SEATS_AVAILABLE",
			3: "STANDING_ROOM_ONLY",
			4: "CRUSHED_STANDING_ROOM_ONLY",
			5: "FULL",
			6: "NOT_ACCEPTING_PASSENGERS",
			7: "NO_DATA_AVAILABLE",
			8: "NOT_BOARDABLE",
		},
	}

	statusValue = struct {
		sync.RWMutex
		m map[string]gtfs.VehiclePosition_OccupancyStatus
	}{
		m: map[string]gtfs.VehiclePosition_OccupancyStatus{
			"EMPTY":                      0,
			"MANY_SEATS_AVAILABLE":       1,
			"FEW_SEATS_AVAILABLE":        2,
			"STANDING_ROOM_ONLY":         3,
			"CRUSHED_STANDING_ROOM_ONLY": 4,
			"FULL":                       5,
			"NOT_ACCEPTING_PASSENGERS":   6,
			"NO_DATA_AVAILABLE":          7,
			"NOT_BOARDABLE":              8,
		},
	}
)

func OccupancyCode(occupancy string) gtfs.VehiclePosition_OccupancyStatus {
	statusValue.RLock()
	o, ok := statusValue.m[occupancy]
	statusValue.RUnlock()
	if !ok {
		return 7 // NO_DATA_AVAILABLE
	}
	return o
}

func OccupancyName(occupancy gtfs.VehiclePosition_OccupancyStatus) string {
	statusName.RLock()
	o, ok := statusName.m[occupancy]
	statusName.RUnlock()
	if !ok {
		return "NO_DATA_AVAILABLE"
	}
	return o
}
