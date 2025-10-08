package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"palback/internal/domain"
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
	ID         string `json:"id"`
	Name       string `json:"name"`
	HasRegions bool   `json:"has_regions"`
}

func (dto *countryDTO) ToModel() model.Country {
	return model.Country{
		ID:         dto.ID,
		Name:       dto.Name,
		HasRegions: dto.HasRegions,
	}
}

// Get Получить информацию об одной стране
func (r *CountryRepo) Get(ctx context.Context, id string) (*model.Country, error) {
	q := `select c.id, c.name, c.has_regions from countries c where c.id = $1`

	var dto countryDTO

	err := r.db.QueryRowContext(ctx, q, id).Scan(&dto.ID, &dto.Name, &dto.HasRegions)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, localErrors.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	country := dto.ToModel()
	return &country, nil
}

// GetAll Получить информацию обо всех странах
func (r *CountryRepo) GetAll(ctx context.Context) ([]model.Country, error) {
	q := `select id, name, has_regions from countries order by name`

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var countries []model.Country

	for rows.Next() {
		var dto countryDTO

		err := rows.Scan(&dto.ID, &dto.Name, &dto.HasRegions)
		if err != nil {
			return nil, err
		}

		countries = append(countries, dto.ToModel())
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return countries, nil
}

func (r *CountryRepo) Create(ctx context.Context, country model.Country) (*model.Country, error) {
	q := `insert into countries (id, name, has_regions) values ($1, $2, $3) returning id`

	var id string
	err := r.db.QueryRowContext(ctx, q, country.ID, country.Name, country.HasRegions).Scan(&id)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "countries_id_key"):
			return nil, domain.ErrCountryAlreadyAdded
		case strings.Contains(err.Error(), "countries_name_key"):
			return nil, domain.ErrCountryNameNotUnique
		default:
			return nil, err
		}
	}

	return &model.Country{
		ID:         id,
		Name:       country.Name,
		HasRegions: country.HasRegions,
	}, nil
}

func (r *CountryRepo) Update(ctx context.Context, id string, country model.Country) error {
	q := `update countries set name=$1, has_regions=$2 where id=$3`

	result, err := r.db.ExecContext(ctx, q, country.Name, country.HasRegions, id)
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
