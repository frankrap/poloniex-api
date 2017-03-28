package main

import (
	"fmt"
	"log"
	"poloniex"
	"poloniex/tradingapi"
	"time"
)

var client *tradingapi.TradingClient

func main() {

	var err error
	client, err = tradingapi.NewTradingClient(poloniex.API_KEY, poloniex.API_SECRET)

	if err != nil {
		log.Fatal(err)
	}

	// printBalances()
	// printCompleteBalances()
	// printDepositAddresses()
	// GenerateNewAddress()
	// printDepositsWithdrawals()
	// printOpenOrders()
	// printAllOpenOrders()
	// printTradeHistory()
	// printAllTradeHistory()
	// printTradesFromOrder()
	// buy()
	// buyFillOrKill()
	// buyImmediateOrCancel()
	// buyPostOnly()
	// sell()
	// sellFillOrKill()
	// sellImmediateOrCancel()
	// sellPostOnly()
	// cancelOrder()
	//TODO move order
	//TODO withdraw
	// printFeeInfo()
	// printAvailableAccountBalances()
	// printAccountBalances()
	printTradableBalances()
}

// Print balances
func printBalances() {

	res, err := client.GetBalances()

	if err != nil {
		log.Fatal(err)
	}

	poloniex.PrettyPrintJson(res)
}

// Print complete balances
func printCompleteBalances() {

	res, err := client.GetCompleteBalances()

	if err != nil {
		log.Fatal(err)
	}

	poloniex.PrettyPrintJson(res)
}

// Print deposit addresses
func printDepositAddresses() {

	res, err := client.GetDepositAddresses()

	if err != nil {
		log.Fatal(err)
	}

	poloniex.PrettyPrintJson(res)
}

// Generate new address for currency BTC
func GenerateNewAddress() {

	currency := "BTC"
	addr, err := client.GenerateNewAddress(currency)

	if err != nil {
		log.Fatal(err)
	}

	toPrint := fmt.Sprintf("New address generated (%s): %s", currency, addr)
	poloniex.PrettyPrintJson(toPrint)
}

// Print deposits and withdrawals that happened the last 20 days
func printDepositsWithdrawals() {

	end := time.Now()
	start := end.Add(-20 * 24 * time.Hour)
	res, err := client.GetDepositsWithdrawals(start, end)

	if err != nil {
		log.Fatal(err)
	}

	poloniex.PrettyPrintJson(res)
}

// Print open orders for BTC_STEEM market
func printOpenOrders() {

	res, err := client.GetOpenOrders("BTC_ETH")

	if err != nil {
		log.Fatal(err)
	}

	poloniex.PrettyPrintJson(res)
}

// Print open orders for all markets
func printAllOpenOrders() {

	res, err := client.GetAllOpenOrders()

	if err != nil {
		log.Fatal(err)
	}

	poloniex.PrettyPrintJson(res)
}

// Print BTC_STEEM trade history that happened the last 20 days
func printTradeHistory() {

	end := time.Now()
	start := end.Add(-20 * 24 * time.Hour)
	res, err := client.GetTradeHistory("BTC_STEEM", start, end)

	if err != nil {
		log.Fatal(err)
	}

	poloniex.PrettyPrintJson(res)
}

// Print trade history for all markets that happened the last 20 days
func printAllTradeHistory() {

	end := time.Now()
	start := end.Add(-20 * 24 * time.Hour)
	res, err := client.GetAllTradeHistory(start, end)

	if err != nil {
		log.Fatal(err)
	}

	poloniex.PrettyPrintJson(res)
}

// Print trade for a given orderId
func printTradesFromOrder() {

	var orderNumber int64 = 258117766006
	res, err := client.GetTradesFromOrder(orderNumber)

	if err != nil {
		log.Fatal(err)
	}

	poloniex.PrettyPrintJson(res)
}

// Place a buy order for 0.01 eth at 0.011btc
func buy() {

	rate, amount := 0.02, 0.01
	res, err := client.Buy("BTC_ETH", rate, amount)

	if err != nil {
		log.Fatal(err)
	}

	poloniex.PrettyPrintJson(res)
}

// Place a buy (fill or kill) order for 0.01 eth at 0.011btc
func buyFillOrKill() {

	rate, amount := 0.011, 0.01
	res, err := client.BuyFillOrKill("BTC_ETH", rate, amount)

	if err != nil {
		log.Fatal(err)
	}

	poloniex.PrettyPrintJson(res)
}

// Place a buy (immediate or cancel) order for 0.01 eth at 0.011btc
func buyImmediateOrCancel() {

	rate, amount := 0.011, 0.01
	res, err := client.BuyImmediateOrCancel("BTC_ETH", rate, amount)

	if err != nil {
		log.Fatal(err)
	}

	poloniex.PrettyPrintJson(res)
}

// Place a buy order (post only) for 0.01 eth at 0.011btc
func buyPostOnly() {

	rate, amount := 0.011, 0.01
	res, err := client.BuyPostOnly("BTC_ETH", rate, amount)

	if err != nil {
		log.Fatal(err)
	}

	poloniex.PrettyPrintJson(res)
}

// Place a sell order for 0.01 eth at 0.011btc
func sell() {

	rate, amount := 0.011, 0.01
	res, err := client.Sell("BTC_ETH", rate, amount)

	if err != nil {
		log.Fatal(err)
	}

	poloniex.PrettyPrintJson(res)
}

// Place a sell (fill or kill) order for 0.01 eth at 0.011btc
func sellFillOrKill() {

	rate, amount := 0.011, 0.01
	res, err := client.SellFillOrKill("BTC_ETH", rate, amount)

	if err != nil {
		log.Fatal(err)
	}

	poloniex.PrettyPrintJson(res)
}

// Place a sell (immediate or cancel) order for 0.01 eth at 0.011btc
func sellImmediateOrCancel() {

	rate, amount := 0.011, 0.01
	res, err := client.SellImmediateOrCancel("BTC_ETH", rate, amount)

	if err != nil {
		log.Fatal(err)
	}

	poloniex.PrettyPrintJson(res)
}

// Place a sell order (post only) for 0.01 eth at 0.011btc
func sellPostOnly() {

	rate, amount := 0.011, 0.01
	res, err := client.SellPostOnly("BTC_ETH", rate, amount)

	if err != nil {
		log.Fatal(err)
	}

	poloniex.PrettyPrintJson(res)
}

// Place a sell order (post only) for 0.01 eth at 0.011btc
func cancelOrder() {

	var orderNumber int64 = 258148121620
	res, err := client.CancelOrder(orderNumber)

	if err != nil {
		log.Fatal(err)
	}

	poloniex.PrettyPrintJson(res)
}

// Print available account balances
func printAvailableAccountBalances() {

	res, err := client.GetAvailableAccountBalances()

	if err != nil {
		log.Fatal(err)
	}

	poloniex.PrettyPrintJson(res)
}

// Print account balances
func printAccountBalances() {

	account := "exchange"
	res, err := client.GetAccountBalances(account)

	if err != nil {
		log.Fatal(err)
	}

	poloniex.PrettyPrintJson(res)
}

// Print fee info
func printFeeInfo() {

	res, err := client.GetFeeInfo()

	if err != nil {
		log.Fatal(err)
	}

	poloniex.PrettyPrintJson(res)
}

// Print tradable balances
func printTradableBalances() {

	res, err := client.GetTradableBalances()

	if err != nil {
		log.Fatal(err)
	}

	poloniex.PrettyPrintJson(res)
}
