package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/sainak/bitsb/domain"
)

var (
	selectUserByIDQuery    = `SELECT id, email, first_name, last_name, access_level, password, last_login, created_at, updated_at FROM users WHERE id=$1`
	selectUserByEmailQuery = `SELECT id, email, first_name, last_name, access_level, password, last_login, created_at, updated_at FROM users WHERE email=$1`
	insertUserQuery        = `INSERT INTO users (email, first_name, last_name, access_level, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	updateUserQuery        = `UPDATE users SET email=$2, first_name=$3, last_name=$4, password=$5, updated_at=$6 WHERE id=$1`
)

type UserRepository struct {
	Conn *sql.DB
}

func New(Conn *sql.DB) domain.UserStorer {
	return &UserRepository{Conn}
}

func (u UserRepository) fetchUser(ctx context.Context, query string, args ...interface{}) (user domain.User, err error) {
	user = domain.User{}

	row := u.Conn.QueryRowContext(ctx, query, args...)
	err = row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Access,
		&user.Password,
		&user.LastLogin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	return
}

func (u UserRepository) SelectUserByID(ctx context.Context, id int64) (user domain.User, err error) {
	return u.fetchUser(ctx, selectUserByIDQuery, id)
}

func (u UserRepository) SelectUserByEmail(ctx context.Context, email string) (user domain.User, err error) {
	return u.fetchUser(ctx, selectUserByEmailQuery, email)
}

func (u UserRepository) InsertUser(ctx context.Context, user *domain.User) (err error) {
	currentTime := time.Now()
	user.CreatedAt = currentTime
	user.UpdatedAt = currentTime

	err = u.Conn.QueryRowContext(
		ctx,
		insertUserQuery,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Access,
		user.Password,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)
	return
}

func (u UserRepository) UpdateUser(ctx context.Context, user *domain.User) (err error) {
	user.UpdatedAt = time.Now()
	result, err := u.Conn.ExecContext(
		ctx,
		updateUserQuery,
		user.ID,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Password,
		user.UpdatedAt,
	)
	if err != nil {
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected != 1 {
		return sql.ErrNoRows
	}
	return
}
