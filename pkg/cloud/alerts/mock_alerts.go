// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mfreeman451/serviceradar/pkg/cloud/alerts (interfaces: AlertService)
//
// Generated by this command:
//
//	mockgen -destination=mock_alerts.go -package=alerts github.com/mfreeman451/serviceradar/pkg/cloud/alerts AlertService
//

// Package alerts is a generated GoMock package.
package alerts

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockAlertService is a mock of AlertService interface.
type MockAlertService struct {
	ctrl     *gomock.Controller
	recorder *MockAlertServiceMockRecorder
	isgomock struct{}
}

// MockAlertServiceMockRecorder is the mock recorder for MockAlertService.
type MockAlertServiceMockRecorder struct {
	mock *MockAlertService
}

// NewMockAlertService creates a new mock instance.
func NewMockAlertService(ctrl *gomock.Controller) *MockAlertService {
	mock := &MockAlertService{ctrl: ctrl}
	mock.recorder = &MockAlertServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAlertService) EXPECT() *MockAlertServiceMockRecorder {
	return m.recorder
}

// Alert mocks base method.
func (m *MockAlertService) Alert(ctx context.Context, alert *WebhookAlert) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Alert", ctx, alert)
	ret0, _ := ret[0].(error)
	return ret0
}

// Alert indicates an expected call of Alert.
func (mr *MockAlertServiceMockRecorder) Alert(ctx, alert any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Alert", reflect.TypeOf((*MockAlertService)(nil).Alert), ctx, alert)
}

// IsEnabled mocks base method.
func (m *MockAlertService) IsEnabled() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsEnabled")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsEnabled indicates an expected call of IsEnabled.
func (mr *MockAlertServiceMockRecorder) IsEnabled() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsEnabled", reflect.TypeOf((*MockAlertService)(nil).IsEnabled))
}
