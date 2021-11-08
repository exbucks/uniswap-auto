package services

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/hirokimoto/uniswap-auto/utils"
)

func Price(eth utils.Crypto, tokens utils.Tokens) (price float64) {
	if eth.Data.Bundles != nil && tokens.Data.Tokens != nil {
		unit, _ := strconv.ParseFloat(eth.Data.Bundles[0].EthPrice, 32)
		amount, _ := strconv.ParseFloat(tokens.Data.Tokens[0].DerivedETH, 32)
		price = unit * amount
	}
	return price
}

func LastPrice(swaps utils.Swaps) (last float64) {
	item := swaps.Data.Swaps[0]
	if item.Amount0In == "0" {
		amountUSD, _ := strconv.ParseFloat(item.AmountUSD, 32)
		amountOut, _ := strconv.ParseFloat(item.Amount0Out, 32)
		last = amountUSD / amountOut
	} else {
		amountUSD, _ := strconv.ParseFloat(item.AmountUSD, 32)
		amountOut, _ := strconv.ParseFloat(item.Amount0In, 32)
		last = amountUSD / amountOut
	}
	return last
}

func MinAndMax(swaps utils.Swaps) (
	min float64,
	max float64,
	minTarget string,
	maxTarget string,
	minTime time.Time,
	maxTime time.Time,
) {
	min = 0
	max = 0
	var _min int64
	var _max int64
	for _, item := range swaps.Data.Swaps {
		if item.Amount0In == "0" {
			amountUSD, _ := strconv.ParseFloat(item.AmountUSD, 32)
			amountOut, _ := strconv.ParseFloat(item.Amount0Out, 32)
			price := amountUSD / amountOut
			if min == 0 || max == 0 {
				min = price
				max = price
			}
			if price < min {
				min = price
				_min, _ = strconv.ParseInt(item.Timestamp, 10, 64)
			}
			if price > max {
				max = price
				_max, _ = strconv.ParseInt(item.Timestamp, 10, 64)
			}
			minTarget = "BUY"
			maxTarget = "BUY"
		} else {
			amountUSD, _ := strconv.ParseFloat(item.AmountUSD, 32)
			amountOut, _ := strconv.ParseFloat(item.Amount0In, 32)
			price := amountUSD / amountOut
			if min == 0 || max == 0 {
				min = price
				max = price
			}
			if price < min {
				min = price
				_min, _ = strconv.ParseInt(item.Timestamp, 10, 64)
			}
			if price > max {
				max = price
				_max, _ = strconv.ParseInt(item.Timestamp, 10, 64)
			}
			minTarget = "SELL"
			maxTarget = "SELL"
		}
	}
	minTime = time.Unix(_min, 0)
	maxTime = time.Unix(_max, 0)
	return min, max, minTarget, maxTarget, minTime, maxTime
}

func PeriodOfSwaps(swaps utils.Swaps) (time.Time, time.Time, time.Duration) {
	first, _ := strconv.ParseInt(swaps.Data.Swaps[0].Timestamp, 10, 64)
	last, _ := strconv.ParseInt(swaps.Data.Swaps[len(swaps.Data.Swaps)-1].Timestamp, 10, 64)
	tf := time.Unix(first, 0)
	tl := time.Unix(last, 0)
	period := tf.Sub(tl)
	return tl, tf, period
}

func findToken(pings <-chan string, id string) {
	var swaps utils.Swaps
	msg := <-pings
	json.Unmarshal([]byte(msg), &swaps)
	min, max, minTarget, maxTarget, minTime, maxTime := MinAndMax(swaps)

	last := LastPrice(swaps)
	fmt.Println(last)

	ts, tl, period := PeriodOfSwaps(swaps)
	if (max-min)/last > 0.1 && period < time.Duration(60*time.Minute) {
		fmt.Println("$$$$$ This is a tradable token! $$$$$")
		fmt.Println("Token ID:", id)
		fmt.Println("Token 0: ", swaps.Data.Swaps[0].Pair.Token0.Name)
		fmt.Println("Token 1: ", swaps.Data.Swaps[0].Pair.Token1.Name)
		fmt.Println("Last price: ", last)
		fmt.Println("Min price: ", min, minTarget, minTime)
		fmt.Println("Max price: ", max, maxTarget, maxTime)
		fmt.Println("Timeframe of 100 swaps: ", period)
		fmt.Println("Start and End time of the above time frame: ", ts, tl)
	}
}

func TradableTokens(wg *sync.WaitGroup, pairs utils.Pairs) {
	defer wg.Done()

	for _, item := range pairs.Data.Pairs {
		c := make(chan string, 1)
		go utils.Post(c, "swaps", item.Id)
		findToken(c, item.Id)
	}
}
