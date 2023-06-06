package remote

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"bitbucket.org/enroute-mobi/ara/gtfs"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/siri/slite"
	"bitbucket.org/enroute-mobi/ara/version"
	"github.com/pkg/errors"
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
	}

	logger.Log.Debugf("Authenticated to OAuth server url: %s, token: %s, token_type: %s, expiring at: %s",
		oauthConfig.TokenURL,
		token.AccessToken,
		token.TokenType,
		token.Expiry.UTC().Format(time.UnixDate))

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
		io.Copy(io.Discard, response.Body)
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
	content, err := io.ReadAll(responseReader)
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

func (c *HTTPClient) SIRILiteStopMonitoringRequest(dest interface{}, stopArea string) (string, error) {
	var rawQuery string
	ctx, cncl := context.WithTimeout(context.Background(), 60*time.Second)
	defer cncl()

	// Prepare URI
	params := url.Values{}
	params.Add("MonitoringRef", stopArea)

	buildUrl, err := URI(c.Url, "", params)
	if err != nil {
		return rawQuery, errors.Wrap(err, "unable to build URI")
	}

	rawQuery, _ = url.PathUnescape(buildUrl.RawQuery)

	// Create http request
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, buildUrl.String(), nil)
	if err != nil {
		return rawQuery, errors.Wrap(err, "cannot create request")
	}

	httpRequest.Header.Add("Accept", "application/json")
	httpRequest.Header.Add("Content-Type", "application/json")
	httpRequest.Header.Set("User-Agent", version.ApplicationName())

	// Send http request
	resp, err := c.httpClient.Do(httpRequest)
	if err != nil {
		return rawQuery, errors.Wrap(err, "cannot proceed request")
	}

	// Attempt to read the body
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return rawQuery, errors.Wrap(err, "cannot read the response")
	}

	// Check empty body
	if len(body) == 0 {
		return rawQuery, errors.Wrap(err, "empty body")
	}

	// Check response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errMsg := string(body)
		errData := &slite.SIRILiteStopMonitoring{}

		if json.Unmarshal(body, &errData); err == nil {
			errMsg = errData.
				Siri.
				ServiceDelivery.
				StopMonitoringDelivery[0].
				ErrorCondition.
				ErrorInformation.
				ErrorText
		}

		return rawQuery, fmt.Errorf("request failed with status %d: %s",
			resp.StatusCode, errMsg)
	}

	// Parse the body
	if dest != nil {
		err := json.NewDecoder(bytes.NewReader(body)).Decode(&dest)
		if err != nil {
			return rawQuery, errors.Wrap(err, "cannot parse server response")
		}
	}

	return rawQuery, nil
}

func URI(baseurl string, path string, params url.Values) (*url.URL, error) {
	base, err := url.Parse(baseurl)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to parse base url %v", baseurl)
	}

	if params == nil {
		params = url.Values{}
	}
	u := base.ResolveReference(&url.URL{Path: path, RawQuery: params.Encode()})
	return u, nil
}
