// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/stackrox/rox/central/cluster/index (interfaces: Indexer)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	v1 "github.com/stackrox/rox/generated/api/v1"
	storage "github.com/stackrox/rox/generated/storage"
	search "github.com/stackrox/rox/pkg/search"
	blevesearch "github.com/stackrox/rox/pkg/search/blevesearch"
	reflect "reflect"
)

// MockIndexer is a mock of Indexer interface
type MockIndexer struct {
	ctrl     *gomock.Controller
	recorder *MockIndexerMockRecorder
}

// MockIndexerMockRecorder is the mock recorder for MockIndexer
type MockIndexerMockRecorder struct {
	mock *MockIndexer
}

// NewMockIndexer creates a new mock instance
func NewMockIndexer(ctrl *gomock.Controller) *MockIndexer {
	mock := &MockIndexer{ctrl: ctrl}
	mock.recorder = &MockIndexerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIndexer) EXPECT() *MockIndexerMockRecorder {
	return m.recorder
}

// AddCluster mocks base method
func (m *MockIndexer) AddCluster(arg0 *storage.Cluster) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCluster", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddCluster indicates an expected call of AddCluster
func (mr *MockIndexerMockRecorder) AddCluster(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCluster", reflect.TypeOf((*MockIndexer)(nil).AddCluster), arg0)
}

// AddClusters mocks base method
func (m *MockIndexer) AddClusters(arg0 []*storage.Cluster) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddClusters", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddClusters indicates an expected call of AddClusters
func (mr *MockIndexerMockRecorder) AddClusters(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddClusters", reflect.TypeOf((*MockIndexer)(nil).AddClusters), arg0)
}

// Count mocks base method
func (m *MockIndexer) Count(arg0 *v1.Query, arg1 ...blevesearch.SearchOption) (int, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Count", varargs...)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count
func (mr *MockIndexerMockRecorder) Count(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockIndexer)(nil).Count), varargs...)
}

// DeleteCluster mocks base method
func (m *MockIndexer) DeleteCluster(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCluster", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCluster indicates an expected call of DeleteCluster
func (mr *MockIndexerMockRecorder) DeleteCluster(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCluster", reflect.TypeOf((*MockIndexer)(nil).DeleteCluster), arg0)
}

// DeleteClusters mocks base method
func (m *MockIndexer) DeleteClusters(arg0 []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteClusters", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteClusters indicates an expected call of DeleteClusters
func (mr *MockIndexerMockRecorder) DeleteClusters(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteClusters", reflect.TypeOf((*MockIndexer)(nil).DeleteClusters), arg0)
}

// MarkInitialIndexingComplete mocks base method
func (m *MockIndexer) MarkInitialIndexingComplete() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MarkInitialIndexingComplete")
	ret0, _ := ret[0].(error)
	return ret0
}

// MarkInitialIndexingComplete indicates an expected call of MarkInitialIndexingComplete
func (mr *MockIndexerMockRecorder) MarkInitialIndexingComplete() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkInitialIndexingComplete", reflect.TypeOf((*MockIndexer)(nil).MarkInitialIndexingComplete))
}

// NeedsInitialIndexing mocks base method
func (m *MockIndexer) NeedsInitialIndexing() (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NeedsInitialIndexing")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NeedsInitialIndexing indicates an expected call of NeedsInitialIndexing
func (mr *MockIndexerMockRecorder) NeedsInitialIndexing() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NeedsInitialIndexing", reflect.TypeOf((*MockIndexer)(nil).NeedsInitialIndexing))
}

// Search mocks base method
func (m *MockIndexer) Search(arg0 *v1.Query, arg1 ...blevesearch.SearchOption) ([]search.Result, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Search", varargs...)
	ret0, _ := ret[0].([]search.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Search indicates an expected call of Search
func (mr *MockIndexerMockRecorder) Search(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Search", reflect.TypeOf((*MockIndexer)(nil).Search), varargs...)
}
