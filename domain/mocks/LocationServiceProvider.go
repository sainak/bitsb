// Code generated by mockery v2.18.0. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/sainak/bitsb/domain"
	mock "github.com/stretchr/testify/mock"

	repo "github.com/sainak/bitsb/utils/repo"
)

// LocationServiceProvider is an autogenerated mock type for the LocationServiceProvider type
type LocationServiceProvider struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, location
func (_m *LocationServiceProvider) Create(ctx context.Context, location *domain.Location) error {
	ret := _m.Called(ctx, location)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Location) error); ok {
		r0 = rf(ctx, location)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: ctx, id
func (_m *LocationServiceProvider) Delete(ctx context.Context, id int64) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetByID provides a mock function with given fields: ctx, id
func (_m *LocationServiceProvider) GetByID(ctx context.Context, id int64) (*domain.Location, error) {
	ret := _m.Called(ctx, id)

	var r0 *domain.Location
	if rf, ok := ret.Get(0).(func(context.Context, int64) *domain.Location); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Location)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListAll provides a mock function with given fields: ctx, cursor, limit, filters
func (_m *LocationServiceProvider) ListAll(ctx context.Context, cursor string, limit int64, filters repo.Filters) ([]*domain.Location, string, error) {
	ret := _m.Called(ctx, cursor, limit, filters)

	var r0 []*domain.Location
	if rf, ok := ret.Get(0).(func(context.Context, string, int64, repo.Filters) []*domain.Location); ok {
		r0 = rf(ctx, cursor, limit, filters)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.Location)
		}
	}

	var r1 string
	if rf, ok := ret.Get(1).(func(context.Context, string, int64, repo.Filters) string); ok {
		r1 = rf(ctx, cursor, limit, filters)
	} else {
		r1 = ret.Get(1).(string)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string, int64, repo.Filters) error); ok {
		r2 = rf(ctx, cursor, limit, filters)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Update provides a mock function with given fields: ctx, location
func (_m *LocationServiceProvider) Update(ctx context.Context, location *domain.Location) error {
	ret := _m.Called(ctx, location)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Location) error); ok {
		r0 = rf(ctx, location)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewLocationServiceProvider interface {
	mock.TestingT
	Cleanup(func())
}

// NewLocationServiceProvider creates a new instance of LocationServiceProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewLocationServiceProvider(t mockConstructorTestingTNewLocationServiceProvider) *LocationServiceProvider {
	mock := &LocationServiceProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
