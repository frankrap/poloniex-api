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

var trollbox Trollbox

var trollboxMu sync.Mutex
var trollboxClosed bool

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
func (client *PushClient) SubscribeTrollbox() (Trollbox, error) {

	trollbox = make(Trollbox)

	handler := func(args []interface{}, kwargs map[string]interface{}) {

		if tbMsg, err := convertArgsToTrollboxMessage(args); err != nil {
			fmt.Printf("convertArgsToTrollboxMessage: %v\n", err)
		} else {

			trollboxMu.Lock()
			if !trollboxClosed {
				trollbox <- tbMsg
			}
			trollboxMu.Unlock()
		}
	}

	if err := client.wampClient.Subscribe(TROLLBOX, nil, handler); err != nil {
		return nil, fmt.Errorf("turnpike.Client.Subscribe: %v", err)
	}

	return trollbox, nil
}

func (client *PushClient) UnsubscribeTrollbox() error {

	if err := client.wampClient.Unsubscribe(TROLLBOX); err != nil {
		return fmt.Errorf("turnpike.Client.Unsuscribe: %v", err)
	}
	trollboxMu.Lock()
	trollboxClosed = true
	trollboxMu.Unlock()

	close(trollbox)

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
