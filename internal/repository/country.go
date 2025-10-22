package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"palback/internal/app"
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
	Weight     int    `json:"weight"`
}

func (dto *countryDTO) ToModel() model.Country {
	return model.Country{
		ID:         dto.ID,
		Name:       dto.Name,
		HasRegions: dto.HasRegions,
		Weight:     dto.Weight,
	}
}

// Get Получить информацию об одной стране
func (r *CountryRepo) Get(ctx context.Context, id string) (*model.Country, error) {
	q := `select c.id, c.name, c.has_regions, c.weight from countries c where c.id = $1`

	var dto countryDTO

	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&dto.ID,
		&dto.Name,
		&dto.HasRegions,
		&dto.Weight,
	)
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
	q := `select id, name, has_regions, weight from countries order by weight desc, name`

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var countries []model.Country

	for rows.Next() {
		var dto countryDTO

		err := rows.Scan(&dto.ID, &dto.Name, &dto.HasRegions, &dto.Weight)
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
	q := `insert into countries (id, name, has_regions, weight) values ($1, $2, $3, $4) returning id`

	var id string
	err := r.db.QueryRowContext(ctx, q,
		country.ID,
		country.Name,
		country.HasRegions,
		country.Weight).Scan(&id)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "countries_id_key"):
			return nil, app.ErrCountryAlreadyAdded
		case strings.Contains(err.Error(), "countries_name_key"):
			return nil, app.ErrCountryNameNotUnique
		default:
			return nil, err
		}
	}

	return &model.Country{
		ID:         id,
		Name:       country.Name,
		HasRegions: country.HasRegions,
		Weight:     country.Weight,
	}, nil
}

func (r *CountryRepo) Update(ctx context.Context, id string, country model.Country) error {
	q := `update countries set name=$1, has_regions=$2, weight=$3 where id=$4`

	result, err := r.db.ExecContext(ctx, q, country.Name, country.HasRegions, country.Weight, id)
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

func (r *CountryRepo) Order(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	squery := make([]string, 0, len(ids))
	for i, id := range ids {
		squery = append(squery, fmt.Sprintf("('%s', %d)", id, len(ids)-i))
	}

	q := `update countries set weight = c2.weight from (values %s) as c2 (id, weight) where c2.id = countries.id`

	_, err := r.db.ExecContext(ctx, fmt.Sprintf(q, strings.Join(squery, ", ")))

	if err != nil {
		return err
	}

	return nil
}
