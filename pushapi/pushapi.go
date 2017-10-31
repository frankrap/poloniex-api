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

	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	turnpike "gopkg.in/jcelliott/turnpike.v2"
)

var (
	conf   *configuration
	logger *logrus.Entry
)

type Client struct {
	wampClientMu sync.RWMutex
	wampClient   *turnpike.Client

	plu *pushLastUpdate
}

type pushLastUpdate struct {
	sync.RWMutex
	lastTimestamp      time.Time
	topicLastTimestamp map[string]time.Time
	subscription       map[string]func() error
}

type configuration struct {
	apiConf `json:"poloniex_push_api"`
}

type apiConf struct {
	WssUri          string `json:"wss_uri"`
	Realm           string `json:"realm"`
	LogLevel        string `json:"log_level"`
	TimeoutSec      int    `json:"timeout_sec"`
	TopicTimeoutMin int    `json:"topic_timeout_min"`
}

func init() {

	customFormatter := new(prefixed.TextFormatter)
	customFormatter.FullTimestamp = true
	customFormatter.ForceColors = true
	customFormatter.ForceFormatting = true
	logrus.SetFormatter(customFormatter)

	logger = logrus.WithField("prefix", "[api:poloniex:pushapi]")

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
		conf.WssUri, nil, nil, nil)

	if err != nil {
		return nil, fmt.Errorf("turnpike.NewWebsocketClient: %v", err)
	}

	_, err = client.JoinRealm(conf.Realm, nil)
	if err != nil {
		return nil, fmt.Errorf("turnpike.Client.JoinRealm: %v", err)
	}

	plu := &pushLastUpdate{
		lastTimestamp:      time.Now(),
		topicLastTimestamp: make(map[string]time.Time),
		subscription:       make(map[string]func() error),
	}

	res := &Client{
		sync.RWMutex{},
		client,
		plu,
	}

	go res.autoReconnect(time.Duration(conf.TimeoutSec) * time.Second)

	return res, nil
}

func (client *Client) autoReconnect(timeout time.Duration) {

	for {
		time.Sleep(timeout)

		client.plu.RLock()
		lastTimestamp := client.plu.lastTimestamp
		client.plu.RUnlock()

		if time.Since(lastTimestamp) > timeout {

			logger.Warn("Auto reconnecting...")
			var err error

			if err = client.Close(); err != nil {
				logger.WithField("error", err).Error(
					"PushClient.autoReconnect: PushClient.Close")
			}

			client.wampClientMu.Lock()
			for {

				time.Sleep(5 * time.Second)
				client.wampClient, err =
					turnpike.NewWebsocketClient(turnpike.JSON, conf.WssUri, nil, nil, nil)

				if err != nil {
					logger.WithField("error", err).Error(
						"PushClient.autoReconnect: turnpike.NewWebsocketClient")
					continue
				}

				_, err = client.wampClient.JoinRealm(conf.Realm, nil)
				if err != nil {
					logger.WithField("error", err).Error(
						"PushClient.autoReconnect: turnpike.Client.JoinRealm")
					continue
				}
				break
			}
			client.wampClientMu.Unlock()

			client.plu.Lock()
			client.plu.lastTimestamp = time.Now()

			logger.Infof("Resubscribing %d topics", len(client.plu.subscription))

			for _, subscribe := range client.plu.subscription {
				if err = subscribe(); err != nil {
					logger.WithField("error", err).Error(
						"PushClient.autoReconnect: subscribe")
				}
			}
			client.plu.Unlock()

			continue
		}

		client.plu.Lock()

		topicTimeout := time.Duration(conf.TopicTimeoutMin) * time.Minute
		topicsError := make(map[string]time.Time)

		for topic, timestamp := range client.plu.topicLastTimestamp {

			if time.Since(timestamp) > topicTimeout {
				topicsError[topic] = timestamp
			}
		}

		for topic, timestamp := range topicsError {

			logger.Infof("%s: no update since %s, resubscribing...",
				topic, time.Since(timestamp))

			if err := client.plu.subscription[topic](); err != nil {
				logger.WithField("error", err).Error(
					"PushClient.autoReconnect: subscribe")
			}
		}
		client.plu.Unlock()
	}
}

func (client *Client) addSubscription(topic string, subscribe func() error) {

	client.plu.Lock()
	defer client.plu.Unlock()

	client.plu.subscription[topic] = subscribe
	client.plu.topicLastTimestamp[topic] = time.Now()
}

func (client *Client) removeSubscription(topic string) {

	client.plu.Lock()
	defer client.plu.Unlock()

	delete(client.plu.subscription, topic)
	delete(client.plu.topicLastTimestamp, topic)
}

func (client *Client) updateTopicTimestamp(topic string) {

	client.plu.Lock()
	defer client.plu.Unlock()

	timestamp := time.Now()
	client.plu.lastTimestamp = timestamp
	client.plu.topicLastTimestamp[topic] = timestamp
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
