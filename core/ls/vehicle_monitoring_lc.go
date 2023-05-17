package ls

import "bitbucket.org/enroute-mobi/ara/model"

type VehicleMonitoringLastChange struct {
	lastState

	LinkDistance float64
	Percentage   float64
	Longitude    float64
	Latitude     float64
	Bearing      float64
}

func (vmlc *VehicleMonitoringLastChange) UpdateState(v *model.Vehicle) {
	vmlc.LinkDistance = v.LinkDistance
	vmlc.Percentage = v.Percentage
	vmlc.Longitude = v.Longitude
	vmlc.Latitude = v.Latitude
	vmlc.Bearing = v.Bearing
}

func NewVehicleMonitoringLastChange(v *model.Vehicle, sub subscription) *VehicleMonitoringLastChange {
	vmlc := &VehicleMonitoringLastChange{}
	vmlc.SetSubscription(sub)
	vmlc.UpdateState(v)
	return vmlc
}

func (vmlc *VehicleMonitoringLastChange) HasChanged(v *model.Vehicle) bool {
	// Check LinkDistance
	if vmlc.LinkDistance != v.LinkDistance {
		return true
	}

	// Check Percentage
	if vmlc.Percentage != v.Percentage {
		return true
	}

	// Check Longitude
	if vmlc.Longitude != v.Longitude {
		return true
	}

	// Check Latitude
	if vmlc.Latitude != v.Latitude {
		return true
	}

	// Check Bearing
	if vmlc.Bearing != v.Bearing {
		return true
	}

	return false
}
