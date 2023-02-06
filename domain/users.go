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
	ID        int64       `json:"id"`
	FirstName string      `json:"first_name"`
	LastName  string      `json:"last_name"`
	Email     string      `json:"email"`
	Password  string      `json:"password"`
	Access    AccessLevel `json:"-"`
	LastLogin null.Time   `json:"last_login"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
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

type UserStorer interface {
	SelectUser(ctx context.Context, id int64) (User, error)
	InsertUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, user *User) error
}
