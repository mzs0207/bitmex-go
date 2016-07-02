package bitmex

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

//Auth - authentication
func (ws *WS) Auth(key, secret string) chan struct{} {
	ws.key = key
	ws.secret = secret

	nonce := ws.Nonce()

	req := fmt.Sprintf("GET/realtime%d", nonce)
	signature := ws.sign(req)

	msg := fmt.Sprintf(
		`{"op": "authKey", "args": ["%s", %d, "%s"]}`,
		key, nonce, signature,
	)

	ch := make(chan struct{})
	ws.Lock()
	ws.chSucc["authKey"] = append(ws.chSucc["authKey"], ch)
	ws.Unlock()

	ws.send(msg)

	return ch
}

func (ws *WS) sign(payload string) string {
	sig := hmac.New(sha256.New, []byte(ws.secret))
	sig.Write([]byte(payload))
	return hex.EncodeToString(sig.Sum(nil))
}

//Nonce - gets next nonce
func (ws *WS) Nonce() int64 {
	ws.nonce++
	return ws.nonce
}

//SubOrder - subscribe to order events
func (ws *WS) SubOrder(ch chan WSOrder, contracts []Contracts) chan struct{} {
	ws.Lock()

	if _, ok := ws.chOrder[ch]; !ok {
		ws.chOrder[ch] = contracts
	} else {
		ws.chOrder[ch] = append(ws.chOrder[ch], contracts...)
	}

	ws.Unlock()

	return ws.subPrivate("order")
}

func (ws *WS) subPrivate(topic string) chan struct{} {

	ch := make(chan struct{})
	ws.Lock()
	ws.chSucc[topic] = append(ws.chSucc[topic], ch)
	ws.Unlock()

	ws.send(`{"op": "subscribe", "args": "` + topic + `"}`)

	return ch

}
