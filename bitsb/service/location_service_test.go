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

	"github.com/sainak/bitsb/apperrors"
	"github.com/sainak/bitsb/bitsb"
	"github.com/sainak/bitsb/mocks"
	"github.com/sainak/bitsb/pkg/repo"
)

type LocationServiceTestSuite struct {
	suite.Suite
	service bitsb.LocationServiceProvider
	repo    *mocks.LocationStorer
}

func TestLocationServiceTestSuite(t *testing.T) {
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

	suite.Run(t, new(LocationServiceTestSuite))
}

func (s *LocationServiceTestSuite) SetupTest() {
	s.repo = mocks.NewLocationStorer(s.T())
	s.service = NewLocationService(s.repo)
}

func (s *LocationServiceTestSuite) TestListAll() {
	t := s.T()

	locations := []*bitsb.Location{
		{1, "abc", time.Now(), time.Now()},
		{2, "def", time.Now(), time.Now()},
		{3, "abd", time.Now(), time.Now()},
		{4, "def", time.Now(), time.Now()},
	}

	t.Run("when location list is successfully retrieved", func(t *testing.T) {
		s.repo.
			On("SelectAll", mock.Anything, "", int64(10), repo.Filters{}).
			Return(locations, "", nil).
			Once()
		list, nextCursor, err := s.service.ListAll(context.Background(), "", int64(10), repo.Filters{})
		require.NoError(t, err)
		require.Equal(t, locations, list)
		require.Equal(t, "", nextCursor)
		s.repo.AssertExpectations(t)
	})

	t.Run("whenlist loction is unsuccessful", func(t *testing.T) {
		s.repo.
			On("SelectAll", mock.Anything, "", int64(10), repo.Filters{}).
			Return([]*bitsb.Location{}, "", fmt.Errorf("error")).
			Once()
		list, nextCursor, err := s.service.ListAll(context.Background(), "", int64(10), repo.Filters{})
		require.Error(t, err)
		require.Empty(t, list)
		require.Equal(t, "", nextCursor)
		s.repo.AssertExpectations(t)
	})
}

func (s *LocationServiceTestSuite) TestGetByID() {
	t := s.T()

	location := &bitsb.Location{ID: 1, Name: "abc", CreatedAt: time.Now(), UpdatedAt: time.Now()}

	t.Run("when location is successfully retrieved", func(t *testing.T) {
		s.repo.
			On("SelectByID", mock.Anything, int64(1)).
			Return(location, nil).
			Once()
		loc, err := s.service.GetByID(context.Background(), int64(1))
		require.NoError(t, err)
		require.Equal(t, location, loc)
		s.repo.AssertExpectations(t)
	})

	t.Run("when location is not found", func(t *testing.T) {
		s.repo.
			On("SelectByID", mock.Anything, int64(1)).
			Return(&bitsb.Location{}, apperrors.ErrNotFound).
			Once()
		loc, err := s.service.GetByID(context.Background(), int64(1))
		require.Error(t, err)
		require.Empty(t, loc)
		s.repo.AssertExpectations(t)
	})

	t.Run("when location is unsuccessful", func(t *testing.T) {
		s.repo.
			On("SelectByID", mock.Anything, int64(1)).
			Return(&bitsb.Location{}, fmt.Errorf("error")).
			Once()
		loc, err := s.service.GetByID(context.Background(), int64(1))
		require.Error(t, err)
		require.Empty(t, loc)
		s.repo.AssertExpectations(t)
	})
}

func (s *LocationServiceTestSuite) TestCreate() {
	t := s.T()

	t.Run("when location create is successfully created", func(t *testing.T) {
		s.repo.
			On("Insert", mock.Anything, mock.Anything).
			Return(nil).
			Once()
		err := s.service.Create(context.Background(), &bitsb.Location{Name: "abc"})
		require.NoError(t, err)
		s.repo.AssertExpectations(t)
	})

	t.Run("when location create is unsuccessful", func(t *testing.T) {
		s.repo.
			On("Insert", mock.Anything, mock.Anything).
			Return(fmt.Errorf("error")).
			Once()
		err := s.service.Create(context.Background(), &bitsb.Location{Name: "abc"})
		require.Error(t, err)
		s.repo.AssertExpectations(t)
	})
}

func (s *LocationServiceTestSuite) TestUpdate() {
	t := s.T()

	t.Run("when location update is successfully updated", func(t *testing.T) {
		s.repo.
			On("Update", mock.Anything, mock.Anything).
			Return(nil).
			Once()
		err := s.service.Update(context.Background(), &bitsb.Location{ID: 1, Name: "abc"})
		require.NoError(t, err)
		s.repo.AssertExpectations(t)
	})

	t.Run("when location update is unsuccessful", func(t *testing.T) {
		s.repo.
			On("Update", mock.Anything, mock.Anything).
			Return(fmt.Errorf("error")).
			Once()
		err := s.service.Update(context.Background(), &bitsb.Location{ID: 1, Name: "abc"})
		require.Error(t, err)
		s.repo.AssertExpectations(t)
	})
}

func (s *LocationServiceTestSuite) TestDelete() {
	t := s.T()

	t.Run("when location delete is successfully deleted", func(t *testing.T) {
		s.repo.
			On("Delete", mock.Anything, int64(1)).
			Return(nil).
			Once()
		err := s.service.Delete(context.Background(), int64(1))
		require.NoError(t, err)
		s.repo.AssertExpectations(t)
	})

	t.Run("when location delete is unsuccessful", func(t *testing.T) {
		s.repo.
			On("Delete", mock.Anything, int64(1)).
			Return(fmt.Errorf("error")).
			Once()
		err := s.service.Delete(context.Background(), int64(1))
		require.Error(t, err)
		s.repo.AssertExpectations(t)
	})
}
