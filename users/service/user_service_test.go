package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gopkg.in/guregu/null.v4"

	"github.com/sainak/bitsb/domain"
	"github.com/sainak/bitsb/domain/mocks"
	"github.com/sainak/bitsb/utils/jwt"
)

type UserServiceTestSuite struct {
	suite.Suite
	service domain.UserServiceProvider
	repo    *mocks.UserStorer
	jwt     *jwt.JWT
}

func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}

func (suite *UserServiceTestSuite) SetupTest() {
	suite.repo = mocks.NewUserStorer(suite.T())
	suite.jwt = jwt.New("test_secret", "24", "5")
	suite.service = New(suite.repo, suite.jwt, 0)
}

func (suite *UserServiceTestSuite) TestLogin() {
	t := suite.T()

	password := "test_pass"
	p, err := hashPassword(password)
	if err != nil {
		t.Fatal(err)
	}
	user := domain.User{
		ID:        1,
		FirstName: "Tester",
		LastName:  "User",
		Email:     "testuser@email.com",
		Password:  p,
		Access:    domain.Passenger,
		LastLogin: null.Time{},
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}
	t.Run("when login is successful", func(t *testing.T) {
		suite.repo.
			On("SelectUserByEmail", mock.Anything, user.Email).
			Return(user, nil)
		creds := &domain.UserLogin{
			Email:    user.Email,
			Password: password,
		}
		token, err := suite.service.Login(context.Background(), creds)
		require.Nil(t, err)
		parsedToken, err := suite.jwt.ParseToken(token.AuthToken)
		require.Nil(t, err)
		require.True(t, parsedToken.Valid)
	})

	t.Run("when password is incorrect", func(t *testing.T) {
		// reuse mock result from previous function
		creds := &domain.UserLogin{
			Email:    user.Email,
			Password: "incorrect_password",
		}
		token, err := suite.service.Login(context.Background(), creds)
		require.Error(t, err)
		require.Zero(t, token)
	})

	t.Run("when user does not exist", func(t *testing.T) {
		suite.repo.
			On("SelectUserByEmail", mock.Anything, "nobody@example.com").
			Return(domain.User{}, fmt.Errorf("record not forund"))
		creds := &domain.UserLogin{
			Email:    "nobody@example.com",
			Password: password,
		}
		token, err := suite.service.Login(context.Background(), creds)
		require.Error(t, err)
		require.Zero(t, token)
	})

}
