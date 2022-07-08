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

const (
	empty                   = "EMPTY"
	manySeatsAvailable      = "MANY_SEATS_AVAILABLE"
	fewSeatsAvailable       = "FEW_SEATS_AVAILABLE"
	standingRoomOnly        = "STANDING_ROOM_ONLY"
	crushedStandingRoomOnly = "CRUSHED_STANDING_ROOM_ONLY"
	full                    = "FULL"
	notAcceptingPassengers  = "NOT_ACCEPTING_PASSENGERS"
	noDataAvailable         = "NO_DATA_AVAILABLE"
	notBoardable            = "NOT_BOARDABLE"
)

func NormalizedOccupancyName(occupancy string) string { // Doesn't do mutch for now, but it felt faster this way
	switch strings.TrimSpace(occupancy) {
	case "":
		return ""
	case empty:
		return empty
	case manySeatsAvailable:
		return manySeatsAvailable
	case fewSeatsAvailable:
		return fewSeatsAvailable
	case standingRoomOnly:
		return standingRoomOnly
	case crushedStandingRoomOnly:
		return crushedStandingRoomOnly
	case full:
		return full
	case notAcceptingPassengers:
		return notAcceptingPassengers
	case notBoardable:
		return notBoardable
	default:
		return noDataAvailable
	}
}

func OccupancyCode(occupancy string) gtfs.VehiclePosition_OccupancyStatus {
	switch occupancy {
	case empty:
		return gtfs.VehiclePosition_OccupancyStatus(0)
	case manySeatsAvailable:
		return gtfs.VehiclePosition_OccupancyStatus(1)
	case fewSeatsAvailable:
		return gtfs.VehiclePosition_OccupancyStatus(2)
	case standingRoomOnly:
		return gtfs.VehiclePosition_OccupancyStatus(3)
	case crushedStandingRoomOnly:
		return gtfs.VehiclePosition_OccupancyStatus(4)
	case full:
		return gtfs.VehiclePosition_OccupancyStatus(5)
	case notAcceptingPassengers:
		return gtfs.VehiclePosition_OccupancyStatus(6)
	case notBoardable:
		return gtfs.VehiclePosition_OccupancyStatus(8)
	default: // NO_DATA_AVAILABLE
		return gtfs.VehiclePosition_OccupancyStatus(7)
	}
}

func OccupancyName(occupancy gtfs.VehiclePosition_OccupancyStatus) string {
	switch occupancy {
	case 0:
		return empty
	case 1:
		return manySeatsAvailable
	case 2:
		return fewSeatsAvailable
	case 3:
		return standingRoomOnly
	case 4:
		return crushedStandingRoomOnly
	case 5:
		return full
	case 6:
		return notAcceptingPassengers
	case 8:
		return notBoardable
	default:
		return noDataAvailable
	}
}
