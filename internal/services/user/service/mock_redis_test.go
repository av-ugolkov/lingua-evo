// Code generated by mockery v2.30.1. DO NOT EDIT.

package service

import (
	context "context"
	time "time"

	mock "github.com/stretchr/testify/mock"
)

// mockRedis is an autogenerated mock type for the redis type
type mockRedis struct {
	mock.Mock
}

// Get provides a mock function with given fields: ctx, key
func (_m *mockRedis) Get(ctx context.Context, key string) (string, error) {
	ret := _m.Called(ctx, key)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (string, error)); ok {
		return rf(ctx, key)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, key)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAccountCode provides a mock function with given fields: ctx, email
func (_m *mockRedis) GetAccountCode(ctx context.Context, email string) (int, error) {
	ret := _m.Called(ctx, email)

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (int, error)); ok {
		return rf(ctx, email)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) int); ok {
		r0 = rf(ctx, email)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTTL provides a mock function with given fields: ct, key
func (_m *mockRedis) GetTTL(ct context.Context, key string) (time.Duration, error) {
	ret := _m.Called(ct, key)

	var r0 time.Duration
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (time.Duration, error)); ok {
		return rf(ct, key)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) time.Duration); ok {
		r0 = rf(ct, key)
	} else {
		r0 = ret.Get(0).(time.Duration)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ct, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetNX provides a mock function with given fields: ctx, key, value, expiration
func (_m *mockRedis) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	ret := _m.Called(ctx, key, value, expiration)

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, interface{}, time.Duration) (bool, error)); ok {
		return rf(ctx, key, value, expiration)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, interface{}, time.Duration) bool); ok {
		r0 = rf(ctx, key, value, expiration)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, interface{}, time.Duration) error); ok {
		r1 = rf(ctx, key, value, expiration)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// newMockRedis creates a new instance of mockRedis. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockRedis(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockRedis {
	mock := &mockRedis{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
