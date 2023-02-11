package domain

import (
	"context"
	"net/http"
	"time"

	"gopkg.in/guregu/null.v4"
)

type AccessLevel int

const (
	Admin     AccessLevel = 1000
	Passenger AccessLevel = 10
)

type User struct {
	ID        int64       `json:"id" db:"id"`
	FirstName string      `json:"first_name"  db:"first_name"`
	LastName  string      `json:"last_name" db:"last_name"`
	Email     string      `json:"email" db:"email"`
	Password  string      `json:"-" db:"password"`
	Access    AccessLevel `json:"-" db:"access_level"`
	LastLogin null.Time   `json:"last_login" db:"last_login"`
	CreatedAt time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt time.Time   `json:"updated_at" db:"updated_at"`
}

type UserLoginForm struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (u UserLoginForm) Bind(r *http.Request) error {
	return nil
}

type UserRegisterForm struct {
	FirstName string `json:"first_name"  binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

func (u UserRegisterForm) Bind(r *http.Request) error {
	return nil
}

type Token struct {
	AuthToken    string `json:"auth_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenFrom struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (r2 RefreshTokenFrom) Bind(r *http.Request) error {
	return nil
}

type UserStorer interface {
	SelectByID(ctx context.Context, id int64) (User, error)
	SelectByEmail(ctx context.Context, email string) (User, error)
	Insert(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
}

type UserServiceProvider interface {
	GetByID(ctx context.Context, id int64) (User, error)
	Login(ctx context.Context, creds *UserLoginForm) (Token, error)
	RefreshToken(token string) (Token, error)
	Signup(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
}
