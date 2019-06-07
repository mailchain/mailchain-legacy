// Copyright 2019 Finobo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package mailbox is a generated GoMock package.
package mailbox

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// mockSigner is a mock of Signer interface
type mockSigner struct {
	ctrl     *gomock.Controller
	recorder *mockSignerMockRecorder
}

// mockSignerMockRecorder is the mock recorder for mockSigner
type mockSignerMockRecorder struct {
	mock *mockSigner
}

// NewmockSigner creates a new mock instance
func newMockSigner(ctrl *gomock.Controller) *mockSigner {
	mock := &mockSigner{ctrl: ctrl}
	mock.recorder = &mockSignerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *mockSigner) EXPECT() *mockSignerMockRecorder {
	return m.recorder
}

// Sign mocks base method
func (m *mockSigner) Sign(opts SignerOpts) (interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Sign", opts)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Sign indicates an expected call of Sign
func (mr *mockSignerMockRecorder) Sign(opts interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sign", reflect.TypeOf((*mockSigner)(nil).Sign), opts)
}

// mockSignerOpts is a mock of SignerOpts interface
type mockSignerOpts struct {
	ctrl     *gomock.Controller
	recorder *mockSignerOptsMockRecorder
}

// mockSignerOptsMockRecorder is the mock recorder for MockSignerOpts
type mockSignerOptsMockRecorder struct {
	mock *mockSignerOpts
}

// newMockSignerOpts creates a new mock instance
func newMockSignerOpts(ctrl *gomock.Controller) *mockSignerOpts {
	mock := &mockSignerOpts{ctrl: ctrl}
	mock.recorder = &mockSignerOptsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *mockSignerOpts) EXPECT() *mockSignerOptsMockRecorder {
	return m.recorder
}
