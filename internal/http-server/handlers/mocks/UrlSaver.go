// Code generated by mockery v2.45.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// UrlSaver is an autogenerated mock type for the UrlSaver type
type UrlSaver struct {
	mock.Mock
}

// AliasExists provides a mock function with given fields: alias
func (_m *UrlSaver) AliasExists(alias string) (bool, error) {
	ret := _m.Called(alias)

	if len(ret) == 0 {
		panic("no return value specified for AliasExists")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (bool, error)); ok {
		return rf(alias)
	}
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(alias)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(alias)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SaveUrl provides a mock function with given fields: urlToSave, alias
func (_m *UrlSaver) SaveUrl(urlToSave string, alias string) (int64, error) {
	ret := _m.Called(urlToSave, alias)

	if len(ret) == 0 {
		panic("no return value specified for SaveUrl")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (int64, error)); ok {
		return rf(urlToSave, alias)
	}
	if rf, ok := ret.Get(0).(func(string, string) int64); ok {
		r0 = rf(urlToSave, alias)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(urlToSave, alias)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewUrlSaver creates a new instance of UrlSaver. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUrlSaver(t interface {
	mock.TestingT
	Cleanup(func())
}) *UrlSaver {
	mock := &UrlSaver{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
