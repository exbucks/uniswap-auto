package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
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

func gql(query map[string]string, target chan string) {
	jsonQuery, _ := json.Marshal(query)
	request, err := http.NewRequest("POST", "https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v2", bytes.NewBuffer(jsonQuery))
	client := &http.Client{Timeout: time.Second * 50}
	response, err := client.Do(request)
	defer response.Body.Close()
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}
	data, _ := ioutil.ReadAll(response.Body)
	target <- string(data)
}

func calcPrice(eth Crypto, tokens Tokens) {
	if eth.Data.Bundles != nil && tokens.Data.Tokens != nil {
		unit, _ := strconv.ParseFloat(eth.Data.Bundles[0].EthPrice, 32)
		amount, _ := strconv.ParseFloat(tokens.Data.Tokens[0].DerivedETH, 32)
		fmt.Println(unit * amount)
	}
}

func main() {
	ethQuery := map[string]string{
		"query": `
	        query bundles {
	            bundles(where: { id: "1" }) {
	                ethPrice
	            }
	        }
	    `,
	}
	xiQuery := map[string]string{
		"query": `
	        query tokens {
	            tokens(where: { id: "0x295b42684f90c77da7ea46336001010f2791ec8c" }) {
	                derivedETH
	                totalLiquidity
	            }
	        }
	    `,
	}
	tradesQuery := map[string]string{
		"query": `
			query swaps {
				swaps(orderBy: timestamp, orderDirection: desc, where:
					{ pair: "0x7a99822968410431edd1ee75dab78866e31caf39" }
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
					}
			}
	    `,
	}

	c1 := make(chan string)
	c2 := make(chan string)
	c3 := make(chan string)

	go func() {
		for {
			gql(ethQuery, c1)
			time.Sleep(time.Second * 2)
		}
	}()

	go func() {
		for {
			gql(xiQuery, c2)
			time.Sleep(time.Second * 3)
		}
	}()

	go func() {
		for {
			gql(tradesQuery, c3)
			time.Sleep(time.Second * 3)
		}
	}()

	var eth Crypto
	var xi Tokens

	go func() {
		for {
			select {
			case msg1 := <-c1:
				json.Unmarshal([]byte(msg1), &eth)
				calcPrice(eth, xi)
			case msg2 := <-c2:
				json.Unmarshal([]byte(msg2), &xi)
				calcPrice(eth, xi)
			case msg3 := <-c3:
				fmt.Println(msg3)
			}
		}
	}()

	var input string
	fmt.Scanln(&input)
}
