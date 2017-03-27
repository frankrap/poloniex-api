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
)

const (
	TRADING_API_URL            = "https://poloniex.com/tradingApi"
	DEFAULT_HTTPCLIENT_TIMEOUT = 10
	MAX_REQUEST_PER_SECOND     = 5
)

type TradingClient struct {
	apiKey     string
	apiSecret  string
	httpClient *http.Client
	throttle   <-chan time.Time
}

type APIError struct {
	Err string `json:"error"`
}

// NewTradingClient returns a newly configured client
func NewTradingClient(apiKey, apiSecret string) (*TradingClient, error) {

	reqInterval := 1000 / MAX_REQUEST_PER_SECOND * time.Millisecond
	client := http.Client{
		Timeout: DEFAULT_HTTPCLIENT_TIMEOUT * time.Second,
	}

	if len(apiKey) == 0 || len(apiSecret) == 0 {
		return nil, errors.New("new trading client: wrong apikey and/or apisecret")
	}

	tc := TradingClient{
		apiKey,
		apiSecret,
		&client,
		time.Tick(reqInterval),
	}

	return &tc, nil
}

// Do prepares and executes api call requests.
func (c *TradingClient) do(form url.Values) ([]byte, error) {

	nonce := time.Now().UnixNano()
	form.Add("nonce", strconv.Itoa(int(nonce)))

	req, err := http.NewRequest("POST", TRADING_API_URL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("new request: %v", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Key", c.apiKey)

	if sig, err := signForm(form, c.apiSecret); err != nil {
		return nil, fmt.Errorf("sign form: %v", err)
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
		return nil, fmt.Errorf("request: %v", res.err)
	}

	defer res.resp.Body.Close()

	body, err := ioutil.ReadAll(res.resp.Body)
	if err != nil {
		return body, fmt.Errorf("readall: %v", err)
	}

	if res.resp.StatusCode != 200 {
		return body, fmt.Errorf("status code: %s", res.resp.Status)
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
		return "", err
	}
	sig := hex.EncodeToString(mac.Sum(nil))

	return sig, nil
}
