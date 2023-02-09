package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"

	"github.com/sainak/bitsb/domain"
	"github.com/sainak/bitsb/domain/errors"
)

type BusRouteRepository struct {
	Conn *sql.DB
}

func NewBusRouteRepository(conn *sql.DB) domain.BusRouteStorer {
	return &BusRouteRepository{conn}
}

func (b *BusRouteRepository) Insert(ctx context.Context, busRoute *domain.BusRoute) (err error) {
	query := `INSERT INTO bus_routes (name, number, start_time, end_time, interval, location_ids, created_at, updated_at)
    	VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`

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
		busRoute.CreatedAt,
		busRoute.UpdatedAt,
	).Scan(&busRoute.ID)
}

func (b *BusRouteRepository) Update(ctx context.Context, busRoute *domain.BusRoute) (err error) {

	query := `UPDATE bus_routes 
				SET name=$2, number=$3, start_time=$4, end_time=$5, interval=$6, location_ids=$7, updated_at=$8
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
		busRoute.UpdatedAt,
	)
	if err != nil {
		return
	}
	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
		err = errors.ErrNotFound
	}
	return
}

func (b *BusRouteRepository) Delete(ctx context.Context, id int64) (err error) {
	query := `DELETE FROM bus_routes WHERE id=$1`

	res, err := b.Conn.ExecContext(ctx, query, id)
	if err != nil {
		return
	}
	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
		err = errors.ErrNotFound
	}
	return
}
