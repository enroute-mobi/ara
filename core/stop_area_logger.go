package core

import (
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
)

type StopAreaLogger struct {
	verbose bool
}

func NewStopAreaLogger(referential *Referential, stopArea *model.StopArea) *StopAreaLogger {
	logger := &StopAreaLogger{}

	// debug or not ?
	verbose := false
	for _, stopAreaObjectId := range referential.LoggerVerboseStopAreas() {
		objectId, ok := stopArea.ObjectID(stopAreaObjectId.Kind())
		if !ok {
			continue
		}
		if objectId.Value() == stopAreaObjectId.Value() {
			verbose = true
			break
		}
	}

	logger.verbose = verbose

	return logger
}

func (stopAreaLogger *StopAreaLogger) Printf(format string, values ...interface{}) {
	if stopAreaLogger.verbose {
		logger.Log.Printf(format, values...)
	} else {
		logger.Log.Debugf(format, values...)
	}
}

func (stopAreaLogger *StopAreaLogger) IsVerbose() bool {
	return stopAreaLogger.verbose
}
