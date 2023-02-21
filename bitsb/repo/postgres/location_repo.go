package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"github.com/sainak/bitsb/apperrors"
	"github.com/sainak/bitsb/bitsb"
	"github.com/sainak/bitsb/pkg/repo"
)

type LocationRepository struct {
	conn *sql.DB
}

func NewLocationRepository(conn *sql.DB) bitsb.LocationStorer {
	return &LocationRepository{conn}
}

func (l LocationRepository) SelectAll(
	ctx context.Context,
	cursor string,
	limit int64,
	filters repo.Filters,
) ([]*bitsb.Location, string, error) {
	q := filters.BuildQuery()
	query := `SELECT id, name, created_at, updated_at 
				FROM locations 
				WHERE created_at < $1`
	if q != "" {
		query += " AND " + q
	}
	query += ` ORDER BY created_at DESC  LIMIT $2;`

	locations := make([]*bitsb.Location, 0, limit)
	decodedCursor, err := repo.DecodeCursor(cursor)
	if err != nil {
		err = apperrors.ErrBadCursor
		return locations, "", err
	}

	rows, err := l.conn.QueryContext(ctx, query, decodedCursor, limit)
	if err != nil {
		return locations, "", err
	}
	defer func(rows *sql.Rows) {
		if err := rows.Close(); err != nil {
			logrus.Error(err)
		}
	}(rows)

	for rows.Next() {
		location := bitsb.Location{}
		err = rows.Scan(
			&location.ID,
			&location.Name,
			&location.CreatedAt,
			&location.UpdatedAt,
		)
		if err != nil {
			return locations, "", err
		}
		locations = append(locations, &location)
	}

	var nextCursor string
	if len(locations) == int(limit) {
		nextCursor = repo.EncodeCursor(locations[len(locations)-1].CreatedAt)
	}

	return locations, nextCursor, nil
}

func (l LocationRepository) SelectByID(ctx context.Context, id int64) (*bitsb.Location, error) {
	query := `SELECT id, name, created_at, updated_at 
				FROM locations 
				WHERE id = $1;`

	row := l.conn.QueryRowContext(ctx, query, id)
	location := &bitsb.Location{}
	err := row.Scan(
		&location.ID,
		&location.Name,
		&location.CreatedAt,
		&location.UpdatedAt,
	)
	return location, err
}

func (l LocationRepository) SelectByIDArray(ctx context.Context, ids []int64) ([]*bitsb.Location, error) {
	query := `
		WITH location_order AS (
			SELECT id, row_number() OVER (ORDER BY position) AS order_position
			FROM unnest(cast($1 as integer[])) WITH ORDINALITY AS t(id, position)
		)
		SELECT  locations.id, name, created_at, updated_at
		FROM locations
			JOIN location_order
				ON locations.id = location_order.id
		ORDER BY location_order.order_position, locations.id
	`
	locations := make([]*bitsb.Location, 0, len(ids))
	rows, err := l.conn.QueryContext(ctx, query, pq.Array(ids))
	if err != nil {
		return locations, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			logrus.Error(err)
		}
	}(rows)

	for rows.Next() {
		location := bitsb.Location{}
		err = rows.Scan(
			&location.ID,
			&location.Name,
			&location.CreatedAt,
			&location.UpdatedAt,
		)
		if err != nil {
			return locations, err
		}
		locations = append(locations, &location)
	}
	return locations, nil
}

func (l LocationRepository) Insert(ctx context.Context, location *bitsb.Location) error {
	query := `INSERT INTO locations (name, created_at, updated_at) VALUES ($1, $2, $3) RETURNING id;`

	currentTime := time.Now()
	location.CreatedAt = currentTime
	location.UpdatedAt = currentTime

	return l.conn.QueryRowContext(
		ctx,
		query,
		location.Name,
		location.CreatedAt,
		location.UpdatedAt,
	).Scan(&location.ID)
}

func (l LocationRepository) Update(ctx context.Context, location *bitsb.Location) error {
	query := `UPDATE locations SET name = $2, updated_at = $3 WHERE id = $1;`

	res, err := l.conn.ExecContext(ctx, query, location.ID, location.Name, location.UpdatedAt)
	if err != nil {
		return err
	}
	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
		err = apperrors.ErrNotFound
	}
	return err
}

func (l LocationRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM locations WHERE id = $1;`

	res, err := l.conn.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
		err = apperrors.ErrNotFound
	}
	return err
}
