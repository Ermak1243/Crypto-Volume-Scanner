// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Logger is an autogenerated mock type for the Logger type
type Logger struct {
	mock.Mock
}

// DPanic provides a mock function with given fields: args
func (_m *Logger) DPanic(args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// DPanicf provides a mock function with given fields: template, args
func (_m *Logger) DPanicf(template string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, template)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// Debug provides a mock function with given fields: args
func (_m *Logger) Debug(args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// Debugf provides a mock function with given fields: template, args
func (_m *Logger) Debugf(template string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, template)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// Error provides a mock function with given fields: args
func (_m *Logger) Error(args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// Errorf provides a mock function with given fields: template, args
func (_m *Logger) Errorf(template string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, template)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// Fatal provides a mock function with given fields: args
func (_m *Logger) Fatal(args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// Fatalf provides a mock function with given fields: template, args
func (_m *Logger) Fatalf(template string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, template)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// Info provides a mock function with given fields: args
func (_m *Logger) Info(args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// Infof provides a mock function with given fields: template, args
func (_m *Logger) Infof(template string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, template)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// InitLogger provides a mock function with given fields:
func (_m *Logger) InitLogger() {
	_m.Called()
}

// Panic provides a mock function with given fields: args
func (_m *Logger) Panic(args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// Panicf provides a mock function with given fields: template, args
func (_m *Logger) Panicf(template string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, template)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// Warn provides a mock function with given fields: args
func (_m *Logger) Warn(args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// Warnf provides a mock function with given fields: template, args
func (_m *Logger) Warnf(template string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, template)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

type mockConstructorTestingTNewLogger interface {
	mock.TestingT
	Cleanup(func())
}

// NewLogger creates a new instance of Logger. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewLogger(t mockConstructorTestingTNewLogger) *Logger {
	mock := &Logger{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
