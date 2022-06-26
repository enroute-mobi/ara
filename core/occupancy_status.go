package core

import "bitbucket.org/enroute-mobi/ara/gtfs"

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

func SIRIOccupancyCode(occupancy string) int32 {
	o, ok := gtfs.VehiclePosition_OccupancyStatus_value[occupancy]
	if !ok {
		return 7 // NO_DATA_AVAILABLE
	}
	return o
}

func SIRIOccupancyName(occupancy int32) string {
	return gtfs.VehiclePosition_OccupancyStatus_name[occupancy]
}
