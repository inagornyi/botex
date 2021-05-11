package kuna

import (
	"testing"

	"github.com/inagornyi/botex/exchange"
)

const (
	apiKey    = "NxwsrobtQcS4UbPwkG5EB1r2mOOyvzb0vCjdwj1J"
	secretKey = "omnyQQYuoyGr2OQDFUGyYqdraQAgJjO4OLj3Unr9"
)

func TestKuna(t *testing.T) {
	k := NewKuna(apiKey, secretKey)
	assert(t, k == nil, "error")
	t.Run("me", me(k))
	t.Run("buy", buy(k))
	t.Run("sell", sell(k))
	t.Run("balances", balances(k))
	t.Run("balance", balance(k))
	t.Run("history", history(k))
	t.Run("orders", orders(k))
	t.Run("ticker", ticker(k))
	t.Run("cancel", cancel(k))
	t.Run("timestamp", ts(k))
}

func me(k *Kuna) func(t *testing.T) {
	return func(t *testing.T) {
		me, err := k.Me()
		assert(t, err != nil, err)
		assert(t, me == nil, "me is nil")
		t.Logf("%#v", me)
	}
}

func buy(k *Kuna) func(t *testing.T) {
	return func(t *testing.T) {
		expected, err := k.Buy(exchange.XRPUAH, 1, 12)
		assert(t, err != nil, err)
		assert(t, expected == nil, "order is nil")
		t.Logf("%#v", expected)
	}
}

func sell(k *Kuna) func(t *testing.T) {
	return func(t *testing.T) {
		expected, err := k.Sell(exchange.XRPUAH, 1, 18)
		assert(t, err != nil, err)
		assert(t, expected == nil, "order is nil")
		t.Logf("%#v", expected)
	}
}

func cancel(k *Kuna) func(t *testing.T) {
	return func(t *testing.T) {
		order, err := k.Cancel(24078245)
		assert(t, err != nil, err)
		assert(t, order == nil, "order is nil")
		t.Logf("%#v", order)
	}
}

func balances(k *Kuna) func(t *testing.T) {
	return func(t *testing.T) {
		balances, err := k.Balances()
		assert(t, err != nil, err)
		assert(t, balances == nil, "balances are nil")
		t.Logf("%#v", balances)
	}
}

func balance(k *Kuna) func(t *testing.T) {
	return func(t *testing.T) {
		balance, err := k.Balance(exchange.XRP)
		assert(t, err != nil, err)
		assert(t, balance == nil, "balance is nil")
		t.Logf("%#v", balance)
	}
}

func history(k *Kuna) func(t *testing.T) {
	return func(t *testing.T) {
		hs, err := k.History(exchange.XRPUAH)
		assert(t, err != nil, err)
		assert(t, hs == nil, "history are nil")
		t.Logf("%#v", hs)
	}
}

func orders(k *Kuna) func(t *testing.T) {
	return func(t *testing.T) {
		os, err := k.Orders(exchange.XRPUAH)
		assert(t, err != nil, err)
		assert(t, os == nil, "orders are nil")
		t.Logf("%#v", os)
	}
}

func ticker(k *Kuna) func(t *testing.T) {
	return func(t *testing.T) {
		ts, err := k.Tickers()
		assert(t, err != nil, err)
		assert(t, ts == nil, "ticker is nil")
		t.Logf("%#v", ts[exchange.XRPUAH].Ticker)
	}
}

func ts(k *Kuna) func(t *testing.T) {
	return func(t *testing.T) {
		ts, err := k.timestamp()
		assert(t, err != nil, err)
		t.Logf("%d", ts)
	}
}

func assert(t *testing.T, b bool, arg interface{}) {
	if b {
		t.Fatal(arg)
	}
}
