// Poloniex trading API implementation.
//
// API Doc: https://poloniex.com/support/api
//
// To use the trading API, you will need to create an API key.
// Please note that there is a default limit of 6 calls per second. If you require more
// than this, please consider optimizing your application using the push API, the "moveOrder"
// command, or the "all" parameter where appropriate. If this is still insufficient, please
// contact support to discuss a limit raise.
//
// All calls to the trading API are sent via HTTP POST to https://poloniex.com/tradingApi
// and must contain the following headers:
//
// Key - Your API key.
// Sign - The query's POST data signed by your key's "secret" according to the HMAC-SHA512 method.
// Additionally, all queries must include a "nonce" POST parameter. The nonce parameter is an integer
// which must always be greater than the previous nonce used.

// All responses from the trading API are in JSON format. In the event of an error, the response will
// always be of the following format:
//
// {"error":"<error message>"}
//
// There are several methods accepted by the trading API, each of which is specified by the "command"
// POST parameter.
package tradingapi

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var (
	conf   *configuration
	logger *logrus.Entry
)

type Client struct {
	apiKey     string
	apiSecret  string
	httpClient *http.Client
	throttle   <-chan time.Time
}

type APIError struct {
	Err string `json:"error"`
}

type configuration struct {
	apiConf `json:"poloniex_trading_api"`
}

type apiConf struct {
	APIUrl               string `json:"api_url"`
	HTTPClientTimeoutSec int    `json:"httpclient_timeout_sec"`
	MaxRequestsSec       int    `json:"max_requests_sec"`
	ApiKey               string `json:"api_key"`
	ApiSecret            string `json:"api_secret"`
	LogLevel             string `json:"log_level"`
}

// Loading configuration
func init() {

	customFormatter := new(prefixed.TextFormatter)
	customFormatter.FullTimestamp = true
	customFormatter.ForceColors = true
	customFormatter.ForceFormatting = true
	logrus.SetFormatter(customFormatter)

	logger = logrus.WithField("prefix", "[api:poloniex:tradingapi]")

	content, err := ioutil.ReadFile("conf.json")

	if err != nil {
		logger.WithField("error", err).Fatal("loading configuration")
	}

	if err := json.Unmarshal(content, &conf); err != nil {
		logger.WithField("error", err).Fatal("loading configuration")
	}

	switch conf.LogLevel {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "fatal":
		logrus.SetLevel(logrus.FatalLevel)
	case "panic":
		logrus.SetLevel(logrus.PanicLevel)
	default:
		logrus.SetLevel(logrus.WarnLevel)
	}
}

// NewClient returns a newly configured client
func NewClient() (*Client, error) {

	reqInterval := 1000 * time.Millisecond / time.Duration(conf.MaxRequestsSec)

	client := http.Client{
		Timeout: time.Duration(conf.HTTPClientTimeoutSec) * time.Second,
	}

	if len(conf.ApiKey) == 0 || len(conf.ApiSecret) == 0 {

		err := errors.New("new trading client: wrong apikey and/or apisecret")
		return nil, err
	}

	tc := Client{
		conf.ApiKey,
		conf.ApiSecret,
		&client,
		time.Tick(reqInterval),
	}

	return &tc, nil
}

// Do prepares and executes api call requests.
func (c *Client) do(form url.Values) ([]byte, error) {

	nonce := time.Now().UnixNano()
	form.Add("nonce", strconv.Itoa(int(nonce)))

	req, err := http.NewRequest("POST",
		conf.APIUrl,
		strings.NewReader(form.Encode()))

	if err != nil {
		return nil, fmt.Errorf("http.NewRequest: %v (API command: %s)",
			err, form.Get("command"))
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Key", c.apiKey)

	if sig, err := signForm(form, c.apiSecret); err != nil {
		return nil, fmt.Errorf("signForm: %v", err)
	} else {
		req.Header.Add("Sign", sig)
	}

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
		return body, fmt.Errorf("Status code: %s (API command: %s)",
			res.resp.Status, form.Get("command"))
	}

	if err := checkAPIError(body); err != nil {
		return nil, err
	}

	return body, nil
}

func checkAPIError(body []byte) error {

	if !strings.Contains(string(body), "\"error\":") {
		return nil
	}

	ae := APIError{}
	if err := json.Unmarshal(body, &ae); err == nil {
		return fmt.Errorf("API error: %s", ae.Err)

	}

	return nil
}

func signForm(form url.Values, apiSecret string) (string, error) {

	mac := hmac.New(sha512.New, []byte(apiSecret))
	_, err := mac.Write([]byte(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("hash.Hash.Write: %v", err)
	}
	sig := hex.EncodeToString(mac.Sum(nil))

	return sig, nil
}
