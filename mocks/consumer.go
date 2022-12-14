// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mqconnector "github.com/Velnbur/mq-connector"
	mock "github.com/stretchr/testify/mock"
)

// Consumer is an autogenerated mock type for the Consumer type
type Consumer struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *Consumer) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Run provides a mock function with given fields: ctx
func (_m *Consumer) Run(ctx context.Context) {
	_m.Called(ctx)
}

// SetContexters provides a mock function with given fields: ctxers
func (_m *Consumer) SetContexters(ctxers ...mqconnector.ContextFunc) mqconnector.Consumer {
	_va := make([]interface{}, len(ctxers))
	for _i := range ctxers {
		_va[_i] = ctxers[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 mqconnector.Consumer
	if rf, ok := ret.Get(0).(func(...mqconnector.ContextFunc) mqconnector.Consumer); ok {
		r0 = rf(ctxers...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(mqconnector.Consumer)
		}
	}

	return r0
}

// SetHandler provides a mock function with given fields: handler
func (_m *Consumer) SetHandler(handler mqconnector.Handler) mqconnector.Consumer {
	ret := _m.Called(handler)

	var r0 mqconnector.Consumer
	if rf, ok := ret.Get(0).(func(mqconnector.Handler) mqconnector.Consumer); ok {
		r0 = rf(handler)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(mqconnector.Consumer)
		}
	}

	return r0
}

type mockConstructorTestingTNewConsumer interface {
	mock.TestingT
	Cleanup(func())
}

// NewConsumer creates a new instance of Consumer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewConsumer(t mockConstructorTestingTNewConsumer) *Consumer {
	mock := &Consumer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
