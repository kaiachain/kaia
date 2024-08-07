// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/kaiachain/kaia/storage/statedb (interfaces: TrieNodeCache)

// Package mock_statedb is a generated GoMock package.
package mock_statedb

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockTrieNodeCache is a mock of TrieNodeCache interface
type MockTrieNodeCache struct {
	ctrl     *gomock.Controller
	recorder *MockTrieNodeCacheMockRecorder
}

// MockTrieNodeCacheMockRecorder is the mock recorder for MockTrieNodeCache
type MockTrieNodeCacheMockRecorder struct {
	mock *MockTrieNodeCache
}

// NewMockTrieNodeCache creates a new mock instance
func NewMockTrieNodeCache(ctrl *gomock.Controller) *MockTrieNodeCache {
	mock := &MockTrieNodeCache{ctrl: ctrl}
	mock.recorder = &MockTrieNodeCacheMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTrieNodeCache) EXPECT() *MockTrieNodeCacheMockRecorder {
	return m.recorder
}

// Close mocks base method
func (m *MockTrieNodeCache) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockTrieNodeCacheMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockTrieNodeCache)(nil).Close))
}

// Get mocks base method
func (m *MockTrieNodeCache) Get(arg0 []byte) []byte {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0)
	ret0, _ := ret[0].([]byte)
	return ret0
}

// Get indicates an expected call of Get
func (mr *MockTrieNodeCacheMockRecorder) Get(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockTrieNodeCache)(nil).Get), arg0)
}

// Has mocks base method
func (m *MockTrieNodeCache) Has(arg0 []byte) ([]byte, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Has", arg0)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// Has indicates an expected call of Has
func (mr *MockTrieNodeCacheMockRecorder) Has(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Has", reflect.TypeOf((*MockTrieNodeCache)(nil).Has), arg0)
}

// SaveToFile mocks base method
func (m *MockTrieNodeCache) SaveToFile(arg0 string, arg1 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveToFile", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveToFile indicates an expected call of SaveToFile
func (mr *MockTrieNodeCacheMockRecorder) SaveToFile(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveToFile", reflect.TypeOf((*MockTrieNodeCache)(nil).SaveToFile), arg0, arg1)
}

// Set mocks base method
func (m *MockTrieNodeCache) Set(arg0, arg1 []byte) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Set", arg0, arg1)
}

// Set indicates an expected call of Set
func (mr *MockTrieNodeCacheMockRecorder) Set(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockTrieNodeCache)(nil).Set), arg0, arg1)
}

// UpdateStats mocks base method
func (m *MockTrieNodeCache) UpdateStats() interface{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateStats")
	ret0, _ := ret[0].(interface{})
	return ret0
}

// UpdateStats indicates an expected call of UpdateStats
func (mr *MockTrieNodeCacheMockRecorder) UpdateStats() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateStats", reflect.TypeOf((*MockTrieNodeCache)(nil).UpdateStats))
}
