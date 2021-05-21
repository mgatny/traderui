package oms

import (
	"fmt"
	"sync"
)

type ClOrdIDGenerator interface {
	Next() string
}

type OrderManager struct {
	sync.RWMutex
	orderID     int
	executionID int
	clOrdID     ClOrdIDGenerator

	orders           map[int]*Order
	crosses          map[int]*Cross
	ordersByClOrdID  map[string]*Order
	crossesByClOrdID map[string]*Cross
	executions       map[int]*Execution
}

func NewOrderManager(idGen ClOrdIDGenerator) *OrderManager {
	return &OrderManager{
		ordersByClOrdID:  make(map[string]*Order),
		orders:           make(map[int]*Order),
		crossesByClOrdID: make(map[string]*Cross),
		crosses:          make(map[int]*Cross),
		executions:       make(map[int]*Execution),
		clOrdID:          idGen,
	}
}

func (om *OrderManager) GetAllOrders() []*Order {
	orders := make([]*Order, 0, len(om.ordersByClOrdID))
	for _, v := range om.ordersByClOrdID {
		orders = append(orders, v)
	}

	return orders
}

func (om *OrderManager) GetAllCrosses() []*Cross {
	crosses := make([]*Cross, 0, len(om.crossesByClOrdID))
	for _, v := range om.crossesByClOrdID {
		crosses = append(crosses, v)
	}

	return crosses
}

func (om *OrderManager) GetAllExecutions() []*Execution {
	executions := make([]*Execution, 0, len(om.executions))
	for _, v := range om.executions {
		executions = append(executions, v)
	}

	return executions
}

func (om *OrderManager) GetOrder(id int) (*Order, error) {
	var err error
	order, ok := om.orders[id]
	if !ok {
		err = fmt.Errorf("could not find order with id %v", id)
	}

	return order, err
}

func (om *OrderManager) GetCross(id int) (*Cross, error) {
	var err error
	cross, ok := om.crosses[id]
	if !ok {
		err = fmt.Errorf("could not find cross with id %v", id)
	}

	return cross, err
}

func (om *OrderManager) GetExecution(id int) (*Execution, error) {
	var err error
	exec, ok := om.executions[id]
	if !ok {
		err = fmt.Errorf("could not find execution with id %v", id)
	}

	return exec, err
}

func (om *OrderManager) GetByClOrdID(clOrdID string) (*Order, error) {
	var err error
	order, ok := om.ordersByClOrdID[clOrdID]
	if !ok {
		err = fmt.Errorf("could not find order with clordid %v", clOrdID)
	}

	return order, err
}

func (om *OrderManager) SaveOrder(order *Order) error {
	order.ID = om.nextOrderID()
	order.ClOrdID = om.clOrdID.Next()

	om.orders[order.ID] = order
	om.ordersByClOrdID[order.ClOrdID] = order

	return nil
}

func (om *OrderManager) SaveCross(cross *Cross) error {
	cross.ID = om.nextOrderID()
	cross.ClOrdID = om.clOrdID.Next()

	om.crosses[cross.ID] = cross
	om.crossesByClOrdID[cross.ClOrdID] = cross

	return nil
}

func (om *OrderManager) SaveExecution(exec *Execution) error {
	exec.ID = om.nextExecutionID()
	om.executions[exec.ID] = exec

	return nil
}

func (om *OrderManager) AssignNextOrderClOrdID(order *Order) string {
	clOrdID := om.clOrdID.Next()
	om.ordersByClOrdID[clOrdID] = order
	return clOrdID
}

func (om *OrderManager) AssignNextCrossClOrdID(cross *Cross) string {
	clOrdID := om.clOrdID.Next()
	om.crossesByClOrdID[clOrdID] = cross
	return clOrdID
}

func (om *OrderManager) nextOrderID() int {
	om.orderID++
	return om.orderID
}

func (om *OrderManager) nextExecutionID() int {
	om.executionID++
	return om.executionID
}
