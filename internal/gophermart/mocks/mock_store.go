package mocks

import (
	context "context"
	storage "github.com/YaNeAndrey/ya-gophermart/internal/gophermart/storage"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockStorageRepo is a mock of StorageRepo interface.
type MockStorageRepo struct {
	ctrl     *gomock.Controller
	recorder *MockStorageRepoMockRecorder
}

// MockStorageRepoMockRecorder is the mock recorder for MockStorageRepo.
type MockStorageRepoMockRecorder struct {
	mock *MockStorageRepo
}

// NewMockStorageRepo creates a new mock instance.
func NewMockStorageRepo(ctrl *gomock.Controller) *MockStorageRepo {
	mock := &MockStorageRepo{ctrl: ctrl}
	mock.recorder = &MockStorageRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorageRepo) EXPECT() *MockStorageRepoMockRecorder {
	return m.recorder
}

// AddNewOrder mocks base method.
func (m *MockStorageRepo) AddNewOrder(arg0 context.Context, arg1, arg2 string) (*storage.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddNewOrder", arg0, arg1, arg2)
	ret0, _ := ret[0].(*storage.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddNewOrder indicates an expected call of AddNewOrder.
func (mr *MockStorageRepoMockRecorder) AddNewOrder(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddNewOrder", reflect.TypeOf((*MockStorageRepo)(nil).AddNewOrder), arg0, arg1, arg2)
}

// AddNewUser mocks base method.
func (m *MockStorageRepo) AddNewUser(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddNewUser", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddNewUser indicates an expected call of AddNewUser.
func (mr *MockStorageRepoMockRecorder) AddNewUser(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddNewUser", reflect.TypeOf((*MockStorageRepo)(nil).AddNewUser), arg0, arg1, arg2)
}

// CheckUserPassword mocks base method.
func (m *MockStorageRepo) CheckUserPassword(arg0 context.Context, arg1, arg2 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckUserPassword", arg0, arg1, arg2)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckUserPassword indicates an expected call of CheckUserPassword.
func (mr *MockStorageRepoMockRecorder) CheckUserPassword(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckUserPassword", reflect.TypeOf((*MockStorageRepo)(nil).CheckUserPassword), arg0, arg1, arg2)
}

// DoRebiting mocks base method.
func (m *MockStorageRepo) DoRebiting(arg0 context.Context, arg1, arg2 string, arg3 float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DoRebiting", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// DoRebiting indicates an expected call of DoRebiting.
func (mr *MockStorageRepoMockRecorder) DoRebiting(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DoRebiting", reflect.TypeOf((*MockStorageRepo)(nil).DoRebiting), arg0, arg1, arg2, arg3)
}

// GetAllNotProcessedOrders mocks base method.
func (m *MockStorageRepo) GetAllNotProcessedOrders(arg0 context.Context) (*[]storage.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllNotProcessedOrders", arg0)
	ret0, _ := ret[0].(*[]storage.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllNotProcessedOrders indicates an expected call of GetAllNotProcessedOrders.
func (mr *MockStorageRepoMockRecorder) GetAllNotProcessedOrders(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllNotProcessedOrders", reflect.TypeOf((*MockStorageRepo)(nil).GetAllNotProcessedOrders), arg0)
}

// GetUserBalance mocks base method.
func (m *MockStorageRepo) GetUserBalance(arg0 context.Context, arg1 string) (*storage.Balance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserBalance", arg0, arg1)
	ret0, _ := ret[0].(*storage.Balance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserBalance indicates an expected call of GetUserBalance.
func (mr *MockStorageRepoMockRecorder) GetUserBalance(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserBalance", reflect.TypeOf((*MockStorageRepo)(nil).GetUserBalance), arg0, arg1)
}

// GetUserOrders mocks base method.
func (m *MockStorageRepo) GetUserOrders(arg0 context.Context, arg1 string) (*[]storage.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserOrders", arg0, arg1)
	ret0, _ := ret[0].(*[]storage.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserOrders indicates an expected call of GetUserOrders.
func (mr *MockStorageRepoMockRecorder) GetUserOrders(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserOrders", reflect.TypeOf((*MockStorageRepo)(nil).GetUserOrders), arg0, arg1)
}

// GetUserWithdrawals mocks base method.
func (m *MockStorageRepo) GetUserWithdrawals(arg0 context.Context, arg1 string) (*[]storage.Withdrawal, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserWithdrawals", arg0, arg1)
	ret0, _ := ret[0].(*[]storage.Withdrawal)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserWithdrawals indicates an expected call of GetUserWithdrawals.
func (mr *MockStorageRepoMockRecorder) GetUserWithdrawals(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserWithdrawals", reflect.TypeOf((*MockStorageRepo)(nil).GetUserWithdrawals), arg0, arg1)
}

// UpdateBalance mocks base method.
func (m *MockStorageRepo) UpdateBalance(arg0 context.Context, arg1 storage.Order) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBalance", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateBalance indicates an expected call of UpdateBalance.
func (mr *MockStorageRepoMockRecorder) UpdateBalance(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBalance", reflect.TypeOf((*MockStorageRepo)(nil).UpdateBalance), arg0, arg1)
}

// UpdateOrder mocks base method.
func (m *MockStorageRepo) UpdateOrder(arg0 context.Context, arg1 storage.Order) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateOrder", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateOrder indicates an expected call of UpdateOrder.
func (mr *MockStorageRepoMockRecorder) UpdateOrder(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOrder", reflect.TypeOf((*MockStorageRepo)(nil).UpdateOrder), arg0, arg1)
}
