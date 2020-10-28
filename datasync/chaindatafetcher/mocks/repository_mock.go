// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/klaytn/klaytn/datasync/chaindatafetcher (interfaces: Repository)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	blockchain "github.com/klaytn/klaytn/blockchain"
	types "github.com/klaytn/klaytn/datasync/chaindatafetcher/types"
)

// MockRepository is a mock of Repository interface
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// HandleChainEvent mocks base method
func (m *MockRepository) HandleChainEvent(arg0 blockchain.ChainEvent, arg1 types.RequestType) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HandleChainEvent", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// HandleChainEvent indicates an expected call of HandleChainEvent
func (mr *MockRepositoryMockRecorder) HandleChainEvent(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleChainEvent", reflect.TypeOf((*MockRepository)(nil).HandleChainEvent), arg0, arg1)
}
