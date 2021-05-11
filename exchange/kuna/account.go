package kuna

import "github.com/inagornyi/botex/exchange"

type Account struct {
	Currency exchange.Currency `json:"currency"`
	Balance  float64           `json:"balance,string"`
	Locked   string            `json:"locked"`
}
