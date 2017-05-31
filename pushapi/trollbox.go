package pushapi

import (
	"fmt"
	"sync"
)

const (
	TROLLBOX = "trollbox"
)

type TrollboxMessage struct {
	TypeMessage   string
	MessageNumber int64
	Username      string
	Message       string
	Reputation    int
}

type Trollbox chan *TrollboxMessage

var (
	trollbox = make(Trollbox)

	trollboxMu           sync.RWMutex
	trollboxUnsubscribed = make(chan struct{})
)

// Poloniex push API implementation of trollbox topic.
//
// API Doc:
// In order to receive new Trollbox messages, subscribe to "trollbox".
//
// Messages will be given in the following format:
//
// [type, messageNumber, username, message, reputation]
//
// Example:
//
// ['trollboxMessage',2094211,'boxOfTroll','Trololol',4]
func (client *Client) SubscribeTrollbox() (Trollbox, error) {

	handler := func(args []interface{}, kwargs map[string]interface{}) {

		client.updateTopicTimestamp(TROLLBOX)

		tbMsg, err := convertArgsToTrollboxMessage(args)
		if err != nil {
			logger.WithField("error", err).Error("convertArgsToTrollboxMessage")
			return
		}

		trollboxMu.RLock()
		select {
		case trollbox <- tbMsg:
		case <-trollboxUnsubscribed:
		}
		trollboxMu.RUnlock()
	}

	subscribe := func() error {

		client.wampClientMu.RLock()
		defer client.wampClientMu.RUnlock()

		if err := client.wampClient.Subscribe(TROLLBOX, nil, handler); err != nil {
			return fmt.Errorf("turnpike.Client.Subscribe: %v", err)
		}
		logger.Infof("Subscribed to: %s", TROLLBOX)

		return nil
	}

	if err := subscribe(); err != nil {
		return nil, err
	}

	trollboxMu.Lock()
	select {
	case <-trollboxUnsubscribed:
		trollboxUnsubscribed = make(chan struct{})
	default:
	}
	trollboxMu.Unlock()

	client.addSubscription(TROLLBOX, subscribe)

	return trollbox, nil
}

func (client *Client) UnsubscribeTrollbox() error {

	client.wampClientMu.RLock()

	if err := client.wampClient.Unsubscribe(TROLLBOX); err != nil {
		return fmt.Errorf("turnpike.Client.Unsuscribe: %v", err)
	}

	client.wampClientMu.RUnlock()

	client.removeSubscription(TROLLBOX)

	trollbox <- nil

	trollboxMu.RLock()
	close(trollboxUnsubscribed)
	trollboxMu.RUnlock()

	return nil
}

func convertArgsToTrollboxMessage(args []interface{}) (*TrollboxMessage, error) {

	var tbMsg = TrollboxMessage{}

	if v, ok := args[0].(string); ok {
		tbMsg.TypeMessage = v
	} else {
		return nil, fmt.Errorf("'TypeMessage' type assertion failed")
	}

	if v, ok := args[1].(float64); ok {
		tbMsg.MessageNumber = int64(v)
	} else {
		return nil, fmt.Errorf("'MessageNumber' type assertion failed")
	}

	if v, ok := args[2].(string); ok {
		tbMsg.Username = v
	} else {
		return nil, fmt.Errorf("'Username' type assertion failed")
	}

	if v, ok := args[3].(string); ok {
		tbMsg.Message = v
	} else {
		return nil, fmt.Errorf("'Message' type assertion failed")
	}

	if len(args) == 5 {
		if v, ok := args[4].(float64); ok {
			tbMsg.Reputation = int(v)
		} else {
			return nil, fmt.Errorf("'Reputation' type assertion failed")
		}
	} else {
		tbMsg.Reputation = -1
	}

	return &tbMsg, nil
}
