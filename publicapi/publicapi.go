// Poloniex public API implementation.
//
// API Doc: https://poloniex.com/support/api/
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
)

const (
	PUBLIC_API_URL             = "https://poloniex.com/public"
	DEFAULT_HTTPCLIENT_TIMEOUT = 10
	MAX_REQUEST_PER_SECOND     = 5
)

type PublicClient struct {
	httpClient *http.Client
	throttle   <-chan time.Time
}

// NewPublicClient returns a newly configured client
func NewPublicClient() *PublicClient {

	reqInterval := 1000 / MAX_REQUEST_PER_SECOND * time.Millisecond
	client := http.Client{
		Timeout: DEFAULT_HTTPCLIENT_TIMEOUT * time.Second,
	}

	return &PublicClient{&client, time.Tick(reqInterval)}
}

// Do prepares and executes api call requests.
func (c *PublicClient) do(method, url, payload string, auth bool) ([]byte, error) {

	req, err := http.NewRequest(method, url, strings.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("new request: %v", err)
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
		return nil, fmt.Errorf("request: %v", res.err)
	}
	defer res.resp.Body.Close()

	body, err := ioutil.ReadAll(res.resp.Body)
	if err != nil {
		return body, fmt.Errorf("readall: %v", err)
	}

	if res.resp.StatusCode != 200 {
		return body, fmt.Errorf("status code: %d", res.resp.Status)
	}

	return body, nil
}

func buildUrl(params map[string]string) string {

	u := PUBLIC_API_URL + "?"

	var parameters []string
	for k, v := range params {
		parameters = append(parameters, k+"="+url.QueryEscape(v))
	}

	return u + strings.Join(parameters, "&")
}
