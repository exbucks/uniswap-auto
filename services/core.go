package services

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/hirokimoto/uniswap-auto/utils"
)

func priceOfSwap(swap utils.Swap) (price float64, target string) {
	amountUSD, _ := strconv.ParseFloat(swap.AmountUSD, 32)
	amountToken, _ := strconv.ParseFloat(swap.Amount0Out, 32)
	price, _ = strconv.ParseFloat(swap.AmountUSD, 32)
	if swap.Amount0In == "0" && swap.Amount1Out == "0" {
		amountToken, _ = strconv.ParseFloat(swap.Amount0Out, 32)
		target = "BUY"
	} else if swap.Amount0Out == "0" && swap.Amount1In == "0" {
		amountToken, _ = strconv.ParseFloat(swap.Amount0In, 32)
		target = "SELL"
	} else if swap.Amount0In != "0" && swap.Amount0Out != "0" {
		amountToken, _ = strconv.ParseFloat(swap.Amount0Out, 32)
		target = "BUY"
	}
	price = amountUSD / amountToken
	return price, target
}

func Price(eth utils.Crypto, tokens utils.Tokens) (price float64) {
	if eth.Data.Bundles != nil && tokens.Data.Tokens != nil {
		unit, _ := strconv.ParseFloat(eth.Data.Bundles[0].EthPrice, 32)
		amount, _ := strconv.ParseFloat(tokens.Data.Tokens[0].DerivedETH, 32)
		price = unit * amount
	}
	return price
}

func LastPrice(swaps utils.Swaps) (price float64) {
	item := swaps.Data.Swaps[0]
	price, _ = priceOfSwap(item)
	return price
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
		price, target := priceOfSwap(item)
		minTarget = target
		maxTarget = target
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

	_, _, period := PeriodOfSwaps(swaps)
	if (max-min)/last > 0.1 && period < time.Duration(6*time.Hour) {
		fmt.Println("$$$$$ This is a tradable token! $$$$$")
		fmt.Println("Token ID:", id)
		fmt.Println("Token 0: ", swaps.Data.Swaps[0].Pair.Token0.Name)
		fmt.Println("Token 1: ", swaps.Data.Swaps[0].Pair.Token1.Name)
		fmt.Println("Last price: ", last)
		fmt.Println("Min price: ", min, minTarget, minTime)
		fmt.Println("Max price: ", max, maxTarget, maxTime)
		fmt.Println("Timeframe of 100 swaps: ", period)
		status := maxTime.After(minTime)
		if status {
			fmt.Println("Status: Safe")
		} else {
			fmt.Println("Status: Dagerous")
		}
	} else {
		fmt.Print(".")
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
