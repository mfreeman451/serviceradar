// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mfreeman451/serviceradar/pkg/cloud (interfaces: DatabaseService,AlertService,MetricsService,APIService,NodeService,CloudService,TransactionService)
//
// Generated by this command:
//
//	mockgen -destination=mock_server.go -package=cloud github.com/mfreeman451/serviceradar/pkg/cloud DatabaseService,AlertService,MetricsService,APIService,NodeService,CloudService,TransactionService
//

// Package cloud is a generated GoMock package.
package cloud

import (
	context "context"
	sql "database/sql"
	reflect "reflect"
	time "time"

	alerts "github.com/mfreeman451/serviceradar/pkg/cloud/alerts"
	api "github.com/mfreeman451/serviceradar/pkg/cloud/api"
	metrics "github.com/mfreeman451/serviceradar/pkg/metrics"
	models "github.com/mfreeman451/serviceradar/pkg/models"
	gomock "go.uber.org/mock/gomock"
)

// MockDatabaseService is a mock of DatabaseService interface.
type MockDatabaseService struct {
	ctrl     *gomock.Controller
	recorder *MockDatabaseServiceMockRecorder
	isgomock struct{}
}

// MockDatabaseServiceMockRecorder is the mock recorder for MockDatabaseService.
type MockDatabaseServiceMockRecorder struct {
	mock *MockDatabaseService
}

// NewMockDatabaseService creates a new mock instance.
func NewMockDatabaseService(ctrl *gomock.Controller) *MockDatabaseService {
	mock := &MockDatabaseService{ctrl: ctrl}
	mock.recorder = &MockDatabaseServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDatabaseService) EXPECT() *MockDatabaseServiceMockRecorder {
	return m.recorder
}

// Begin mocks base method.
func (m *MockDatabaseService) Begin() (*sql.Tx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Begin")
	ret0, _ := ret[0].(*sql.Tx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Begin indicates an expected call of Begin.
func (mr *MockDatabaseServiceMockRecorder) Begin() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Begin", reflect.TypeOf((*MockDatabaseService)(nil).Begin))
}

// CleanOldData mocks base method.
func (m *MockDatabaseService) CleanOldData(age time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CleanOldData", age)
	ret0, _ := ret[0].(error)
	return ret0
}

// CleanOldData indicates an expected call of CleanOldData.
func (mr *MockDatabaseServiceMockRecorder) CleanOldData(age any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CleanOldData", reflect.TypeOf((*MockDatabaseService)(nil).CleanOldData), age)
}

// Close mocks base method.
func (m *MockDatabaseService) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockDatabaseServiceMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockDatabaseService)(nil).Close))
}

// Exec mocks base method.
func (m *MockDatabaseService) Exec(query string, args ...any) (sql.Result, error) {
	m.ctrl.T.Helper()
	varargs := []any{query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Exec", varargs...)
	ret0, _ := ret[0].(sql.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exec indicates an expected call of Exec.
func (mr *MockDatabaseServiceMockRecorder) Exec(query any, args ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exec", reflect.TypeOf((*MockDatabaseService)(nil).Exec), varargs...)
}

// GetNodeHistoryPoints mocks base method.
func (m *MockDatabaseService) GetNodeHistoryPoints(nodeID string, limit int) ([]NodeHistoryPoint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNodeHistoryPoints", nodeID, limit)
	ret0, _ := ret[0].([]NodeHistoryPoint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNodeHistoryPoints indicates an expected call of GetNodeHistoryPoints.
func (mr *MockDatabaseServiceMockRecorder) GetNodeHistoryPoints(nodeID, limit any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNodeHistoryPoints", reflect.TypeOf((*MockDatabaseService)(nil).GetNodeHistoryPoints), nodeID, limit)
}

// Query mocks base method.
func (m *MockDatabaseService) Query(query string, args ...any) (*sql.Rows, error) {
	m.ctrl.T.Helper()
	varargs := []any{query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Query", varargs...)
	ret0, _ := ret[0].(*sql.Rows)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Query indicates an expected call of Query.
func (mr *MockDatabaseServiceMockRecorder) Query(query any, args ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Query", reflect.TypeOf((*MockDatabaseService)(nil).Query), varargs...)
}

// QueryRow mocks base method.
func (m *MockDatabaseService) QueryRow(query string, args ...any) *sql.Row {
	m.ctrl.T.Helper()
	varargs := []any{query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "QueryRow", varargs...)
	ret0, _ := ret[0].(*sql.Row)
	return ret0
}

// QueryRow indicates an expected call of QueryRow.
func (mr *MockDatabaseServiceMockRecorder) QueryRow(query any, args ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryRow", reflect.TypeOf((*MockDatabaseService)(nil).QueryRow), varargs...)
}

// UpdateNodeStatus mocks base method.
func (m *MockDatabaseService) UpdateNodeStatus(status *NodeStatus) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateNodeStatus", status)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateNodeStatus indicates an expected call of UpdateNodeStatus.
func (mr *MockDatabaseServiceMockRecorder) UpdateNodeStatus(status any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateNodeStatus", reflect.TypeOf((*MockDatabaseService)(nil).UpdateNodeStatus), status)
}

// UpdateServiceStatus mocks base method.
func (m *MockDatabaseService) UpdateServiceStatus(status *ServiceStatus) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateServiceStatus", status)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateServiceStatus indicates an expected call of UpdateServiceStatus.
func (mr *MockDatabaseServiceMockRecorder) UpdateServiceStatus(status any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateServiceStatus", reflect.TypeOf((*MockDatabaseService)(nil).UpdateServiceStatus), status)
}

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
func (m *MockAlertService) Alert(ctx context.Context, alert *alerts.WebhookAlert) error {
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

// MockMetricsService is a mock of MetricsService interface.
type MockMetricsService struct {
	ctrl     *gomock.Controller
	recorder *MockMetricsServiceMockRecorder
	isgomock struct{}
}

// MockMetricsServiceMockRecorder is the mock recorder for MockMetricsService.
type MockMetricsServiceMockRecorder struct {
	mock *MockMetricsService
}

// NewMockMetricsService creates a new mock instance.
func NewMockMetricsService(ctrl *gomock.Controller) *MockMetricsService {
	mock := &MockMetricsService{ctrl: ctrl}
	mock.recorder = &MockMetricsServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMetricsService) EXPECT() *MockMetricsServiceMockRecorder {
	return m.recorder
}

// AddMetric mocks base method.
func (m *MockMetricsService) AddMetric(nodeID string, timestamp time.Time, value int64, serviceType string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddMetric", nodeID, timestamp, value, serviceType)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddMetric indicates an expected call of AddMetric.
func (mr *MockMetricsServiceMockRecorder) AddMetric(nodeID, timestamp, value, serviceType any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddMetric", reflect.TypeOf((*MockMetricsService)(nil).AddMetric), nodeID, timestamp, value, serviceType)
}

// CleanupStaleNodes mocks base method.
func (m *MockMetricsService) CleanupStaleNodes(age time.Duration) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CleanupStaleNodes", age)
}

// CleanupStaleNodes indicates an expected call of CleanupStaleNodes.
func (mr *MockMetricsServiceMockRecorder) CleanupStaleNodes(age any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CleanupStaleNodes", reflect.TypeOf((*MockMetricsService)(nil).CleanupStaleNodes), age)
}

// GetMetrics mocks base method.
func (m *MockMetricsService) GetMetrics(nodeID string) []models.MetricPoint {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMetrics", nodeID)
	ret0, _ := ret[0].([]models.MetricPoint)
	return ret0
}

// GetMetrics indicates an expected call of GetMetrics.
func (mr *MockMetricsServiceMockRecorder) GetMetrics(nodeID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMetrics", reflect.TypeOf((*MockMetricsService)(nil).GetMetrics), nodeID)
}

// MockAPIService is a mock of APIService interface.
type MockAPIService struct {
	ctrl     *gomock.Controller
	recorder *MockAPIServiceMockRecorder
	isgomock struct{}
}

// MockAPIServiceMockRecorder is the mock recorder for MockAPIService.
type MockAPIServiceMockRecorder struct {
	mock *MockAPIService
}

// NewMockAPIService creates a new mock instance.
func NewMockAPIService(ctrl *gomock.Controller) *MockAPIService {
	mock := &MockAPIService{ctrl: ctrl}
	mock.recorder = &MockAPIServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAPIService) EXPECT() *MockAPIServiceMockRecorder {
	return m.recorder
}

// SetKnownPollers mocks base method.
func (m *MockAPIService) SetKnownPollers(knownPollers []string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetKnownPollers", knownPollers)
}

// SetKnownPollers indicates an expected call of SetKnownPollers.
func (mr *MockAPIServiceMockRecorder) SetKnownPollers(knownPollers any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetKnownPollers", reflect.TypeOf((*MockAPIService)(nil).SetKnownPollers), knownPollers)
}

// SetNodeHistoryHandler mocks base method.
func (m *MockAPIService) SetNodeHistoryHandler(handler func(string) ([]api.NodeHistoryPoint, error)) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetNodeHistoryHandler", handler)
}

// SetNodeHistoryHandler indicates an expected call of SetNodeHistoryHandler.
func (mr *MockAPIServiceMockRecorder) SetNodeHistoryHandler(handler any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetNodeHistoryHandler", reflect.TypeOf((*MockAPIService)(nil).SetNodeHistoryHandler), handler)
}

// Start mocks base method.
func (m *MockAPIService) Start(addr string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start", addr)
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start.
func (mr *MockAPIServiceMockRecorder) Start(addr any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockAPIService)(nil).Start), addr)
}

// UpdateNodeStatus mocks base method.
func (m *MockAPIService) UpdateNodeStatus(nodeID string, status *api.NodeStatus) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpdateNodeStatus", nodeID, status)
}

// UpdateNodeStatus indicates an expected call of UpdateNodeStatus.
func (mr *MockAPIServiceMockRecorder) UpdateNodeStatus(nodeID, status any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateNodeStatus", reflect.TypeOf((*MockAPIService)(nil).UpdateNodeStatus), nodeID, status)
}

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
func (m *MockNodeService) GetNodeHistory(nodeID string, limit int) ([]NodeHistoryPoint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNodeHistory", nodeID, limit)
	ret0, _ := ret[0].([]NodeHistoryPoint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNodeHistory indicates an expected call of GetNodeHistory.
func (mr *MockNodeServiceMockRecorder) GetNodeHistory(nodeID, limit any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNodeHistory", reflect.TypeOf((*MockNodeService)(nil).GetNodeHistory), nodeID, limit)
}

// GetNodeStatus mocks base method.
func (m *MockNodeService) GetNodeStatus(nodeID string) (*NodeStatus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNodeStatus", nodeID)
	ret0, _ := ret[0].(*NodeStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNodeStatus indicates an expected call of GetNodeStatus.
func (mr *MockNodeServiceMockRecorder) GetNodeStatus(nodeID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNodeStatus", reflect.TypeOf((*MockNodeService)(nil).GetNodeStatus), nodeID)
}

// UpdateNodeStatus mocks base method.
func (m *MockNodeService) UpdateNodeStatus(nodeID string, status *NodeStatus) error {
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
func (m *MockCloudService) ReportStatus(ctx context.Context, nodeID string, status *NodeStatus) error {
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

// MockTransactionService is a mock of TransactionService interface.
type MockTransactionService struct {
	ctrl     *gomock.Controller
	recorder *MockTransactionServiceMockRecorder
	isgomock struct{}
}

// MockTransactionServiceMockRecorder is the mock recorder for MockTransactionService.
type MockTransactionServiceMockRecorder struct {
	mock *MockTransactionService
}

// NewMockTransactionService creates a new mock instance.
func NewMockTransactionService(ctrl *gomock.Controller) *MockTransactionService {
	mock := &MockTransactionService{ctrl: ctrl}
	mock.recorder = &MockTransactionServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransactionService) EXPECT() *MockTransactionServiceMockRecorder {
	return m.recorder
}

// Begin mocks base method.
func (m *MockTransactionService) Begin() (*sql.Tx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Begin")
	ret0, _ := ret[0].(*sql.Tx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Begin indicates an expected call of Begin.
func (mr *MockTransactionServiceMockRecorder) Begin() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Begin", reflect.TypeOf((*MockTransactionService)(nil).Begin))
}

// Commit mocks base method.
func (m *MockTransactionService) Commit() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Commit")
	ret0, _ := ret[0].(error)
	return ret0
}

// Commit indicates an expected call of Commit.
func (mr *MockTransactionServiceMockRecorder) Commit() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Commit", reflect.TypeOf((*MockTransactionService)(nil).Commit))
}

// Exec mocks base method.
func (m *MockTransactionService) Exec(query string, args ...any) (sql.Result, error) {
	m.ctrl.T.Helper()
	varargs := []any{query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Exec", varargs...)
	ret0, _ := ret[0].(sql.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exec indicates an expected call of Exec.
func (mr *MockTransactionServiceMockRecorder) Exec(query any, args ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exec", reflect.TypeOf((*MockTransactionService)(nil).Exec), varargs...)
}

// Query mocks base method.
func (m *MockTransactionService) Query(query string, args ...any) (*sql.Rows, error) {
	m.ctrl.T.Helper()
	varargs := []any{query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Query", varargs...)
	ret0, _ := ret[0].(*sql.Rows)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Query indicates an expected call of Query.
func (mr *MockTransactionServiceMockRecorder) Query(query any, args ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Query", reflect.TypeOf((*MockTransactionService)(nil).Query), varargs...)
}

// QueryRow mocks base method.
func (m *MockTransactionService) QueryRow(query string, args ...any) *sql.Row {
	m.ctrl.T.Helper()
	varargs := []any{query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "QueryRow", varargs...)
	ret0, _ := ret[0].(*sql.Row)
	return ret0
}

// QueryRow indicates an expected call of QueryRow.
func (mr *MockTransactionServiceMockRecorder) QueryRow(query any, args ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryRow", reflect.TypeOf((*MockTransactionService)(nil).QueryRow), varargs...)
}

// Rollback mocks base method.
func (m *MockTransactionService) Rollback() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Rollback")
	ret0, _ := ret[0].(error)
	return ret0
}

// Rollback indicates an expected call of Rollback.
func (mr *MockTransactionServiceMockRecorder) Rollback() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Rollback", reflect.TypeOf((*MockTransactionService)(nil).Rollback))
}
