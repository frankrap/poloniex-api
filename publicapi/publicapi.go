// Poloniex public API implementation.
//
// API Doc: https://poloniex.com/support/api
//
// Please note that making more than 6 calls per second to the public API, or repeatedly and
// needlessly fetching excessive amounts of data, can result in your IP being banned.
//
// There are six public methods, all of which take HTTP GET requests and return output in JSON format.
package publicapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

var conf *configuration

type PublicClient struct {
	httpClient *http.Client
	throttle   <-chan time.Time
}

type configuration struct {
	PublicAPI struct {
		PublicAPIUrl         string `json:"public_api_url"`
		HTTPClientTimeoutSec int    `json:"httpclient_timeout_sec"`
		MaxRequestsSec       int    `json:"max_requests_sec"`
		LogLevel             string `json:"log_level"`
	} `json:"public_api"`
}

func init() {

	content, err := ioutil.ReadFile("conf.json")

	if err != nil {
		log.WithField("error", err).Fatal("loading configuration")
	}

	if err := json.Unmarshal(content, &conf); err != nil {
		log.WithField("error", err).Fatal("loading configuration")
	}

	switch conf.PublicAPI.LogLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "panic":
		log.SetLevel(log.PanicLevel)
	default:
		log.SetLevel(log.WarnLevel)
	}
}

// NewPublicClient returns a newly configured client
func NewPublicClient() *PublicClient {

	reqInterval := 1000 * time.Millisecond /
		time.Duration(conf.PublicAPI.MaxRequestsSec)

	client := http.Client{
		Timeout: time.Duration(conf.PublicAPI.HTTPClientTimeoutSec) *
			time.Second,
	}

	return &PublicClient{&client, time.Tick(reqInterval)}
}

// Do prepares and executes api call requests.
func (c *PublicClient) do(params map[string]string) ([]byte, error) {

	url := buildUrl(params)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest: %v (API command: %s)",
			err, params["command"])
	}

	req.Header.Add("Accept", "application/json")

	type result struct {
		resp *http.Response
		err  error
	}

	done := make(chan result)
	go func() {
		<-c.throttle
		resp, err := c.httpClient.Do(req)
		done <- result{resp, err}
	}()
	res := <-done

	if res.err != nil {
		return nil, fmt.Errorf("http.Client.Do: %v", res.err)
	}

	defer res.resp.Body.Close()

	body, err := ioutil.ReadAll(res.resp.Body)
	if err != nil {
		return body, fmt.Errorf("ioutil.readAll: %v", err)
	}

	if res.resp.StatusCode != 200 {
		return body, fmt.Errorf("status code: %s (API command: %s)",
			res.resp.Status, params["command"])
	}

	return body, nil
}

func buildUrl(params map[string]string) string {

	u := conf.PublicAPI.PublicAPIUrl + "?"

	var parameters []string
	for k, v := range params {
		parameters = append(parameters, k+"="+url.QueryEscape(v))
	}

	return u + strings.Join(parameters, "&")
}
