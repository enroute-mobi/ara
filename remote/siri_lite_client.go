package remote

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"bitbucket.org/enroute-mobi/ara/siri/slite"
	"bitbucket.org/enroute-mobi/ara/version"
	"github.com/pkg/errors"
)

type SIRILiteClient struct {
	HTTPClientUrls

	httpClient *HTTPClient
}

type siriLiteClientArguments struct {
	params           url.Values
	requestType      requestType
	expectedResponse string
	destination      interface{}
}

func NewSIRILiteClient(c *HTTPClient) *SIRILiteClient {
	return &SIRILiteClient{
		httpClient: c,
	}
}

func (client *SIRILiteClient) remoteClient() *http.Client {
	return client.httpClient.HTTPClient()
}

func (c *SIRILiteClient) prepareAndSendRequest(args siriLiteClientArguments) error {
	dest := args.destination

	ctx, cncl := context.WithTimeout(context.Background(), getTimeOut(args.requestType))
	defer cncl()

	buildUrl, err := URI(c.getURL(args.requestType), "", args.params)
	if err != nil {
		return errors.Wrap(err, "unable to build URI")
	}

	// Create http request
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, buildUrl.String(), nil)
	if err != nil {
		return errors.Wrap(err, "cannot create request")
	}

	httpRequest.Header.Add("Accept", "application/json")
	httpRequest.Header.Add("Content-Type", "application/json")
	httpRequest.Header.Set("User-Agent", version.ApplicationName())

	// Send http request
	resp, err := c.remoteClient().Do(httpRequest)
	if err != nil {
		return errors.Wrap(err, "cannot proceed request")
	}

	// Attempt to read the body
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "cannot read the response")
	}

	// Check empty body
	if len(body) == 0 {
		return errors.Wrap(err, "empty body")
	}

	// Check response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return c.handleUnsuccessfullStatusCode(resp.StatusCode, args.expectedResponse, body)
	}

	// rewrite the payload
	rewriteBody, err := slite.RewriteValues(body)
	if err != nil {
		return errors.Wrap(err, "cannot rewrite the payload")
	}

	// Parse the body
	err = json.NewDecoder(bytes.NewReader(rewriteBody)).Decode(&dest)
	if err != nil {
		return errors.Wrap(err, "cannot parse server response")
	}

	return nil
}

func (client *SIRILiteClient) StopMonitoring(stopArea string) (*slite.SIRILiteStopMonitoring, error) {
	params := url.Values{}
	params.Add("MonitoringRef", stopArea)

	dest := &slite.SIRILiteStopMonitoring{}

	err := client.prepareAndSendRequest(siriLiteClientArguments{
		params:           params,
		requestType:      NOTIFICATION,
		expectedResponse: "StopMonitoringDelivery",
		destination:      dest,
	})
	if err != nil {
		return nil, err
	}

	return dest, nil
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

func (client *SIRILiteClient) getURL(requestType requestType) string {
	switch requestType {
	case SUBSCRIPTION:
		if client.httpClient.SubscriptionsUrl != "" {
			return client.httpClient.SubscriptionsUrl
		}
	case NOTIFICATION:
		if client.httpClient.NotificationsUrl != "" {
			return client.httpClient.NotificationsUrl
		}
	}
	return client.httpClient.Url
}

func (client *SIRILiteClient) handleUnsuccessfullStatusCode(statusCode int, expectedResponse string, body []byte) error {
	switch expectedResponse {
	case "StopMonitoringDelivery":
		errData := slite.SIRILiteStopMonitoring{}
		var errMsg string
		if err := json.Unmarshal(body, &errData); err == nil {
			errMsg = errData.
				Siri.
				ServiceDelivery.
				StopMonitoringDelivery[0].
				ErrorCondition.
				ErrorInformation.
				ErrorText
			return fmt.Errorf("request failed with status %d: %s", statusCode, errMsg)
		} else {
			// cannot parse the response
			return fmt.Errorf("request failed with status %d: %s", statusCode, strings.Replace(string(body), "\n", "", -1))
		}
	}
	return fmt.Errorf("request failed with status %d: ", statusCode)
}
