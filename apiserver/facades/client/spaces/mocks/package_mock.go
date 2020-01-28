// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/apiserver/facades/client/spaces (interfaces: Backing,BlockChecker,Machine)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	set "github.com/juju/collections/set"
	networkingcommon "github.com/juju/juju/apiserver/common/networkingcommon"
	spaces "github.com/juju/juju/apiserver/facades/client/spaces"
	controller "github.com/juju/juju/controller"
	constraints "github.com/juju/juju/core/constraints"
	network "github.com/juju/juju/core/network"
	settings "github.com/juju/juju/core/settings"
	environs "github.com/juju/juju/environs"
	config "github.com/juju/juju/environs/config"
	names_v3 "gopkg.in/juju/names.v3"
	reflect "reflect"
)

// MockBacking is a mock of Backing interface
type MockBacking struct {
	ctrl     *gomock.Controller
	recorder *MockBackingMockRecorder
}

// MockBackingMockRecorder is the mock recorder for MockBacking
type MockBackingMockRecorder struct {
	mock *MockBacking
}

// NewMockBacking creates a new mock instance
func NewMockBacking(ctrl *gomock.Controller) *MockBacking {
	mock := &MockBacking{ctrl: ctrl}
	mock.recorder = &MockBackingMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockBacking) EXPECT() *MockBackingMockRecorder {
	return m.recorder
}

// AddSpace mocks base method
func (m *MockBacking) AddSpace(arg0 string, arg1 network.Id, arg2 []string, arg3 bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddSpace", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddSpace indicates an expected call of AddSpace
func (mr *MockBackingMockRecorder) AddSpace(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddSpace", reflect.TypeOf((*MockBacking)(nil).AddSpace), arg0, arg1, arg2, arg3)
}

// AllEndpointBindings mocks base method
func (m *MockBacking) AllEndpointBindings() ([]spaces.ApplicationEndpointBindingsShim, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllEndpointBindings")
	ret0, _ := ret[0].([]spaces.ApplicationEndpointBindingsShim)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllEndpointBindings indicates an expected call of AllEndpointBindings
func (mr *MockBackingMockRecorder) AllEndpointBindings() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllEndpointBindings", reflect.TypeOf((*MockBacking)(nil).AllEndpointBindings))
}

// AllMachines mocks base method
func (m *MockBacking) AllMachines() ([]spaces.Machine, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllMachines")
	ret0, _ := ret[0].([]spaces.Machine)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllMachines indicates an expected call of AllMachines
func (mr *MockBackingMockRecorder) AllMachines() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllMachines", reflect.TypeOf((*MockBacking)(nil).AllMachines))
}

// AllSpaces mocks base method
func (m *MockBacking) AllSpaces() ([]networkingcommon.BackingSpace, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllSpaces")
	ret0, _ := ret[0].([]networkingcommon.BackingSpace)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllSpaces indicates an expected call of AllSpaces
func (mr *MockBackingMockRecorder) AllSpaces() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllSpaces", reflect.TypeOf((*MockBacking)(nil).AllSpaces))
}

// CloudSpec mocks base method
func (m *MockBacking) CloudSpec() (environs.CloudSpec, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloudSpec")
	ret0, _ := ret[0].(environs.CloudSpec)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CloudSpec indicates an expected call of CloudSpec
func (mr *MockBackingMockRecorder) CloudSpec() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloudSpec", reflect.TypeOf((*MockBacking)(nil).CloudSpec))
}

// ConstraintsBySpace mocks base method
func (m *MockBacking) ConstraintsBySpace(arg0 string) (map[string]constraints.Value, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConstraintsBySpace", arg0)
	ret0, _ := ret[0].(map[string]constraints.Value)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ConstraintsBySpace indicates an expected call of ConstraintsBySpace
func (mr *MockBackingMockRecorder) ConstraintsBySpace(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConstraintsBySpace", reflect.TypeOf((*MockBacking)(nil).ConstraintsBySpace), arg0)
}

// ControllerConfig mocks base method
func (m *MockBacking) ControllerConfig() (controller.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ControllerConfig")
	ret0, _ := ret[0].(controller.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ControllerConfig indicates an expected call of ControllerConfig
func (mr *MockBackingMockRecorder) ControllerConfig() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ControllerConfig", reflect.TypeOf((*MockBacking)(nil).ControllerConfig))
}

// ModelConfig mocks base method
func (m *MockBacking) ModelConfig() (*config.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ModelConfig")
	ret0, _ := ret[0].(*config.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ModelConfig indicates an expected call of ModelConfig
func (mr *MockBackingMockRecorder) ModelConfig() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ModelConfig", reflect.TypeOf((*MockBacking)(nil).ModelConfig))
}

// ModelTag mocks base method
func (m *MockBacking) ModelTag() names_v3.ModelTag {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ModelTag")
	ret0, _ := ret[0].(names_v3.ModelTag)
	return ret0
}

// ModelTag indicates an expected call of ModelTag
func (mr *MockBackingMockRecorder) ModelTag() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ModelTag", reflect.TypeOf((*MockBacking)(nil).ModelTag))
}

// ReloadSpaces mocks base method
func (m *MockBacking) ReloadSpaces(arg0 environs.BootstrapEnviron) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReloadSpaces", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReloadSpaces indicates an expected call of ReloadSpaces
func (mr *MockBackingMockRecorder) ReloadSpaces(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReloadSpaces", reflect.TypeOf((*MockBacking)(nil).ReloadSpaces), arg0)
}

// RenameSpace mocks base method
func (m *MockBacking) RenameSpace(arg0 settings.ItemChanges, arg1 map[string]constraints.Value, arg2, arg3 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RenameSpace", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// RenameSpace indicates an expected call of RenameSpace
func (mr *MockBackingMockRecorder) RenameSpace(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RenameSpace", reflect.TypeOf((*MockBacking)(nil).RenameSpace), arg0, arg1, arg2, arg3)
}

// SpaceByName mocks base method
func (m *MockBacking) SpaceByName(arg0 string) (networkingcommon.BackingSpace, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SpaceByName", arg0)
	ret0, _ := ret[0].(networkingcommon.BackingSpace)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SpaceByName indicates an expected call of SpaceByName
func (mr *MockBackingMockRecorder) SpaceByName(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SpaceByName", reflect.TypeOf((*MockBacking)(nil).SpaceByName), arg0)
}

// SubnetByCIDR mocks base method
func (m *MockBacking) SubnetByCIDR(arg0 string) (networkingcommon.BackingSubnet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubnetByCIDR", arg0)
	ret0, _ := ret[0].(networkingcommon.BackingSubnet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SubnetByCIDR indicates an expected call of SubnetByCIDR
func (mr *MockBackingMockRecorder) SubnetByCIDR(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubnetByCIDR", reflect.TypeOf((*MockBacking)(nil).SubnetByCIDR), arg0)
}

// MockBlockChecker is a mock of BlockChecker interface
type MockBlockChecker struct {
	ctrl     *gomock.Controller
	recorder *MockBlockCheckerMockRecorder
}

// MockBlockCheckerMockRecorder is the mock recorder for MockBlockChecker
type MockBlockCheckerMockRecorder struct {
	mock *MockBlockChecker
}

// NewMockBlockChecker creates a new mock instance
func NewMockBlockChecker(ctrl *gomock.Controller) *MockBlockChecker {
	mock := &MockBlockChecker{ctrl: ctrl}
	mock.recorder = &MockBlockCheckerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockBlockChecker) EXPECT() *MockBlockCheckerMockRecorder {
	return m.recorder
}

// ChangeAllowed mocks base method
func (m *MockBlockChecker) ChangeAllowed() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeAllowed")
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangeAllowed indicates an expected call of ChangeAllowed
func (mr *MockBlockCheckerMockRecorder) ChangeAllowed() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeAllowed", reflect.TypeOf((*MockBlockChecker)(nil).ChangeAllowed))
}

// RemoveAllowed mocks base method
func (m *MockBlockChecker) RemoveAllowed() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveAllowed")
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveAllowed indicates an expected call of RemoveAllowed
func (mr *MockBlockCheckerMockRecorder) RemoveAllowed() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveAllowed", reflect.TypeOf((*MockBlockChecker)(nil).RemoveAllowed))
}

// MockMachine is a mock of Machine interface
type MockMachine struct {
	ctrl     *gomock.Controller
	recorder *MockMachineMockRecorder
}

// MockMachineMockRecorder is the mock recorder for MockMachine
type MockMachineMockRecorder struct {
	mock *MockMachine
}

// NewMockMachine creates a new mock instance
func NewMockMachine(ctrl *gomock.Controller) *MockMachine {
	mock := &MockMachine{ctrl: ctrl}
	mock.recorder = &MockMachineMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMachine) EXPECT() *MockMachineMockRecorder {
	return m.recorder
}

// AllSpaces mocks base method
func (m *MockMachine) AllSpaces() (set.Strings, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllSpaces")
	ret0, _ := ret[0].(set.Strings)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllSpaces indicates an expected call of AllSpaces
func (mr *MockMachineMockRecorder) AllSpaces() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllSpaces", reflect.TypeOf((*MockMachine)(nil).AllSpaces))
}
