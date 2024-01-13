// Code generated by MockGen. DO NOT EDIT.
// Source: usecase/authtoken/interface.go
//
// Generated by this command:
//
//	mockgen -source usecase/authtoken/interface.go -typed -destination internal/mocks/ucmocks/authtoken_usecase.go -mock_names=Usecase=MockAuthTokenUsecase -package ucmocks
//

// Package ucmocks is a generated GoMock package.
package ucmocks

import (
	context "context"
	reflect "reflect"

	authtoken "github.com/psi59/payhere-assignment/usecase/authtoken"
	gomock "go.uber.org/mock/gomock"
)

// MockAuthTokenUsecase is a mock of Usecase interface.
type MockAuthTokenUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockAuthTokenUsecaseMockRecorder
}

// MockAuthTokenUsecaseMockRecorder is the mock recorder for MockAuthTokenUsecase.
type MockAuthTokenUsecaseMockRecorder struct {
	mock *MockAuthTokenUsecase
}

// NewMockAuthTokenUsecase creates a new mock instance.
func NewMockAuthTokenUsecase(ctrl *gomock.Controller) *MockAuthTokenUsecase {
	mock := &MockAuthTokenUsecase{ctrl: ctrl}
	mock.recorder = &MockAuthTokenUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthTokenUsecase) EXPECT() *MockAuthTokenUsecaseMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockAuthTokenUsecase) Create(c context.Context, input *authtoken.CreateInput) (*authtoken.CreateOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", c, input)
	ret0, _ := ret[0].(*authtoken.CreateOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockAuthTokenUsecaseMockRecorder) Create(c, input any) *MockAuthTokenUsecaseCreateCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockAuthTokenUsecase)(nil).Create), c, input)
	return &MockAuthTokenUsecaseCreateCall{Call: call}
}

// MockAuthTokenUsecaseCreateCall wrap *gomock.Call
type MockAuthTokenUsecaseCreateCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c_2 *MockAuthTokenUsecaseCreateCall) Return(arg0 *authtoken.CreateOutput, arg1 error) *MockAuthTokenUsecaseCreateCall {
	c_2.Call = c_2.Call.Return(arg0, arg1)
	return c_2
}

// Do rewrite *gomock.Call.Do
func (c_2 *MockAuthTokenUsecaseCreateCall) Do(f func(context.Context, *authtoken.CreateInput) (*authtoken.CreateOutput, error)) *MockAuthTokenUsecaseCreateCall {
	c_2.Call = c_2.Call.Do(f)
	return c_2
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c_2 *MockAuthTokenUsecaseCreateCall) DoAndReturn(f func(context.Context, *authtoken.CreateInput) (*authtoken.CreateOutput, error)) *MockAuthTokenUsecaseCreateCall {
	c_2.Call = c_2.Call.DoAndReturn(f)
	return c_2
}

// RegisterToBlacklist mocks base method.
func (m *MockAuthTokenUsecase) RegisterToBlacklist(c context.Context, input *authtoken.RegisterToBlacklistInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterToBlacklist", c, input)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterToBlacklist indicates an expected call of RegisterToBlacklist.
func (mr *MockAuthTokenUsecaseMockRecorder) RegisterToBlacklist(c, input any) *MockAuthTokenUsecaseRegisterToBlacklistCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterToBlacklist", reflect.TypeOf((*MockAuthTokenUsecase)(nil).RegisterToBlacklist), c, input)
	return &MockAuthTokenUsecaseRegisterToBlacklistCall{Call: call}
}

// MockAuthTokenUsecaseRegisterToBlacklistCall wrap *gomock.Call
type MockAuthTokenUsecaseRegisterToBlacklistCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c_2 *MockAuthTokenUsecaseRegisterToBlacklistCall) Return(arg0 error) *MockAuthTokenUsecaseRegisterToBlacklistCall {
	c_2.Call = c_2.Call.Return(arg0)
	return c_2
}

// Do rewrite *gomock.Call.Do
func (c_2 *MockAuthTokenUsecaseRegisterToBlacklistCall) Do(f func(context.Context, *authtoken.RegisterToBlacklistInput) error) *MockAuthTokenUsecaseRegisterToBlacklistCall {
	c_2.Call = c_2.Call.Do(f)
	return c_2
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c_2 *MockAuthTokenUsecaseRegisterToBlacklistCall) DoAndReturn(f func(context.Context, *authtoken.RegisterToBlacklistInput) error) *MockAuthTokenUsecaseRegisterToBlacklistCall {
	c_2.Call = c_2.Call.DoAndReturn(f)
	return c_2
}

// Verify mocks base method.
func (m *MockAuthTokenUsecase) Verify(c context.Context, input *authtoken.VerifyInput) (*authtoken.VerifyOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Verify", c, input)
	ret0, _ := ret[0].(*authtoken.VerifyOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Verify indicates an expected call of Verify.
func (mr *MockAuthTokenUsecaseMockRecorder) Verify(c, input any) *MockAuthTokenUsecaseVerifyCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Verify", reflect.TypeOf((*MockAuthTokenUsecase)(nil).Verify), c, input)
	return &MockAuthTokenUsecaseVerifyCall{Call: call}
}

// MockAuthTokenUsecaseVerifyCall wrap *gomock.Call
type MockAuthTokenUsecaseVerifyCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c_2 *MockAuthTokenUsecaseVerifyCall) Return(arg0 *authtoken.VerifyOutput, arg1 error) *MockAuthTokenUsecaseVerifyCall {
	c_2.Call = c_2.Call.Return(arg0, arg1)
	return c_2
}

// Do rewrite *gomock.Call.Do
func (c_2 *MockAuthTokenUsecaseVerifyCall) Do(f func(context.Context, *authtoken.VerifyInput) (*authtoken.VerifyOutput, error)) *MockAuthTokenUsecaseVerifyCall {
	c_2.Call = c_2.Call.Do(f)
	return c_2
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c_2 *MockAuthTokenUsecaseVerifyCall) DoAndReturn(f func(context.Context, *authtoken.VerifyInput) (*authtoken.VerifyOutput, error)) *MockAuthTokenUsecaseVerifyCall {
	c_2.Call = c_2.Call.DoAndReturn(f)
	return c_2
}