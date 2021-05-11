package kuna

type Tickers struct {
	At     uint64  `json:"at"`
	Ticker *Ticker `json:"ticker"`
}

type Ticker struct {
	Buy   float64     `json:"buy,string"`
	Sell  float64     `json:"sell,string"`
	Low   float64     `json:"low,string"`
	High  float64     `json:"high,string"`
	Last  float64     `json:"last,string"`
	Vol   float64     `json:"vol,string"`
	Price interface{} `json:"price"` // bug on the side of kuna
}
