// Code generated by mockery v1.0.0
package ammomock

import http "net/http"
import mock "github.com/stretchr/testify/mock"
import netsample "github.com/yandex/pandora/core/aggregator/netsample"

// Ammo is an autogenerated mock type for the Ammo type
type Ammo struct {
	mock.Mock
}

// Id provides a mock function with given fields:
func (_m *Ammo) Id() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// Request provides a mock function with given fields:
func (_m *Ammo) Request() (*http.Request, *netsample.Sample) {
	ret := _m.Called()

	var r0 *http.Request
	if rf, ok := ret.Get(0).(func() *http.Request); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*http.Request)
		}
	}

	var r1 *netsample.Sample
	if rf, ok := ret.Get(1).(func() *netsample.Sample); ok {
		r1 = rf()
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*netsample.Sample)
		}
	}

	return r0, r1
}
