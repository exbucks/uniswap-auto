package services

import (
	"strconv"
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
