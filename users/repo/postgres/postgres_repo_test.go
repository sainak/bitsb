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

func (suite *UserRepositoryTestSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	if err != nil {
		suite.T().Fatal(err)
	}
	suite.db = db
	suite.mock = mock
	suite.repo = New(db)
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}

func (suite *UserRepositoryTestSuite) TestSelectUser() {
	t := suite.T()

	user := &domain.User{
		ID:        1,
		FirstName: "Jhon",
		LastName:  "Doe",
		Email:     "jhon.doe@example.com",
		Password:  "test_password",
		Access:    domain.Admin,
	}

	t.Run("when select by user id is successful", func(t *testing.T) {
		suite.mock.ExpectQuery("SELECT (.+) FROM users").
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
		res, err := suite.repo.SelectUserByID(context.Background(), user.ID)
		assert.Nil(t, err)
		assert.Equal(t, user, &res)
	})

	t.Run("when select by user email is successful", func(t *testing.T) {
		suite.mock.ExpectQuery("SELECT (.+) FROM users").
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
		res, err := suite.repo.SelectUserByEmail(context.Background(), user.Email)
		assert.Nil(t, err)
		assert.Equal(t, user, &res)
	})

	t.Run("when select by user id is not successful", func(t *testing.T) {
		suite.mock.ExpectQuery("SELECT (.+) FROM users").
			WithArgs(user.ID).
			WillReturnError(sql.ErrNoRows)
		_, err := suite.repo.SelectUserByID(context.Background(), user.ID)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})
}

func (suite *UserRepositoryTestSuite) TestInsertUser() {
	t := suite.T()

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
		suite.mock.ExpectQuery("INSERT INTO").
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
		err = suite.repo.InsertUser(context.Background(), user)
		assert.Nil(t, err)
		assert.Equal(t, int64(1), user.ID)
		assert.Equal(t, user.CreatedAt, time.Now())
	})

	t.Run("when insert is not successful", func(t *testing.T) {
		suite.mock.ExpectQuery("INSERT INTO").
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
		err := suite.repo.InsertUser(context.Background(), user)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})
}

func (suite *UserRepositoryTestSuite) TestUpdateUser() {
	t := suite.T()

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
		suite.mock.ExpectExec("UPDATE users SET (.+)").
			WithArgs(
				user.ID,
				user.Email,
				user.FirstName,
				user.LastName,
				user.Password,
				time.Now(),
			).
			WillReturnResult(
				sqlmock.NewResult(1, 1),
			)
		err := suite.repo.UpdateUser(context.Background(), user)
		assert.Nil(t, err)
		assert.Equal(t, user.UpdatedAt, time.Now())
	})

	t.Run("when update is performed on invalid user", func(t *testing.T) {
		suite.mock.ExpectExec("UPDATE users SET (.+)").
			WithArgs(
				user.ID,
				user.Email,
				user.FirstName,
				user.LastName,
				user.Password,
				time.Now(),
			).
			WillReturnResult(
				sqlmock.NewResult(0, 0),
			)
		err := suite.repo.UpdateUser(context.Background(), user)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})

	t.Run("when update is not successful", func(t *testing.T) {
		suite.mock.ExpectExec("UPDATE users SET (.+)").
			WithArgs(
				user.ID,
				user.Email,
				user.FirstName,
				user.LastName,
				user.Password,
				time.Now(),
			).
			WillReturnError(sql.ErrNoRows)
		err := suite.repo.UpdateUser(context.Background(), user)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})
}
