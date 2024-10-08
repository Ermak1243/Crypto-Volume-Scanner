// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	context "context"
	models "cvs/internal/models"

	mock "github.com/stretchr/testify/mock"
)

// UserPairsRepository is an autogenerated mock type for the UserPairsRepository type
type UserPairsRepository struct {
	mock.Mock
}

// Add provides a mock function with given fields: ctx, pairData
func (_m *UserPairsRepository) Add(ctx context.Context, pairData models.UserPairs) error {
	ret := _m.Called(ctx, pairData)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, models.UserPairs) error); ok {
		r0 = rf(ctx, pairData)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeletePair provides a mock function with given fields: ctx, pairData
func (_m *UserPairsRepository) DeletePair(ctx context.Context, pairData models.UserPairs) error {
	ret := _m.Called(ctx, pairData)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, models.UserPairs) error); ok {
		r0 = rf(ctx, pairData)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAllUserPairs provides a mock function with given fields: ctx, userID
func (_m *UserPairsRepository) GetAllUserPairs(ctx context.Context, userID int) ([]models.UserPairs, error) {
	ret := _m.Called(ctx, userID)

	var r0 []models.UserPairs
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) ([]models.UserPairs, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) []models.UserPairs); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.UserPairs)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPairsByExchange provides a mock function with given fields: ctx, exchange
func (_m *UserPairsRepository) GetPairsByExchange(ctx context.Context, exchange string) ([]string, error) {
	ret := _m.Called(ctx, exchange)

	var r0 []string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]string, error)); ok {
		return rf(ctx, exchange)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []string); ok {
		r0 = rf(ctx, exchange)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, exchange)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateExactValue provides a mock function with given fields: ctx, pairData
func (_m *UserPairsRepository) UpdateExactValue(ctx context.Context, pairData models.UserPairs) error {
	ret := _m.Called(ctx, pairData)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, models.UserPairs) error); ok {
		r0 = rf(ctx, pairData)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewUserPairsRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewUserPairsRepository creates a new instance of UserPairsRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUserPairsRepository(t mockConstructorTestingTNewUserPairsRepository) *UserPairsRepository {
	mock := &UserPairsRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
