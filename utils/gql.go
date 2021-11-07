package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func request(query map[string]string, target chan string) {
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

func ETHQuery(target chan string) {
	query := Query("bundles", "")
	request(query, target)
}

func XIQuery(target chan string) {
	query := Query("tokens", "0x295b42684f90c77da7ea46336001010f2791ec8c")
	request(query, target)
}

func TradesQuery(target chan string) {
	query := Query("swaps", "0x7a99822968410431edd1ee75dab78866e31caf39")
	request(query, target)
}
