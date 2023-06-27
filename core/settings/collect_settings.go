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

	includedSA    collection
	excludedSA    collection
	includedLines collection
	excludedLines collection
}

type collection map[string]struct{}

func (c collection) include(s string) (ok bool) {
	_, ok = c[s]
	return
}

// Returns true if at least one element of the argument map isn't in the collection
func (c collection) atLeastOneNotInCollection(m map[string]struct{}) bool {
	for k := range m {
		if _, ok := c[k]; !ok {
			return true
		}
	}
	return false
}

func (cs *CollectSettings) Empty() bool {
	return len(cs.includedSA) == 0 && len(cs.excludedSA) == 0 && len(cs.includedLines) == 0 && len(cs.excludedLines) == 0 && !cs.UseDiscoveredSA && !cs.UseDiscoveredLines
}

func (cs *CollectSettings) IncludeStop(s string) CollectStatus {
	if cs.includedSA.include(s) {
		return CAN_COLLECT
	}
	return COLLECT_UNKNOWN
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
	if cs.includedLines.include(lineId) {
		return CAN_COLLECT
	}
	return COLLECT_UNKNOWN
}

func (cs *CollectSettings) ExcludeLine(lineId string) CollectStatus {
	if cs.excludedLines.include(lineId) {
		return CANNOT_COLLECT
	}
	return COLLECT_UNKNOWN
}

// Returns true if all lines are excluded
func (cs *CollectSettings) ExcludeAllLines(ls map[string]struct{}) bool {
	if len(cs.excludedLines) == 0 {
		return false
	}
	return !cs.excludedLines.atLeastOneNotInCollection(ls)
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

func (cs *CollectSettings) CanCollectLines(lineIds map[string]struct{}) CollectStatus {
	for line := range lineIds {
		canCollect := cs.CanCollectLine(line)

		if canCollect != COLLECT_UNKNOWN {
			return canCollect
		}
	}
	return COLLECT_UNKNOWN
}
