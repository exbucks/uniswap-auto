package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
