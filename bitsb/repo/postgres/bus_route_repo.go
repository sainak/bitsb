package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"

	"github.com/sainak/bitsb/domain"
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
