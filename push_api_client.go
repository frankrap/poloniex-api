package poloniex

import (
	"fmt"

	turnpike "gopkg.in/jcelliott/turnpike.v2"
)

const (
	WSSURI = "wss://api.poloniex.com"
	REALM  = "realm1"
)

type PushClient struct {
	wampClient *turnpike.Client
}

func NewPushClient() (*PushClient, error) {

	turnpike.Debug()
	client, err := turnpike.NewWebsocketClient(turnpike.JSON, WSSURI, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("new websocket client: %v", err)
	}

	_, err = client.JoinRealm(REALM, nil)
	if err != nil {
		return nil, fmt.Errorf("join realm: %v", err)
	}

	return &PushClient{client}, nil
}
