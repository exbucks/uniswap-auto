package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/c-bata/go-prompt"

	"github.com/hirokimoto/uniswap-auto/services"
	"github.com/hirokimoto/uniswap-auto/utils"
)

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "1.", Description: " Track price of a token"},
		{Text: "2.", Description: " Track the token if it is tradable now"},
		{Text: "3.", Description: " Find a token which is a stable"},
		{Text: "4.", Description: " Find a tradable token"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func trackSwap(pings <-chan string) {
	msg := <-pings
	var swaps utils.Swaps
	json.Unmarshal([]byte(msg), &swaps)

	min, max, minTarget, maxTarget, minTime, maxTime := services.MinAndMax(swaps)
	fmt.Println("Min price: ", min, minTarget, minTime)
	fmt.Println("Max price: ", max, maxTarget, maxTime)

	last := services.LastPrice(swaps)
	fmt.Println("Last price: ", last)

	ts, tl, period := services.PeriodOfSwaps(swaps)
	fmt.Println("Timeframe of 100 swaps: ", period)
	fmt.Println("Start and End time of the above time frame: ", ts, tl)
	if (max-min)/last > 0.5 {
		fmt.Println("$$$$$ This is a tradable token! $$$$$")
	}
}

func trackPairs(pings <-chan string) {
	msg := <-pings
	var pairs utils.Pairs

	json.Unmarshal([]byte(msg), &pairs)

	var wg sync.WaitGroup
	wg.Add(len(pairs.Data.Pairs))
	go services.TradableTokens(&wg, pairs)
	wg.Wait()
}

func main() {
	fmt.Println("Please select what you are going to do.")
	t := prompt.Input("> ", completer)

	if t == "1." {
		c1 := make(chan string)
		c2 := make(chan string)

		var token string
		var eth utils.Crypto
		var xi utils.Tokens

		fmt.Print("Please enter your token: ")
		fmt.Scanf("%s", &token)
		fmt.Println("Token address: ", token)

		go func() {
			for {
				utils.Post(c1, "bundles", "")
				time.Sleep(time.Second * 2)
			}
		}()
		go func() {
			for {
				utils.Post(c2, "tokens", "0x295b42684f90c77da7ea46336001010f2791ec8c")
				time.Sleep(time.Second * 2)
			}
		}()

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

				}
			}
		}()
	} else if t == "2." {
		var pair string
		fmt.Print("Please enter your token: ")
		fmt.Scanf("%s", &pair)
		fmt.Println("Token address: ", pair)

		go func() {
			for {
				c3 := make(chan string)
				go utils.Post(c3, "swaps", "0x7a99822968410431edd1ee75dab78866e31caf39")
				trackSwap(c3)
			}
		}()
	} else if t == "3." {
		fmt.Println("hello, Yourself")
	} else if t == "4." {
		go func() {
			for {
				c4 := make(chan string)
				go utils.Post(c4, "pairs", "")
				trackPairs(c4)
			}
		}()
	}
	var input string
	fmt.Scanln(&input)
}
