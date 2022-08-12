// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	types "github.com/cosmos/cosmos-sdk/types"
	mock "github.com/stretchr/testify/mock"
)

// EarnHooks is an autogenerated mock type for the EarnHooks type
type EarnHooks struct {
	mock.Mock
}

// AfterVaultDepositCreated provides a mock function with given fields: ctx, vaultDenom, depositor, sharedOwned
func (_m *EarnHooks) AfterVaultDepositCreated(ctx types.Context, vaultDenom string, depositor types.AccAddress, sharedOwned types.Dec) {
	_m.Called(ctx, vaultDenom, depositor, sharedOwned)
}

// BeforeVaultDepositModified provides a mock function with given fields: ctx, vaultDenom, depositor, sharedOwned
func (_m *EarnHooks) BeforeVaultDepositModified(ctx types.Context, vaultDenom string, depositor types.AccAddress, sharedOwned types.Dec) {
	_m.Called(ctx, vaultDenom, depositor, sharedOwned)
}

type mockConstructorTestingTNewEarnHooks interface {
	mock.TestingT
	Cleanup(func())
}

// NewEarnHooks creates a new instance of EarnHooks. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewEarnHooks(t mockConstructorTestingTNewEarnHooks) *EarnHooks {
	mock := &EarnHooks{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
