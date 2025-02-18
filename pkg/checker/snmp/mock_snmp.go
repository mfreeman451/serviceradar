// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mfreeman451/serviceradar/pkg/checker/snmp (interfaces: Collector,Aggregator,Service,CollectorFactory,AggregatorFactory,SNMPClient,SNMPManager,DataStore)
//
// Generated by this command:
//
//	mockgen -destination=mock_snmp.go -package=snmp github.com/mfreeman451/serviceradar/pkg/checker/snmp Collector,Aggregator,Service,CollectorFactory,AggregatorFactory,SNMPClient,SNMPManager,DataStore
//

// Package snmp is a generated GoMock package.
package snmp

import (
	context "context"
	reflect "reflect"
	time "time"

	db "github.com/mfreeman451/serviceradar/pkg/db"
	gomock "go.uber.org/mock/gomock"
)

// MockCollector is a mock of Collector interface.
type MockCollector struct {
	ctrl     *gomock.Controller
	recorder *MockCollectorMockRecorder
	isgomock struct{}
}

// MockCollectorMockRecorder is the mock recorder for MockCollector.
type MockCollectorMockRecorder struct {
	mock *MockCollector
}

// NewMockCollector creates a new mock instance.
func NewMockCollector(ctrl *gomock.Controller) *MockCollector {
	mock := &MockCollector{ctrl: ctrl}
	mock.recorder = &MockCollectorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCollector) EXPECT() *MockCollectorMockRecorder {
	return m.recorder
}

// GetResults mocks base method.
func (m *MockCollector) GetResults() <-chan DataPoint {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetResults")
	ret0, _ := ret[0].(<-chan DataPoint)
	return ret0
}

// GetResults indicates an expected call of GetResults.
func (mr *MockCollectorMockRecorder) GetResults() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetResults", reflect.TypeOf((*MockCollector)(nil).GetResults))
}

// GetStatus mocks base method.
func (m *MockCollector) GetStatus() TargetStatus {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStatus")
	ret0, _ := ret[0].(TargetStatus)
	return ret0
}

// GetStatus indicates an expected call of GetStatus.
func (mr *MockCollectorMockRecorder) GetStatus() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStatus", reflect.TypeOf((*MockCollector)(nil).GetStatus))
}

// Start mocks base method.
func (m *MockCollector) Start(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start.
func (mr *MockCollectorMockRecorder) Start(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockCollector)(nil).Start), ctx)
}

// Stop mocks base method.
func (m *MockCollector) Stop() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stop")
	ret0, _ := ret[0].(error)
	return ret0
}

// Stop indicates an expected call of Stop.
func (mr *MockCollectorMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockCollector)(nil).Stop))
}

// MockAggregator is a mock of Aggregator interface.
type MockAggregator struct {
	ctrl     *gomock.Controller
	recorder *MockAggregatorMockRecorder
	isgomock struct{}
}

// MockAggregatorMockRecorder is the mock recorder for MockAggregator.
type MockAggregatorMockRecorder struct {
	mock *MockAggregator
}

// NewMockAggregator creates a new mock instance.
func NewMockAggregator(ctrl *gomock.Controller) *MockAggregator {
	mock := &MockAggregator{ctrl: ctrl}
	mock.recorder = &MockAggregatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAggregator) EXPECT() *MockAggregatorMockRecorder {
	return m.recorder
}

// AddPoint mocks base method.
func (m *MockAggregator) AddPoint(point *DataPoint) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddPoint", point)
}

// AddPoint indicates an expected call of AddPoint.
func (mr *MockAggregatorMockRecorder) AddPoint(point any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddPoint", reflect.TypeOf((*MockAggregator)(nil).AddPoint), point)
}

// GetAggregatedData mocks base method.
func (m *MockAggregator) GetAggregatedData(oidName string, interval Interval) (*DataPoint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAggregatedData", oidName, interval)
	ret0, _ := ret[0].(*DataPoint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAggregatedData indicates an expected call of GetAggregatedData.
func (mr *MockAggregatorMockRecorder) GetAggregatedData(oidName, interval any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAggregatedData", reflect.TypeOf((*MockAggregator)(nil).GetAggregatedData), oidName, interval)
}

// Reset mocks base method.
func (m *MockAggregator) Reset() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Reset")
}

// Reset indicates an expected call of Reset.
func (mr *MockAggregatorMockRecorder) Reset() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reset", reflect.TypeOf((*MockAggregator)(nil).Reset))
}

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
	isgomock struct{}
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// AddTarget mocks base method.
func (m *MockService) AddTarget(ctx context.Context, target *Target) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddTarget", ctx, target)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddTarget indicates an expected call of AddTarget.
func (mr *MockServiceMockRecorder) AddTarget(ctx, target any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddTarget", reflect.TypeOf((*MockService)(nil).AddTarget), ctx, target)
}

// GetStatus mocks base method.
func (m *MockService) GetStatus(arg0 context.Context) (map[string]TargetStatus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStatus", arg0)
	ret0, _ := ret[0].(map[string]TargetStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStatus indicates an expected call of GetStatus.
func (mr *MockServiceMockRecorder) GetStatus(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStatus", reflect.TypeOf((*MockService)(nil).GetStatus), arg0)
}

// RemoveTarget mocks base method.
func (m *MockService) RemoveTarget(targetName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveTarget", targetName)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveTarget indicates an expected call of RemoveTarget.
func (mr *MockServiceMockRecorder) RemoveTarget(targetName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveTarget", reflect.TypeOf((*MockService)(nil).RemoveTarget), targetName)
}

// Start mocks base method.
func (m *MockService) Start(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start.
func (mr *MockServiceMockRecorder) Start(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockService)(nil).Start), ctx)
}

// Stop mocks base method.
func (m *MockService) Stop() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stop")
	ret0, _ := ret[0].(error)
	return ret0
}

// Stop indicates an expected call of Stop.
func (mr *MockServiceMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockService)(nil).Stop))
}

// MockCollectorFactory is a mock of CollectorFactory interface.
type MockCollectorFactory struct {
	ctrl     *gomock.Controller
	recorder *MockCollectorFactoryMockRecorder
	isgomock struct{}
}

// MockCollectorFactoryMockRecorder is the mock recorder for MockCollectorFactory.
type MockCollectorFactoryMockRecorder struct {
	mock *MockCollectorFactory
}

// NewMockCollectorFactory creates a new mock instance.
func NewMockCollectorFactory(ctrl *gomock.Controller) *MockCollectorFactory {
	mock := &MockCollectorFactory{ctrl: ctrl}
	mock.recorder = &MockCollectorFactoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCollectorFactory) EXPECT() *MockCollectorFactoryMockRecorder {
	return m.recorder
}

// CreateCollector mocks base method.
func (m *MockCollectorFactory) CreateCollector(target *Target) (Collector, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCollector", target)
	ret0, _ := ret[0].(Collector)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCollector indicates an expected call of CreateCollector.
func (mr *MockCollectorFactoryMockRecorder) CreateCollector(target any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCollector", reflect.TypeOf((*MockCollectorFactory)(nil).CreateCollector), target)
}

// MockAggregatorFactory is a mock of AggregatorFactory interface.
type MockAggregatorFactory struct {
	ctrl     *gomock.Controller
	recorder *MockAggregatorFactoryMockRecorder
	isgomock struct{}
}

// MockAggregatorFactoryMockRecorder is the mock recorder for MockAggregatorFactory.
type MockAggregatorFactoryMockRecorder struct {
	mock *MockAggregatorFactory
}

// NewMockAggregatorFactory creates a new mock instance.
func NewMockAggregatorFactory(ctrl *gomock.Controller) *MockAggregatorFactory {
	mock := &MockAggregatorFactory{ctrl: ctrl}
	mock.recorder = &MockAggregatorFactoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAggregatorFactory) EXPECT() *MockAggregatorFactoryMockRecorder {
	return m.recorder
}

// CreateAggregator mocks base method.
func (m *MockAggregatorFactory) CreateAggregator(interval time.Duration, maxPoints int) (Aggregator, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAggregator", interval, maxPoints)
	ret0, _ := ret[0].(Aggregator)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAggregator indicates an expected call of CreateAggregator.
func (mr *MockAggregatorFactoryMockRecorder) CreateAggregator(interval, maxPoints any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAggregator", reflect.TypeOf((*MockAggregatorFactory)(nil).CreateAggregator), interval, maxPoints)
}

// MockSNMPClient is a mock of SNMPClient interface.
type MockSNMPClient struct {
	ctrl     *gomock.Controller
	recorder *MockSNMPClientMockRecorder
	isgomock struct{}
}

// MockSNMPClientMockRecorder is the mock recorder for MockSNMPClient.
type MockSNMPClientMockRecorder struct {
	mock *MockSNMPClient
}

// NewMockSNMPClient creates a new mock instance.
func NewMockSNMPClient(ctrl *gomock.Controller) *MockSNMPClient {
	mock := &MockSNMPClient{ctrl: ctrl}
	mock.recorder = &MockSNMPClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSNMPClient) EXPECT() *MockSNMPClientMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockSNMPClient) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockSNMPClientMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockSNMPClient)(nil).Close))
}

// Connect mocks base method.
func (m *MockSNMPClient) Connect() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Connect")
	ret0, _ := ret[0].(error)
	return ret0
}

// Connect indicates an expected call of Connect.
func (mr *MockSNMPClientMockRecorder) Connect() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Connect", reflect.TypeOf((*MockSNMPClient)(nil).Connect))
}

// Get mocks base method.
func (m *MockSNMPClient) Get(oids []string) (map[string]any, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", oids)
	ret0, _ := ret[0].(map[string]any)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockSNMPClientMockRecorder) Get(oids any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockSNMPClient)(nil).Get), oids)
}

// MockSNMPManager is a mock of SNMPManager interface.
type MockSNMPManager struct {
	ctrl     *gomock.Controller
	recorder *MockSNMPManagerMockRecorder
	isgomock struct{}
}

// MockSNMPManagerMockRecorder is the mock recorder for MockSNMPManager.
type MockSNMPManagerMockRecorder struct {
	mock *MockSNMPManager
}

// NewMockSNMPManager creates a new mock instance.
func NewMockSNMPManager(ctrl *gomock.Controller) *MockSNMPManager {
	mock := &MockSNMPManager{ctrl: ctrl}
	mock.recorder = &MockSNMPManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSNMPManager) EXPECT() *MockSNMPManagerMockRecorder {
	return m.recorder
}

// GetSNMPMetrics mocks base method.
func (m *MockSNMPManager) GetSNMPMetrics(nodeID string, startTime, endTime time.Time) ([]db.SNMPMetric, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSNMPMetrics", nodeID, startTime, endTime)
	ret0, _ := ret[0].([]db.SNMPMetric)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSNMPMetrics indicates an expected call of GetSNMPMetrics.
func (mr *MockSNMPManagerMockRecorder) GetSNMPMetrics(nodeID, startTime, endTime any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSNMPMetrics", reflect.TypeOf((*MockSNMPManager)(nil).GetSNMPMetrics), nodeID, startTime, endTime)
}

// MockDataStore is a mock of DataStore interface.
type MockDataStore struct {
	ctrl     *gomock.Controller
	recorder *MockDataStoreMockRecorder
	isgomock struct{}
}

// MockDataStoreMockRecorder is the mock recorder for MockDataStore.
type MockDataStoreMockRecorder struct {
	mock *MockDataStore
}

// NewMockDataStore creates a new mock instance.
func NewMockDataStore(ctrl *gomock.Controller) *MockDataStore {
	mock := &MockDataStore{ctrl: ctrl}
	mock.recorder = &MockDataStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDataStore) EXPECT() *MockDataStoreMockRecorder {
	return m.recorder
}

// Cleanup mocks base method.
func (m *MockDataStore) Cleanup(age time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Cleanup", age)
	ret0, _ := ret[0].(error)
	return ret0
}

// Cleanup indicates an expected call of Cleanup.
func (mr *MockDataStoreMockRecorder) Cleanup(age any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Cleanup", reflect.TypeOf((*MockDataStore)(nil).Cleanup), age)
}

// Query mocks base method.
func (m *MockDataStore) Query(filter DataFilter) ([]DataPoint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Query", filter)
	ret0, _ := ret[0].([]DataPoint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Query indicates an expected call of Query.
func (mr *MockDataStoreMockRecorder) Query(filter any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Query", reflect.TypeOf((*MockDataStore)(nil).Query), filter)
}

// Store mocks base method.
func (m *MockDataStore) Store(point DataPoint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Store", point)
	ret0, _ := ret[0].(error)
	return ret0
}

// Store indicates an expected call of Store.
func (mr *MockDataStoreMockRecorder) Store(point any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Store", reflect.TypeOf((*MockDataStore)(nil).Store), point)
}
