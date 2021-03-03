package monitoring

import (
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
)

func HandleHttpPanic(response http.ResponseWriter) {
	err := recover()

	if err != nil {
		sentry.CurrentHub().Recover(err)
		sentry.Flush(time.Second * 5)
		http.Error(response, "Internal error", http.StatusInternalServerError)
	}
}

func HandlePanic() {
	err := recover()

	if err != nil {
		sentry.CurrentHub().Recover(err)
		sentry.Flush(time.Second * 5)
	}
}
