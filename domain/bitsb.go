package domain

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/sainak/bitsb/utils/repo"
)

type Location struct {
	ID        int64     `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"createdAt"`
	UpdatedAt time.Time `json:"updated_at" db:"updatedAt"`
}

type LocationForm struct {
	Name string `json:"name"`
}

func (l LocationForm) Bind(r *http.Request) error {
	if l.Name == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}

type (
	LocationStorer interface {
		SelectAll(ctx context.Context, cursor string, limit int64, filters repo.Filters) (locations []*Location, nextCursor string, err error)
		SelectByID(ctx context.Context, id int64) (location *Location, err error)
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
