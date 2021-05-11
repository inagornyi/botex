package kuna

import "github.com/inagornyi/botex/exchange"

type Balance struct {
	Currency exchange.Currency `json:"currency"`
	Balance  float64           `json:"balance,string"`
}
