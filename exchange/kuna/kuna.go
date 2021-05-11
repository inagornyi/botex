package kuna

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/inagornyi/botex/exchange"
	"github.com/inagornyi/botex/httpex"
)

const (
	apiUrl        = "https://kuna.io"
	apiVersion    = "api/v2"
	members       = "members/me"
	newOrder      = "orders"
	cancelOrder   = "order/delete"
	historyOrders = "trades/my"
	activeOrders  = "orders"
	tickers       = "tickers"
	timestamp     = "timestamp"
)

const (
	name      = "KUNA"
	accessKey = "access_key"
	tonce     = "tonce"
)

type Kuna struct {
	http      *httpex.Client
	apiKey    string
	secretKey string
}

func NewKuna(apiKey, secretKey string) *Kuna {
	return &Kuna{
		httpex.NewClient().Throttle(3 * time.Second).Timeout(1 * time.Minute),
		apiKey,
		secretKey,
	}
}

func (k Kuna) Name() string {
	return name
}

func (k Kuna) Me() (*Me, error) {
	var me Me
	if err := k.send(httpex.GET, &me, members); err != nil {
		return nil, err
	}
	return &me, nil
}

func (k Kuna) Buy(pair exchange.Pair, volume, price float64) (*Order, error) {
	return k.newOrder("buy", pair, volume, price)
}

func (k Kuna) Sell(pair exchange.Pair, volume, price float64) (*Order, error) {
	return k.newOrder("sell", pair, volume, price)
}

func (k Kuna) Cancel(id int64) (*Order, error) {
	var order Order
	if err := k.send(httpex.POST, &order, cancelOrder, "id", fmt.Sprintf("%d", id)); err != nil {
		return nil, err
	}
	return &order, nil
}

func (k Kuna) Balance(currency exchange.Currency) (*Balance, error) {
	bs, err := k.Balances()
	if err != nil {
		return nil, err
	}
	return bs[currency], nil
}

func (k Kuna) Balances() (Balances, error) {
	me, err := k.Me()
	if err != nil {
		return nil, err
	}

	bs := Balances{}
	for _, acc := range me.Accounts {
		bs[acc.Currency] = &Balance{
			Currency: acc.Currency,
			Balance:  acc.Balance,
		}
	}
	return bs, nil
}

func (k Kuna) History(pair exchange.Pair) ([]History, error) {
	var hs []History
	if err := k.send(httpex.GET, &hs, historyOrders, "market", string(pair)); err != nil {
		return nil, err
	}
	return hs, nil
}

func (k Kuna) Orders(pair exchange.Pair) ([]Order, error) {
	var os []Order
	if err := k.send(httpex.GET, &os, activeOrders, "market", string(pair)); err != nil {
		return nil, err
	}
	return os, nil
}

func (k Kuna) Tickers() (map[exchange.Pair]Tickers, error) {
	var ts map[exchange.Pair]Tickers
	if err := k.send(httpex.GET, &ts, tickers); err != nil {
		return nil, err
	}
	return ts, nil
}

func (k Kuna) newOrder(side string, pair exchange.Pair, volume, price float64) (*Order, error) {
	var order Order
	if err := k.send(httpex.POST, &order, newOrder,
		"market", string(pair), "price", fmt.Sprintf("%f", price), "side", side, "volume", fmt.Sprintf("%f", volume),
	); err != nil {
		return nil, err
	}
	return &order, nil
}

func (k Kuna) timestamp() (int64, error) {
	bs, err := k.http.Header(header()).Do(httpex.GET, fmt.Sprintf("%s/%s/%s", apiUrl, apiVersion, timestamp))
	if err != nil {
		return 0, err
	}
	var kunaError KunaError
	json.Unmarshal(bs, &kunaError)
	if kunaError.Error != nil {
		return 0, kunaError.Error
	}
	timestamp, err := strconv.ParseInt(string(bs), 0, 64)
	if err != nil {
		return 0, err
	}
	return timestamp, nil
}

func (k Kuna) send(method string, obj interface{}, endpoint string, args ...string) error {
	ts, err := k.timestamp()
	if err != nil {
		return err
	}
	value := value(args...)
	bs, err := k.http.
		Header(header()).
		Value(value).
		Do(genURL(method, fmt.Sprintf("/%s/%s", apiVersion, endpoint), k.apiKey, k.secretKey, ts, value))
	if err != nil {
		return err
	}

	var kunaError KunaError
	json.Unmarshal(bs, &kunaError)
	if kunaError.Error != nil {
		return kunaError.Error
	}

	if err = json.Unmarshal(bs, &obj); err != nil {
		return err
	}
	return nil
}

func header() http.Header {
	header := http.Header{}
	header.Add("Content-Type", "application/x-www-form-urlencoded")
	return header
}

func value(args ...string) url.Values {
	value := url.Values{}
	if len(args)%2 == 0 {
		for index := 0; index < len(args); index += 2 {
			value.Add(args[index], args[index+1])
		}
	}
	return value
}

// Generate URL for private API request
func genURL(method, url, apiKey, secretKey string, timestamp int64, args url.Values) (string, string) {
	args.Add(accessKey, apiKey)
	args.Add(tonce, fmt.Sprintf("%d000", timestamp))
	query := args.Encode()
	query = strings.Trim(query, "&")
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(fmt.Sprintf("%s|%s|%s", method, url, query)))
	secret := hex.EncodeToString(h.Sum(nil))
	url = fmt.Sprintf("%s?%s&signature=%s", url, query, secret)
	return method, fmt.Sprintf("%s%s", apiUrl, url)
}
