package poloniex

import (
	"fmt"

	turnpike "gopkg.in/jcelliott/turnpike.v2"
)

const (
	WSSURI = "wss://api.poloniex.com"
	REALM  = "realm1"
	TICKER = "ticker"
)

func NewClient() (*turnpike.Client, error) {

	turnpike.Debug()
	client, err := turnpike.NewWebsocketClient(turnpike.JSON, WSSURI, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("newclient: %v", err)
	}

	_, err = client.JoinRealm(REALM, nil)
	if err != nil {
		return nil, fmt.Errorf("newclient: %v", err)
	}

	return client, nil
}
