package model

import (
	"strings"

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

func NormalizedOccupancyName(occupancy string) string { // Doesn't do mutch for now, but it felt faster this way
	switch strings.TrimSpace(occupancy) {
	case "":
		return ""
	case "EMPTY":
		return "EMPTY"
	case "MANY_SEATS_AVAILABLE":
		return "MANY_SEATS_AVAILABLE"
	case "FEW_SEATS_AVAILABLE":
		return "FEW_SEATS_AVAILABLE"
	case "STANDING_ROOM_ONLY":
		return "STANDING_ROOM_ONLY"
	case "CRUSHED_STANDING_ROOM_ONLY":
		return "CRUSHED_STANDING_ROOM_ONLY"
	case "FULL":
		return "FULL"
	case "NOT_ACCEPTING_PASSENGERS":
		return "NOT_ACCEPTING_PASSENGERS"
	case "NOT_BOARDABLE":
		return "NOT_BOARDABLE"
	default: // NO_DATA_AVAILABLE
		return "NO_DATA_AVAILABLE"
	}
}

func OccupancyCode(occupancy string) gtfs.VehiclePosition_OccupancyStatus {
	switch occupancy {
	case "EMPTY":
		return gtfs.VehiclePosition_OccupancyStatus(0)
	case "MANY_SEATS_AVAILABLE":
		return gtfs.VehiclePosition_OccupancyStatus(1)
	case "FEW_SEATS_AVAILABLE":
		return gtfs.VehiclePosition_OccupancyStatus(2)
	case "STANDING_ROOM_ONLY":
		return gtfs.VehiclePosition_OccupancyStatus(3)
	case "CRUSHED_STANDING_ROOM_ONLY":
		return gtfs.VehiclePosition_OccupancyStatus(4)
	case "FULL":
		return gtfs.VehiclePosition_OccupancyStatus(5)
	case "NOT_ACCEPTING_PASSENGERS":
		return gtfs.VehiclePosition_OccupancyStatus(6)
	case "NOT_BOARDABLE":
		return gtfs.VehiclePosition_OccupancyStatus(8)
	default: // NO_DATA_AVAILABLE
		return gtfs.VehiclePosition_OccupancyStatus(7)
	}
}

func OccupancyName(occupancy gtfs.VehiclePosition_OccupancyStatus) string {
	switch occupancy {
	case 0:
		return "EMPTY"
	case 1:
		return "MANY_SEATS_AVAILABLE"
	case 2:
		return "FEW_SEATS_AVAILABLE"
	case 3:
		return "STANDING_ROOM_ONLY"
	case 4:
		return "CRUSHED_STANDING_ROOM_ONLY"
	case 5:
		return "FULL"
	case 6:
		return "NOT_ACCEPTING_PASSENGERS"
	case 8:
		return "NOT_BOARDABLE"
	default: // 7
		return "NO_DATA_AVAILABLE"
	}
}
