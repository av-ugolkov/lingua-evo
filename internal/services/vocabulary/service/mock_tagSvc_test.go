// Code generated by mockery v2.30.1. DO NOT EDIT.

package service

import (
	context "context"

	tag "github.com/av-ugolkov/lingua-evo/internal/services/tag"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// mockTagSvc is an autogenerated mock type for the tagSvc type
type mockTagSvc struct {
	mock.Mock
}

// AddTags provides a mock function with given fields: ctx, tags
func (_m *mockTagSvc) AddTags(ctx context.Context, tags []tag.Tag) ([]uuid.UUID, error) {
	ret := _m.Called(ctx, tags)

	var r0 []uuid.UUID
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []tag.Tag) ([]uuid.UUID, error)); ok {
		return rf(ctx, tags)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []tag.Tag) []uuid.UUID); ok {
		r0 = rf(ctx, tags)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]uuid.UUID)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []tag.Tag) error); ok {
		r1 = rf(ctx, tags)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// newMockTagSvc creates a new instance of mockTagSvc. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockTagSvc(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockTagSvc {
	mock := &mockTagSvc{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
