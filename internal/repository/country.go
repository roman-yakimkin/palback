package repository

import (
	"context"
	"database/sql"
	"errors"
	"palback/internal/domain/model"
	localErrors "palback/internal/pkg/errors"
)

type CountryRepo struct {
	db *sql.DB
}

func NewCountryRepo(db *sql.DB) *CountryRepo {
	return &CountryRepo{
		db: db,
	}
}

type countryDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (dto *countryDTO) ToCountry() model.Country {
	return model.Country{
		ID:   dto.ID,
		Name: dto.Name,
	}
}

// Get Получить информацию об одной стране
func (r *CountryRepo) Get(ctx context.Context, id string) (*model.Country, error) {
	q := `select c.id, c.name from countries c where c.id = $1`

	var dto countryDTO

	err := r.db.QueryRowContext(ctx, q, id).Scan(&dto.ID, &dto.Name)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, localErrors.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	country := dto.ToCountry()
	return &country, nil
}

// GetAll Получить информацию обо всех странах
func (r *CountryRepo) GetAll(ctx context.Context) ([]model.Country, error) {
	q := `select id, name from countries order by name`

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var countries []model.Country

	for rows.Next() {
		var dto countryDTO

		err := rows.Scan(&dto.ID, &dto.Name)
		if err != nil {
			return nil, err
		}

		countries = append(countries, dto.ToCountry())
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return countries, nil
}

func (r *CountryRepo) Post(ctx context.Context, country model.Country) (*model.Country, error) {
	q := `insert into countries (id, name) values ($1, $2) returning id`

	var id string
	err := r.db.QueryRowContext(ctx, q, country.ID, country.Name).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &model.Country{
		ID:   id,
		Name: country.Name,
	}, nil
}

func (r *CountryRepo) Put(ctx context.Context, id string, country model.Country) error {
	q := `update countries set name=$1 where id=$2`

	result, err := r.db.ExecContext(ctx, q, country.Name, id)
	if err != nil {
		return err
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return localErrors.ErrNotFound
	}

	return nil
}

func (r *CountryRepo) Delete(ctx context.Context, id string) error {
	q := `delete from countries where id=$1`

	result, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return err
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return localErrors.ErrNotFound
	}

	return nil
}
