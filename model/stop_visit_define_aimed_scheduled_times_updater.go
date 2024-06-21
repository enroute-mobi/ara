package model

func NewStopVisitDefineAimedScheduledTimesUpdater(sm *SelectMacro) (updater, error) {
	return func(mi ModelInstance) error {
		sv := mi.(*StopVisit)
		sv.Schedules.SetDefaultAimedTimes()
		return nil
	}, nil
}
