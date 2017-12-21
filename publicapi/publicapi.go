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
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var (
	logger *logrus.Entry
	apiUrl = "https://poloniex.com/public"
	maxRequestsSec = 5
	httpClientTimeoutSec = 10
)

type Client struct {
	httpClient *http.Client
	throttle   <-chan time.Time
}

type configuration struct {
	apiConf `json:"poloniex_public_api"`
}

type apiConf struct {
	APIUrl               string `json:"api_url"`
	HTTPClientTimeoutSec int    `json:"httpclient_timeout_sec"`
	MaxRequestsSec       int    `json:"max_requests_sec"`
	LogLevel             string `json:"log_level"`
}

func init() {

	customFormatter := new(prefixed.TextFormatter)
	customFormatter.FullTimestamp = true
	customFormatter.ForceColors = true
	customFormatter.ForceFormatting = true
	logrus.SetFormatter(customFormatter)

	logger = logrus.WithField("prefix", "[api:poloniex:publicapi]")

	logrus.SetLevel(logrus.InfoLevel)
}

// NewClient returns a newly configured client
func NewClient() *Client {

	reqInterval := 1000 * time.Millisecond / time.Duration(maxRequestsSec)

	client := http.Client{
		Timeout: time.Duration(httpClientTimeoutSec) * time.Second,
	}

	return &Client{&client, time.Tick(reqInterval)}
}

// Do prepares and executes api call requests.
func (c *Client) do(params map[string]string) ([]byte, error) {

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

	u := apiUrl

	var parameters []string
	for k, v := range params {
		parameters = append(parameters, k+"="+url.QueryEscape(v))
	}

	if len(parameters) > 0 {
		return u + "?" + strings.Join(parameters, "&")
	}

	return u + strings.Join(parameters, "&")
}
