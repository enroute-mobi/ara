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

	"bitbucket.org/enroute-mobi/ara/gtfs"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/version"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"google.golang.org/protobuf/proto"
)

type HTTPClientOptions struct {
	SiriEnvelopeType string
	OAuth            *HTTPClientOAuth
	Urls             HTTPClientUrls
}

type HTTPClientOAuth struct {
	ClientID     string
	ClientSecret string
	TokenURL     string
}

type HTTPClientUrls struct {
	Url              string
	SubscriptionsUrl string
	NotificationsUrl string
}

type HTTPClient struct {
	HTTPClientUrls

	httpClient *http.Client
	siriClient *SIRIClient
}

func NewHTTPClient(opts HTTPClientOptions) *HTTPClient {
	c := &HTTPClient{
		HTTPClientUrls: opts.Urls,
		httpClient:     httpClient(opts),
	}
	sc := NewSIRIClient(c, opts.SiriEnvelopeType)
	c.siriClient = sc

	return c
}

func httpClient(opts HTTPClientOptions) (c *http.Client) {
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
	c = &http.Client{
		Timeout:   60 * time.Second,
		Transport: netTransport,
	}

	if opts.OAuth == nil {
		return c
	}

	oauthConfig := clientcredentials.Config{
		ClientID:     opts.OAuth.ClientID,
		ClientSecret: opts.OAuth.ClientSecret,
		TokenURL:     opts.OAuth.TokenURL,
	}

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, c)

	token, err := oauthConfig.Token(ctx)

	if err != nil {
		logger.Log.Printf("Could not authenticate with OAuth: %s", err)
		return c // Return default server if we can't authenticate with OAuth
	} else {
		logger.Log.Debugf("Authenticated to OAuth server url: %s, token_type: %s, expiring at: %s",
			oauthConfig.TokenURL,
			token.TokenType,
			token.Expiry.UTC().Format(time.UnixDate))
	}

	return oauthConfig.Client(ctx)
}

func (c *HTTPClient) SetURLs(urls HTTPClientUrls) {
	c.HTTPClientUrls = urls
}

func (c *HTTPClient) SIRIClient() *SIRIClient {
	return c.siriClient
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

	// ARA-878
	// if response.Header.Get("Content-Type") != "application/x-protobuf" {
	// 	return nil, NewGtfsError(fmt.Sprintf("HTTP Content-Type %v", response.Header.Get("Content-Type")))
	// }

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
		return nil, NewGtfsError("Empty Body")
	}

	feed := &gtfs.FeedMessage{}
	err = proto.Unmarshal(content, feed)
	if err != nil {
		return nil, NewGtfsError(fmt.Sprintf("Error while unmarshalling: %v", err))
	}

	return feed, nil
}
