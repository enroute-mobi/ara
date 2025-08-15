package settings

const (
	COLLECT_UNKNOWN CollectStatus = iota
	CAN_COLLECT
	CANNOT_COLLECT
)

type CollectStatus int64

type CollectSettings struct {
	UseDiscoveredSA    bool
	UseDiscoveredLines bool

	includedSA         collection
	excludedSA         collection
	includedLines      collection
	excludedLines      collection
	includedFacilities collection
	excludedFacilities collection
}

type collection map[string]struct{}

func (c collection) include(s string) (ok bool) {
	_, ok = c[s]
	return
}

func (cs *CollectSettings) Empty() bool {
	return len(cs.includedSA) == 0 && len(cs.excludedSA) == 0 && len(cs.includedLines) == 0 && len(cs.excludedLines) == 0 && !cs.UseDiscoveredSA && !cs.UseDiscoveredLines
}

func (cs *CollectSettings) IncludeStop(s string) CollectStatus {
	if len(cs.includedSA) == 0 {
		return COLLECT_UNKNOWN
	}

	if cs.includedSA.include(s) {
		return CAN_COLLECT
	}

	return CANNOT_COLLECT
}

func (cs *CollectSettings) ExcludeStop(s string) CollectStatus {
	if cs.excludedSA.include(s) {
		return CANNOT_COLLECT
	}
	return COLLECT_UNKNOWN
}

func (cs *CollectSettings) CanCollectStop(s string) CollectStatus {
	canCollect := cs.IncludeStop(s)
	if canCollect != COLLECT_UNKNOWN {
		return canCollect
	}

	return cs.ExcludeStop(s)
}

func (cs *CollectSettings) IncludeLine(lineId string) CollectStatus {
	if len(cs.includedLines) == 0 {
		return COLLECT_UNKNOWN
	}

	if cs.includedLines.include(lineId) {
		return CAN_COLLECT
	}

	return CANNOT_COLLECT
}

func (cs *CollectSettings) ExcludeLine(lineId string) CollectStatus {
	if cs.excludedLines.include(lineId) {
		return CANNOT_COLLECT
	}
	return COLLECT_UNKNOWN
}

func (cs *CollectSettings) IncludeFacility(facilityId string) CollectStatus {
	if len(cs.includedFacilities) == 0 {
		return COLLECT_UNKNOWN
	}

	if cs.includedFacilities.include(facilityId) {
		return CAN_COLLECT
	}

	return CANNOT_COLLECT
}

func (cs *CollectSettings) ExcludeFacility(facilityId string) CollectStatus {
	if cs.excludedFacilities.include(facilityId) {
		return CANNOT_COLLECT
	}
	return COLLECT_UNKNOWN
}

func (cs *CollectSettings) CanCollectLine(lineId string) CollectStatus {
	canCollect := cs.IncludeLine(lineId)
	if canCollect != COLLECT_UNKNOWN {
		return canCollect
	}

	canCollect = cs.ExcludeLine(lineId)
	if canCollect != COLLECT_UNKNOWN {
		return canCollect
	}
	return COLLECT_UNKNOWN
}

func (cs *CollectSettings) CanCollectFacility(facilityId string) CollectStatus {
	canCollect := cs.IncludeFacility(facilityId)
	if canCollect != COLLECT_UNKNOWN {
		return canCollect
	}

	canCollect = cs.ExcludeFacility(facilityId)
	if canCollect != COLLECT_UNKNOWN {
		return canCollect
	}
	return COLLECT_UNKNOWN
}

func (cs *CollectSettings) CanCollectLines(lineIds map[string]struct{}) CollectStatus {
	for line := range lineIds {
		canCollect := cs.CanCollectLine(line)

		if canCollect != COLLECT_UNKNOWN {
			return canCollect
		}
	}
	return COLLECT_UNKNOWN
}
