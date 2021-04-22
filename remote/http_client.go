package remote

import (
	"net/http"
	"time"
)

type HTTPClientUrls struct {
	Url              string
	SubscriptionsUrl string
	NotificationsUrl string
}

type HTTPClient struct {
	HTTPClientUrls

	httpClient *http.Client
	soapClient *SOAPClient
}

func NewHTTPClient(urls HTTPClientUrls) *HTTPClient {
	// Customize the Transport based on DefaultTransport
	// DefaultTransport for reference
	// var DefaultTransport RoundTripper = &Transport{
	// 	Proxy: ProxyFromEnvironment,
	// 	DialContext: (&net.Dialer{
	// 		Timeout:   30 * time.Second,
	// 		KeepAlive: 30 * time.Second,
	// 	}).DialContext,
	// 	ForceAttemptHTTP2:     true,
	// 	MaxIdleConns:          100,
	// 	IdleConnTimeout:       90 * time.Second,
	// 	TLSHandshakeTimeout:   10 * time.Second,
	// 	ExpectContinueTimeout: 1 * time.Second,
	// }

	netTransport := http.DefaultTransport.(*http.Transport).Clone()
	netTransport.MaxConnsPerHost = 30
	netTransport.MaxIdleConnsPerHost = 10
	netTransport.TLSHandshakeTimeout = 5 * time.Second

	// set a long default time for safety, but we use context for request specific timeouts
	httpClient := &http.Client{
		Timeout:   60 * time.Second,
		Transport: netTransport,
	}

	c := &HTTPClient{
		HTTPClientUrls: urls,
		httpClient:     httpClient,
	}
	sc := NewSOAPClient(c)
	c.soapClient = sc

	return c
}

func (c *HTTPClient) SOAPClient() *SOAPClient {
	return c.soapClient
}

func (c *HTTPClient) HTTPClient() *http.Client {
	return c.httpClient
}
