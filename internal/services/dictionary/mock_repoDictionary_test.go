// Code generated by mockery v2.30.1. DO NOT EDIT.

package dictionary

import (
	context "context"

	uuid "github.com/google/uuid"
	mock "github.com/stretchr/testify/mock"
)

// mockRepoDictionary is an autogenerated mock type for the repoDictionary type
type mockRepoDictionary struct {
	mock.Mock
}

// AddWords provides a mock function with given fields: ctx, words
func (_m *mockRepoDictionary) AddWords(ctx context.Context, words []DictWord) ([]DictWord, error) {
	ret := _m.Called(ctx, words)

	var r0 []DictWord
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []DictWord) ([]DictWord, error)); ok {
		return rf(ctx, words)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []DictWord) []DictWord); ok {
		r0 = rf(ctx, words)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]DictWord)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []DictWord) error); ok {
		r1 = rf(ctx, words)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteWordByText provides a mock function with given fields: ctx, word
func (_m *mockRepoDictionary) DeleteWordByText(ctx context.Context, word *DictWord) error {
	ret := _m.Called(ctx, word)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *DictWord) error); ok {
		r0 = rf(ctx, word)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindWords provides a mock function with given fields: ctx, w
func (_m *mockRepoDictionary) FindWords(ctx context.Context, w *DictWord) ([]uuid.UUID, error) {
	ret := _m.Called(ctx, w)

	var r0 []uuid.UUID
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *DictWord) ([]uuid.UUID, error)); ok {
		return rf(ctx, w)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *DictWord) []uuid.UUID); ok {
		r0 = rf(ctx, w)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]uuid.UUID)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *DictWord) error); ok {
		r1 = rf(ctx, w)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDictionary provides a mock function with given fields: ctx, langCode, search, page, itemsPerPage
func (_m *mockRepoDictionary) GetDictionary(ctx context.Context, langCode string, search string, page int, itemsPerPage int) ([]DictWord, error) {
	ret := _m.Called(ctx, langCode, search, page, itemsPerPage)

	var r0 []DictWord
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, int, int) ([]DictWord, error)); ok {
		return rf(ctx, langCode, search, page, itemsPerPage)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, int, int) []DictWord); ok {
		r0 = rf(ctx, langCode, search, page, itemsPerPage)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]DictWord)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, int, int) error); ok {
		r1 = rf(ctx, langCode, search, page, itemsPerPage)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPronunciation provides a mock function with given fields: ctx, text, langCode
func (_m *mockRepoDictionary) GetPronunciation(ctx context.Context, text string, langCode string) (string, error) {
	ret := _m.Called(ctx, text, langCode)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (string, error)); ok {
		return rf(ctx, text, langCode)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) string); ok {
		r0 = rf(ctx, text, langCode)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, text, langCode)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRandomWord provides a mock function with given fields: ctx, langCode
func (_m *mockRepoDictionary) GetRandomWord(ctx context.Context, langCode string) (DictWord, error) {
	ret := _m.Called(ctx, langCode)

	var r0 DictWord
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (DictWord, error)); ok {
		return rf(ctx, langCode)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) DictWord); ok {
		r0 = rf(ctx, langCode)
	} else {
		r0 = ret.Get(0).(DictWord)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, langCode)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetWords provides a mock function with given fields: ctx, ids
func (_m *mockRepoDictionary) GetWords(ctx context.Context, ids []uuid.UUID) ([]DictWord, error) {
	ret := _m.Called(ctx, ids)

	var r0 []DictWord
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []uuid.UUID) ([]DictWord, error)); ok {
		return rf(ctx, ids)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []uuid.UUID) []DictWord); ok {
		r0 = rf(ctx, ids)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]DictWord)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []uuid.UUID) error); ok {
		r1 = rf(ctx, ids)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetWordsByText provides a mock function with given fields: ctx, words
func (_m *mockRepoDictionary) GetWordsByText(ctx context.Context, words []DictWord) ([]DictWord, error) {
	ret := _m.Called(ctx, words)

	var r0 []DictWord
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []DictWord) ([]DictWord, error)); ok {
		return rf(ctx, words)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []DictWord) []DictWord); ok {
		r0 = rf(ctx, words)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]DictWord)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []DictWord) error); ok {
		r1 = rf(ctx, words)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateWord provides a mock function with given fields: ctx, w
func (_m *mockRepoDictionary) UpdateWord(ctx context.Context, w *DictWord) error {
	ret := _m.Called(ctx, w)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *DictWord) error); ok {
		r0 = rf(ctx, w)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// newMockRepoDictionary creates a new instance of mockRepoDictionary. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockRepoDictionary(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockRepoDictionary {
	mock := &mockRepoDictionary{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
