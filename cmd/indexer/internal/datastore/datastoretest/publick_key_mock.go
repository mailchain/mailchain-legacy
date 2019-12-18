// Code generated by MockGen. DO NOT EDIT.
// Source: public_key.go

// Package datastoretest is a generated GoMock package.
package datastoretest

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	crypto "github.com/mailchain/mailchain/crypto"
	reflect "reflect"
)

// MockPublicKeyStore is a mock of PublicKeyStore interface
type MockPublicKeyStore struct {
	ctrl     *gomock.Controller
	recorder *MockPublicKeyStoreMockRecorder
}

// MockPublicKeyStoreMockRecorder is the mock recorder for MockPublicKeyStore
type MockPublicKeyStoreMockRecorder struct {
	mock *MockPublicKeyStore
}

// NewMockPublicKeyStore creates a new mock instance
func NewMockPublicKeyStore(ctrl *gomock.Controller) *MockPublicKeyStore {
	mock := &MockPublicKeyStore{ctrl: ctrl}
	mock.recorder = &MockPublicKeyStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPublicKeyStore) EXPECT() *MockPublicKeyStoreMockRecorder {
	return m.recorder
}

// PutPublicKey mocks base method
func (m *MockPublicKeyStore) PutPublicKey(ctx context.Context, protocol, network string, address []byte, pubKey crypto.PublicKey) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PutPublicKey", ctx, protocol, network, address, pubKey)
	ret0, _ := ret[0].(error)
	return ret0
}

// PutPublicKey indicates an expected call of PutPublicKey
func (mr *MockPublicKeyStoreMockRecorder) PutPublicKey(ctx, protocol, network, address, pubKey interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutPublicKey", reflect.TypeOf((*MockPublicKeyStore)(nil).PutPublicKey), ctx, protocol, network, address, pubKey)
}

// GetPublicKey mocks base method
func (m *MockPublicKeyStore) GetPublicKey(ctx context.Context, protocol, network string, address []byte) (crypto.PublicKey, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPublicKey", ctx, protocol, network, address)
	ret0, _ := ret[0].(crypto.PublicKey)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPublicKey indicates an expected call of GetPublicKey
func (mr *MockPublicKeyStoreMockRecorder) GetPublicKey(ctx, protocol, network, address interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPublicKey", reflect.TypeOf((*MockPublicKeyStore)(nil).GetPublicKey), ctx, protocol, network, address)
}
