package kuna

import (
	"github.com/inagornyi/botex/exchange"
	"time"
)

type History struct {
	ID        int64         `json:"id"`
	Price     float64       `json:"price,string"`
	Volume    float64       `json:"volume,string"`
	Funds     string        `json:"funds"`
	Market    exchange.Pair `json:"market"`
	CreatedAt time.Time     `json:"created_at"`
	Side      string        `json:"side"`
	OrderID   int64         `json:"order_id"`
}
