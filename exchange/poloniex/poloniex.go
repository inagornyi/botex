package poloniex

import (
	"github.com/inagornyi/botex/httpex"
	"time"
)

/*
Public HTTP Endpoint: https://poloniex.com/public
Private HTTP Endpoint: https://poloniex.com/tradingApi
*/

const (
	apiUrl = "https://poloniex.com"
	pubEp  = "public"
	pvtEp  = "tradingApi"
)

type Poloniex struct {
	http      *httpex.Client
	apiKey    string
	secretKey string
}

func NewPoloniex(apiKey, secretKey string) *Poloniex {
	return &Poloniex{
		httpex.NewClient().Throttle(3 * time.Second),
		apiKey,
		secretKey,
	}
}

//func (p Poloniex) Balances() (Balances, error) {
//
//	return bs, nil
//}
