package bitmex

import (
	"encoding/json"
	"time"
)

type WSTrade struct {
	Size            float64   `json:"size"`
	Price           float64   `json:"price"`
	ForeignNotional float64   `json:"foreignNotional"`
	GrossValue      float64   `json:"grossValue"`
	HomeNotional    float64   `json:"homeNotional"`
	Symbol          string    `json:"symbol"`
	TickDirection   string    `json:"tickDirection"`
	Side            string    `json:"side"`
	TradeMatchID    string    `json:"trdMatchID"`
	Timestamp       time.Time `json:"timestamp"`
}

type wsData struct {
	Table       string            `json:"table"`
	Action      string            `json:"action"`
	Keys        []string          `json:"keys"`
	Attributes  map[string]string `json:"attributes"`
	Types       map[string]string `json:"keys"`
	ForeignKeys map[string]string `json:"foreignKeys"`
	Data        json.RawMessage
}

type wsSuccess struct {
	Success   bool              `json:"success"`
	Subscribe string            `json:"subscribe"`
	Request   map[string]string `json:"request"`
}

type wsInfo struct {
	Info      string    `json:"info"`
	Version   string    `json:"version"`
	Time      time.Time `json:"timestamp"`
	Docs      string    `json:"docs"`
	Heartbeat bool      `json:"heartbeatEnabled"`
}

type wsError struct {
	Error string `json:"error"`
}

type Trade struct {
}
