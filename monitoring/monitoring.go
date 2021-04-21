package monitoring

import (
	"net/http"
	"time"

	"bitbucket.org/enroute-mobi/ara/logger"
	"github.com/getsentry/sentry-go"
)

func HandleHttpPanic(response http.ResponseWriter) {
	err := recover()

	if err != nil {
		ReportError(err)
		http.Error(response, "Internal error", http.StatusInternalServerError)
	}
}

func HandlePanic() {
	err := recover()

	if err != nil {
		ReportError(err)
	}
}

func ReportError(err interface{}) {
	logger.Log.Printf("Error in processing: %v", err)
	sentry.CurrentHub().Recover(err)
	sentry.Flush(time.Second * 5)
}
