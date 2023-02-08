package postgres

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/undefinedlabs/go-mpatch"
	"gopkg.in/guregu/null.v4"

	"github.com/sainak/bitsb/domain"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	db   *sql.DB
	mock sqlmock.Sqlmock
	repo domain.UserStorer
}

func (s *UserRepositoryTestSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	if err != nil {
		s.T().Fatal(err)
	}
	s.db = db
	s.mock = mock
	s.repo = NewUserRepository(db)
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}

func (s *UserRepositoryTestSuite) TestSelectUser() {
	t := s.T()

	user := &domain.User{
		ID:        1,
		FirstName: "Jhon",
		LastName:  "Doe",
		Email:     "jhon.doe@example.com",
		Password:  "test_password",
		Access:    domain.Admin,
	}

	t.Run("when select by user id is successful", func(t *testing.T) {
		s.mock.ExpectQuery("SELECT (.+) FROM users").
			WithArgs(user.ID).
			WillReturnRows(sqlmock.
				NewRows([]string{
					"id",
					"email",
					"first_name",
					"last_name",
					"access_level",
					"password",
					"last_login",
					"created_at",
					"updated_at",
				}).
				AddRow(
					user.ID,
					user.Email,
					user.FirstName,
					user.LastName,
					user.Access,
					user.Password,
					user.LastLogin,
					user.CreatedAt,
					user.UpdatedAt,
				),
			)
		res, err := s.repo.SelectByID(context.Background(), user.ID)
		assert.Nil(t, err)
		assert.Equal(t, user, &res)
	})

	t.Run("when select by user email is successful", func(t *testing.T) {
		s.mock.ExpectQuery("SELECT (.+) FROM users").
			WithArgs(user.Email).
			WillReturnRows(sqlmock.
				NewRows([]string{
					"id",
					"email",
					"first_name",
					"last_name",
					"access_level",
					"password",
					"last_login",
					"created_at",
					"updated_at",
				}).
				AddRow(
					user.ID,
					user.Email,
					user.FirstName,
					user.LastName,
					user.Access,
					user.Password,
					user.LastLogin,
					user.CreatedAt,
					user.UpdatedAt,
				),
			)
		res, err := s.repo.SelectByEmail(context.Background(), user.Email)
		assert.Nil(t, err)
		assert.Equal(t, user, &res)
	})

	t.Run("when select by user id is not successful", func(t *testing.T) {
		s.mock.ExpectQuery("SELECT (.+) FROM users").
			WithArgs(user.ID).
			WillReturnError(sql.ErrNoRows)
		_, err := s.repo.SelectByID(context.Background(), user.ID)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})
}

func (s *UserRepositoryTestSuite) TestInsertUser() {
	t := s.T()

	patch, err := mpatch.PatchMethod(time.Now, func() time.Time {
		return time.Date(2020, 11, 01, 00, 00, 00, 0, time.UTC)
	})
	if err != nil {
		t.Fatal(err)
	}
	defer func(patch *mpatch.Patch) {
		err := patch.Unpatch()
		if err != nil {
			t.Fatal(err)
		}
	}(patch)

	user := &domain.User{
		ID:        1,
		FirstName: "Jhon",
		LastName:  "Doe",
		Email:     "jhon.doe@example.com",
		Password:  "test_password",
		Access:    domain.Admin,
	}

	t.Run("when insert is successful", func(t *testing.T) {
		s.mock.ExpectQuery("INSERT INTO").
			WithArgs(
				user.Email,
				user.FirstName,
				user.LastName,
				user.Access,
				user.Password,
				time.Now(),
				time.Now(),
			).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).AddRow(user.ID),
			)
		err = s.repo.Insert(context.Background(), user)
		assert.Nil(t, err)
		assert.Equal(t, int64(1), user.ID)
		assert.Equal(t, user.CreatedAt, time.Now())
	})

	t.Run("when insert is not successful", func(t *testing.T) {
		s.mock.ExpectQuery("INSERT INTO").
			WithArgs(
				user.Email,
				user.FirstName,
				user.LastName,
				user.Access,
				user.Password,
				time.Now(),
				time.Now(),
			).
			WillReturnError(sql.ErrNoRows)
		err := s.repo.Insert(context.Background(), user)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})
}

func (s *UserRepositoryTestSuite) TestUpdateUser() {
	t := s.T()

	patch, err := mpatch.PatchMethod(time.Now, func() time.Time {
		return time.Date(2020, 11, 01, 00, 00, 00, 0, time.UTC)
	})
	if err != nil {
		t.Fatal(err)
	}
	defer func(patch *mpatch.Patch) {
		err := patch.Unpatch()
		if err != nil {
			t.Fatal(err)
		}
	}(patch)

	user := &domain.User{
		ID:        0,
		FirstName: "",
		LastName:  "",
		Email:     "",
		Password:  "",
		Access:    0,
		LastLogin: null.Time{},
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	t.Run("when update is successful", func(t *testing.T) {
		s.mock.ExpectExec("UPDATE users SET (.+)").
			WithArgs(
				user.ID,
				user.Email,
				user.FirstName,
				user.LastName,
				user.Password,
				user.LastLogin,
				time.Now(),
			).
			WillReturnResult(
				sqlmock.NewResult(1, 1),
			)
		err := s.repo.Update(context.Background(), user)
		assert.Nil(t, err)
		assert.Equal(t, user.UpdatedAt, time.Now())
	})

	t.Run("when update is performed on invalid user", func(t *testing.T) {
		s.mock.ExpectExec("UPDATE users SET (.+)").
			WithArgs(
				user.ID,
				user.Email,
				user.FirstName,
				user.LastName,
				user.Password,
				user.LastLogin,
				time.Now(),
			).
			WillReturnResult(
				sqlmock.NewResult(0, 0),
			)
		err := s.repo.Update(context.Background(), user)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})

	t.Run("when update is not successful", func(t *testing.T) {
		s.mock.ExpectExec("UPDATE users SET (.+)").
			WithArgs(
				user.ID,
				user.Email,
				user.FirstName,
				user.LastName,
				user.Password,
				user.LastLogin,
				time.Now(),
			).
			WillReturnError(sql.ErrNoRows)
		err := s.repo.Update(context.Background(), user)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})
}
