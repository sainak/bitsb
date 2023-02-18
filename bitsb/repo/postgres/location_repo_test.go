package postgres

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/undefinedlabs/go-mpatch"

	"github.com/sainak/bitsb/domain"
	"github.com/sainak/bitsb/pkg/repo"
)

type LocationRepositoryTestSuite struct {
	suite.Suite
	db   *sql.DB
	mock sqlmock.Sqlmock
	repo domain.LocationStorer
}

func TestLocationRepositoryTestSuite(t *testing.T) {
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

	suite.Run(t, new(LocationRepositoryTestSuite))
}

func (s *LocationRepositoryTestSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	if err != nil {
		s.T().Fatal(err)
	}
	s.db = db
	s.mock = mock
	s.repo = NewLocationRepository(db)
}

func (s *LocationRepositoryTestSuite) TestSelectAll() {
	t := s.T()

	t.Run("when select all locations is successful", func(t *testing.T) {
		s.mock.ExpectQuery("SELECT (.+) FROM locations").
			WillReturnRows(sqlmock.
				NewRows([]string{
					"id",
					"name",
					"created_at",
					"updated_at",
				}).
				AddRow(
					1,
					"Test Location",
					time.Now(),
					time.Now(),
				).
				AddRow(
					2,
					"Test Location 2",
					time.Now(),
					time.Now(),
				),
			)

		got, cursor, err := s.repo.SelectAll(context.Background(), "", int64(10), repo.Filters{})
		require.NoError(t, err)
		require.NotEmpty(t, got)
		require.Equal(t, "", cursor)
	})

	t.Run("when select all locations is not successful", func(t *testing.T) {
		s.mock.ExpectQuery("SELECT (.+) FROM locations").
			WillReturnError(sql.ErrNoRows)

		got, cursor, err := s.repo.SelectAll(context.Background(), "", int64(10), repo.Filters{})
		require.Error(t, err)
		require.Empty(t, got)
		require.Equal(t, "", cursor)
	})
}

func (s *LocationRepositoryTestSuite) TestSelectByID() {
	t := s.T()

	location := &domain.Location{
		ID:   1,
		Name: "Test Location",
	}

	t.Run("when select by location id is successful", func(t *testing.T) {
		s.mock.ExpectQuery("SELECT (.+) FROM locations").
			WithArgs(location.ID).
			WillReturnRows(sqlmock.
				NewRows([]string{
					"id",
					"name",
					"created_at",
					"updated_at",
				}).
				AddRow(
					location.ID,
					location.Name,
					location.CreatedAt,
					location.UpdatedAt,
				))

		got, err := s.repo.SelectByID(context.Background(), location.ID)
		require.NoError(t, err)
		require.Equal(t, location, got)
	})

	t.Run("when select by location id is not successful", func(t *testing.T) {
		s.mock.ExpectQuery("SELECT (.+) FROM locations").
			WithArgs(int64(2)).
			WillReturnError(sql.ErrNoRows)

		got, err := s.repo.SelectByID(context.Background(), int64(2))
		require.Error(t, err)
		require.Empty(t, got)
	})
}

func (s *LocationRepositoryTestSuite) SelectByIDArray() {
	t := s.T()

	location := &domain.Location{
		ID:   1,
		Name: "Test Location",
	}

	t.Run("when select by location id array is successful", func(t *testing.T) {
		s.mock.ExpectQuery("SELECT (.+) FROM locations").
			WithArgs(location.ID).
			WillReturnRows(sqlmock.
				NewRows([]string{
					"id",
					"name",
					"created_at",
					"updated_at",
				}).
				AddRow(
					location.ID,
					location.Name,
					location.CreatedAt,
					location.UpdatedAt,
				))

		got, err := s.repo.SelectByIDArray(context.Background(), []int64{location.ID})
		require.NoError(t, err)
		require.Equal(t, []*domain.Location{location}, got)
	})

	t.Run("when select by location id array is not successful", func(t *testing.T) {
		s.mock.ExpectQuery("SELECT (.+) FROM locations").
			WithArgs([]int64{2}).
			WillReturnError(sql.ErrNoRows)

		got, err := s.repo.SelectByIDArray(context.Background(), []int64{2})
		require.Error(t, err)
		require.Empty(t, got)
	})
}

func (s *LocationRepositoryTestSuite) TestInsert() {
	t := s.T()

	location := &domain.Location{
		Name: "Test Location",
	}

	t.Run("when insert location is successful", func(t *testing.T) {
		s.mock.ExpectQuery("INSERT INTO locations").
			WithArgs(
				location.Name,
				time.Now(),
				time.Now(),
			).
			WillReturnRows(sqlmock.
				NewRows([]string{"id"}).
				AddRow(location.ID),
			)
		s.mock.ExpectCommit()

		err := s.repo.Insert(context.Background(), location)
		require.NoError(t, err)
	})

	t.Run("when insert location is not successful", func(t *testing.T) {
		s.mock.ExpectQuery("INSERT INTO locations (.+) RETURNING id").
			WithArgs(location.Name).
			WillReturnError(sql.ErrNoRows)
		s.mock.ExpectRollback()

		err := s.repo.Insert(context.Background(), location)
		require.Error(t, err)
	})
}

func (s *LocationRepositoryTestSuite) TestUpdate() {
	t := s.T()

	location := &domain.Location{
		ID:        1,
		Name:      "Test Location",
		UpdatedAt: time.Now(),
	}

	t.Run("when update location is successful", func(t *testing.T) {
		s.mock.ExpectExec("UPDATE locations").
			WithArgs(location.ID, location.Name, time.Now()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectCommit()

		err := s.repo.Update(context.Background(), location)
		require.NoError(t, err)
	})

	t.Run("when update location is not successful", func(t *testing.T) {
		s.mock.ExpectExec("UPDATE locations").
			WithArgs(location.ID, location.Name, time.Now()).
			WillReturnError(sql.ErrNoRows)
		s.mock.ExpectRollback()

		err := s.repo.Update(context.Background(), location)
		require.Error(t, err)
	})

	t.Run("when no rows affected", func(t *testing.T) {
		lc := domain.Location{
			ID:        2,
			Name:      "Test Location",
			UpdatedAt: time.Now(),
		}
		s.mock.ExpectExec("UPDATE locations").
			WithArgs(lc.ID, lc.Name, time.Now()).
			WillReturnResult(sqlmock.NewResult(0, 0))
		s.mock.ExpectCommit()

		err := s.repo.Update(context.Background(), &lc)
		require.Error(t, err)
	})
}
