package service

import (
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/guregu/null.v4"

	"github.com/sainak/bitsb/domain"
	"github.com/sainak/bitsb/domain/errors"
	"github.com/sainak/bitsb/pkg/jwt"
	"github.com/sainak/bitsb/pkg/utils"
)

type UserService struct {
	repo domain.UserStorer
	jwt  *jwt.JWT
}

func NewUserService(repo domain.UserStorer, jwtInstance *jwt.JWT) domain.UserServiceProvider {
	return &UserService{
		repo: repo,
		jwt:  jwtInstance,
	}
}

func (u UserService) Login(ctx context.Context, creds *domain.UserLoginForm) (domain.Token, error) {
	token := domain.Token{}

	user, err := u.repo.SelectByEmail(ctx, creds.Email)
	if err != nil {
		// user not found
		err = errors.ErrInvalidCredentials
		return token, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
	if err != nil {
		// wrong password
		err = errors.ErrInvalidCredentials
		return token, err
	}

	// update last login
	user.LastLogin = null.TimeFrom(time.Now())
	err = u.repo.Update(ctx, &user)
	if err != nil {
		return token, err
	}

	token.AuthToken, err = u.jwt.CreateToken(user.ID)
	if err != nil {
		return token, err
	}
	token.RefreshToken, err = u.jwt.CreateRefreshToken(user.ID)
	if err != nil {
		return token, err
	}
	return token, nil
}

func (u UserService) RefreshToken(refreshToken string) (domain.Token, error) {
	token, err := u.jwt.RefreshToken(refreshToken)
	if err != nil {
		err = errors.ErrInvalidToken
	}
	return domain.Token{AuthToken: token}, err
}

func (u UserService) GetByID(ctx context.Context, id int64) (domain.User, error) {
	return u.repo.SelectByID(ctx, id)
}

func (u UserService) Signup(ctx context.Context, user *domain.User) error {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	return u.repo.Insert(ctx, user)
}

func (u UserService) Update(ctx context.Context, user *domain.User) error {
	return u.repo.Update(ctx, user)
}
