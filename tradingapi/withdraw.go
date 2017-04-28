package tradingapi

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

type Withdrawal struct {
	Response string `json:"response"`
}

// Poloniex trading API implementation of withdraw command.
//
// API Doc:
// Immediately places a withdrawal for a given currency, with no email confirmation. In order
// to use this method, the withdrawal privilege must be enabled for your API key. Required POST
// parameters are "currency", "amount", and "address". For XMR withdrawals, you may optionally
// specify "paymentId".
//
// Sample output:
//
//  {
//    "response": "Withdrew 2398 NXT."
//  }
func (client *Client) Withdraw(currency string, amount float64, address string) (*Withdrawal, error) {

	postParameters := url.Values{}
	postParameters.Add("command", "withdraw")
	postParameters.Add("currency", currency)
	postParameters.Add("amount", strconv.FormatFloat(amount, 'f', -1, 64))
	postParameters.Add("address", address)

	resp, err := client.do(postParameters)
	if err != nil {
		return nil, fmt.Errorf("TradingClient.do: %v", err)
	}

	fmt.Println(string(resp))
	res := Withdrawal{}

	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %v", err)
	}

	return &res, nil
}

// WithdrawWithPaymentId withdraw for currency with special id parameter (XMR, XRP ...)
func (client *Client) WithdrawWithPaymentId(currency string, amount float64, address, paymentId string) (*Withdrawal, error) {

	postParameters := url.Values{}
	postParameters.Add("command", "withdraw")
	postParameters.Add("currency", currency)
	postParameters.Add("amount", strconv.FormatFloat(amount, 'f', -1, 64))
	postParameters.Add("address", address)
	postParameters.Add("paymentId", paymentId)

	resp, err := client.do(postParameters)
	if err != nil {
		return nil, fmt.Errorf("TradingClient.do %v", err)
	}

	res := Withdrawal{}

	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %v", err)
	}

	return &res, nil
}
