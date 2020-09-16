// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ortymid/market/market (interfaces: Interface)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	market "github.com/ortymid/market/market"
	reflect "reflect"
)

// MockMarket is a mock of Interface interface
type MockMarket struct {
	ctrl     *gomock.Controller
	recorder *MockMarketMockRecorder
}

// MockMarketMockRecorder is the mock recorder for MockMarket
type MockMarketMockRecorder struct {
	mock *MockMarket
}

// NewMockMarket creates a new mock instance
func NewMockMarket(ctrl *gomock.Controller) *MockMarket {
	mock := &MockMarket{ctrl: ctrl}
	mock.recorder = &MockMarketMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMarket) EXPECT() *MockMarketMockRecorder {
	return m.recorder
}

// AddProduct mocks base method
func (m *MockMarket) AddProduct(arg0 context.Context, arg1 market.AddProductRequest, arg2 string) (*market.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddProduct", arg0, arg1, arg2)
	ret0, _ := ret[0].(*market.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddProduct indicates an expected call of AddProduct
func (mr *MockMarketMockRecorder) AddProduct(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddProduct", reflect.TypeOf((*MockMarket)(nil).AddProduct), arg0, arg1, arg2)
}

// DeleteProduct mocks base method
func (m *MockMarket) DeleteProduct(arg0 int, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteProduct", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteProduct indicates an expected call of DeleteProduct
func (mr *MockMarketMockRecorder) DeleteProduct(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteProduct", reflect.TypeOf((*MockMarket)(nil).DeleteProduct), arg0, arg1)
}

// EditProduct mocks base method
func (m *MockMarket) EditProduct(arg0 context.Context, arg1 market.EditProductRequest, arg2 string) (*market.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EditProduct", arg0, arg1, arg2)
	ret0, _ := ret[0].(*market.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EditProduct indicates an expected call of EditProduct
func (mr *MockMarketMockRecorder) EditProduct(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditProduct", reflect.TypeOf((*MockMarket)(nil).EditProduct), arg0, arg1, arg2)
}

// Product mocks base method
func (m *MockMarket) Product(arg0 context.Context, arg1 int) (*market.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Product", arg0, arg1)
	ret0, _ := ret[0].(*market.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Product indicates an expected call of Product
func (mr *MockMarketMockRecorder) Product(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Product", reflect.TypeOf((*MockMarket)(nil).Product), arg0, arg1)
}

// Products mocks base method
func (m *MockMarket) Products(arg0 context.Context, arg1, arg2 int) ([]*market.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Products", arg0, arg1, arg2)
	ret0, _ := ret[0].([]*market.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Products indicates an expected call of Products
func (mr *MockMarketMockRecorder) Products(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Products", reflect.TypeOf((*MockMarket)(nil).Products), arg0, arg1, arg2)
}

// ReplaceProduct mocks base method
func (m *MockMarket) ReplaceProduct(arg0 *market.Product, arg1 string) (*market.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReplaceProduct", arg0, arg1)
	ret0, _ := ret[0].(*market.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReplaceProduct indicates an expected call of ReplaceProduct
func (mr *MockMarketMockRecorder) ReplaceProduct(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReplaceProduct", reflect.TypeOf((*MockMarket)(nil).ReplaceProduct), arg0, arg1)
}
