// Code generated by MockGen. DO NOT EDIT.
// Source: pubkey.go

// Package mailboxtest is a generated GoMock package.
package mailboxtest

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockPubKeyFinder is a mock of PubKeyFinder interface
type MockPubKeyFinder struct {
	ctrl     *gomock.Controller
	recorder *MockPubKeyFinderMockRecorder
}

// MockPubKeyFinderMockRecorder is the mock recorder for MockPubKeyFinder
type MockPubKeyFinderMockRecorder struct {
	mock *MockPubKeyFinder
}

// NewMockPubKeyFinder creates a new mock instance
func NewMockPubKeyFinder(ctrl *gomock.Controller) *MockPubKeyFinder {
	mock := &MockPubKeyFinder{ctrl: ctrl}
	mock.recorder = &MockPubKeyFinderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPubKeyFinder) EXPECT() *MockPubKeyFinderMockRecorder {
	return m.recorder
}

// PublicKeyFromAddress mocks base method
func (m *MockPubKeyFinder) PublicKeyFromAddress(ctx context.Context, protocol, network string, address []byte) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PublicKeyFromAddress", ctx, protocol, network, address)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PublicKeyFromAddress indicates an expected call of PublicKeyFromAddress
func (mr *MockPubKeyFinderMockRecorder) PublicKeyFromAddress(ctx, protocol, network, address interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublicKeyFromAddress", reflect.TypeOf((*MockPubKeyFinder)(nil).PublicKeyFromAddress), ctx, protocol, network, address)
}
