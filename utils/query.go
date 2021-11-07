package utils

import (
	"fmt"
)

type Crypto struct {
	Data struct {
		Bundles []struct {
			EthPrice string `json:"ethPrice"`
		} `json:"bundles"`
	} `json:"data"`
}

type Tokens struct {
	Data struct {
		Tokens []struct {
			DerivedETH     string `json:"derivedETH"`
			TotalLiquidity string `json:"totalLiquidity"`
		} `json:"tokens"`
	} `json:"data"`
}

type Swaps struct {
	Data struct {
		Swaps []struct {
			Amount0In  string `json:"amount0In"`
			Amount0Out string `json:"amount0Out"`
			Amount1In  string `json:"amount1In"`
			Amount1Out string `json:"amount1Out"`
			AmountUSD  string `json:"amountUSD"`
			Id         string `json:"id"`
			Pair       struct {
				Token0 struct {
					Symbol string `json:"symbol"`
				} `json:"token0"`
				Token1 struct {
					Symbol string `json:"symbol"`
				} `json:"token1"`
			} `json:"pair"`
			Timestamp string `json:"timestamp"`
			To        string `json:"to"`
		}
	}
}

func Query(target string, id string) map[string]string {
	var query map[string]string
	switch target {
	case "bundles":
		query = map[string]string{
			"query": `
				query bundles {
					bundles(where: { id: "1" }) {
						ethPrice
					}
				}
			`,
		}
		break
	case "tokens":
		sub := fmt.Sprintf(`
			query tokens {
				tokens(where: { id: %s }) {
					derivedETH
					totalLiquidity
				}
			}
		`, id)
		query = map[string]string{"query": sub}
		break
	case "swaps":
		sub := fmt.Sprintf(`
			query swaps {
				swaps(orderBy: timestamp, orderDirection: desc, where:
					{ pair: %s }
				) {
					pair {
						token0 {
							symbol
						}
						token1 {
							symbol
						}
					}
					amount0In
					amount0Out
					amount1In
					amount1Out
					amountUSD
					to
					timestamp
					id
				}
			}
		`, id)
		query = map[string]string{"query": sub}
		break
	default:
	}
	return query
}
