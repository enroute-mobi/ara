package rah

import (
	"net/http"
	"strings"
)

type RemoteAddressHandler struct{}

func (*RemoteAddressHandler) HandleRemoteAddress(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); len(xff) > 0 {
		return strings.Split(xff, ",")[0]
	}
	return strings.Split(r.RemoteAddr, ":")[0]
}
