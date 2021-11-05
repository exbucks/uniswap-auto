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
		fmt.Println(unit * amount)
	}
}
