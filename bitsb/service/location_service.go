package service

import (
	"context"

	"github.com/sainak/bitsb/bitsb"
	"github.com/sainak/bitsb/pkg/repo"
)

type LocationService struct {
	repo bitsb.LocationStorer
}

func NewLocationService(r bitsb.LocationStorer) bitsb.LocationServiceProvider {
	return &LocationService{
		repo: r,
	}
}

func (l LocationService) ListAll(
	ctx context.Context,
	cursor string,
	limit int64,
	filters repo.Filters,
) ([]*bitsb.Location, string, error) {
	return l.repo.SelectAll(ctx, cursor, limit, filters)
}

func (l LocationService) GetByID(ctx context.Context, id int64) (*bitsb.Location, error) {
	return l.repo.SelectByID(ctx, id)
}

func (l LocationService) Create(ctx context.Context, location *bitsb.Location) error {
	return l.repo.Insert(ctx, location)
}

func (l LocationService) Update(ctx context.Context, location *bitsb.Location) error {
	return l.repo.Update(ctx, location)
}

func (l LocationService) Delete(ctx context.Context, id int64) error {
	return l.repo.Delete(ctx, id)
}
