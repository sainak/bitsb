package service

import (
	"context"

	"github.com/sainak/bitsb/domain"
	"github.com/sainak/bitsb/utils/repo"
)

type CompanyService struct {
	repo domain.CompanyStorer
}

func NewCompanyService(repo domain.CompanyStorer) domain.CompanyServiceProvider {
	return &CompanyService{
		repo: repo,
	}
}

func (c CompanyService) ListAll(
	ctx context.Context,
	cursor string,
	limit int64,
	filters repo.Filters,
) (companies []*domain.Company, nextCursor string, err error) {
	return c.repo.SelectAll(ctx, cursor, limit, filters)
}

func (c CompanyService) GetByID(ctx context.Context, id int64) (company *domain.Company, err error) {
	return c.repo.SelectByID(ctx, id)
}

func (c CompanyService) Create(ctx context.Context, company *domain.Company) (err error) {
	return c.repo.Insert(ctx, company)
}

func (c CompanyService) Update(ctx context.Context, company *domain.Company) (err error) {
	return c.repo.Update(ctx, company)
}

func (c CompanyService) Delete(ctx context.Context, id int64) (err error) {
	return c.repo.Delete(ctx, id)
}
