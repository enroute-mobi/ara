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
	for _, stopAreaCode := range referential.LoggerVerboseStopAreas() {
		code, ok := stopArea.Code(stopAreaCode.CodeSpace())
		if !ok {
			continue
		}
		if code.Value() == stopAreaCode.Value() {
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
