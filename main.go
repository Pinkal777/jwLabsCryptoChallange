package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const RateURI string = "https://api.coinbase.com/v2/exchange-rates?currency=USD"

var (
	usd float64
)

type Response struct {
	Data struct {
		Currancy string            `json:"currency"`
		Rates    map[string]string `json:"rates"`
	} `json:"data"`
}

type InvestmentAllocation struct {
	BTC string `json:"BTC"`
	ETH string `json:"ETH"`
}

type API struct {
	Client *http.Client
	URL    string
}

// this function get latest Crypto Exchange rate from supplied URL
func (api *API) GetLatestExchangeRate() (float64, float64, error) {
	resp, err := api.Client.Get(api.URL)
	if err != nil {
		return -1, -1, fmt.Errorf("no response from server: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body) //response body is []byte
	if err != nil {
		return -1, -1, fmt.Errorf("can not read response body: %w", err)
	}

	var result Response
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte
		return -1, -1, fmt.Errorf("can not unmarshal json response: %w", err)
	}

	bitcoinValue, _ := strconv.ParseFloat(strings.TrimSpace(result.Data.Rates["BTC"]), 64)
	etheriumValue, _ := strconv.ParseFloat(strings.TrimSpace(result.Data.Rates["ETH"]), 64)

	return bitcoinValue, etheriumValue, nil
}

// This Function will provide crypto split 70%/30%(BTC/ETH) for Provided input in USD.
func GenerateCryptoPortfolio(usd float64) (string, error) {
	api := &API{
		Client: &http.Client{
			Timeout: 15 * time.Second,
		},
		URL: RateURI,
	}

	btcVal, ethVal, err := api.GetLatestExchangeRate()
	if err != nil {
		return "", fmt.Errorf("not able to fetch latest exchange rate from server: %w", err)
	}

	btcPercent := 0.7
	ethPercent := 0.3

	btcInvestmentFund := usd * btcPercent
	ethInvestmentFund := usd * ethPercent

	totalBTC := btcInvestmentFund * btcVal
	totalETH := ethInvestmentFund * ethVal

	allocation := &InvestmentAllocation{
		BTC: strconv.FormatFloat(totalBTC, 'f', -1, 64),
		ETH: strconv.FormatFloat(totalETH, 'f', -1, 64),
	}

	alloc, err := json.Marshal(allocation)
	if err != nil {
		return "", fmt.Errorf("allocation cannot marshall to json: %w", err)
	}
	return string(alloc), nil
}

func main() {
	flag.Float64Var(&usd, "USD", 0.0, "Amount in usd To Invest in crypto")
	flag.Parse()
	if usd < 0.0 {
		fmt.Println("please provide valid amount(in USD) you want to invest in crypto.")
		log.Fatal("i.e.  go run .\\main.go -USD 100.25")
	}

	investment, err := GenerateCryptoPortfolio(usd)
	if err != nil {
		log.Fatal("investment portfolio cannot generated.", err)
	}
	fmt.Printf("\nwith %.2f$ you can buy/invest following cryptos \n\n", usd)
	fmt.Println(string(investment))
}
