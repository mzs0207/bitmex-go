package bitmex

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
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
		nonce:  time.Now().Unix(),
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
