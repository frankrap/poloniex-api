package tradingapi

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// Poloniex trading API implementation of buy and sell command.
//
// API Doc:
// Places a limit buy or sell order in a given market. Required POST parameters are "currencyPair",
// "rate", and "amount". If successful, the method will return the order number.
//
// You may optionally set "fillOrKill", "immediateOrCancel", "postOnly" to 1. A fill-or-kill
// order will either fill in its entirety or be completely aborted. An immediate-or-cancel
// order can be partially or completely filled, but any portion of the order that cannot be
// filled immediately will be canceled rather than left on the order book. A post-only order
// will only be placed if no portion of it fills immediately; this guarantees you will never
// pay the taker fee on any part of the order that fills.
//
// Sample output:
//
//  {
//    "orderNumber": "31226040",
//    "resultingTrades": [
//      {
//        "amount": "338.8732",
//        "date": "2014-10-18 23:03:21",
//        "rate": "0.00000173",
//        "total": "0.00058625",
//        "tradeID": "16164",
//        "type": "buy"
//      }
//    ], ...
//    amountUnfilled: "332.23"
//  }
type BuyOrSellOrder struct {
	OrderNumber     int64            `json:"orderNumber,string"`
	ResultingTrades []ResultingTrade `json:"resultingTrades"`
	AmountUnfilled  float64          `json:"amountUnfilled,string"` // Only for ImmediateOrCancel option
}

type ResultingTrade struct {
	Amount    float64 `json:"amount,string"`
	Date      int64   // Unix timestamp
	Rate      float64 `json:"rate,string"`
	Total     float64 `json:"total,string"`
	TradeId   int64   `json:"tradeID,string"`
	TypeOrder string  `json:"type"`
}

func (client *TradingClient) BuyFillOrKill(currencyPair string, rate, amount float64) (*BuyOrSellOrder, error) {
	return client.buyOrSell("buy", currencyPair, rate, amount, "fillOrKill")
}

func (client *TradingClient) BuyImmediateOrCancel(currencyPair string, rate, amount float64) (*BuyOrSellOrder, error) {
	return client.buyOrSell("buy", currencyPair, rate, amount, "immediateOrCancel")
}

func (client *TradingClient) BuyPostOnly(currencyPair string, rate, amount float64) (*BuyOrSellOrder, error) {
	return client.buyOrSell("buy", currencyPair, rate, amount, "postOnly")
}

func (client *TradingClient) Buy(currencyPair string, rate, amount float64) (*BuyOrSellOrder, error) {
	return client.buyOrSell("buy", currencyPair, rate, amount, "")
}

func (client *TradingClient) SellFillOrKill(currencyPair string, rate, amount float64) (*BuyOrSellOrder, error) {
	return client.buyOrSell("sell", currencyPair, rate, amount, "fillOrKill")
}

func (client *TradingClient) SellImmediateOrCancel(currencyPair string, rate, amount float64) (*BuyOrSellOrder, error) {
	return client.buyOrSell("sell", currencyPair, rate, amount, "immediateOrCancel")
}

func (client *TradingClient) SellPostOnly(currencyPair string, rate, amount float64) (*BuyOrSellOrder, error) {
	return client.buyOrSell("sell", currencyPair, rate, amount, "postOnly")
}

func (client *TradingClient) Sell(currencyPair string, rate, amount float64) (*BuyOrSellOrder, error) {
	return client.buyOrSell("sell", currencyPair, rate, amount, "")
}

func (client *TradingClient) buyOrSell(command, currencyPair string, rate, amount float64, option string) (*BuyOrSellOrder, error) {

	postParameters := url.Values{}
	postParameters.Add("command", command)
	postParameters.Add("currencyPair", currencyPair)
	postParameters.Add("rate", strconv.FormatFloat(rate, 'f', -1, 64))
	postParameters.Add("amount", strconv.FormatFloat(amount, 'f', -1, 64))

	if option != "" {
		postParameters.Add(option, "1")
	}

	resp, err := client.do(postParameters)
	if err != nil {
		return nil, fmt.Errorf("do: %v", err)
	}

	res := BuyOrSellOrder{}

	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, fmt.Errorf("json unmarshal: %v", err)
	}

	return &res, nil
}

func (r *ResultingTrade) UnmarshalJSON(data []byte) error {

	type alias ResultingTrade
	aux := struct {
		Date string `json:"Date"`
		*alias
	}{
		alias: (*alias)(r),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("unmarshal aux: %v", err)
	}

	if timestamp, err := time.Parse("2006-01-02 15:04:05", aux.Date); err != nil {
		return fmt.Errorf("timestamp conversion: %v", err)
	} else {
		r.Date = int64(timestamp.Unix())
	}

	return nil
}
