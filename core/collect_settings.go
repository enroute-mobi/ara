package core

type CollectSettings struct {
	UseDiscovered bool

	includedSA    collection
	includedLines collection
	excludedSA    collection
}

type collection []string

func (c collection) include(s string) bool {
	for i := range c {
		if c[i] == s {
			return true
		}
	}
	return false
}

// Returns true if any element of the collection is in the argument map
func (c collection) isInMap(m map[string]struct{}) bool {
	for i := range c {
		if _, ok := m[c[i]]; ok {
			return true
		}
	}
	return false
}

func (cs *CollectSettings) Empty() bool {
	return len(cs.includedSA) == 0 && len(cs.includedLines) == 0 && len(cs.excludedSA) == 0 && !cs.UseDiscovered
}

func (cc *CollectSettings) IncludeStop(s string) bool {
	return cc.includedSA.include(s)
}

func (cc *CollectSettings) ExcludeStop(s string) bool {
	return cc.excludedSA.include(s)
}

func (cc *CollectSettings) IncludeLine(s string) bool {
	if len(cc.includedLines) == 0 {
		return false
	}
	if cc.includedLines.include(s) {
		return true
	}
	return false
}

func (cc *CollectSettings) CanCollectLines(lineIds map[string]struct{}) bool {
	if len(cc.includedLines) == 0 {
		return false
	}
	return cc.includedLines.isInMap(lineIds)
}
