package model

func NewDefineSituationAffectsUpdater(sm *SelectMacro) (updater, error) {
	return func(mi ModelInstance) error {
		s := mi.(*Situation)
		affects := make(map[ModelId]Affect)
		for _, a := range s.Affects {
			affects[a.GetId()] = a
		}
		for _, c := range s.Consequences {
			for _, a := range c.Affects {
				affects[a.GetId()] = a
			}
		}
		newAffects := make([]Affect, 0, len(affects))
		for _, v := range affects {
			newAffects = append(newAffects, v)
		}
		s.Affects = newAffects

		return nil
	}, nil
}
