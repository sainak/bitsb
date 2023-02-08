package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/sainak/bitsb/domain"
	"github.com/sainak/bitsb/domain/errors"
	"github.com/sainak/bitsb/utils/repo"
)

type CompanyRepository struct {
	Conn *sql.DB
}

func NewCompanyRepository(conn *sql.DB) domain.CompanyStorer {
	return &CompanyRepository{conn}
}

func (l CompanyRepository) SelectAll(
	ctx context.Context,
	cursor string,
	limit int64,
	filters repo.Filters,
) (companies []*domain.Company, nextCursor string, err error) {
	query := `SELECT id, name, created_at, updated_at, location_id
				FROM companies
				WHERE created_at < $1`
	q := filters.BuildQuery()
	if q != "" {
		query += " AND " + q
	}
	query += ` ORDER BY created_at DESC  LIMIT $2;`

	companies = make([]*domain.Company, 0, limit)
	decodedCursor, err := repo.DecodeCursor(cursor)
	if err != nil {
		err = errors.ErrBadCursor
		return
	}

	rows, err := l.Conn.QueryContext(ctx, query, decodedCursor, limit)
	if err != nil {
		return
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			logrus.Error(err)
		}
	}(rows)

	for rows.Next() {
		company := domain.Company{}
		err = rows.Scan(
			&company.ID,
			&company.Name,
			&company.CreatedAt,
			&company.UpdatedAt,
			&company.LocationID,
			&company.Location.Name,
		)
		if err != nil {
			return
		}
		companies = append(companies, &company)
	}

	if len(companies) == int(limit) {
		nextCursor = repo.EncodeCursor(companies[len(companies)-1].CreatedAt)
	}

	return
}

func (l CompanyRepository) SelectByID(ctx context.Context, id int64) (company *domain.Company, err error) {
	query := `SELECT id, name, created_at, updated_at, location_id
				FROM companies 
				WHERE id = $1;`

	row := l.Conn.QueryRowContext(ctx, query, id)
	company = &domain.Company{}
	err = row.Scan(
		&company.ID,
		&company.Name,
		&company.CreatedAt,
		&company.UpdatedAt,
		&company.LocationID,
	)
	if err != nil {
		return
	}
	return
}

func (l CompanyRepository) Insert(ctx context.Context, company *domain.Company) error {
	query := `INSERT INTO companies (name, created_at, updated_at, location_id) 
			VALUES ($1, $2, $3, $4) 
			RETURNING id, (SELECT l.name FROM locations l WHERE l.id = location_id);`
	currentTime := time.Now()
	company.CreatedAt = currentTime
	company.UpdatedAt = currentTime

	return l.Conn.QueryRowContext(
		ctx,
		query,
		company.Name,
		company.CreatedAt,
		company.UpdatedAt,
		company.LocationID,
	).Scan(&company.ID, &company.Location.Name)
}

func (l CompanyRepository) Update(ctx context.Context, company *domain.Company) (err error) {
	query := `UPDATE companies SET name = $2, updated_at = $3 WHERE id = $1;`

	res, err := l.Conn.ExecContext(ctx, query, company.ID, company.Name, company.UpdatedAt)
	if err != nil {
		return
	}
	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
		err = errors.ErrNotFound
	}
	return
}

func (l CompanyRepository) Delete(ctx context.Context, id int64) (err error) {
	query := `DELETE FROM companies WHERE id = $1;`

	res, err := l.Conn.ExecContext(ctx, query, id)
	if err != nil {
		return
	}
	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
		err = errors.ErrNotFound
	}
	return
}
