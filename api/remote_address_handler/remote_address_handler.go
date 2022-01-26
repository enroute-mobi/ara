package remote_address_handler

import (
	"net/http"
	"strings"
)

type RemoteAddressHandler struct{}

func (_ *RemoteAddressHandler) HandleRemoteAddress(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); len(xff) > 0 {
		return strings.Split(xff, ",")[0]
	}
	return strings.Split(r.RemoteAddr, ":")[0]
}
