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

type BusRouteRepository struct {
	Conn *sql.DB
}

func NewBusRouteRepository(conn *sql.DB) bitsb.BusRouteStorer {
	return &BusRouteRepository{conn}
}

func (b *BusRouteRepository) SelectAll(
	ctx context.Context,
	cursor string,
	limit int64,
	locations []int64,
) ([]*bitsb.BusRoute, string, error) {
	query := `SELECT id, name, number, start_time, end_time, interval, location_ids, min_price, max_price, created_at, updated_at  FROM bus_routes 
				WHERE created_at < $1 ORDER BY created_at DESC LIMIT $2;`

	queryWithStops := `SELECT  id, name, number, start_time, end_time, interval, location_ids, min_price, max_price, created_at, updated_at FROM bus_routes
    						WHERE location_ids @> cast($3 as int[]) AND created_at < $1 ORDER BY created_at DESC LIMIT $2;`

	decodedCursor, err := repo.DecodeCursor(cursor)
	if err != nil {
		err = apperrors.ErrBadCursor
		return []*bitsb.BusRoute{}, "", err
	}

	var rows *sql.Rows
	if len(locations) == 0 {
		rows, err = b.Conn.QueryContext(ctx, query, decodedCursor, limit)
	} else {
		rows, err = b.Conn.QueryContext(ctx, queryWithStops, decodedCursor, limit, pq.Array(locations))
	}

	if err != nil {
		return []*bitsb.BusRoute{}, "", err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			logrus.Error(err)
		}
	}(rows)

	busRoutes := make([]*bitsb.BusRoute, 0, limit)
	for rows.Next() {
		busRoute := bitsb.BusRoute{}
		err = rows.Scan(
			&busRoute.ID,
			&busRoute.Name,
			&busRoute.Number,
			&busRoute.StartTime,
			&busRoute.EndTime,
			&busRoute.Interval,
			pq.Array(&busRoute.LocationIDS),
			&busRoute.MinPrice,
			&busRoute.MaxPrice,
			&busRoute.CreatedAt,
			&busRoute.UpdatedAt,
		)
		if err != nil {
			return []*bitsb.BusRoute{}, "", err
		}
		busRoutes = append(busRoutes, &busRoute)
	}

	var nextCursor string
	if len(busRoutes) > 0 {
		nextCursor = repo.EncodeCursor(busRoutes[len(busRoutes)-1].CreatedAt)
	}
	return busRoutes, nextCursor, nil
}

func (b *BusRouteRepository) SelectByID(ctx context.Context, id int64) (*bitsb.BusRoute, error) {
	query := `SELECT id, name, number, start_time, end_time, interval, location_ids, min_price, max_price, created_at, updated_at FROM bus_routes WHERE id=$1;`
	busRoute := &bitsb.BusRoute{}
	err := b.Conn.QueryRowContext(ctx, query, id).Scan(
		&busRoute.ID,
		&busRoute.Name,
		&busRoute.Number,
		&busRoute.StartTime,
		&busRoute.EndTime,
		&busRoute.Interval,
		pq.Array(&busRoute.LocationIDS),
		&busRoute.MinPrice,
		&busRoute.MaxPrice,
		&busRoute.CreatedAt,
		&busRoute.UpdatedAt,
	)
	return busRoute, err
}

func (b *BusRouteRepository) Insert(ctx context.Context, busRoute *bitsb.BusRoute) error {
	query := `INSERT INTO bus_routes (name, number, start_time, end_time, interval, location_ids, min_price, max_price, created_at, updated_at)
    	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`

	currentTime := time.Now()
	busRoute.CreatedAt = currentTime
	busRoute.UpdatedAt = currentTime

	return b.Conn.QueryRowContext(
		ctx,
		query,
		busRoute.Name,
		busRoute.Number,
		busRoute.StartTime,
		busRoute.EndTime,
		busRoute.Interval,
		pq.Array(busRoute.LocationIDS),
		busRoute.MinPrice,
		busRoute.MaxPrice,
		busRoute.CreatedAt,
		busRoute.UpdatedAt,
	).Scan(&busRoute.ID)
}

func (b *BusRouteRepository) Update(ctx context.Context, busRoute *bitsb.BusRoute) error {
	query := `UPDATE bus_routes 
				SET name=$2, number=$3, start_time=$4, end_time=$5, interval=$6, location_ids=$7, min_price=$8, max_price=$9, updated_at=$10
				WHERE id=$1`

	busRoute.UpdatedAt = time.Now()

	res, err := b.Conn.ExecContext(
		ctx,
		query,
		busRoute.ID,
		busRoute.Name,
		busRoute.Number,
		busRoute.StartTime,
		busRoute.EndTime,
		busRoute.Interval,
		pq.Array(busRoute.LocationIDS),
		busRoute.MinPrice,
		busRoute.MaxPrice,
		busRoute.UpdatedAt,
	)
	if err != nil {
		return err
	}
	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
		err = apperrors.ErrNotFound
	}
	return err
}

func (b *BusRouteRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM bus_routes WHERE id=$1`

	res, err := b.Conn.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
		err = apperrors.ErrNotFound
	}
	return err
}
