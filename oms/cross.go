package oms

import (
	"errors"

	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/quickfix"
	"github.com/shopspring/decimal"
)

//Cross is the cross type
type Cross struct {
	ID                  int                      `json:"id"`
	CrossType           enum.CrossType           `json:"cross_type"`
	CrossPrioritization enum.CrossPrioritization `json:"cross_prioritization"`
	SessionID           quickfix.SessionID       `json:"-"`
	ClOrdID             string                   `json:"clord_id"`
	OrderID             string                   `json:"order_id"`
	Symbol              string                   `json:"symbol"`
	QuantityDecimal     decimal.Decimal          `json:"-"`
	Quantity            string                   `json:"quantity"`
	Account             string                   `json:"account"`
	Session             string                   `json:"session_id"`
	Side                enum.Side                `json:"side"`
	OrdType             enum.OrdType             `json:"ord_type"`
	PriceDecimal        decimal.Decimal          `json:"-"`
	Price               string                   `json:"price"`
	StopPriceDecimal    decimal.Decimal          `json:"-"`
	StopPrice           string                   `json:"stop_price"`
	Closed              string                   `json:"closed"`
	Open                string                   `json:"open"`
	AvgPx               string                   `json:"avg_px"`
	SecurityType        enum.SecurityType        `json:"security_type"`
	SecurityDesc        string                   `json:"security_desc"`
	SecurityID          string                   `json:"security_id"`
	SecurityIDSource    enum.SecurityIDSource    `json:"security_id_source"`
	ExecInst            string                   `json:"exec_inst"`
	MaturityMonthYear   string                   `json:"maturity_month_year"`
	MaturityDay         int                      `json:"maturity_day"`
	PutOrCall           enum.PutOrCall           `json:"put_or_call"`
	StrikePrice         string                   `json:"strike_price"`
	StrikePriceDecimal  decimal.Decimal          `json:"-"`
	BuyClOrdID          string                   `json:"buy_clord_id"`
	BuyAccount          string                   `json:"buy_account"`
	BuyDesignation      string                   `json:"buy_designation"`
	SellClOrdID         string                   `json:"sell_clord_id"`
	SellAccount         string                   `json:"sell_account"`
	SellDesignation     string                   `json:"sell_designation"`
	SenderSubID         string                   `json:"sender_sub_id"`
}

//Init initialized computed fields on cross from user input
func (cross *Cross) Init() error {
	var err error
	if cross.QuantityDecimal, err = decimal.NewFromString(cross.Quantity); err != nil {
		return errors.New("Invalid Qty")
	}

	if cross.StrikePrice != "" {
		if cross.StrikePriceDecimal, err = decimal.NewFromString(cross.StrikePrice); err != nil {
			return errors.New("Invalid StrikePrice")
		}
	}

	switch cross.OrdType {
	case enum.OrdType_LIMIT, enum.OrdType_STOP_LIMIT:
		if cross.PriceDecimal, err = decimal.NewFromString(cross.Price); err != nil {
			return errors.New("Invalid Price")
		}
	}

	switch cross.OrdType {
	case enum.OrdType_STOP, enum.OrdType_STOP_LIMIT:
		if cross.StopPriceDecimal, err = decimal.NewFromString(cross.StopPrice); err != nil {
			return errors.New("Invalid StopPrice")
		}
	}

	return nil
}
