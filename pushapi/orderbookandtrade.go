package pushapi

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type MarketUpdates struct {
	Sequence int64
	Updates  []*MarketUpdate
}

type MarketUpdate struct {
	Data       interface{}
	TypeUpdate string `json:"type"`
}

type OrderBookModify struct {
	Rate      float64 `json:"rate,string"`
	TypeOrder string  `json:"type"`
	Amount    float64 `json:"amount,string"`
}

type OrderBookRemove struct {
	Rate      float64 `json:"rate,string"`
	TypeOrder string  `json:"type"`
}

type NewTrade struct {
	TradeId   int64   `json:"tradeID,string"`
	Rate      float64 `json:"rate,string"`
	Amount    float64 `json:"amount,string"`
	Date      int64   // Unix timestamp
	Total     float64 `json:"total,string"`
	TypeOrder string  `json:"type"`
}

type MarketUpdater chan *MarketUpdates

type marketUpdaterInfo struct {
	updater MarketUpdater

	mu           sync.RWMutex
	unsubscribed chan struct{}
}

var (
	mutex              sync.RWMutex
	marketUpdaterInfos = make(map[string]*marketUpdaterInfo)
)

// Poloniex push API implementation of order book and trade topics.
//
// API Doc:
// To receive order book and trade updates, subscribe to the desired currencyPair,
// e.g. "BTC_XMR".
//
// There are two types of order book updates:
//
//  [
//    {
//      data: {
//        rate: '0.00300888',
//        type: 'bid',
//        amount: '3.32349029'
//      },
//      type: 'orderBookModify'
//    }
//  ]
//
//  [
//    {
//      data: {
//        rate: '0.00311164',
//        type: 'ask'
//      },
//      type: 'orderBookRemove'
//    }
//  ]
//
// Updates of type orderBookModify can be either additions to the order book
// or changes to existing entries. The value of 'amount' indicates the new total
// amount on the books at the given rate — in other words, it replaces any previous
// value, rather than indicates an adjustment to a previous value.
//
// Trade history updates are provided in the following format:
//
//  [
//    {
//      data: {
//        tradeID: '364476',
//        rate: '0.00300888',
//        amount: '0.03580906',
//        date: '2014-10-07 21:51:20',
//        total: '0.00010775',
//        type: 'sell'
//      },
//      type: 'newTrade'
//    }
//  ]
// The dictionary portion of each market message ("kwargs" in the Node.js example)
// will contain a sequence number with the key "seq". In order to keep your order
// book consistent, you will need to ensure that messages are applied in the order
// of their sequence numbers, even if they arrive out of order. In some markets, if
// there is no update for more than 1 second, a heartbeat message consisting of an
// empty argument list and the latest sequence number will be sent. These will go
// out once per second, but if there is no update for more than 60 seconds, the
// heartbeat interval will be reduced to 8 seconds until the next update.
//
// Several order book and trade history updates will often arrive in a single message.
// Be sure to loop through the entire array, otherwise you will miss some updates.
func (client *Client) SubscribeMarket(currencyPair string) (MarketUpdater, error) {

	handler := func(args []interface{}, kwargs map[string]interface{}) {

		client.updateTopicTimestamp(currencyPair)

		seq, ok := kwargs["seq"].(float64)
		if !ok {
			logger.WithField("error", "'seq' type assertion failed").Error(
				"convertArgsToMarketUpdateSlice")
			return
		}

		if len(args) == 0 {
			// Heartbeat
			// int64(seq)
			return
		}

		updates, err := convertArgsToMarketUpdateSlice(args)
		if err != nil {
			logger.WithField("error", err).Error("convertArgsToMarketUpdateSlice")
			return
		}

		mutex.RLock()
		updaterInfo := marketUpdaterInfos[currencyPair]
		mutex.RUnlock()

		updaterInfo.mu.RLock()
		select {
		case updaterInfo.updater <- &MarketUpdates{int64(seq), updates}:
		case <-updaterInfo.unsubscribed:
		}
		updaterInfo.mu.RUnlock()
	}

	subscribe := func() error {

		client.wampClientMu.RLock()
		defer client.wampClientMu.RUnlock()

		if err := client.wampClient.Subscribe(currencyPair, nil, handler); err != nil {
			return fmt.Errorf("turnpike.Client.Subscribe: %v", err)
		}
		logger.Infof("Subscribed to: %s", currencyPair)

		return nil
	}

	mutex.Lock()
	updaterInfo, ok := marketUpdaterInfos[currencyPair]
	if !ok {
		updaterInfo = &marketUpdaterInfo{
			updater:      make(MarketUpdater),
			unsubscribed: make(chan struct{}),
			mu:           sync.RWMutex{},
		}
		marketUpdaterInfos[currencyPair] = updaterInfo
	}
	mutex.Unlock()

	if err := subscribe(); err != nil {
		mutex.Lock()
		delete(marketUpdaterInfos, currencyPair)
		mutex.Unlock()
		return nil, err
	}

	updaterInfo.mu.Lock()
	select {
	case <-updaterInfo.unsubscribed:
		updaterInfo.unsubscribed = make(chan struct{})
	default:
	}
	updaterInfo.mu.Unlock()

	client.addSubscription(currencyPair, subscribe)

	return updaterInfo.updater, nil
}

func (client *Client) UnsubscribeMarket(currencyPair string) error {

	client.wampClientMu.RLock()

	if err := client.wampClient.Unsubscribe(currencyPair); err != nil {
		return fmt.Errorf("turnpike.Client.Unsuscribe: %v", err)
	}

	client.wampClientMu.RUnlock()

	client.removeSubscription(currencyPair)

	mutex.RLock()
	updaterInfo := marketUpdaterInfos[currencyPair]
	mutex.RUnlock()

	updaterInfo.updater <- nil

	updaterInfo.mu.RLock()
	close(updaterInfo.unsubscribed)
	updaterInfo.mu.RUnlock()

	return nil
}

func convertArgsToMarketUpdateSlice(args []interface{}) ([]*MarketUpdate, error) {

	res := make([]*MarketUpdate, len(args))

	for i, val := range args {

		strjson, err := json.Marshal(val)
		if err != nil {
			return nil, fmt.Errorf("json.Marshal: %v", err)
		}

		var dataField json.RawMessage

		marketUpdate := MarketUpdate{
			Data: &dataField,
		}

		if err := json.Unmarshal(strjson, &marketUpdate); err != nil {
			return nil, fmt.Errorf("json.Unmarshal: %v", err)
		}

		switch marketUpdate.TypeUpdate {
		case "orderBookModify":
			obm := OrderBookModify{}
			if err := json.Unmarshal(dataField, &obm); err != nil {
				return nil, fmt.Errorf("json.Unmarshal: %v", err)
			}
			marketUpdate.Data = &obm
		case "orderBookRemove":
			obr := OrderBookRemove{}
			if err := json.Unmarshal(dataField, &obr); err != nil {
				return nil, fmt.Errorf("json.Unmarshal: %v", err)
			}
			marketUpdate.Data = &obr
		case "newTrade":
			nt := NewTrade{}
			if err := json.Unmarshal(dataField, &nt); err != nil {
				return nil, fmt.Errorf("json.Unmarshal: %v", err)
			}
			marketUpdate.Data = &nt
		}

		res[i] = &marketUpdate
	}

	return res, nil
}

func (n *NewTrade) UnmarshalJSON(data []byte) error {

	type alias NewTrade
	aux := struct {
		Date string `json:"Date"`
		*alias
	}{
		alias: (*alias)(n),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("json.Unmarshal: %v", err)
	}

	if timestamp, err := time.Parse("2006-01-02 15:04:05", aux.Date); err != nil {
		return fmt.Errorf("time.Parse: %v", err)
	} else {
		n.Date = int64(timestamp.Unix())
	}

	return nil
}
