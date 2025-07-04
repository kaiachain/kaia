// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/kaiachain/kaia/kaiax/builder (interfaces: BuilderModule)

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	rpc "github.com/kaiachain/kaia/networks/rpc"
)

// MockBuilderModule is a mock of BuilderModule interface.
type MockBuilderModule struct {
	ctrl     *gomock.Controller
	recorder *MockBuilderModuleMockRecorder
}

// MockBuilderModuleMockRecorder is the mock recorder for MockBuilderModule.
type MockBuilderModuleMockRecorder struct {
	mock *MockBuilderModule
}

// NewMockBuilderModule creates a new mock instance.
func NewMockBuilderModule(ctrl *gomock.Controller) *MockBuilderModule {
	mock := &MockBuilderModule{ctrl: ctrl}
	mock.recorder = &MockBuilderModuleMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBuilderModule) EXPECT() *MockBuilderModuleMockRecorder {
	return m.recorder
}

// APIs mocks base method.
func (m *MockBuilderModule) APIs() []rpc.API {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "APIs")
	ret0, _ := ret[0].([]rpc.API)
	return ret0
}

// APIs indicates an expected call of APIs.
func (mr *MockBuilderModuleMockRecorder) APIs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "APIs", reflect.TypeOf((*MockBuilderModule)(nil).APIs))
}

// Start mocks base method.
func (m *MockBuilderModule) Start() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start")
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start.
func (mr *MockBuilderModuleMockRecorder) Start() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockBuilderModule)(nil).Start))
}

// Stop mocks base method.
func (m *MockBuilderModule) Stop() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Stop")
}

// Stop indicates an expected call of Stop.
func (mr *MockBuilderModuleMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockBuilderModule)(nil).Stop))
}
