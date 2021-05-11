package bot

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/inagornyi/botex/exchange"
	"github.com/inagornyi/botex/exchange/kuna"
)

type Bot struct {
	purchaseStep     float64
	quantityStep     float64
	profitPercentage float64
	volume           float64
	fee              float64
	e                exchanger
}

type order struct {
	id      *int64
	price   float64
	volume  float64
	history []order
	bought  float64
}

type exchanger interface {
	Tickers() (map[exchange.Pair]kuna.Tickers, error)
	Buy(pair exchange.Pair, volume, price float64) (*kuna.Order, error)
	Sell(pair exchange.Pair, volume, price float64) (*kuna.Order, error)
	Orders(pair exchange.Pair) ([]kuna.Order, error)
	Cancel(id int64) (*kuna.Order, error)
	Balance(currency exchange.Currency) (*kuna.Balance, error)
	History(pair exchange.Pair) ([]kuna.History, error)
}

func NewBot(e exchanger) *Bot {
	return &Bot{
		purchaseStep:     0.002,
		quantityStep:     0.5,
		profitPercentage: 0.01,
		volume:           0.001,
		fee:              0.0025,
		e:                e,
	}
}

func (b *Bot) Trade() error {
	ts, err := b.e.Tickers()
	if err != nil {
		fmt.Printf("[ERROR] - error: %s\n", err)
		return err
	}

	bo := order{
		price:  ts[exchange.ETHUAH].Ticker.Sell,
		volume: b.volume,
	}

	so := order{}

	for {
		if bo.id == nil {
			o, err := b.e.Buy(exchange.ETHUAH, bo.volume, bo.price)
			if err != nil {
				fmt.Printf("[BUY][ERROR] - error: %s\n", err)
				continue
			}

			if o == nil {
				fmt.Printf("[BUY][ERROR] - error: unknown error\n")
				continue
			}

			bo.id = &o.ID
			bo.price = o.Price
			bo.volume = o.Volume

			fmt.Printf("[BUY][BOUGHT] - id: %d, price: %f, fee: %f, volume: %f\n",
				o.ID, o.Price, b.fee, o.Volume)
		}

		// get orders
		os, err := b.e.Orders(exchange.ETHUAH)
		if err != nil {
			fmt.Printf("[ORDERS][ERROR] - error: %s\n", err)
			continue
		}

		// check order
		co := func(o order, os []kuna.Order) error {
			for _, ho := range os {
				if *o.id == ho.ID {
					return nil
				}

			}
			return errors.New("order not exists")
		}

		// check order for sell
		if so.id != nil {
			err = co(so, os)
			if err != nil {
				break
			}
		}

		hos, err := b.e.History(exchange.ETHUAH)
		if err != nil {
			continue
		}

		chos := func(o order, hos []kuna.History) error {
			for _, ho := range hos {
				if *o.id == ho.OrderID {
					return nil
				}

			}
			return errors.New("order not exists")
		}

		// check history orders
		err = chos(bo, hos)
		if err != nil {
			continue
		}

		// check order for buy
		err = co(bo, os)
		if err == nil {
			continue
		}

		// append order
		ao := func(o order) {
			for _, ho := range o.history {
				if o.id == ho.id {
					return
				}
			}
			//fmt.Printf("[HISTOTY][DEBUG] json: %+v\n", o)
			bo.history = append(bo.history, order{
				id:     o.id,
				price:  o.price,
				volume: o.volume,
			})
		}
		ao(bo)

		if so.id != nil {
			o, err := b.e.Cancel(*so.id)
			if err != nil {
				continue
			}
			fmt.Printf("[CANCEL][SELL] - id: %d, obj: %+v\n", o.ID, o)
			so.bought = o.ExecutedVolume
			so.id = nil
		}

		if so.id == nil {
			var amount, price, volume float64
			var iterations int
			for _, ho := range bo.history {
				volume += ho.volume - (ho.volume * b.fee)
				amount += ho.price * ho.volume
				price += ho.price
				iterations++
			}

			amount = amount / float64(iterations)
			price = price / float64(iterations)
			price = amount / ((amount / price) - ((amount / price) * b.profitPercentage))

			volumeStr := fmt.Sprintf("%f", volume)

			ss := strings.Split(volumeStr, ".")
			var zeros int
			var numbers int
			for _, s := range ss[1] {
				if s == '\u0030' {
					zeros++
				} else {
					numbers++
				}
				if numbers == 2 {
					break
				}
			}
			volumeStr = ss[0] + "." + ss[1][0:(zeros+numbers)]
			volume, err := strconv.ParseFloat(volumeStr, 64)

			o, err := b.e.Sell(exchange.ETHUAH, volume, price)
			if err != nil {
				fmt.Printf("[SELL][ERROR] - error: %s\n", err)
				continue
			}

			so.id = &o.ID
			so.price = o.Price
			so.volume = o.Volume

			fmt.Printf("[SELL][SOLD] - id: %d, price: %f, volume: %f, obj: %+v\n", o.ID, o.Price, o.Volume, o)
		}

		bo.id = nil
		bo.price = bo.price - bo.price*b.purchaseStep
		bo.volume = bo.volume + bo.volume*b.quantityStep
	}

	os, err := b.e.Orders(exchange.ETHUAH)
	if err != nil {
		fmt.Printf("[ORDERS][ERROR] - error: %s\n", err)
		return err
	}
	for _, o := range os {
		_, err := b.e.Cancel(o.ID)
		if err != nil {
			fmt.Printf("[CANCEL][ERROR] - error: %s\n", err)
		}
	}

	balance, err := b.e.Balance(exchange.ETH)
	if err != nil {
		return err
	}

	ts, err = b.e.Tickers()
	if err != nil {
		return err
	}

	_, err = b.e.Sell(exchange.ETHUAH, balance.Balance, ts[exchange.ETHUAH].Ticker.Sell)
	if err != nil {
		return err
	}

	return nil
}
