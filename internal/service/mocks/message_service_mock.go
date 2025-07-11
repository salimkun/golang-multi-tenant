// Code generated by MockGen. DO NOT EDIT.
// Source: message_service.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockMessageServiceInterface is a mock of MessageServiceInterface interface.
type MockMessageServiceInterface struct {
	ctrl     *gomock.Controller
	recorder *MockMessageServiceInterfaceMockRecorder
}

// MockMessageServiceInterfaceMockRecorder is the mock recorder for MockMessageServiceInterface.
type MockMessageServiceInterfaceMockRecorder struct {
	mock *MockMessageServiceInterface
}

// NewMockMessageServiceInterface creates a new mock instance.
func NewMockMessageServiceInterface(ctrl *gomock.Controller) *MockMessageServiceInterface {
	mock := &MockMessageServiceInterface{ctrl: ctrl}
	mock.recorder = &MockMessageServiceInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMessageServiceInterface) EXPECT() *MockMessageServiceInterfaceMockRecorder {
	return m.recorder
}

// FetchMessages mocks base method.
func (m *MockMessageServiceInterface) FetchMessages(tenantID, cursor string, limit int) ([]map[string]interface{}, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchMessages", tenantID, cursor, limit)
	ret0, _ := ret[0].([]map[string]interface{})
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// FetchMessages indicates an expected call of FetchMessages.
func (mr *MockMessageServiceInterfaceMockRecorder) FetchMessages(tenantID, cursor, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchMessages", reflect.TypeOf((*MockMessageServiceInterface)(nil).FetchMessages), tenantID, cursor, limit)
}

// PublishToTenantQueue mocks base method.
func (m *MockMessageServiceInterface) PublishToTenantQueue(tenantID string, payload map[string]interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PublishToTenantQueue", tenantID, payload)
	ret0, _ := ret[0].(error)
	return ret0
}

// PublishToTenantQueue indicates an expected call of PublishToTenantQueue.
func (mr *MockMessageServiceInterfaceMockRecorder) PublishToTenantQueue(tenantID, payload interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublishToTenantQueue", reflect.TypeOf((*MockMessageServiceInterface)(nil).PublishToTenantQueue), tenantID, payload)
}
