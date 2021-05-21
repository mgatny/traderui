package basic

import (
	"errors"
	"strconv"
	"time"

	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/field"
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/traderui/oms"
	"github.com/quickfixgo/traderui/secmaster"

	fix40nos "github.com/quickfixgo/fix40/newordersingle"
	fix41nos "github.com/quickfixgo/fix41/newordersingle"
	fix42nos "github.com/quickfixgo/fix42/newordersingle"
	fix43nos "github.com/quickfixgo/fix43/newordersingle"
	fix44nos "github.com/quickfixgo/fix44/newordersingle"
	fix50nos "github.com/quickfixgo/fix50/newordersingle"

	fix50cross "github.com/quickfixgo/fix50/newordercross"

	fix42cxl "github.com/quickfixgo/fix42/ordercancelrequest"
)

//FIXFactory builds vanilla fix messages, implements traderui.fixFactory
type FIXFactory struct{}

func (FIXFactory) NewOrderSingle(order oms.Order) (msg quickfix.Messagable, err error) {
	switch order.SessionID.BeginString {
	case quickfix.BeginStringFIX40:
		msg, err = nos40(order)
	case quickfix.BeginStringFIX41:
		msg, err = nos41(order)
	case quickfix.BeginStringFIX42:
		msg, err = nos42(order)
	case quickfix.BeginStringFIX43:
		msg, err = nos43(order)
	case quickfix.BeginStringFIX44:
		msg, err = nos44(order)
	case quickfix.BeginStringFIXT11:
		msg, err = nos50(order)
	default:
		err = errors.New("unhandled BeginString")
	}

	return
}

func (FIXFactory) NewOrderCross(cross oms.Cross) (msg quickfix.Messagable, err error) {
	switch cross.SessionID.BeginString {
	case quickfix.BeginStringFIXT11:
		msg, err = cross50(cross)
	default:
		err = errors.New("unhandled BeginString")
	}

	return
}

func (FIXFactory) OrderCancelRequest(order oms.Order, clOrdID string) (msg quickfix.Messagable, err error) {
	switch order.SessionID.BeginString {
	case quickfix.BeginStringFIX42:
		msg, err = cxl42(order, clOrdID)
	default:
		err = errors.New("uhandled BeginString")
	}

	return
}

func (FIXFactory) SecurityDefinitionRequest(req secmaster.SecurityDefinitionRequest) (msg quickfix.Messagable, err error) {
	err = errors.New("not implemented")
	return
}

func populateOrder(genMessage quickfix.Messagable, ord oms.Order) (quickfix.Messagable, error) {
	msg := genMessage.ToMessage()

	switch ord.OrdType {
	case enum.OrdType_LIMIT, enum.OrdType_STOP_LIMIT:
		msg.Body.Set(field.NewPrice(ord.PriceDecimal, 2))
	}

	switch ord.OrdType {
	case enum.OrdType_STOP, enum.OrdType_STOP_LIMIT:
		msg.Body.Set(field.NewStopPx(ord.StopPriceDecimal, 2))
	}

	return msg, nil
}

func populateCross(genMessage quickfix.Messagable, cross oms.Cross) (quickfix.Messagable, error) {
	msg := genMessage.ToMessage()

	switch cross.OrdType {
	case enum.OrdType_LIMIT, enum.OrdType_STOP_LIMIT:
		msg.Body.Set(field.NewPrice(cross.PriceDecimal, 2))
	}

	switch cross.OrdType {
	case enum.OrdType_STOP, enum.OrdType_STOP_LIMIT:
		msg.Body.Set(field.NewStopPx(cross.StopPriceDecimal, 2))
	}

	return msg, nil
}

func nos40(ord oms.Order) (quickfix.Messagable, error) {
	nos := fix40nos.New(
		field.NewClOrdID(ord.ClOrdID),
		field.NewHandlInst("1"),
		field.NewSymbol(ord.Symbol),
		field.NewSide(ord.Side),
		field.NewOrderQty(ord.QuantityDecimal, 0),
		field.NewOrdType(ord.OrdType),
	)

	return populateOrder(nos, ord)
}

func nos41(ord oms.Order) (quickfix.Messagable, error) {
	nos := fix41nos.New(
		field.NewClOrdID(ord.ClOrdID),
		field.NewHandlInst("1"),
		field.NewSymbol(ord.Symbol),
		field.NewSide(ord.Side),
		field.NewOrdType(ord.OrdType),
	)
	nos.Set(field.NewOrderQty(ord.QuantityDecimal, 0))

	return populateOrder(nos, ord)
}

func nos42(ord oms.Order) (quickfix.Messagable, error) {
	nos := fix42nos.New(
		field.NewClOrdID(ord.ClOrdID),
		field.NewHandlInst("1"),
		field.NewSymbol(ord.Symbol),
		field.NewSide(ord.Side),
		field.NewTransactTime(time.Now()),
		field.NewOrdType(ord.OrdType),
	)
	nos.Set(field.NewOrderQty(ord.QuantityDecimal, 0))

	return populateOrder(nos, ord)
}

func cxl42(ord oms.Order, clOrdID string) (quickfix.Messagable, error) {
	cxl := fix42cxl.New(
		field.NewOrigClOrdID(ord.ClOrdID),
		field.NewClOrdID(clOrdID),
		field.NewSymbol(ord.Symbol),
		field.NewSide(ord.Side),
		field.NewTransactTime(time.Now()),
	)

	return cxl, nil
}

func nos43(ord oms.Order) (quickfix.Messagable, error) {
	nos := fix43nos.New(
		field.NewClOrdID(ord.ClOrdID),
		field.NewHandlInst("1"),
		field.NewSide(ord.Side),
		field.NewTransactTime(time.Now()),
		field.NewOrdType(ord.OrdType),
	)
	nos.Set(field.NewSymbol(ord.Symbol))
	nos.Set(field.NewOrderQty(ord.QuantityDecimal, 0))

	return populateOrder(nos, ord)
}

func nos44(ord oms.Order) (quickfix.Messagable, error) {
	nos := fix44nos.New(
		field.NewClOrdID(ord.ClOrdID),
		field.NewSide(ord.Side),
		field.NewTransactTime(time.Now()),
		field.NewOrdType(ord.OrdType),
	)
	nos.Set(field.NewSymbol(ord.Symbol))
	nos.Set(field.NewHandlInst("1"))
	nos.Set(field.NewOrderQty(ord.QuantityDecimal, 0))

	return populateOrder(nos, ord)
}

func nos50(ord oms.Order) (quickfix.Messagable, error) {
	nos := fix50nos.New(
		field.NewClOrdID(ord.ClOrdID),
		field.NewSide(ord.Side),
		field.NewTransactTime(time.Now()),
		field.NewOrdType(ord.OrdType),
	)
	nos.Set(field.NewHandlInst("1"))
	nos.Set(field.NewOrderQty(ord.QuantityDecimal, 0))
	nos.Set(field.NewSecurityType(ord.SecurityType))

	if ord.Account != "" {
		nos.Set(field.NewAccount(ord.Account))
	}
	if ord.Symbol != "" {
		nos.Set(field.NewSymbol(ord.Symbol))
	}
	if ord.SecurityID != "" {
		nos.Set(field.NewSecurityID(ord.SecurityID))
		nos.Set(field.NewSecurityIDSource(ord.SecurityIDSource))
	}
	if ord.ExecInst != "" {
		nos.Set(field.NewExecInst(enum.ExecInst(ord.ExecInst)))
	}

	return populateOrder(nos, ord)
}

func cross50(cross oms.Cross) (quickfix.Messagable, error) {
	msg := fix50cross.New(
		field.NewCrossID(strconv.Itoa(cross.ID)),
		field.NewCrossType(cross.CrossType),
		field.NewCrossPrioritization(cross.CrossPrioritization),
		field.NewTransactTime(time.Now()),
		field.NewOrdType(cross.OrdType),
	)
	msg.Set(field.NewHandlInst("1"))
	msg.Set(field.NewOrderQty(cross.QuantityDecimal, 0))
	msg.Set(field.NewSecurityType(cross.SecurityType))

	if cross.Symbol != "" {
		msg.Set(field.NewSymbol(cross.Symbol))
	}
	if cross.SecurityID != "" {
		msg.Set(field.NewSecurityID(cross.SecurityID))
		msg.Set(field.NewSecurityIDSource(cross.SecurityIDSource))
	}
	if cross.ExecInst != "" {
		msg.Set(field.NewExecInst(enum.ExecInst(cross.ExecInst)))
	}

	noSides := fix50cross.NewNoSidesRepeatingGroup()
	buy := noSides.Add()
	buy.Set(field.NewSide(enum.Side_BUY))
	buy.Set(field.NewDesignation("BuyDesignation"))
	buy.Set(field.NewClOrdID("BuyClOrdID"))
	buy.Set(field.NewAccount("BuyAccount"))
	sell := noSides.Add()
	sell.Set(field.NewSide(enum.Side_SELL))
	sell.Set(field.NewDesignation("SellDesignation"))
	sell.Set(field.NewClOrdID("SellClOrdID"))
	sell.Set(field.NewAccount("SellAccount"))
	msg.SetNoSides(noSides)

	return populateCross(msg, cross)
}
