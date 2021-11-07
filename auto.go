package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hirokimoto/uniswap-auto/services"
	"github.com/hirokimoto/uniswap-auto/utils"
)

func ethQuery(target chan string) {
	for {
		utils.Post(target, "bundles", "")
		time.Sleep(time.Second * 2)
	}
}

func xiQuery(target chan string) {
	for {
		utils.Post(target, "tokens", "0x295b42684f90c77da7ea46336001010f2791ec8c")
		time.Sleep(time.Second * 3)
	}
}

func tradesQuery(target chan string) {
	for {
		utils.Post(target, "swaps", "0x7a99822968410431edd1ee75dab78866e31caf39")
		time.Sleep(time.Second * 3)
	}
}

func main() {
	c1 := make(chan string)
	c2 := make(chan string)
	c3 := make(chan string)

	go ethQuery(c1)
	go xiQuery(c2)
	go tradesQuery(c3)

	var eth utils.Crypto
	var xi utils.Tokens
	var swaps utils.Swaps

	go func() {
		for {
			select {
			case msg1 := <-c1:
				json.Unmarshal([]byte(msg1), &eth)
				price := services.Price(eth, xi)
				fmt.Println("Current price: ", price)
			case msg2 := <-c2:
				json.Unmarshal([]byte(msg2), &xi)
				price := services.Price(eth, xi)
				fmt.Println("Current price: ", price)
			case msg3 := <-c3:
				json.Unmarshal([]byte(msg3), &swaps)
				last := services.LastPrice(swaps)
				min, max, minTarget, maxTarget, minTime, maxTime := services.MinAndMax(swaps)

				ts := time.Unix(minTime, 0)
				te := time.Unix(maxTime, 0)

				fmt.Println("Last price: ", last)
				fmt.Println("Min price: ", min, minTarget, ts)
				fmt.Println("Max price: ", max, maxTarget, te)
			}
		}
	}()

	var input string
	fmt.Scanln(&input)
}
