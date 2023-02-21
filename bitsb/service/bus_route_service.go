package service

import (
	"context"

	"github.com/sainak/bitsb/apperrors"
	"github.com/sainak/bitsb/bitsb"
	"github.com/sainak/bitsb/pkg/utils"
)

type BusRouteService struct {
	repo         bitsb.BusRouteStorer
	locationRepo bitsb.LocationStorer
}

func NewBusRouteService(r bitsb.BusRouteStorer, l bitsb.LocationStorer) bitsb.BusRouteServiceProvider {
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
) ([]*bitsb.BusRoute, string, error) {
	return b.repo.SelectAll(ctx, cursor, limit, locations)
}

func (b *BusRouteService) GetByID(ctx context.Context, id int64) (*bitsb.BusRoute, error) {
	busRoute, err := b.repo.SelectByID(ctx, id)
	if err != nil {
		return &bitsb.BusRoute{}, err
	}
	locations, err := b.locationRepo.SelectByIDArray(ctx, busRoute.LocationIDS)
	if err != nil {
		return &bitsb.BusRoute{}, err
	}
	for _, l := range locations {
		loc := &bitsb.LocationForm{Name: l.Name}
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
		return 0, apperrors.ErrInvalidLocation
	}

	price := utils.Min(int64(distance)*busRoute.MinPrice, busRoute.MaxPrice)
	return price, err
}

func (b *BusRouteService) Create(ctx context.Context, busRoute *bitsb.BusRoute) error {
	return b.repo.Insert(ctx, busRoute)
}

func (b *BusRouteService) Update(ctx context.Context, busRoute *bitsb.BusRoute) error {
	return b.repo.Update(ctx, busRoute)
}

func (b *BusRouteService) Delete(ctx context.Context, id int64) error {
	return b.repo.Delete(ctx, id)
}
