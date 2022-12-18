package model

import (
	"strings"
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
	Undefined               = ""
	Unknown                 = "unknown"
	Empty                   = "empty"
	ManySeatsAvailable      = "manySeatsAvailable"
	FewSeatsAvailable       = "fewSeatsAvailable"
	StandingRoomOnly        = "standingRoomOnly"
	CrushedStandingRoomOnly = "crushedStandingRoomOnly"
	Full                    = "full"
	NotAcceptingPassengers  = "notAcceptingPassengers"
	StandingAvailable       = "standingAvailable"
	SeatsAvailable          = "seatsAvailable"
)

func NormalizedOccupancyName(occupancy string) string {
	switch strings.TrimSpace(occupancy) {
	case Unknown:
		return Unknown
	case Empty:
		return Empty
	case ManySeatsAvailable:
		return ManySeatsAvailable
	case FewSeatsAvailable:
		return FewSeatsAvailable
	case StandingRoomOnly:
		return StandingRoomOnly
	case CrushedStandingRoomOnly:
		return CrushedStandingRoomOnly
	case Full:
		return Full
	case NotAcceptingPassengers:
		return NotAcceptingPassengers
	case StandingAvailable:
		return StandingRoomOnly
	case SeatsAvailable:
		return ManySeatsAvailable
	default:
		return Undefined
	}
}
