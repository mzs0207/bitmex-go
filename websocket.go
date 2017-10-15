package bitmex

import (
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/apex/log"

	"golang.org/x/net/websocket"
)

const wsURL = "wss://www.bitmex.com/realtime"

//WS - websocket connection object
type WS struct {
	sync.Mutex
	conn   *websocket.Conn
	log    *log.Logger
	nonce  int64
	key    string
	secret string
	chSucc map[string][]chan struct{}
	quit   chan struct{}

	// channels subscribed to different contracts

	chTrade    map[chan WSTrade][]Contract
	chQuote    map[chan WSQuote][]Contract
	chOrder    map[chan WSOrder][]Contract
	chPosition map[chan WSPosition][]Contract
}

//NewWS - creates new websocket object
func NewWS() *WS {
	return &WS{
		nonce:      time.Now().Unix(),
		quit:       make(chan struct{}),
		chTrade:    make(map[chan WSTrade][]Contract, 0),
		chQuote:    make(map[chan WSQuote][]Contract, 0),
		chOrder:    make(map[chan WSOrder][]Contract, 0),
		chPosition: make(map[chan WSPosition][]Contract, 0),
		chSucc:     make(map[string][]chan struct{}, 0),
	}
}

//Connect - connects
func (ws *WS) Connect() error {
	conn, err := websocket.Dial(wsURL, "", "http://localhost/")

	if err != nil {
		return err
	}

	log.Info("Connected")

	ws.conn = conn

	go ws.read()

	return nil
}

//Disconnect - Disconnects from websocket
func (ws *WS) Disconnect() {
	log.Info("Disconnecting")
	close(ws.quit)
	ws.conn.Close()
	//TODO Close all channels
	return
}

func (ws *WS) read() {
	for {
		// TODO []byte
		var msg string

		err := websocket.Message.Receive(ws.conn, &msg)
		if err != nil {
			select {
			case <-ws.quit:
				return
			default:
				log.Fatalf("WS error: %v", err)
			}
		}

		log.Debugf("Raw: %v", msg)

		switch {
		case strings.HasPrefix(msg, `{"success"`):
			var success wsSuccess
			json.Unmarshal([]byte(msg), &success)
			log.Debugf("Success: %#v", success)

			if channels, found := ws.chSucc[success.Request["op"]]; found {
				for _, ch := range channels {
					select {
					case ch <- struct{}{}:
					default:
					}
				}
			}

			if channels, found := ws.chSucc[success.Request["args"]]; found {
				for _, ch := range channels {
					select {
					case ch <- struct{}{}:
					default:
					}
				}
			}

		case strings.HasPrefix(msg, `{"info"`):
			var info wsInfo
			json.Unmarshal([]byte(msg), &info)
			log.Infof("Info: %v", info)

		case strings.Contains(msg, `{"table"`):
			var table wsData
			json.Unmarshal([]byte(msg), &table)
			log.Debugf("Table: %#v", table)

			switch table.Table {

			case "trade":
				var trades []WSTrade
				json.Unmarshal(table.Data, &trades)

				log.Debugf("Trades: %#v", trades)

				for _, one := range trades {
					ws.trade(one)
				}

			case "quote":
				var quotes []WSQuote
				json.Unmarshal(table.Data, &quotes)

				log.Debugf("Quotes: %#v", quotes)

				for _, one := range quotes {
					ws.quote(one)
				}

			case "order":
				var orders []WSOrder
				json.Unmarshal(table.Data, &orders)

				log.Debugf("Orders: %#v", orders)

				for _, one := range orders {
					ws.order(one)
				}

			case "position":
				var positions []WSPosition
				json.Unmarshal(table.Data, &positions)

				log.Debugf("Positions: %#v", positions)

				for _, one := range positions {
					ws.position(one)
				}
			}
		default:
			ws.fatal(errors.New("Unkown WS message"))
		}
	}
}

func (ws *WS) sendTrade(ch chan WSTrade, trade WSTrade) {
	select {
	case ch <- trade:
		log.Debugf("Trade sent: %#v - %#v", ch, trade)
	default:
		log.Debugf("Trade channel busy: %#v", ch)
	}
}

func (ws *WS) sendOrder(ch chan WSOrder, order WSOrder) {
	select {
	case ch <- order:
		log.Debugf("Order sent: %#v - %#v", ch, order)
	default:
		log.Debugf("Order channel busy: %#v", ch)
	}
}

func (ws *WS) sendQuote(ch chan WSQuote, quote WSQuote) {
	select {
	case ch <- quote:
		log.Debugf("Quote sent: %#v - %#v", ch, quote)
	default:
		log.Debugf("Quote channel busy: %#v", ch)
	}
}

func (ws *WS) sendPosition(ch chan WSPosition, position WSPosition) {
	select {
	case ch <- position:
		log.Debugf("Position sent: %#v - %#v", ch, position)
	default:
		log.Debugf("Position channel busy: %#v", ch)
	}
}

func (ws *WS) trade(trade WSTrade) {
	for ch, symbols := range ws.chTrade {
		// All
		if len(symbols) == 0 {
			ws.sendTrade(ch, trade)
			continue
		}

		// Filtered
		for _, oneSymbol := range symbols {
			if oneSymbol == Contract(trade.Symbol) {
				ws.sendTrade(ch, trade)
			}
		}
	}
}

func (ws *WS) order(order WSOrder) {
	for ch, symbols := range ws.chOrder {
		// All
		if len(symbols) == 0 {
			ws.sendOrder(ch, order)
			continue
		}

		// Filtered
		for _, oneSymbol := range symbols {
			if oneSymbol == Contract(order.Symbol) {
				ws.sendOrder(ch, order)
			}
		}
	}
}

func (ws *WS) position(position WSPosition) {
	for ch, symbols := range ws.chPosition {
		// All
		if len(symbols) == 0 {
			ws.sendPosition(ch, position)
			continue
		}

		// Filtered
		for _, oneSymbol := range symbols {
			if oneSymbol == Contract(position.Symbol) {
				ws.sendPosition(ch, position)
			}
		}
	}
}

func (ws *WS) quote(quote WSQuote) {
	for ch, symbols := range ws.chQuote {
		// All
		if len(symbols) == 0 {
			ws.sendQuote(ch, quote)
			continue
		}

		// Filtered
		for _, oneSymbol := range symbols {
			if oneSymbol == Contract(quote.Symbol) {
				ws.sendQuote(ch, quote)
			}
		}
	}
}

//Writing to WS
func (ws *WS) send(msg string) {
	defer ws.Unlock()

	log.Debugf("Writing WS: %#v", string(msg))
	ws.Lock()

	if _, err := ws.conn.Write([]byte(msg)); err != nil {
		ws.fatal(err)
	}
}

func (ws *WS) fatal(err error) {
	ws.Disconnect()
	log.Fatalf("%v", err.Error())
}

//SubTrade - subscribes channel to trades
func (ws *WS) SubTrade(ch chan WSTrade, contract []Contract) {
	ws.Lock()

	if _, ok := ws.chTrade[ch]; !ok {
		ws.chTrade[ch] = contract
	} else {
		ws.chTrade[ch] = append(ws.chTrade[ch], contract...)
	}

	ws.Unlock()

	for _, one := range contract {
		ws.send(`{"op": "subscribe", "args": "trade:` + string(one) + `"}`)
	}
}

//SubQuote - subscribes to quotes
func (ws *WS) SubQuote(ch chan WSQuote, contract []Contract) {

	ws.Lock()

	if _, ok := ws.chQuote[ch]; !ok {
		ws.chQuote[ch] = contract
	} else {
		ws.chQuote[ch] = append(ws.chQuote[ch], contract...)
	}

	ws.Unlock()

	for _, one := range contract {
		ws.send(`{"op": "subscribe", "args": "quote:` + string(one) + `"}`)
	}
}
