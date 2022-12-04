package network

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type binancePriceResponse struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price,string"`
}

func Play_HTTP_Request() {
	var res binancePriceResponse

	res = getTokenPrice("MATIC", "USDT")
	log.Printf("%+v", res)
	res = getTokenPrice("AVAX", "USDT")
	log.Printf("%+v", res)
}

func getTokenPrice(tokenA string, tokenB string) binancePriceResponse {
	resp, err := http.Get(fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%s%s", strings.ToUpper(tokenA), strings.ToUpper(tokenB)))
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Errorf("%+v", err))
	}

	if resp.StatusCode != 200 {
		log.Panicf("response status %s - message: %s", resp.Status, string(bytes))
	}

	var priceData binancePriceResponse
	err = json.Unmarshal(bytes, &priceData)
	if err != nil {
		log.Panic(err)
	}
	return priceData
}
