package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/sainak/bitsb/domain"
	"github.com/sainak/bitsb/domain/errors"
	"github.com/sainak/bitsb/utils/repo"
)

type LocationRepository struct {
	Conn *sql.DB
}

func NewLocationRepository(conn *sql.DB) domain.LocationStorer {
	return &LocationRepository{conn}
}

func (l LocationRepository) SelectAll(
	ctx context.Context,
	cursor string,
	limit int64,
	filters repo.Filters,
) (locations []*domain.Location, nextCursor string, err error) {
	q := filters.BuildQuery()
	query := `SELECT id, name, created_at, updated_at 
				FROM locations 
				WHERE created_at < $1`
	if q != "" {
		query += " AND " + q
	}
	query += ` ORDER BY created_at DESC  LIMIT $2;`

	locations = make([]*domain.Location, 0, limit)
	decodedCursor, err := repo.DecodeCursor(cursor)
	if err != nil {
		err = errors.ErrBadCursor
		return
	}

	rows, err := l.Conn.QueryContext(ctx, query, decodedCursor, limit)
	if err != nil {
		return
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			logrus.Error(err)
		}
	}(rows)

	for rows.Next() {
		location := domain.Location{}
		err = rows.Scan(
			&location.ID,
			&location.Name,
			&location.CreatedAt,
			&location.UpdatedAt,
		)
		if err != nil {
			return
		}
		locations = append(locations, &location)
	}

	if len(locations) == int(limit) {
		nextCursor = repo.EncodeCursor(locations[len(locations)-1].CreatedAt)
	}

	return
}

func (l LocationRepository) SelectByID(ctx context.Context, id int64) (location *domain.Location, err error) {
	query := `SELECT id, name, created_at, updated_at 
				FROM locations 
				WHERE id = $1;`

	row := l.Conn.QueryRowContext(ctx, query, id)
	location = &domain.Location{}
	err = row.Scan(
		&location.ID,
		&location.Name,
		&location.CreatedAt,
		&location.UpdatedAt,
	)
	if err != nil {
		return
	}
	return
}

func (l LocationRepository) Insert(ctx context.Context, location *domain.Location) error {
	query := `INSERT INTO locations (name, created_at, updated_at) VALUES ($1, $2, $3) RETURNING id;`

	currentTime := time.Now()
	location.CreatedAt = currentTime
	location.UpdatedAt = currentTime

	return l.Conn.QueryRowContext(
		ctx,
		query,
		location.Name,
		location.CreatedAt,
		location.UpdatedAt,
	).Scan(&location.ID)
}

func (l LocationRepository) Update(ctx context.Context, location *domain.Location) (err error) {
	query := `UPDATE locations SET name = $2, updated_at = $3 WHERE id = $1;`

	res, err := l.Conn.ExecContext(ctx, query, location.ID, location.Name, location.UpdatedAt)
	if err != nil {
		return
	}
	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
		err = errors.ErrNotFound
	}
	return
}

func (l LocationRepository) Delete(ctx context.Context, id int64) (err error) {
	query := `DELETE FROM locations WHERE id = $1;`

	res, err := l.Conn.ExecContext(ctx, query, id)
	if err != nil {
		return
	}
	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
		err = errors.ErrNotFound
	}
	return
}
