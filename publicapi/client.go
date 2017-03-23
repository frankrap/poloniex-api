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

func buildUrl(params map[string]string) string {

	u := PUBLIC_API_URL + "?"

	var parameters []string
	for k, v := range params {
		parameters = append(parameters, k+"="+url.QueryEscape(v))
	}

	return u + strings.Join(parameters, "&")
}

func NewPublicClient() *PublicClient {

	reqInterval := 1000 / MAX_REQUEST_PER_SECOND * time.Millisecond
	client := http.Client{
		Timeout: DEFAULT_HTTPCLIENT_TIMEOUT * time.Second,
	}

	return &PublicClient{&client, time.Tick(reqInterval)}
}

func (c *PublicClient) do(method, url, payload string, auth bool) ([]byte, error) {

	req, err := http.NewRequest(method, url, strings.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("new request: %v", err)
	}

	if method == "POST" || method == "PUT" {
		req.Header.Add("Content-Type", "application/json;charset=utf-8")
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

// Auth
// if authNeeded {
//  if len(c.apiKey) == 0 || len(c.apiSecret) == 0 {
//      return nil, errors.New("You need to set API Key and API Secret to call this method")
//  }
//  nonce := time.Now().UnixNano()
//  q := req.URL.Query()
//  q.Set("apikey", c.apiKey)
//  q.Set("nonce", fmt.Sprintf("%d", nonce))
//  req.URL.RawQuery = q.Encode()
//  mac := hmac.New(sha512.New, []byte(c.apiSecret))
//  _, err = mac.Write([]byte(req.URL.String()))
//  sig := hex.EncodeToString(mac.Sum(nil))
//  req.Header.Add("apisign", sig)
// }
