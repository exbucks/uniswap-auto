package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

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
	jsonETH, _ := json.Marshal(ethQuery)
	jsonXI, _ := json.Marshal(xiQuery)
	requestETH, err := http.NewRequest("POST", "https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v2", bytes.NewBuffer(jsonETH))
	requestXI, err := http.NewRequest("POST", "https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v2", bytes.NewBuffer(jsonXI))
	client := &http.Client{Timeout: time.Second * 10}
	responseETH, err := client.Do(requestETH)
	responseXI, err := client.Do(requestXI)
	defer responseETH.Body.Close()
	defer requestXI.Body.Close()
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}
	dataETH, _ := ioutil.ReadAll(responseETH.Body)
	dataXI, _ := ioutil.ReadAll(responseXI.Body)
	fmt.Println(string(dataETH))
	fmt.Println((string(dataXI)))
}
