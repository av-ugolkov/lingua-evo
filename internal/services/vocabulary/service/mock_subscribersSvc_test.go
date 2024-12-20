// Code generated by mockery v2.30.1. DO NOT EDIT.

package service

import (
	context "context"

	uuid "github.com/google/uuid"
	mock "github.com/stretchr/testify/mock"
)

// mockSubscribersSvc is an autogenerated mock type for the subscribersSvc type
type mockSubscribersSvc struct {
	mock.Mock
}

// Check provides a mock function with given fields: ctx, uid, subID
func (_m *mockSubscribersSvc) Check(ctx context.Context, uid uuid.UUID, subID uuid.UUID) (bool, error) {
	ret := _m.Called(ctx, uid, subID)

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID) (bool, error)); ok {
		return rf(ctx, uid, subID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID) bool); ok {
		r0 = rf(ctx, uid, subID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, uuid.UUID) error); ok {
		r1 = rf(ctx, uid, subID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// newMockSubscribersSvc creates a new instance of mockSubscribersSvc. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockSubscribersSvc(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockSubscribersSvc {
	mock := &mockSubscribersSvc{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
