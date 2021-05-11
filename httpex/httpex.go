package httpex

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

const (
	POST    = "POST"
	GET     = "GET"
	HEAD    = "HEAD"
	PUT     = "PUT"
	DELETE  = "DELETE"
	PATCH   = "PATCH"
	OPTIONS = "OPTIONS"
)

type Client struct {
	httpClient  *http.Client
	httpTimeout time.Duration
	throttle    <-chan time.Time
	headers     http.Header
	values      url.Values
	debug       bool
}

// NewClient return a new exchange HTTP client
func NewClient() *Client {
	return &Client{
		httpTimeout: 30 * time.Second,
		throttle:    time.Tick(100 * time.Millisecond),
		httpClient:  &http.Client{},
	}
}

func (c *Client) Timeout(timeout time.Duration) *Client {
	c.httpTimeout = timeout
	return c
}

func (c *Client) Throttle(throttle time.Duration) *Client {
	c.throttle = time.Tick(throttle)
	return c
}

func (c *Client) Header(headers http.Header) *Client {
	c.headers = headers
	return c
}

func (c *Client) Value(values url.Values) *Client {
	c.values = values
	return c
}

func (c *Client) Debug(enable bool) *Client {
	c.debug = enable
	return c
}

// do prepare and process HTTP request to exchange API
func (c *Client) Do(method string, url string) ([]byte, error) {
	respCh := make(chan []byte)
	errCh := make(chan error)
	<-c.throttle
	go c.makeReq(method, url, c.headers, c.values, respCh, errCh)
	return <-respCh, <-errCh
}

func (c *Client) makeReq(method string, url string, headers http.Header, values url.Values, respCh chan<- []byte, errCh chan<- error) {
	var body []byte
	connectTimer := time.NewTimer(c.httpTimeout)
	req, err := http.NewRequest(method, url, strings.NewReader(values.Encode()))
	if err != nil {
		respCh <- body
		errCh <- err
		return
	}
	if headers != nil {
		req.Header = headers
	}
	resp, err := c.doTimeoutRequest(connectTimer, req)
	if err != nil {
		respCh <- body
		errCh <- err
		return
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		respCh <- nil
		errCh <- err
		return
	}
	respCh <- body
	errCh <- nil
	close(respCh)
	close(errCh)
}

// doTimeoutRequest do a HTTP request with timeout
func (c *Client) doTimeoutRequest(timer *time.Timer, req *http.Request) (*http.Response, error) {
	// Do the request in the background so we can check the timeout
	type result struct {
		resp *http.Response
		err  error
	}

	done := make(chan result, 1)
	go func() {
		if c.debug {
			c.dumpRequest(req)
		}
		resp, err := c.httpClient.Do(req)
		if c.debug {
			c.dumpResponse(resp)
		}
		done <- result{resp, err}
	}()

	// Wait for the read or the timeout
	select {
	case r := <-done:
		return r.resp, r.err
	case <-timer.C:
		return nil, errors.New("timeout on reading data from exchange API")
	}
}

func (c *Client) dumpRequest(r *http.Request) {
	if r == nil {
		log.Print("dumpReq ok: <nil>")
		return
	}

	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Print("dumpReq err:", err)
	} else {
		log.Print("dumpReq ok:", string(dump))
	}
}

func (c *Client) dumpResponse(r *http.Response) {
	if r == nil {
		log.Print("dumpResponse ok: <nil>")
		return
	}

	dump, err := httputil.DumpResponse(r, true)
	if err != nil {
		log.Print("dumpResponse err:", err)
	} else {
		log.Print("dumpResponse ok:", string(dump))
	}
}
