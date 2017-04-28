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
	"strconv"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	turnpike "gopkg.in/jcelliott/turnpike.v2"
)

var (
	conf   *configuration
	logger *logrus.Entry
)

type Client struct {
	wampClientMu sync.RWMutex
	wampClient   *turnpike.Client
	subscription map[string]func() error

	mc *msgCount
}

type msgCount struct {
	sync.Mutex
	count         uint64
	lastTimestamp time.Time
}

type configuration struct {
	apiConf `json:"poloniex_push_api"`
}

type apiConf struct {
	WssUri     string `json:"wss_uri"`
	Realm      string `json:"realm"`
	LogLevel   string `json:"log_level"`
	TimeoutSec int    `json:"timeout_sec"`
}

func init() {

	customFormatter := new(logrus.TextFormatter)
	customFormatter.FullTimestamp = true
	logrus.SetFormatter(customFormatter)

	logger = logrus.WithField("context", "[api:poloniex:pushapi]")

	content, err := ioutil.ReadFile("conf.json")

	if err != nil {
		logger.WithField("error", err).Fatal("loading configuration")
	}

	if err := json.Unmarshal(content, &conf); err != nil {
		logger.WithField("error", err).Fatal("loading configuration")
	}

	switch conf.LogLevel {
	case "debug":
		turnpike.Debug()
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

func NewClient() (*Client, error) {

	client, err := turnpike.NewWebsocketClient(turnpike.JSON,
		conf.WssUri, nil, nil)

	if err != nil {
		return nil, fmt.Errorf("turnpike.NewWebsocketClient: %v", err)
	}

	_, err = client.JoinRealm(conf.Realm, nil)
	if err != nil {
		return nil, fmt.Errorf("turnpike.Client.JoinRealm: %v", err)
	}

	res := &Client{
		sync.RWMutex{},
		client,
		make(map[string]func() error), &msgCount{},
	}

	go res.autoReconnect(time.Duration(conf.TimeoutSec) * time.Second)

	return res, nil
}

func (client *Client) autoReconnect(timeout time.Duration) {

	for {

		time.Sleep(timeout)

		client.mc.Lock()
		count := client.mc.count
		lastTimestamp := client.mc.lastTimestamp
		client.mc.Unlock()

		if count > 0 && time.Since(lastTimestamp) > timeout {

			logger.Warn("Auto reconnecting...")
			var err error

			if err = client.Close(); err != nil {
				logger.WithField("error", err).Error("PushClient.autoReconnect: PushClient.Close")
			}

			client.wampClientMu.Lock()

			for {

				time.Sleep(5 * time.Second)
				client.wampClient, err = turnpike.NewWebsocketClient(turnpike.JSON, conf.WssUri, nil, nil)

				if err != nil {
					logger.WithField("error", err).Error("PushClient.autoReconnect: turnpike.NewWebsocketClient")
					continue
				}

				_, err = client.wampClient.JoinRealm(conf.Realm, nil)
				if err != nil {
					logger.WithField("error", err).Error("PushClient.autoReconnect: turnpike.Client.JoinRealm")
					continue
				}

				client.mc.Lock()
				client.mc.count = 0
				client.mc.Unlock()
				break
			}

			subscribes := make(map[string]func() error)
			for topic, subscribe := range client.subscription {
				subscribes[topic] = subscribe
			}
			client.wampClientMu.Unlock()

			logger.WithField("subscriptions", subscribes).Infof("Resubscribing %d topics", len(subscribes))

			for _, subscribe := range subscribes {
				if err = subscribe(); err != nil {
					logger.WithField("error", err).Error("PushClient.autoReconnect: subscribe")
				}
			}
		}
	}
}

func (client *Client) addSubscription(topic string, subscribe func() error) {

	client.wampClientMu.Lock()
	defer client.wampClientMu.Unlock()

	client.subscription[topic] = subscribe
}

func (client *Client) removeSubscription(topic string) {

	client.wampClientMu.Lock()
	defer client.wampClientMu.Unlock()

	delete(client.subscription, topic)
}

func (client *Client) updateMsgCount() {

	client.mc.Lock()
	defer client.mc.Unlock()

	client.mc.count++
	client.mc.lastTimestamp = time.Now()
}

func (client *Client) Close() error {

	client.wampClientMu.RLock()
	defer client.wampClientMu.RUnlock()

	if err := client.wampClient.Close(); err != nil {
		return fmt.Errorf("turnpike.Client.Close: %v", err)
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
