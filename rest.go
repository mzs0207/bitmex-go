package bitmex

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	uuid "github.com/satori/go.uuid"
)

const (
	endpoint   = "https://www.bitmex.com"
	apiVersion = "/api/v1"
)

// REST API object
type REST struct {
	client      *http.Client
	key, secret string
	nonce       int64
}

//NewREST REST Bitmex object
func NewREST() *REST {
	tr := &http.Transport{
		MaxIdleConns:    1,
		IdleConnTimeout: 60 * time.Second,
	}

	return &REST{
		client: &http.Client{Transport: tr},
		key:    os.Getenv("BITMEX_KEY"),
		secret: os.Getenv("BITMEX_SECRET"),
		nonce:  time.Now().UnixNano() / int64(time.Millisecond),
	}
}

// Auth func
func (r *REST) Auth(key, secret string) {
	r.key, r.secret = key, secret
}

//Send order func
func (r *REST) Send(order *Order) error {
	body, err := json.Marshal(order)
	req, err := r.request("POST", "/order", body)
	if err != nil {
		return err
	}

	_, err = r.client.Do(req)
	return err
}

//OrderSend 发送订单 .
func (r *REST) OrderSend(order *Order) (Order, error) {
	o := Order{}
	body, err := json.Marshal(order)
	req, err := r.request("POST", "/order", body)
	if err != nil {
		return o, err
	}
	resp, err := r.client.Do(req)
	if err != nil {
		return o, err
	}
	defer resp.Body.Close()

	respbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return o, err
	}
	fmt.Println(string(respbody))
	err = json.Unmarshal(respbody, &o)
	return o, nil
}

// Order 生成订单的基础方法.
func (r *REST) Order(symbol string, price float64, amount float64, side, orderType string, postOnly bool) (Order, error) {
	o := NewOrder(Contract(symbol))
	o.Price = price
	o.OrderQty = amount
	o.Side = side
	o.OrdType = orderType
	if postOnly {
		o.ExecInst = ParticipateDoNotInitiate
	}

	return r.OrderSend(o)

}

// LimitOrder 限价单.
func (r *REST) LimitOrder(symbol string, price float64, amount float64, side string, postOnly bool) (Order, error) {
	return r.Order(symbol, price, amount, side, Limit, postOnly)
}

// LimitBuyOrder 限价单买.
func (r *REST) LimitBuyOrder(symbol string, price float64, amount float64, postOnly bool) (Order, error) {
	return r.LimitOrder(symbol, price, amount, "Buy", postOnly)
}

// LimitSellOrder 限价单买.
func (r *REST) LimitSellOrder(symbol string, price float64, amount float64, postOnly bool) (Order, error) {
	return r.LimitOrder(symbol, price, amount, "Sell", postOnly)
}

// MarketOrder 市价单.
func (r *REST) MarketOrder(symbol string, price float64, amount float64, side string) (Order, error) {
	return r.Order(symbol, price, amount, side, Market, false)
}

// MarketBuyOrder 市价买单.
func (r *REST) MarketBuyOrder(symbol string, price float64, amount float64) (Order, error) {
	return r.Order(symbol, price, amount, "Buy", Market, false)
}

// MarketSellOrder 市价买单.
func (r *REST) MarketSellOrder(symbol string, price float64, amount float64) (Order, error) {
	return r.Order(symbol, price, amount, "Sell", Market, false)
}

// CancelOrder 取消订单.
func (r *REST) CancelOrder(orderID uuid.UUID) error {
	o := Order{}
	o.OrderID = orderID
	body, err := json.Marshal(o)
	req, err := r.request("DELETE", "/order", body)
	if err != nil {
		return err
	}
	resp, err := r.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(respbody))
	return nil
}

// ModifyOrder 修改订单.
func (r *REST) ModifyOrder(order Order) (Order, error) {
	o := Order{}
	body, err := json.Marshal(order)
	req, err := r.request("PUT", "/order", body)
	if err != nil {
		return o, err
	}
	resp, err := r.client.Do(req)
	if err != nil {
		return o, err
	}
	defer resp.Body.Close()

	respbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return o, err
	}
	fmt.Println(string(respbody))
	err = json.Unmarshal(respbody, &o)
	return o, nil
}

func (r *REST) getNonce() int64 {
	r.nonce++
	return r.nonce
}

func (r *REST) request(method, url string, body []byte) (*http.Request, error) {

	if method == "GET" {
		// TODO
		return nil, nil
	}
	req, err := http.NewRequest(
		method, endpoint+apiVersion+url, bytes.NewReader(body),
	)

	if err != nil {
		return nil, err
	}

	nonce := r.getNonce()
	sig := signature(r.secret, method, url, nonce, body)

	req.Header.Add("Content-Length", strconv.Itoa(len(body)))
	req.Header.Add("Content-Type", "application/json")

	req.Header.Add("api-nonce", strconv.FormatInt(nonce, 10))
	req.Header.Add("api-key", r.key)
	req.Header.Add("api-signature", sig)

	return req, nil
}

// hex(HMAC_SHA256(apiSecret, verb + path + nonce + data))
func signature(secret, verb, path string, nonce int64, data []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	var buf bytes.Buffer
	buf.WriteString(verb)
	buf.WriteString(apiVersion + path)
	fmt.Fprintf(&buf, "%d", nonce)
	buf.Write(data)
	mac.Write(buf.Bytes())
	return hex.EncodeToString(mac.Sum(nil))
}
