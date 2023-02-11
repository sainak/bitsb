package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/sainak/bitsb/domain"
)

type UserRepository struct {
	conn *sql.DB
}

func NewUserRepository(conn *sql.DB) domain.UserStorer {
	return &UserRepository{conn}
}

func (u UserRepository) fetchUser(ctx context.Context, query string, args ...interface{}) (domain.User, error) {
	user := domain.User{}
	err := u.conn.QueryRowContext(ctx, query, args...).Scan(
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
	return user, err
}

func (u UserRepository) SelectByID(ctx context.Context, id int64) (domain.User, error) {
	query := `SELECT id, email, first_name, last_name, access_level, password, last_login, created_at, updated_at 
				FROM users 
				WHERE id=$1`
	return u.fetchUser(ctx, query, id)
}

func (u UserRepository) SelectByEmail(ctx context.Context, email string) (domain.User, error) {
	query := `SELECT id, email, first_name, last_name, access_level, password, last_login, created_at, updated_at 
				FROM users 
				WHERE email=$1`
	return u.fetchUser(ctx, query, email)
}

func (u UserRepository) Insert(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (email, first_name, last_name, access_level, password, created_at, updated_at) 
				VALUES ($1, $2, $3, $4, $5, $6, $7) 
				RETURNING id`

	currentTime := time.Now()
	user.CreatedAt = currentTime
	user.UpdatedAt = currentTime

	err := u.conn.QueryRowContext(
		ctx,
		query,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Access,
		user.Password,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)
	return err
}

func (u UserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `UPDATE users 
				SET email=$2, first_name=$3, last_name=$4, password=$5, last_login=$6, updated_at=$7 
				WHERE id=$1`
	user.UpdatedAt = time.Now()
	result, err := u.conn.ExecContext(
		ctx,
		query,
		user.ID,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Password,
		user.LastLogin,
		user.UpdatedAt,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected != 1 {
		return sql.ErrNoRows
	}
	return err
}
