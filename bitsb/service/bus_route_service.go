package service

import (
	"context"

	"github.com/sainak/bitsb/domain"
)

type BusRouteService struct {
	repo         domain.BusRouteStorer
	locationRepo domain.LocationStorer
}

func NewBusRouteService(r domain.BusRouteStorer, l domain.LocationStorer) domain.BusRouteServiceProvider {
	return &BusRouteService{
		repo:         r,
		locationRepo: l,
	}
}

func (b *BusRouteService) Create(ctx context.Context, busRoute *domain.BusRoute) (err error) {
	err = b.repo.Insert(ctx, busRoute)
	if err != nil {
		return err
	}
	locations, err := b.locationRepo.SelectByIDArray(ctx, busRoute.LocationIDS)
	if err != nil {
		return
	}
	busRoute.Locations = locations
	return
}

func (b *BusRouteService) Update(ctx context.Context, busRoute *domain.BusRoute) (err error) {
	err = b.repo.Update(ctx, busRoute)
	if err != nil {
		return err
	}
	locations, err := b.locationRepo.SelectByIDArray(ctx, busRoute.LocationIDS)
	if err != nil {
		return
	}
	busRoute.Locations = locations
	return
}

func (b *BusRouteService) Delete(ctx context.Context, id int64) (err error) {
	return b.repo.Delete(ctx, id)
}
