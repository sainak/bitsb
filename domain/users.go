package domain

import (
	"context"
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
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRegister struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type Token struct {
	AuthToken    string `json:"auth_token"`
	RefreshToken string `json:"refresh_token"`
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
