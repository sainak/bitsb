package domain

import (
	"context"
	"net/http"
	"time"

	"gopkg.in/guregu/null.v4"
)

type AccessLevel int

const (
	Admin     AccessLevel = 1
	Passenger AccessLevel = 1000
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

type UserLogin struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (u UserLogin) Bind(r *http.Request) error {
	return nil
}

type UserRegister struct {
	FirstName string `json:"first_name"  binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

func (u UserRegister) Bind(r *http.Request) error {
	return nil
}

type Token struct {
	AuthToken    string `json:"auth_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshToken struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (r2 RefreshToken) Bind(r *http.Request) error {
	return nil
}

type UserStorer interface {
	SelectUserByID(ctx context.Context, id int64) (user User, err error)
	SelectUserByEmail(ctx context.Context, email string) (user User, err error)
	InsertUser(ctx context.Context, user *User) (err error)
	UpdateUser(ctx context.Context, user *User) (err error)
}

type UserServiceProvider interface {
	GetUserByID(ctx context.Context, id int64) (user User, err error)
	Login(ctx context.Context, creds *UserLogin) (token Token, err error)
	RefreshToken(ctx context.Context, token string) (newToken Token, err error)
	Signup(ctx context.Context, user *User) (err error)
	UpdateUser(ctx context.Context, user *User) (err error)
}
