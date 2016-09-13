package oms

import (
	"errors"

	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/quickfix/enum"
	"github.com/shopspring/decimal"
)

//Order is the order type
type Order struct {
	ID                 int                `json:"id"`
	SessionID          quickfix.SessionID `json:"-"`
	ClOrdID            string             `json:"clord_id"`
	Symbol             string             `json:"symbol"`
	QuantityDecimal    decimal.Decimal    `json:"-"`
	Quantity           string             `json:"quantity"`
	Account            string             `json:"account"`
	Session            string             `json:"session_id"`
	Side               string             `json:"side"`
	OrdType            string             `json:"ord_type"`
	PriceDecimal       decimal.Decimal    `json:"-"`
	Price              string             `json:"price"`
	StopPriceDecimal   decimal.Decimal    `json:"-"`
	StopPrice          string             `json:"stop_price"`
	Closed             string             `json:"closed"`
	Open               string             `json:"open"`
	AvgPx              string             `json:"avg_px"`
	SecurityType       string             `json:"security_type"`
	MaturityMonthYear  string             `json:"maturity_month_year"`
	MaturityDay        int                `json:"maturity_day"`
	PutOrCall          int                `json:"put_or_call"`
	StrikePrice        string             `json:"strike_price"`
	StrikePriceDecimal decimal.Decimal    `json:"-"`
}

//Init initialized computed fields on order from user input
func (order *Order) Init() error {
	if qty, err := decimal.NewFromString(order.Quantity); err == nil {
		order.QuantityDecimal = qty
	} else {
		return errors.New("Invalid Qty")
	}

	switch order.OrdType {
	case enum.OrdType_LIMIT, enum.OrdType_STOP_LIMIT:
		if price, err := decimal.NewFromString(order.Price); err == nil {
			order.PriceDecimal = price
		} else {
			return errors.New("Invalid Price")
		}
	}

	switch order.OrdType {
	case enum.OrdType_STOP, enum.OrdType_STOP_LIMIT:
		if stopPrice, err := decimal.NewFromString(order.StopPrice); err == nil {
			order.StopPriceDecimal = stopPrice
		} else {
			return errors.New("Invalid StopPrice")
		}
	}

	return nil
}