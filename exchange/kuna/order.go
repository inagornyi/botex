package kuna

import "time"

type Order struct {
	ID              int64     `json:"id"`
	Side            string    `json:"side"`
	OrdType         string    `json:"order_type"`
	Price           float64   `json:"price,string"`
	AvgPrice        string    `json:"avg_price"`
	State           string    `json:"state"`
	Market          string    `json:"market"`
	CreatedAt       time.Time `json:"created_at"`
	Volume          float64   `json:"volume,string"`
	RemainingVolume float64   `json:"remaining_volume,string"`
	ExecutedVolume  float64   `json:"executed_volume,string"`
	TradesCount     int       `json:"trades_count"`
}
