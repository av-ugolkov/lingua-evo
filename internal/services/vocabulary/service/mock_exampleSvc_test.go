// Code generated by mockery v2.30.1. DO NOT EDIT.

package service

import (
	context "context"

	example "github.com/av-ugolkov/lingua-evo/internal/services/example"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// mockExampleSvc is an autogenerated mock type for the exampleSvc type
type mockExampleSvc struct {
	mock.Mock
}

// AddExamples provides a mock function with given fields: ctx, examples, langCode
func (_m *mockExampleSvc) AddExamples(ctx context.Context, examples []example.Example, langCode string) ([]uuid.UUID, error) {
	ret := _m.Called(ctx, examples, langCode)

	var r0 []uuid.UUID
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []example.Example, string) ([]uuid.UUID, error)); ok {
		return rf(ctx, examples, langCode)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []example.Example, string) []uuid.UUID); ok {
		r0 = rf(ctx, examples, langCode)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]uuid.UUID)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []example.Example, string) error); ok {
		r1 = rf(ctx, examples, langCode)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetExamples provides a mock function with given fields: ctx, exampleIDs
func (_m *mockExampleSvc) GetExamples(ctx context.Context, exampleIDs []uuid.UUID) ([]example.Example, error) {
	ret := _m.Called(ctx, exampleIDs)

	var r0 []example.Example
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []uuid.UUID) ([]example.Example, error)); ok {
		return rf(ctx, exampleIDs)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []uuid.UUID) []example.Example); ok {
		r0 = rf(ctx, exampleIDs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]example.Example)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []uuid.UUID) error); ok {
		r1 = rf(ctx, exampleIDs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// newMockExampleSvc creates a new instance of mockExampleSvc. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockExampleSvc(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockExampleSvc {
	mock := &mockExampleSvc{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
