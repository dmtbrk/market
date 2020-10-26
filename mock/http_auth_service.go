// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ortymid/market/http (interfaces: AuthService)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	user "github.com/ortymid/market/market/user"
	http "net/http"
	reflect "reflect"
)

// HTTPAuthService is a mock of AuthService interface
type HTTPAuthService struct {
	ctrl     *gomock.Controller
	recorder *HTTPAuthServiceMockRecorder
}

// HTTPAuthServiceMockRecorder is the mock recorder for HTTPAuthService
type HTTPAuthServiceMockRecorder struct {
	mock *HTTPAuthService
}

// NewHTTPAuthService creates a new mock instance
func NewHTTPAuthService(ctrl *gomock.Controller) *HTTPAuthService {
	mock := &HTTPAuthService{ctrl: ctrl}
	mock.recorder = &HTTPAuthServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *HTTPAuthService) EXPECT() *HTTPAuthServiceMockRecorder {
	return m.recorder
}

// Authorize mocks base method
func (m *HTTPAuthService) Authorize(arg0 context.Context, arg1 *http.Request) (*user.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Authorize", arg0, arg1)
	ret0, _ := ret[0].(*user.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Authorize indicates an expected call of Authorize
func (mr *HTTPAuthServiceMockRecorder) Authorize(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Authorize", reflect.TypeOf((*HTTPAuthService)(nil).Authorize), arg0, arg1)
}
