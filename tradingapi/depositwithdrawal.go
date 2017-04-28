package tradingapi

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

type DepositsWithdrawals struct {
	Deposits    []*DepositHistory
	Withdrawals []*WithdrawalHistory
}

type DepositHistory struct {
	Currency      string  `json:"currency"`
	Address       string  `json:"address"`
	Amount        float64 `json:"amount,string"`
	Confirmations int     `json:"confirmations"`
	TxId          string  `json:"txid"`
	Timestamp     int64   `json:"timestamp"`
	Status        string  `json:"status"`
}

type WithdrawalHistory struct {
	WithdrawalNumber int64   `json:"withdrawalNumber"`
	Currency         string  `json:"currency"`
	Address          string  `json:"address"`
	Amount           float64 `json:"amount,string"`
	Timestamp        int64   `json:"timestamp"`
	Status           string  `json:"status"`
	IpAddress        string  `json:"ipAddress"`
}

// Poloniex trading API implementation of returnDepositsWithdrawals command.
//
// API Doc:
// Returns your deposit and withdrawal history within a range, specified by the "start"
// and "end" POST parameters, both of which should be given as UNIX timestamps.
//
// Withdrawals status: "COMPLETE" || "COMPLETE: ERROR" || "COMPLETE: txid" || "PENDING"
//
// Sample output:
//  {
//    "deposits": [
//      {
//        "currency": "BTC",
//        "address":"...",
//        "amount": "0.01006132",
//        "confirmations":10,
//        "txid":"17f819a91369a9ff6c4a34216d434597cfc1b4a3d0489b46bd6f924137a47701",
//        "timestamp":1399305798,
//        "status":"COMPLETE"
//      }, ...
//    ],
//    "withdrawals": [
//      {
//        "withdrawalNumber": 134933,
//        "currency": "BTC",
//        "address": "1N2i5n8DwTGzUq2Vmn9TUL8J1vdr1XBDFg",
//        "amount": "5.00010000",
//        "timestamp": 1399267904,
//        "status": "COMPLETE: 36e483efa6aff9fd53a235177579d98451c4eb237c210e66cd2b9a2d4a988f8e",
//        "ipAddress": "100.100.100.100"
//      }, ...
//    ]
//  }
func (client *Client) GetDepositsWithdrawals(start, end time.Time) (*DepositsWithdrawals, error) {

	postParameters := url.Values{}
	postParameters.Add("command", "returnDepositsWithdrawals")
	postParameters.Add("start", strconv.Itoa(int(start.Unix())))
	postParameters.Add("end", strconv.Itoa(int(end.Unix())))

	resp, err := client.do(postParameters)
	if err != nil {
		return nil, fmt.Errorf("TradingClient.do: %v", err)
	}

	res := DepositsWithdrawals{}

	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %v", err)
	}

	return &res, nil
}
