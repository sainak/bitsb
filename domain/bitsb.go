package domain

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/sainak/bitsb/utils/repo"
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
		SelectAll(ctx context.Context, cursor string, limit int64, filters repo.Filters) (locations []*Location, nextCursor string, err error)
		SelectByID(ctx context.Context, id int64) (location *Location, err error)
		SelectByIDArray(ctx context.Context, ids []int64) (locations []*Location, err error)
		Insert(ctx context.Context, location *Location) (err error)
		Update(ctx context.Context, location *Location) (err error)
		Delete(ctx context.Context, id int64) (err error)
	}
	LocationServiceProvider interface {
		ListAll(ctx context.Context, cursor string, limit int64, filters repo.Filters) (locations []*Location, nextCursor string, err error)
		GetByID(ctx context.Context, id int64) (location *Location, err error)
		Create(ctx context.Context, location *Location) (err error)
		Update(ctx context.Context, location *Location) (err error)
		Delete(ctx context.Context, id int64) (err error)
	}
)

// ---- BusRoute ----

type BusRoute struct {
	ID          int64       `json:"id" db:"id"`
	Name        string      `json:"name" db:"name"`
	Number      string      `json:"number" db:"number"`
	StartTime   time.Time   `json:"start_time" db:"start_time"`
	EndTime     time.Time   `json:"end_time" db:"end_time"`
	Interval    int64       `json:"interval" db:"interval"`
	LocationIDS []int64     `json:"location_ids" db:"locations"`
	CreatedAt   time.Time   `json:"created_at" db:"createdAt"`
	UpdatedAt   time.Time   `json:"updated_at" db:"updatedAt"`
	Locations   []*Location `json:"locations"`
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
	LocationIDS []int64   `json:"location_ids"`
}

func (b *BusRouteForm) Bind(r *http.Request) error {
	var errors []string
	if b.Name == "" {
		errors = append(errors, "'name' is required")
	}
	if b.Number == "" {
		errors = append(errors, "'number' is required")
	}
	if b.StartTime.IsZero() {
		errors = append(errors, "'start_time' is required")
	}
	if b.EndTime.IsZero() {
		errors = append(errors, "'end_time' is required")
	}
	if b.Interval == 0 {
		errors = append(errors, "'interval' is required")
	}
	if len(b.LocationIDS) == 0 {
		errors = append(errors, "'location_ids' is required")
	} else if len(b.LocationIDS) < 2 {
		errors = append(errors, "'location_ids' should have atleast 2 stops")
	} else if len(b.LocationIDS) > 10 {
		errors = append(errors, "'location_ids' should have atmost 10 stops")
	}
	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, ", "))
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
		Insert(ctx context.Context, busRoute *BusRoute) (err error)
		Update(ctx context.Context, busRoute *BusRoute) (err error)
		Delete(ctx context.Context, id int64) (err error)
	}

	BusRouteServiceProvider interface {
		Create(ctx context.Context, busRoute *BusRoute) (err error)
		Update(ctx context.Context, busRoute *BusRoute) (err error)
		Delete(ctx context.Context, id int64) (err error)
	}
)
