// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import io "io"
import jlaffayeftp "github.com/jlaffaye/ftp"
import mock "github.com/stretchr/testify/mock"

// ServerConnexion is an autogenerated mock type for the ServerConnexion type
type ServerConnexion struct {
	mock.Mock
}

// Delete provides a mock function with given fields: oath
func (_m *ServerConnexion) Delete(oath string) error {
	ret := _m.Called(oath)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(oath)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FileSize provides a mock function with given fields: path
func (_m *ServerConnexion) FileSize(path string) (int64, error) {
	ret := _m.Called(path)

	var r0 int64
	if rf, ok := ret.Get(0).(func(string) int64); ok {
		r0 = rf(path)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: sPath
func (_m *ServerConnexion) List(sPath string) ([]*jlaffayeftp.Entry, error) {
	ret := _m.Called(sPath)

	var r0 []*jlaffayeftp.Entry
	if rf, ok := ret.Get(0).(func(string) []*jlaffayeftp.Entry); ok {
		r0 = rf(sPath)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*jlaffayeftp.Entry)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(sPath)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Login provides a mock function with given fields: sUsername, sPwd
func (_m *ServerConnexion) Login(sUsername string, sPwd string) error {
	ret := _m.Called(sUsername, sPwd)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(sUsername, sPwd)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Quit provides a mock function with given fields:
func (_m *ServerConnexion) Quit() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Retr provides a mock function with given fields: path
func (_m *ServerConnexion) Retr(path string) (*jlaffayeftp.Response, error) {
	ret := _m.Called(path)

	var r0 *jlaffayeftp.Response
	if rf, ok := ret.Get(0).(func(string) *jlaffayeftp.Response); ok {
		r0 = rf(path)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*jlaffayeftp.Response)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Stor provides a mock function with given fields: path, oReader
func (_m *ServerConnexion) Stor(path string, oReader io.Reader) error {
	ret := _m.Called(path, oReader)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, io.Reader) error); ok {
		r0 = rf(path, oReader)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
