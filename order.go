package bitmex

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// Order types
const (
	Market                    = "Market"
	Limit                     = "Limit"
	Stop                      = "Stop"
	StopLimit                 = "StopLimit"
	MarketIfTouched           = "MarketIfTouched"
	LimitIfTouched            = "LimitIfTouched"
	MarketWithLeftOverAsLimit = "MarketWithLeftOverAsLimit"
	Pegged                    = "Pegged"
)

// TimeInForce types
const (
	Day               = "Day"
	GoodTillCancel    = "GoodTillCancel"
	ImmediateOrCancel = "ImmediateOrCancel"
	FillOrKill        = "FillOrKill"
)

// PegPriceType types
const (
	LastPeg         = "LastPeg"
	MidPricePeg     = "MidPricePeg"
	MarketPeg       = "MarketPeg"
	PrimaryPeg      = "PrimaryPeg"
	TrailingStopPeg = "TrailingStopPeg"
)

// Execution instructions
const (
	ParticipateDoNotInitiate = "ParticipateDoNotInitiate"
	AllOrNone                = "AllOrNone"
	MarkPrice                = "MarkPrice"
	IndexPrice               = "IndexPrice"
	LastPrice                = "LastPrice"
	Close                    = "Close"
	ReduceOnly               = "ReduceOnly"
	Fixed                    = "Fixed"
)

// Order type
type Order struct {
	Account               float64   `json:"account,omitempty"`
	AvgPx                 float64   `json:"avgPx,omitempty"`
	ClOrdID               string    `json:"clOrdID,omitempty"`
	ClOrdLinkID           string    `json:"clOrdLinkID,omitempty"`
	ContingencyType       string    `json:"contingencyType,omitempty"`
	CumQty                float64   `json:"cumQty,omitempty"`
	Currency              Contract  `json:"currency,omitempty"`
	DisplayQty            float64   `json:"displayQty,omitempty"`
	ExDestination         string    `json:"exDestination,omitempty"`
	ExecInst              string    `json:"execInst,omitempty"`
	LeavesQty             float64   `json:"leavesQty,omitempty"`
	MultiLegReportingType string    `json:"multiLegReportingType,omitempty"`
	OrderID               uuid.UUID `json:"orderID,omitempty"`
	OrderQty              float64   `json:"orderQty,omitempty"`
	OrdRejReason          string    `json:"ordRejReason,omitempty"`
	OrdStatus             string    `json:"ordStatus,omitempty"`
	OrdType               string    `json:"ordType,omitempty"`
	PegOffsetValue        float64   `json:"pegOffsetValue,omitempty"`
	PegPriceType          string    `json:"pegPriceType,omitempty"`
	Price                 float64   `json:"price,omitempty"`
	SettlCurrency         Contract  `json:"settlCurrency,omitempty"`
	Side                  string    `json:"side,omitempty"`
	SimpleCumQty          float64   `json:"simpleCumQty,omitempty"`
	SimpleLeavesQty       float64   `json:"simpleLeavesQty,omitempty"`
	SimpleOrderQty        float64   `json:"simpleOrderQty,omitempty"`
	StopPx                float64   `json:"stopPx,omitempty"`
	Symbol                Contract  `json:"symbol,omitempty"`
	Text                  string    `json:"text,omitempty"`
	TimeInForce           string    `json:"timeInForce,omitempty"`
	Timestamp             time.Time `json:"timestamp,omitempty"`
	TransactTime          time.Time `json:"transactTime,omitempty"`
	Triggered             string    `json:"triggered,omitempty"`
	WorkingIndicator      bool      `json:"workingIndicator,omitempty"`
}

//NewOrder constructor
func NewOrder(symbol Contract) *Order {
	return &Order{Symbol: symbol}
}

//NewOrderMarket order
func NewOrderMarket(symbol Contract, qty float64) *Order {
	return &Order{
		Symbol:   symbol,
		OrderQty: qty,
	}
}
