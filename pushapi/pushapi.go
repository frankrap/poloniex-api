// Poloniex push API implementation.
//
// API Doc: https://poloniex.com/support/api
//
// The best way to get public data updates on markets is via the push API,
// which pushes live ticker, order book, trade, and Trollbox updates over
// WebSockets using the WAMP protocol. In order to use the push API,
// connect to wss://api.poloniex.com and subscribe to the desired feed.
package pushapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	turnpike "gopkg.in/jcelliott/turnpike.v2"
)

var conf *configuration

type PushClient struct {
	wampClient *turnpike.Client
}

type configuration struct {
	PushAPI struct {
		WssUri string `json:"wss_uri"`
		Realm  string `json:"realm"`
	} `json:"push_api"`
}

func init() {

	content, err := ioutil.ReadFile("conf.json")

	if err != nil {
		log.Fatalf("loading configuration: %v", err)
	}

	if err := json.Unmarshal(content, &conf); err != nil {
		log.Fatalf("loading configuration: %v", err)
	}
}

func NewPushClient() (*PushClient, error) {

	turnpike.Debug()
	client, err := turnpike.NewWebsocketClient(turnpike.JSON,
		conf.PushAPI.WssUri, nil, nil)

	if err != nil {
		return nil, fmt.Errorf("turnpike.NewWebsocketClient: %v", err)
	}

	_, err = client.JoinRealm(conf.PushAPI.Realm, nil)
	if err != nil {
		return nil, fmt.Errorf("turnpike.Client.JoinRealm: %v", err)
	}

	return &PushClient{client}, nil
}

func (client *PushClient) LeaveRealm() error {

	if err := client.wampClient.LeaveRealm(); err != nil {
		return fmt.Errorf("turnpike.Client.LeaveRealm: %v", err)
	}
	return nil
}

func (client *PushClient) Close() error {

	if err := client.wampClient.Close(); err != nil {
		return fmt.Errorf("turnpike.Client.LeaveRealm: %v", err)
	}
	return nil
}

func convertStringToFloat(arg interface{}) (float64, error) {

	if v, ok := arg.(string); ok {

		val, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, fmt.Errorf("strconv.ParseFloat: %v", err)
		}
		return val, nil

	} else {
		return 0, fmt.Errorf("type assertion failed: %v", arg)
	}
}
