package psettings

import "bitbucket.org/enroute-mobi/ara/logger"

type CollectSettings struct {
	UseDiscovered bool

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

// Returns true if any element of the collection is in the argument map
func (c collection) isInMap(m map[string]struct{}) bool {
	for i := range c {
		if _, ok := m[i]; ok {
			return true
		}
	}
	return false
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
	return len(cs.includedSA) == 0 && len(cs.excludedSA) == 0 && len(cs.includedLines) == 0 && len(cs.excludedLines) == 0 && !cs.UseDiscovered
}

func (cs *CollectSettings) IncludeStop(s string) bool {
	if len(cs.includedSA) == 0 {
		return true
	}
	return cs.includedSA.include(s)
}

func (cs *CollectSettings) ExcludeStop(s string) bool {
	return cs.excludedSA.include(s)
}

func (cs *CollectSettings) CanCollectStop(s string) bool {
	// logger.Log.Printf("Can collect stop: %v %v", cs.IncludeStop(s), !cs.ExcludeStop(s))
	return cs.IncludeStop(s) && !cs.ExcludeStop(s)
}

func (cs *CollectSettings) IncludeLine(s string) bool {
	if len(cs.includedLines) == 0 {
		return true
	}
	return cs.includedLines.include(s)
}

func (cs *CollectSettings) ExcludeLine(s string) bool {
	return cs.excludedLines.include(s)
}

// Returns true if all lines are excluded
func (cs *CollectSettings) ExcludeAllLines(ls map[string]struct{}) bool {
	if len(cs.excludedLines) == 0 {
		return false
	}
	return !cs.excludedLines.atLeastOneNotInCollection(ls)
}

func (cs *CollectSettings) CanCollectLine(s string) bool {
	return cs.IncludeLine(s) && !cs.ExcludeLine(s)
}

// Return true if we can collect any of the lines passed in argument
func (cs *CollectSettings) canCollectLines(ls map[string]struct{}) bool {
	return len(cs.includedLines) == 0 || cs.includedLines.isInMap(ls)
}

// Return true if we can collect the lines and don't exclude at least one
func (cs *CollectSettings) CanCollectLines(ls map[string]struct{}) bool {
	logger.Log.Printf("Can collect lines: %v %v", cs.canCollectLines(ls), cs.excludedLines.atLeastOneNotInCollection(ls))

	return cs.canCollectLines(ls) && !cs.ExcludeAllLines(ls)
}
