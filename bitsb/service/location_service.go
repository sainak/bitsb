package service

import (
	"context"

	"github.com/sainak/bitsb/domain"
	"github.com/sainak/bitsb/pkg/repo"
)

type LocationService struct {
	repo domain.LocationStorer
}

func NewLocationService(r domain.LocationStorer) domain.LocationServiceProvider {
	return &LocationService{
		repo: r,
	}
}

func (l LocationService) ListAll(
	ctx context.Context,
	cursor string,
	limit int64,
	filters repo.Filters,
) ([]*domain.Location, string, error) {
	return l.repo.SelectAll(ctx, cursor, limit, filters)
}

func (l LocationService) GetByID(ctx context.Context, id int64) (*domain.Location, error) {
	return l.repo.SelectByID(ctx, id)
}

func (l LocationService) Create(ctx context.Context, location *domain.Location) error {
	return l.repo.Insert(ctx, location)
}

func (l LocationService) Update(ctx context.Context, location *domain.Location) error {
	return l.repo.Update(ctx, location)
}

func (l LocationService) Delete(ctx context.Context, id int64) error {
	return l.repo.Delete(ctx, id)
}
