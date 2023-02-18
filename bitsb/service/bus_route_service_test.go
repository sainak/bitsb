package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/undefinedlabs/go-mpatch"

	"github.com/sainak/bitsb/domain"
	"github.com/sainak/bitsb/domain/mocks"
)

type BusRouteServiceTestSuite struct {
	suite.Suite
	service      domain.BusRouteServiceProvider
	repo         *mocks.BusRouteStorer
	locationRepo *mocks.LocationStorer
}

func TestBusRouteServiceTestSuite(t *testing.T) {
	patch, err := mpatch.PatchMethod(time.Now, func() time.Time {
		return time.Date(2020, 11, 01, 00, 00, 00, 0, time.UTC)
	})
	if err != nil {
		t.Fatal(err)
	}
	defer func(patch *mpatch.Patch) {
		if err := patch.Unpatch(); err != nil {
			t.Fatal(err)
		}
	}(patch)

	suite.Run(t, new(BusRouteServiceTestSuite))
}

func (s *BusRouteServiceTestSuite) SetupTest() {
	s.repo = mocks.NewBusRouteStorer(s.T())
	s.locationRepo = mocks.NewLocationStorer(s.T())
	s.service = NewBusRouteService(s.repo, s.locationRepo)
}

func (s *BusRouteServiceTestSuite) TestListAll() {
	t := s.T()

	busRoutes := []*domain.BusRoute{
		{
			ID:   1,
			Name: "Test Route 1",
		},
		{
			ID:   2,
			Name: "Test Route 2",
		},
	}

	t.Run("when list all routes is successful", func(t *testing.T) {
		s.repo.
			On("SelectAll", mock.Anything, "", int64(10), []int64{}).
			Return(busRoutes, "", nil)

		routes, cursor, err := s.service.ListAll(context.Background(), "", int64(10), []int64{})
		require.NoError(t, err)
		require.Equal(t, busRoutes, routes)
		require.Equal(t, "", cursor)
	})

	t.Run("when list all routes is unsuccessful", func(t *testing.T) {
		s.repo.
			On("SelectAll", mock.Anything, "awd342", int64(10), []int64{}).
			Return([]*domain.BusRoute{}, "", fmt.Errorf("error"))

		routes, cursor, err := s.service.ListAll(context.Background(), "awd342", int64(10), []int64{})
		require.Error(t, err)
		require.Empty(t, routes)
		require.Equal(t, "", cursor)
	})
}

func (s *BusRouteServiceTestSuite) TestGetByID() {
	t := s.T()

	busRoute := &domain.BusRoute{
		ID:          1,
		Name:        "Test Route 1",
		MinPrice:    3,
		MaxPrice:    10,
		LocationIDS: []int64{1, 2},
	}

	locationDetails := []*domain.Location{
		{ID: 1, Name: "location 1"},
		{ID: 2, Name: "location 2"},
	}

	locations := []*domain.LocationForm{
		{"location 1"},
		{"location 2"},
	}

	routeWithLoc := &domain.BusRoute{
		ID:          1,
		Name:        "Test Route 1",
		MinPrice:    3,
		MaxPrice:    10,
		LocationIDS: []int64{1, 2},
		Locations:   locations,
	}

	t.Run("when get route by id is successful", func(t *testing.T) {
		s.repo.
			On("SelectByID", mock.Anything, int64(1)).
			Return(busRoute, nil)

		s.locationRepo.
			On("SelectByIDArray", mock.Anything, busRoute.LocationIDS).
			Return(locationDetails, nil)

		route, err := s.service.GetByID(context.Background(), int64(1))
		require.NoError(t, err)
		require.Equal(t, routeWithLoc, route)
	})

	t.Run("when get route by id is unsuccessful", func(t *testing.T) {
		s.repo.
			On("SelectByID", mock.Anything, int64(2)).
			Return(nil, fmt.Errorf("error"))

		route, err := s.service.GetByID(context.Background(), int64(2))
		require.Error(t, err)
		require.Empty(t, route)
	})

	t.Run("when get route by id is unsuccessful with invalid locations", func(t *testing.T) {
		invalidBusRoute := &domain.BusRoute{
			ID:          3,
			Name:        "Test Route 1",
			MinPrice:    3,
			MaxPrice:    10,
			LocationIDS: []int64{999, 33},
		}
		s.repo.
			On("SelectByID", mock.Anything, invalidBusRoute.ID).
			Return(invalidBusRoute, nil)
		s.locationRepo.
			On("SelectByIDArray", mock.Anything, invalidBusRoute.LocationIDS).
			Return(nil, fmt.Errorf("error"))

		route, err := s.service.GetByID(context.Background(), invalidBusRoute.ID)
		require.Error(t, err)
		require.Empty(t, route)
	})
}

func (s *BusRouteServiceTestSuite) TestCalculateTicketPrice() {
	t := s.T()

	busRoute := &domain.BusRoute{
		ID:          1,
		Name:        "Test Route 1",
		MinPrice:    3,
		MaxPrice:    10,
		LocationIDS: []int64{1, 2, 3, 5, 7, 8},
	}

	t.Run("when calculate ticket price is successful for 5 stops", func(t *testing.T) {
		var start, end int64 = 1, 7
		s.repo.
			On("SelectByID", mock.Anything, int64(1)).
			Return(busRoute, nil)

		price, err := s.service.CalculateTicketPrice(context.Background(), busRoute.ID, start, end)
		require.NoError(t, err)
		require.Equal(t, int64(10), price)
	})

	t.Run("when calculate ticket price is successful for 1 stop", func(t *testing.T) {
		var start, end int64 = 1, 2
		s.repo.
			On("SelectByID", mock.Anything, int64(1)).
			Return(busRoute, nil)

		price, err := s.service.CalculateTicketPrice(context.Background(), busRoute.ID, start, end)
		require.NoError(t, err)
		require.Equal(t, int64(3), price)
	})

	t.Run("claculete ticket price for invalid route", func(t *testing.T) {
		var start, end int64 = 1, 2
		s.repo.
			On("SelectByID", mock.Anything, int64(3)).
			Return(nil, fmt.Errorf("error"))

		price, err := s.service.CalculateTicketPrice(context.Background(), int64(3), start, end)
		require.Error(t, err)
		require.Equal(t, int64(0), price)
	})

	t.Run("claculete ticket price for invalid start location", func(t *testing.T) {
		var start, end int64 = 10, 2

		price, err := s.service.CalculateTicketPrice(context.Background(), busRoute.ID, start, end)
		require.Error(t, err)
		require.Equal(t, int64(0), price)
	})
}

func (s *BusRouteServiceTestSuite) TestCreate() {
	t := s.T()

	busRoute := &domain.BusRoute{
		ID:          1,
		Name:        "Test Route 1",
		MinPrice:    3,
		MaxPrice:    10,
		LocationIDS: []int64{1, 2, 3, 5, 7, 8},
	}

	t.Run("when create route is successful", func(t *testing.T) {
		s.repo.
			On("Insert", mock.Anything, busRoute).
			Return(nil)

		err := s.service.Create(context.Background(), busRoute)
		require.NoError(t, err)
	})

	t.Run("when create route is unsuccessful", func(t *testing.T) {
		s.repo.
			On("Insert", mock.Anything, &domain.BusRoute{}).
			Return(fmt.Errorf("error"))

		err := s.service.Create(context.Background(), &domain.BusRoute{})
		require.Error(t, err)
	})
}

func (s *BusRouteServiceTestSuite) TestUpdate() {
	t := s.T()

	busRoute := &domain.BusRoute{
		ID:          1,
		Name:        "Test Route 1",
		MinPrice:    3,
		MaxPrice:    10,
		LocationIDS: []int64{1, 2, 3, 5, 7, 8},
	}

	t.Run("when update route is successful", func(t *testing.T) {
		s.repo.
			On("Update", mock.Anything, busRoute).
			Return(nil)

		err := s.service.Update(context.Background(), busRoute)
		require.NoError(t, err)
	})

	t.Run("when update route is unsuccessful", func(t *testing.T) {
		s.repo.
			On("Update", mock.Anything, &domain.BusRoute{}).
			Return(fmt.Errorf("error"))

		err := s.service.Update(context.Background(), &domain.BusRoute{})
		require.Error(t, err)
	})
}

func (s *BusRouteServiceTestSuite) TestDelete() {
	t := s.T()

	t.Run("when delete route is successful", func(t *testing.T) {
		s.repo.
			On("Delete", mock.Anything, int64(1)).
			Return(nil)

		err := s.service.Delete(context.Background(), int64(1))
		require.NoError(t, err)
	})

	t.Run("when delete route is unsuccessful", func(t *testing.T) {
		s.repo.
			On("Delete", mock.Anything, int64(4)).
			Return(fmt.Errorf("error"))

		err := s.service.Delete(context.Background(), int64(4))
		require.Error(t, err)
	})
}
