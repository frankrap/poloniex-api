package tradingapi

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	poloniex "github.com/joemocquant/poloniex-api"
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
	OrderNumber     int64                     `json:"orderNumber,string"`
	ResultingTrades []poloniex.ResultingTrade `json:"resultingTrades"`
	AmountUnfilled  float64                   `json:"amountUnfilled,string"` // Only for ImmediateOrCancel option
}

func (client *Client) BuyFillOrKill(currencyPair string, rate, amount float64) (*BuyOrSellOrder, error) {
	return client.buyOrSell("buy", currencyPair, rate, amount, "fillOrKill")
}

func (client *Client) BuyImmediateOrCancel(currencyPair string, rate, amount float64) (*BuyOrSellOrder, error) {
	return client.buyOrSell("buy", currencyPair, rate, amount, "immediateOrCancel")
}

func (client *Client) BuyPostOnly(currencyPair string, rate, amount float64) (*BuyOrSellOrder, error) {
	return client.buyOrSell("buy", currencyPair, rate, amount, "postOnly")
}

func (client *Client) Buy(currencyPair string, rate, amount float64) (*BuyOrSellOrder, error) {
	return client.buyOrSell("buy", currencyPair, rate, amount, "")
}

func (client *Client) SellFillOrKill(currencyPair string, rate, amount float64) (*BuyOrSellOrder, error) {
	return client.buyOrSell("sell", currencyPair, rate, amount, "fillOrKill")
}

func (client *Client) SellImmediateOrCancel(currencyPair string, rate, amount float64) (*BuyOrSellOrder, error) {
	return client.buyOrSell("sell", currencyPair, rate, amount, "immediateOrCancel")
}

func (client *Client) SellPostOnly(currencyPair string, rate, amount float64) (*BuyOrSellOrder, error) {
	return client.buyOrSell("sell", currencyPair, rate, amount, "postOnly")
}

func (client *Client) Sell(currencyPair string, rate, amount float64) (*BuyOrSellOrder, error) {
	return client.buyOrSell("sell", currencyPair, rate, amount, "")
}

func (client *Client) buyOrSell(command, currencyPair string, rate, amount float64, option string) (*BuyOrSellOrder, error) {

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
		return nil, fmt.Errorf("TradingClient.do: %v", err)
	}

	res := BuyOrSellOrder{}

	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %v", err)
	}

	return &res, nil
}
