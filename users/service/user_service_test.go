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
	"gopkg.in/guregu/null.v4"

	"github.com/sainak/bitsb/mocks"
	"github.com/sainak/bitsb/pkg/jwt"
	"github.com/sainak/bitsb/pkg/utils"
	"github.com/sainak/bitsb/users"
)

type UserServiceTestSuite struct {
	suite.Suite
	service users.UserServiceProvider
	repo    *mocks.UserStorer
	jwt     *jwt.JWT
}

func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}

func (s *UserServiceTestSuite) SetupTest() {
	s.repo = mocks.NewUserStorer(s.T())
	s.jwt = jwt.New("test_secret", "24", "5")
	s.service = NewUserService(s.repo, s.jwt)
}

func (s *UserServiceTestSuite) TestLogin() {
	t := s.T()

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

	password := "test_pass"
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		t.Fatal(err)
	}
	user := users.User{
		ID:        1,
		FirstName: "Tester",
		LastName:  "User",
		Email:     "testuser@email.com",
		Password:  hashedPassword,
		Access:    users.Passenger,
		LastLogin: null.TimeFrom(time.Now()),
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	t.Run("when login is successful", func(t *testing.T) {
		s.repo.
			On("SelectByEmail", mock.Anything, user.Email).
			Return(user, nil)
		s.repo.
			On("Update", mock.Anything, &user).
			Return(nil)
		creds := &users.UserLoginForm{
			Email:    user.Email,
			Password: password,
		}
		token, err := s.service.Login(context.Background(), creds)
		require.Nil(t, err)
		parsedToken, err := s.jwt.ParseToken(token.AuthToken)
		require.Nil(t, err)
		require.True(t, parsedToken.Valid)
	})

	t.Run("when password is incorrect", func(t *testing.T) {
		// reuse mock result from previous function
		creds := &users.UserLoginForm{
			Email:    user.Email,
			Password: "incorrect_password",
		}
		token, err := s.service.Login(context.Background(), creds)
		require.Error(t, err)
		require.Zero(t, token)
	})

	t.Run("when user does not exist", func(t *testing.T) {
		s.repo.
			On("SelectByEmail", mock.Anything, "nobody@example.com").
			Return(users.User{}, fmt.Errorf("record not forund"))
		creds := &users.UserLoginForm{
			Email:    "nobody@example.com",
			Password: password,
		}
		token, err := s.service.Login(context.Background(), creds)
		require.Error(t, err)
		require.Zero(t, token)
	})
}
