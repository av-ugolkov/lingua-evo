// Code generated by mockery v2.30.1. DO NOT EDIT.

package service

import (
	context "context"

	user "github.com/av-ugolkov/lingua-evo/internal/services/user"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// mockUserRepo is an autogenerated mock type for the userRepo type
type mockUserRepo struct {
	mock.Mock
}

// AddGoogleUser provides a mock function with given fields: ctx, userCreate
func (_m *mockUserRepo) AddGoogleUser(ctx context.Context, userCreate user.GoogleUser) (uuid.UUID, error) {
	ret := _m.Called(ctx, userCreate)

	var r0 uuid.UUID
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, user.GoogleUser) (uuid.UUID, error)); ok {
		return rf(ctx, userCreate)
	}
	if rf, ok := ret.Get(0).(func(context.Context, user.GoogleUser) uuid.UUID); ok {
		r0 = rf(ctx, userCreate)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(uuid.UUID)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, user.GoogleUser) error); ok {
		r1 = rf(ctx, userCreate)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AddUser provides a mock function with given fields: ctx, u, pswHash
func (_m *mockUserRepo) AddUser(ctx context.Context, u *user.User, pswHash string) (uuid.UUID, error) {
	ret := _m.Called(ctx, u, pswHash)

	var r0 uuid.UUID
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *user.User, string) (uuid.UUID, error)); ok {
		return rf(ctx, u, pswHash)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *user.User, string) uuid.UUID); ok {
		r0 = rf(ctx, u, pswHash)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(uuid.UUID)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *user.User, string) error); ok {
		r1 = rf(ctx, u, pswHash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AddUserData provides a mock function with given fields: ctx, data
func (_m *mockUserRepo) AddUserData(ctx context.Context, data user.UserData) error {
	ret := _m.Called(ctx, data)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, user.UserData) error); ok {
		r0 = rf(ctx, data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AddUserNewsletters provides a mock function with given fields: ctx, data
func (_m *mockUserRepo) AddUserNewsletters(ctx context.Context, data user.UserNewsletters) error {
	ret := _m.Called(ctx, data)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, user.UserNewsletters) error); ok {
		r0 = rf(ctx, data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetPswHash provides a mock function with given fields: ctx, uid
func (_m *mockUserRepo) GetPswHash(ctx context.Context, uid uuid.UUID) (string, error) {
	ret := _m.Called(ctx, uid)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (string, error)); ok {
		return rf(ctx, uid)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) string); ok {
		r0 = rf(ctx, uid)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, uid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserByEmail provides a mock function with given fields: ctx, email
func (_m *mockUserRepo) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	ret := _m.Called(ctx, email)

	var r0 *user.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*user.User, error)); ok {
		return rf(ctx, email)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *user.User); ok {
		r0 = rf(ctx, email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*user.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserByGoogleID provides a mock function with given fields: ctx, email
func (_m *mockUserRepo) GetUserByGoogleID(ctx context.Context, email string) (*user.User, error) {
	ret := _m.Called(ctx, email)

	var r0 *user.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*user.User, error)); ok {
		return rf(ctx, email)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *user.User); ok {
		r0 = rf(ctx, email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*user.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserByID provides a mock function with given fields: ctx, uid
func (_m *mockUserRepo) GetUserByID(ctx context.Context, uid uuid.UUID) (*user.User, error) {
	ret := _m.Called(ctx, uid)

	var r0 *user.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (*user.User, error)); ok {
		return rf(ctx, uid)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) *user.User); ok {
		r0 = rf(ctx, uid)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*user.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, uid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserByNickname provides a mock function with given fields: ctx, name
func (_m *mockUserRepo) GetUserByNickname(ctx context.Context, name string) (*user.User, error) {
	ret := _m.Called(ctx, name)

	var r0 *user.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*user.User, error)); ok {
		return rf(ctx, name)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *user.User); ok {
		r0 = rf(ctx, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*user.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserData provides a mock function with given fields: ctx, uid
func (_m *mockUserRepo) GetUserData(ctx context.Context, uid uuid.UUID) (*user.UserData, error) {
	ret := _m.Called(ctx, uid)

	var r0 *user.UserData
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (*user.UserData, error)); ok {
		return rf(ctx, uid)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) *user.UserData); ok {
		r0 = rf(ctx, uid)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*user.UserData)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, uid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserSubscriptions provides a mock function with given fields: ctx, uid
func (_m *mockUserRepo) GetUserSubscriptions(ctx context.Context, uid uuid.UUID) ([]user.Subscriptions, error) {
	ret := _m.Called(ctx, uid)

	var r0 []user.Subscriptions
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) ([]user.Subscriptions, error)); ok {
		return rf(ctx, uid)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) []user.Subscriptions); ok {
		r0 = rf(ctx, uid)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]user.Subscriptions)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, uid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUsers provides a mock function with given fields: ctx, page, perPage, sort, order, search
func (_m *mockUserRepo) GetUsers(ctx context.Context, page int, perPage int, sort int, order int, search string) ([]user.User, int, error) {
	ret := _m.Called(ctx, page, perPage, sort, order, search)

	var r0 []user.User
	var r1 int
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, int, int, int, int, string) ([]user.User, int, error)); ok {
		return rf(ctx, page, perPage, sort, order, search)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int, int, int, int, string) []user.User); ok {
		r0 = rf(ctx, page, perPage, sort, order, search)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]user.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int, int, int, int, string) int); ok {
		r1 = rf(ctx, page, perPage, sort, order, search)
	} else {
		r1 = ret.Get(1).(int)
	}

	if rf, ok := ret.Get(2).(func(context.Context, int, int, int, int, string) error); ok {
		r2 = rf(ctx, page, perPage, sort, order, search)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// RemoveUser provides a mock function with given fields: ctx, uid
func (_m *mockUserRepo) RemoveUser(ctx context.Context, uid uuid.UUID) error {
	ret := _m.Called(ctx, uid)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) error); ok {
		r0 = rf(ctx, uid)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateEmail provides a mock function with given fields: ctx, uid, newEmail
func (_m *mockUserRepo) UpdateEmail(ctx context.Context, uid uuid.UUID, newEmail string) error {
	ret := _m.Called(ctx, uid, newEmail)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, string) error); ok {
		r0 = rf(ctx, uid, newEmail)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateNickname provides a mock function with given fields: ctx, uid, newNickname
func (_m *mockUserRepo) UpdateNickname(ctx context.Context, uid uuid.UUID, newNickname string) error {
	ret := _m.Called(ctx, uid, newNickname)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, string) error); ok {
		r0 = rf(ctx, uid, newNickname)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdatePsw provides a mock function with given fields: ctx, uid, hashPsw
func (_m *mockUserRepo) UpdatePsw(ctx context.Context, uid uuid.UUID, hashPsw string) error {
	ret := _m.Called(ctx, uid, hashPsw)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, string) error); ok {
		r0 = rf(ctx, uid, hashPsw)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateVisitedAt provides a mock function with given fields: ctx, uid
func (_m *mockUserRepo) UpdateVisitedAt(ctx context.Context, uid uuid.UUID) error {
	ret := _m.Called(ctx, uid)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) error); ok {
		r0 = rf(ctx, uid)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// newMockUserRepo creates a new instance of mockUserRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockUserRepo(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockUserRepo {
	mock := &mockUserRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
