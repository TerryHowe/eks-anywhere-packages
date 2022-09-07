// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/artifacts/puller.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockPuller is a mock of Puller interface.
type MockPuller struct {
	ctrl     *gomock.Controller
	recorder *MockPullerMockRecorder
}

// MockPullerMockRecorder is the mock recorder for MockPuller.
type MockPullerMockRecorder struct {
	mock *MockPuller
}

// NewMockPuller creates a new mock instance.
func NewMockPuller(ctrl *gomock.Controller) *MockPuller {
	mock := &MockPuller{ctrl: ctrl}
	mock.recorder = &MockPullerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPuller) EXPECT() *MockPullerMockRecorder {
	return m.recorder
}

// Pull mocks base method.
func (m *MockPuller) Pull(ctx context.Context, ref string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Pull", ctx, ref)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Pull indicates an expected call of Pull.
func (mr *MockPullerMockRecorder) Pull(ctx, ref interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Pull", reflect.TypeOf((*MockPuller)(nil).Pull), ctx, ref)
}