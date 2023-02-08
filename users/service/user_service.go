package service

import (
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/sainak/bitsb/domain"
	"github.com/sainak/bitsb/domain/errors"
	"github.com/sainak/bitsb/utils/jwt"
)

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

type UserService struct {
	repo           domain.UserStorer
	jwt            *jwt.JWT
	contextTimeout time.Duration
}

func NewUserService(repo domain.UserStorer, jwtInstance *jwt.JWT, timeout time.Duration) domain.UserServiceProvider {
	return &UserService{
		repo:           repo,
		jwt:            jwtInstance,
		contextTimeout: timeout,
	}
}

func (u UserService) Login(ctx context.Context, creds *domain.UserLoginForm) (token domain.Token, err error) {
	token = domain.Token{}
	user, err := u.repo.SelectByEmail(ctx, creds.Email)
	if err != nil {
		// user not found
		err = errors.ErrInvalidCredentials
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
	if err != nil {
		// wrong password
		err = errors.ErrInvalidCredentials
		return
	}

	token.AuthToken, err = u.jwt.CreateToken(user.ID)
	if err != nil {
		return
	}
	token.RefreshToken, err = u.jwt.CreateRefreshToken(user.ID)
	if err != nil {
		return
	}
	return
}

func (u UserService) RefreshToken(ctx context.Context, refreshToken string) (newToken domain.Token, err error) {
	newToken = domain.Token{}
	newToken.AuthToken, err = u.jwt.RefreshToken(refreshToken)
	return
}

func (u UserService) GetByID(ctx context.Context, id int64) (user domain.User, err error) {
	return u.repo.SelectByID(ctx, id)
}

func (u UserService) Signup(ctx context.Context, user *domain.User) (err error) {
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return
	}
	user.Password = hashedPassword
	return u.repo.Insert(ctx, user)
}

func (u UserService) Update(ctx context.Context, user *domain.User) (err error) {
	return u.repo.Update(ctx, user)
}
