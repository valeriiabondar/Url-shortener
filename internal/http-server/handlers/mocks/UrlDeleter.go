// Code generated by mockery v2.46.3. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// UrlDeleter is an autogenerated mock type for the UrlDeleter type
type UrlDeleter struct {
	mock.Mock
}

// DeleteUrl provides a mock function with given fields: alias
func (_m *UrlDeleter) DeleteUrl(alias string) error {
	ret := _m.Called(alias)

	if len(ret) == 0 {
		panic("no return value specified for DeleteUrl")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(alias)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewUrlDeleter creates a new instance of UrlDeleter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUrlDeleter(t interface {
	mock.TestingT
	Cleanup(func())
}) *UrlDeleter {
	mock := &UrlDeleter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}