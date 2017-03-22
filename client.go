package poloniex

import (
	"fmt"
	"net/url"
	"strings"

	turnpike "gopkg.in/jcelliott/turnpike.v2"
)

const (
	PUBLIC_API_URL = "https://poloniex.com/public"
	WSSURI         = "wss://api.poloniex.com"
	REALM          = "realm1"
	TICKER         = "ticker"
)

func NewClient() (*turnpike.Client, error) {

	turnpike.Debug()
	client, err := turnpike.NewWebsocketClient(turnpike.JSON, WSSURI, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("new websocket client: %v", err)
	}

	_, err = client.JoinRealm(REALM, nil)
	if err != nil {
		return nil, fmt.Errorf("join realm: %v", err)
	}

	return client, nil
}

func buildUrl(params map[string]string) string {

	u := PUBLIC_API_URL + "?"

	var parameters []string
	for k, v := range params {
		parameters = append(parameters, k+"="+url.QueryEscape(v))
	}

	return u + strings.Join(parameters, "&")
}
