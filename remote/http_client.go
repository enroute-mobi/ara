package remote

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"bitbucket.org/enroute-mobi/ara/version"
	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"github.com/golang/protobuf/proto"
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

func (c *HTTPClient) SetURLs(urls HTTPClientUrls) {
	c.HTTPClientUrls = urls
}

func (c *HTTPClient) SOAPClient() *SOAPClient {
	return c.soapClient
}

func (c *HTTPClient) HTTPClient() *http.Client {
	return c.httpClient
}

func (c *HTTPClient) GTFSRequest() (*gtfs.FeedMessage, error) {
	ctx, cncl := context.WithTimeout(context.Background(), 60*time.Second)
	defer cncl()

	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, c.Url, nil)
	if err != nil {
		return nil, err
	}
	// httpRequest.SetBasicAuth(username, password)
	httpRequest.Header.Set("Accept-Encoding", "gzip, deflate")
	httpRequest.Header.Set("User-Agent", version.ApplicationName())

	// Send http request
	response, err := c.httpClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	defer func() {
		io.Copy(ioutil.Discard, response.Body)
		response.Body.Close()
	}()

	// Check response status
	if response.StatusCode != http.StatusOK {
		return nil, NewGtfsError(fmt.Sprintf("HTTP status %v", strconv.Itoa(response.StatusCode)))
	}

	if response.Header.Get("Content-Type") != "application/x-protobuf" {
		return nil, NewGtfsError(fmt.Sprintf("HTTP Content-Type %v", response.Header.Get("Content-Type")))
	}

	// Check if response is gzip
	var responseReader io.Reader
	if response.Header.Get("Content-Encoding") == "gzip" {
		gzipReader, err := gzip.NewReader(response.Body)
		if err != nil {
			return nil, NewGtfsError(fmt.Sprintf("GTFS response error: can't unzip response: %v", err))
		}
		defer gzipReader.Close()
		responseReader = gzipReader
	} else {
		responseReader = response.Body
	}

	// Attempt to read the body
	content, err := ioutil.ReadAll(responseReader)
	if err != nil {
		return nil, NewGtfsError(fmt.Sprintf("Error while reading body: %v", err))
	}
	if len(content) == 0 {
		return nil, NewGtfsError(fmt.Sprintf("Empty Body"))
	}

	feed := &gtfs.FeedMessage{}
	err = proto.Unmarshal(content, feed)
	if err != nil {
		return nil, NewGtfsError(fmt.Sprintf("Error while unmarshalling: %v", err))
	}

	return feed, nil
}
