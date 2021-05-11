package main

import (
	"fmt"
	"log"

	"github.com/inagornyi/botex/bot"
	"github.com/inagornyi/botex/exchange"
	"github.com/inagornyi/botex/exchange/kuna"
)

const (
	apiKey    = ""
	secretKey = ""
)

func main() {
	client := kuna.NewKuna(apiKey, secretKey)

	b, err := client.Balance(exchange.UAH)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n[UAH] balance: %f\n\n", b.Balance)

	bot := bot.NewBot(client)
	err = bot.Trade()
	if err != nil {
		log.Fatal(err)
	}

	b, err = client.Balance(exchange.UAH)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n[UAH] balance: %f\n", b.Balance)
}
