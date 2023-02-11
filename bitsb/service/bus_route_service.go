package service

import (
	"context"

	"github.com/sainak/bitsb/domain"
	"github.com/sainak/bitsb/domain/errors"
	"github.com/sainak/bitsb/pkg/utils"
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

func (b *BusRouteService) ListAll(
	ctx context.Context,
	cursor string,
	limit int64,
	locations []int64,
) ([]*domain.BusRoute, string, error) {
	return b.repo.SelectAll(ctx, cursor, limit, locations)
}

func (b *BusRouteService) GetByID(ctx context.Context, id int64) (*domain.BusRoute, error) {
	busRoute, err := b.repo.SelectByID(ctx, id)
	if err != nil {
		return &domain.BusRoute{}, err
	}
	locations, err := b.locationRepo.SelectByIDArray(ctx, busRoute.LocationIDS)
	if err != nil {
		return &domain.BusRoute{}, err
	}
	for _, l := range locations {
		loc := &domain.LocationForm{Name: l.Name}
		busRoute.Locations = append(busRoute.Locations, loc)
	}
	return busRoute, err
}

func (b *BusRouteService) CalculateTicketPrice(ctx context.Context, id, start, end int64) (int64, error) {
	busRoute, err := b.repo.SelectByID(ctx, id)
	if err != nil {
		return 0, err
	}
	startIndex := utils.IndexOf(busRoute.LocationIDS, start)
	endIndex := utils.IndexOf(busRoute.LocationIDS, end)
	distance := utils.Abs(endIndex - startIndex)
	if startIndex == -1 || endIndex == -1 || distance == 0 {
		return 0, errors.ErrInvalidLocation
	}

	price := utils.Min(int64(distance)*busRoute.MinPrice, busRoute.MaxPrice)
	return price, err
}

func (b *BusRouteService) Create(ctx context.Context, busRoute *domain.BusRoute) error {
	return b.repo.Insert(ctx, busRoute)
}

func (b *BusRouteService) Update(ctx context.Context, busRoute *domain.BusRoute) error {
	return b.repo.Update(ctx, busRoute)
}

func (b *BusRouteService) Delete(ctx context.Context, id int64) error {
	return b.repo.Delete(ctx, id)
}
