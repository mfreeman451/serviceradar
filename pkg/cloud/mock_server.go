// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mfreeman451/serviceradar/pkg/cloud (interfaces: NodeService,CloudService)
//
// Generated by this command:
//
//	mockgen -destination=mock_server.go -package=cloud github.com/mfreeman451/serviceradar/pkg/cloud NodeService,CloudService
//

// Package cloud is a generated GoMock package.
package cloud

import (
	context "context"
	reflect "reflect"

	api "github.com/mfreeman451/serviceradar/pkg/cloud/api"
	metrics "github.com/mfreeman451/serviceradar/pkg/metrics"
	gomock "go.uber.org/mock/gomock"
)

// MockNodeService is a mock of NodeService interface.
type MockNodeService struct {
	ctrl     *gomock.Controller
	recorder *MockNodeServiceMockRecorder
	isgomock struct{}
}

// MockNodeServiceMockRecorder is the mock recorder for MockNodeService.
type MockNodeServiceMockRecorder struct {
	mock *MockNodeService
}

// NewMockNodeService creates a new mock instance.
func NewMockNodeService(ctrl *gomock.Controller) *MockNodeService {
	mock := &MockNodeService{ctrl: ctrl}
	mock.recorder = &MockNodeServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNodeService) EXPECT() *MockNodeServiceMockRecorder {
	return m.recorder
}

// CheckNodeHealth mocks base method.
func (m *MockNodeService) CheckNodeHealth(nodeID string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckNodeHealth", nodeID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckNodeHealth indicates an expected call of CheckNodeHealth.
func (mr *MockNodeServiceMockRecorder) CheckNodeHealth(nodeID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckNodeHealth", reflect.TypeOf((*MockNodeService)(nil).CheckNodeHealth), nodeID)
}

// GetNodeHistory mocks base method.
func (m *MockNodeService) GetNodeHistory(nodeID string, limit int) ([]api.NodeHistoryPoint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNodeHistory", nodeID, limit)
	ret0, _ := ret[0].([]api.NodeHistoryPoint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNodeHistory indicates an expected call of GetNodeHistory.
func (mr *MockNodeServiceMockRecorder) GetNodeHistory(nodeID, limit any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNodeHistory", reflect.TypeOf((*MockNodeService)(nil).GetNodeHistory), nodeID, limit)
}

// GetNodeStatus mocks base method.
func (m *MockNodeService) GetNodeStatus(nodeID string) (*api.NodeStatus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNodeStatus", nodeID)
	ret0, _ := ret[0].(*api.NodeStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNodeStatus indicates an expected call of GetNodeStatus.
func (mr *MockNodeServiceMockRecorder) GetNodeStatus(nodeID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNodeStatus", reflect.TypeOf((*MockNodeService)(nil).GetNodeStatus), nodeID)
}

// UpdateNodeStatus mocks base method.
func (m *MockNodeService) UpdateNodeStatus(nodeID string, status *api.NodeStatus) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateNodeStatus", nodeID, status)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateNodeStatus indicates an expected call of UpdateNodeStatus.
func (mr *MockNodeServiceMockRecorder) UpdateNodeStatus(nodeID, status any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateNodeStatus", reflect.TypeOf((*MockNodeService)(nil).UpdateNodeStatus), nodeID, status)
}

// MockCloudService is a mock of CloudService interface.
type MockCloudService struct {
	ctrl     *gomock.Controller
	recorder *MockCloudServiceMockRecorder
	isgomock struct{}
}

// MockCloudServiceMockRecorder is the mock recorder for MockCloudService.
type MockCloudServiceMockRecorder struct {
	mock *MockCloudService
}

// NewMockCloudService creates a new mock instance.
func NewMockCloudService(ctrl *gomock.Controller) *MockCloudService {
	mock := &MockCloudService{ctrl: ctrl}
	mock.recorder = &MockCloudServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCloudService) EXPECT() *MockCloudServiceMockRecorder {
	return m.recorder
}

// GetMetricsManager mocks base method.
func (m *MockCloudService) GetMetricsManager() metrics.MetricCollector {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMetricsManager")
	ret0, _ := ret[0].(metrics.MetricCollector)
	return ret0
}

// GetMetricsManager indicates an expected call of GetMetricsManager.
func (mr *MockCloudServiceMockRecorder) GetMetricsManager() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMetricsManager", reflect.TypeOf((*MockCloudService)(nil).GetMetricsManager))
}

// ReportStatus mocks base method.
func (m *MockCloudService) ReportStatus(ctx context.Context, nodeID string, status *api.NodeStatus) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReportStatus", ctx, nodeID, status)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReportStatus indicates an expected call of ReportStatus.
func (mr *MockCloudServiceMockRecorder) ReportStatus(ctx, nodeID, status any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReportStatus", reflect.TypeOf((*MockCloudService)(nil).ReportStatus), ctx, nodeID, status)
}

// Start mocks base method.
func (m *MockCloudService) Start(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start.
func (mr *MockCloudServiceMockRecorder) Start(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockCloudService)(nil).Start), ctx)
}

// Stop mocks base method.
func (m *MockCloudService) Stop(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stop", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Stop indicates an expected call of Stop.
func (mr *MockCloudServiceMockRecorder) Stop(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockCloudService)(nil).Stop), ctx)
}
