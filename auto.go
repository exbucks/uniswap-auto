package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hirokimoto/uniswap-auto/services"
	"github.com/hirokimoto/uniswap-auto/utils"
)

func main() {
	c1 := make(chan string)
	c2 := make(chan string)
	c3 := make(chan string)

	go func() {
		for {
			utils.ETHQuery(c1)
			time.Sleep(time.Second * 2)
		}
	}()

	go func() {
		for {
			utils.XIQuery(c2)
			time.Sleep(time.Second * 3)
		}
	}()

	go func() {
		for {
			utils.TradesQuery(c3)
			time.Sleep(time.Second * 3)
		}
	}()

	var eth utils.Crypto
	var xi utils.Tokens
	var swaps utils.Swaps

	go func() {
		for {
			select {
			case msg1 := <-c1:
				json.Unmarshal([]byte(msg1), &eth)
				services.Price(eth, xi)
			case msg2 := <-c2:
				json.Unmarshal([]byte(msg2), &xi)
				services.Price(eth, xi)
			case msg3 := <-c3:
				json.Unmarshal([]byte(msg3), &swaps)
				fmt.Println(msg3)
			}
		}
	}()

	var input string
	fmt.Scanln(&input)
}
