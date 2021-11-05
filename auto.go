package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func gql(query map[string]string, target chan string) {
	jsonQuery, _ := json.Marshal(query)
	request, err := http.NewRequest("POST", "https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v2", bytes.NewBuffer(jsonQuery))
	client := &http.Client{Timeout: time.Second * 10}
	response, err := client.Do(request)
	defer response.Body.Close()
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}
	data, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(data))
	target <- string(data)
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

	c1 := make(chan string)
	c2 := make(chan string)

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
			select {
			case msg1 := <-c1:
				fmt.Println(msg1)
			case msg2 := <-c2:
				fmt.Println(msg2)
			}
		}
	}()

	var input string
	fmt.Scanln(&input)
}
