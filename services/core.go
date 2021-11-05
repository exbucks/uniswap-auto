package services

import (
	"fmt"
	"strconv"

	"github.com/hirokimoto/uniswap-auto/utils"
)

func Price(eth utils.Crypto, tokens utils.Tokens) {
	if eth.Data.Bundles != nil && tokens.Data.Tokens != nil {
		unit, _ := strconv.ParseFloat(eth.Data.Bundles[0].EthPrice, 32)
		amount, _ := strconv.ParseFloat(tokens.Data.Tokens[0].DerivedETH, 32)
		fmt.Println("Price: ", unit*amount)
	}
}

func MinAndMax(swaps utils.Swaps) (min float64, max float64) {
	for _, item := range swaps.Data.Swaps {
		amountUSD, _ := strconv.ParseFloat(item.AmountUSD, 32)
		// amountOut, _ := strconv.ParseFloat(item.Amount0Out, 32)
		// price := amountUSD / amountOut
		if amountUSD < min {
			min = amountUSD
		}
		if amountUSD > max {
			max = amountUSD
		}
	}
	return min, max
}
