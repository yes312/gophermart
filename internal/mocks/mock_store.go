// Code generated by MockGen. DO NOT EDIT.
// Source: gophermart/internal/database (interfaces: StoragerDB)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	db "gophermart/internal/database"
	models "gophermart/internal/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockStoragerDB is a mock of StoragerDB interface.
type MockStoragerDB struct {
	ctrl     *gomock.Controller
	recorder *MockStoragerDBMockRecorder
}

// MockStoragerDBMockRecorder is the mock recorder for MockStoragerDB.
type MockStoragerDBMockRecorder struct {
	mock *MockStoragerDB
}

// NewMockStoragerDB creates a new mock instance.
func NewMockStoragerDB(ctrl *gomock.Controller) *MockStoragerDB {
	mock := &MockStoragerDB{ctrl: ctrl}
	mock.recorder = &MockStoragerDBMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStoragerDB) EXPECT() *MockStoragerDBMockRecorder {
	return m.recorder
}

// AddOrder mocks base method.
func (m *MockStoragerDB) AddOrder(arg0 context.Context, arg1, arg2 string) db.DBOperation {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddOrder", arg0, arg1, arg2)
	ret0, _ := ret[0].(db.DBOperation)
	return ret0
}

// AddOrder indicates an expected call of AddOrder.
func (mr *MockStoragerDBMockRecorder) AddOrder(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddOrder", reflect.TypeOf((*MockStoragerDB)(nil).AddOrder), arg0, arg1, arg2)
}

// AddUser mocks base method.
func (m *MockStoragerDB) AddUser(arg0 context.Context, arg1, arg2 string) db.DBOperation {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUser", arg0, arg1, arg2)
	ret0, _ := ret[0].(db.DBOperation)
	return ret0
}

// AddUser indicates an expected call of AddUser.
func (mr *MockStoragerDBMockRecorder) AddUser(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUser", reflect.TypeOf((*MockStoragerDB)(nil).AddUser), arg0, arg1, arg2)
}

// Close mocks base method.
func (m *MockStoragerDB) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockStoragerDBMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockStoragerDB)(nil).Close))
}

// GetBalance mocks base method.
func (m *MockStoragerDB) GetBalance(arg0 context.Context, arg1 string) db.DBOperation {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBalance", arg0, arg1)
	ret0, _ := ret[0].(db.DBOperation)
	return ret0
}

// GetBalance indicates an expected call of GetBalance.
func (mr *MockStoragerDBMockRecorder) GetBalance(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBalance", reflect.TypeOf((*MockStoragerDB)(nil).GetBalance), arg0, arg1)
}

// GetNewProcessedOrders mocks base method.
func (m *MockStoragerDB) GetNewProcessedOrders(arg0 context.Context) db.DBOperation {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNewProcessedOrders", arg0)
	ret0, _ := ret[0].(db.DBOperation)
	return ret0
}

// GetNewProcessedOrders indicates an expected call of GetNewProcessedOrders.
func (mr *MockStoragerDBMockRecorder) GetNewProcessedOrders(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNewProcessedOrders", reflect.TypeOf((*MockStoragerDB)(nil).GetNewProcessedOrders), arg0)
}

// GetOrders mocks base method.
func (m *MockStoragerDB) GetOrders(arg0 context.Context, arg1 string) db.DBOperation {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrders", arg0, arg1)
	ret0, _ := ret[0].(db.DBOperation)
	return ret0
}

// GetOrders indicates an expected call of GetOrders.
func (mr *MockStoragerDBMockRecorder) GetOrders(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrders", reflect.TypeOf((*MockStoragerDB)(nil).GetOrders), arg0, arg1)
}

// GetUser mocks base method.
func (m *MockStoragerDB) GetUser(arg0 context.Context, arg1 string) db.DBOperation {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", arg0, arg1)
	ret0, _ := ret[0].(db.DBOperation)
	return ret0
}

// GetUser indicates an expected call of GetUser.
func (mr *MockStoragerDBMockRecorder) GetUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockStoragerDB)(nil).GetUser), arg0, arg1)
}

// GetWithdrawals mocks base method.
func (m *MockStoragerDB) GetWithdrawals(arg0 context.Context, arg1 string) db.DBOperation {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWithdrawals", arg0, arg1)
	ret0, _ := ret[0].(db.DBOperation)
	return ret0
}

// GetWithdrawals indicates an expected call of GetWithdrawals.
func (mr *MockStoragerDBMockRecorder) GetWithdrawals(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWithdrawals", reflect.TypeOf((*MockStoragerDB)(nil).GetWithdrawals), arg0, arg1)
}

// PutStatuses mocks base method.
func (m *MockStoragerDB) PutStatuses(arg0 context.Context, arg1 *[]models.OrderStatusNew) db.DBOperation {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PutStatuses", arg0, arg1)
	ret0, _ := ret[0].(db.DBOperation)
	return ret0
}

// PutStatuses indicates an expected call of PutStatuses.
func (mr *MockStoragerDBMockRecorder) PutStatuses(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutStatuses", reflect.TypeOf((*MockStoragerDB)(nil).PutStatuses), arg0, arg1)
}

// WithRetry mocks base method.
func (m *MockStoragerDB) WithRetry(arg0 context.Context, arg1 db.DBOperation) (interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithRetry", arg0, arg1)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WithRetry indicates an expected call of WithRetry.
func (mr *MockStoragerDBMockRecorder) WithRetry(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithRetry", reflect.TypeOf((*MockStoragerDB)(nil).WithRetry), arg0, arg1)
}

// WithdrawBalance mocks base method.
func (m *MockStoragerDB) WithdrawBalance(arg0 context.Context, arg1 string, arg2 models.OrderSum) db.DBOperation {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithdrawBalance", arg0, arg1, arg2)
	ret0, _ := ret[0].(db.DBOperation)
	return ret0
}

// WithdrawBalance indicates an expected call of WithdrawBalance.
func (mr *MockStoragerDBMockRecorder) WithdrawBalance(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithdrawBalance", reflect.TypeOf((*MockStoragerDB)(nil).WithdrawBalance), arg0, arg1, arg2)
}
