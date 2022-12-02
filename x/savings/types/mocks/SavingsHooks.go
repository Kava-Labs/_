// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	savingstypes "github.com/kava-labs/kava/x/savings/types"
	mock "github.com/stretchr/testify/mock"

	types "github.com/cosmos/cosmos-sdk/types"
)

// SavingsHooks is an autogenerated mock type for the SavingsHooks type
type SavingsHooks struct {
	mock.Mock
}

// AfterSavingsDepositCreated provides a mock function with given fields: ctx, deposit
func (_m *SavingsHooks) AfterSavingsDepositCreated(ctx types.Context, deposit savingstypes.Deposit) {
	_m.Called(ctx, deposit)
}

// BeforeSavingsDepositModified provides a mock function with given fields: ctx, deposit, incomingDenoms
func (_m *SavingsHooks) BeforeSavingsDepositModified(ctx types.Context, deposit savingstypes.Deposit, incomingDenoms []string) {
	_m.Called(ctx, deposit, incomingDenoms)
}

type mockConstructorTestingTNewSavingsHooks interface {
	mock.TestingT
	Cleanup(func())
}

// NewSavingsHooks creates a new instance of SavingsHooks. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewSavingsHooks(t mockConstructorTestingTNewSavingsHooks) *SavingsHooks {
	mock := &SavingsHooks{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}