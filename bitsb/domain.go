package bitsb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/sainak/bitsb/pkg/repo"
)

// ---- Location ----

type Location struct {
	ID        int64     `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"createdAt"`
	UpdatedAt time.Time `json:"updated_at" db:"updatedAt"`
}

type LocationForm struct {
	Name string `json:"name"`
}

func (l *LocationForm) Bind(r *http.Request) error {
	if l.Name == "" {
		return fmt.Errorf("'name' is required")
	}
	return nil
}

type (
	LocationStorer interface {
		SelectAll(ctx context.Context, cursor string, limit int64, filters repo.Filters) ([]*Location, string, error)
		SelectByID(ctx context.Context, id int64) (*Location, error)
		SelectByIDArray(ctx context.Context, ids []int64) ([]*Location, error)
		Insert(ctx context.Context, location *Location) error
		Update(ctx context.Context, location *Location) error
		Delete(ctx context.Context, id int64) error
	}
	LocationServiceProvider interface {
		ListAll(ctx context.Context, cursor string, limit int64, filters repo.Filters) ([]*Location, string, error)
		GetByID(ctx context.Context, id int64) (*Location, error)
		Create(ctx context.Context, location *Location) error
		Update(ctx context.Context, location *Location) error
		Delete(ctx context.Context, id int64) error
	}
)

// ---- BusRoute ----

type BusRoute struct {
	ID          int64           `json:"id" db:"id"`
	Name        string          `json:"name" db:"name"`
	Number      string          `json:"number" db:"number"`
	StartTime   time.Time       `json:"start_time" db:"start_time"`
	EndTime     time.Time       `json:"end_time" db:"end_time"`
	Interval    int64           `json:"interval" db:"interval"`
	LocationIDS []int64         `json:"location_ids" db:"locations"`
	MinPrice    int64           `json:"min_price"`
	MaxPrice    int64           `json:"max_price"`
	CreatedAt   time.Time       `json:"created_at" db:"createdAt"`
	UpdatedAt   time.Time       `json:"updated_at" db:"updatedAt"`
	Locations   []*LocationForm `json:"stops,omitempty"`
}

func (b *BusRoute) MarshalJSON() ([]byte, error) {
	startTime := b.StartTime.Format("15:04")
	endTime := b.EndTime.Format("15:04")
	type Alias BusRoute
	return json.Marshal(&struct {
		*Alias
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
	}{
		Alias:     (*Alias)(b),
		StartTime: startTime,
		EndTime:   endTime,
	})
}

type BusRouteForm struct {
	Name        string    `json:"name"`
	Number      string    `json:"number"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Interval    int64     `json:"interval"`
	MinPrice    int64     `json:"min_price"`
	MaxPrice    int64     `json:"max_price"`
	LocationIDS []int64   `json:"location_ids"`
}

func (b *BusRouteForm) Bind(r *http.Request) error {
	const MaxStops = 10
	var errs []string
	if b.Name == "" {
		errs = append(errs, "'name' is required")
	}
	if b.Number == "" {
		errs = append(errs, "'number' is required")
	}
	if b.StartTime.IsZero() {
		errs = append(errs, "'start_time' is required")
	}
	if b.EndTime.IsZero() {
		errs = append(errs, "'end_time' is required")
	}
	if b.Interval == 0 {
		errs = append(errs, "'interval' is required")
	}
	if len(b.LocationIDS) == 0 {
		errs = append(errs, "'location_ids' is required")
	} else if len(b.LocationIDS) < 2 {
		errs = append(errs, "'location_ids' should have atleast 2 stops")
	} else if len(b.LocationIDS) > MaxStops {
		errs = append(errs, fmt.Sprintf("'location_ids' should have atmost %d stops", MaxStops))
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, ", "))
	}
	return nil
}

func (b *BusRouteForm) UnmarshalJSON(data []byte) error {
	type Alias BusRouteForm
	aux := &struct {
		*Alias
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
	}{
		Alias: (*Alias)(b),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	b.StartTime, _ = time.Parse("15:04", aux.StartTime)
	b.EndTime, _ = time.Parse("15:04", aux.EndTime)
	return nil
}

type (
	BusRouteStorer interface {
		SelectAll(ctx context.Context, cursor string, limit int64, locations []int64) ([]*BusRoute, string, error)
		SelectByID(ctx context.Context, id int64) (*BusRoute, error)
		Insert(ctx context.Context, busRoute *BusRoute) error
		Update(ctx context.Context, busRoute *BusRoute) error
		Delete(ctx context.Context, id int64) error
	}

	BusRouteServiceProvider interface {
		ListAll(ctx context.Context, cursor string, limit int64, locations []int64) ([]*BusRoute, string, error)
		GetByID(ctx context.Context, id int64) (*BusRoute, error)
		CalculateTicketPrice(ctx context.Context, id, start, end int64) (int64, error)
		Create(ctx context.Context, busRoute *BusRoute) error
		Update(ctx context.Context, busRoute *BusRoute) error
		Delete(ctx context.Context, id int64) error
	}
)
