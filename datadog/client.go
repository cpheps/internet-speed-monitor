// Package datadog publishes speedtest results to DataDog
package datadog

import (
	"fmt"
	"os"

	"github.com/cpheps/internet-speed-monitor/speedtest"
	"gopkg.in/zorkian/go-datadog-api.v2"
)

const (
	apiKeyEnv = "DD_API_KEY"
	appKeyEnv = "DD_APP_KEY"
)

// Client is a wrapper around the DataDog client
type Client struct {
	client *datadog.Client
}

// NewClient creates and returns a new client. Uses the ENVs to look up API and App keys.
// The ENVs are:
//  DD_API_KEY - your API key
//  DD_APP_KEY - your app key
func NewClient() (*Client, error) {
	apiKey, appKey, err := loadKeys()
	if err != nil {
		return nil, err
	}

	return &Client{
		client: datadog.NewClient(apiKey, appKey),
	}, nil
}

// SendTestResults sends the results of a speed test as metrics to datadog
func (c Client) SendTestResults(results *speedtest.Results, classifier string) error {
	metricSeries := make([]datadog.Metric, 0, 3)

	unixTimestamp := float64(results.Timstamp.Unix())

	// Upload metric
	metricSeries = append(metricSeries, datadog.Metric{
		Metric: datadog.String("upload_speed"),
		Points: []datadog.DataPoint{{&unixTimestamp, &results.Upload}},
		Type:   datadog.String("gauge"),
		Host:   &classifier,
		Unit:   datadog.String("bit"),
	})

	// Download metric
	metricSeries = append(metricSeries, datadog.Metric{
		Metric: datadog.String("download_speed"),
		Points: []datadog.DataPoint{{&unixTimestamp, &results.Download}},
		Type:   datadog.String("gauge"),
		Host:   &classifier,
		Unit:   datadog.String("bit"),
	})

	// Ping metric
	metricSeries = append(metricSeries, datadog.Metric{
		Metric: datadog.String("ping"),
		Points: []datadog.DataPoint{{&unixTimestamp, &results.Download}},
		Type:   datadog.String("gauge"),
		Host:   &classifier,
		Unit:   datadog.String("millisecond"),
	})

	return c.client.PostMetrics(metricSeries)
}

func loadKeys() (apiKey, appKey string, err error) {
	apiKey, ok := os.LookupEnv(apiKeyEnv)
	if !ok {
		err = fmt.Errorf("could not find ENV %s", apiKeyEnv)
		return
	}

	appKey, ok = os.LookupEnv(appKeyEnv)
	if !ok {
		err = fmt.Errorf("could not find ENV %s", appKeyEnv)
		return
	}

	return
}
