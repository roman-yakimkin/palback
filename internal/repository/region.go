package repository

import (
	"context"
	"database/sql"
	"errors"
	"palback/internal/domain"
	"palback/internal/domain/model"
	localErrors "palback/internal/pkg/errors"
	"strings"
)

type RegionRepo struct {
	db *sql.DB
}

func NewRegionRepo(db *sql.DB) *RegionRepo {
	return &RegionRepo{
		db: db,
	}
}

type regionDTO struct {
	ID        int    `json:"id"`
	CountryID string `json:"country_id"`
	Name      string `json:"name"`
}

func (dto *regionDTO) ToModel() model.Region {
	return model.Region{
		ID:        dto.ID,
		CountryID: dto.CountryID,
		Name:      dto.Name,
	}
}

func (r *RegionRepo) Get(ctx context.Context, id int) (*model.Region, error) {
	q := `select r.id, r.country_id, r.name from regions r where r.id = $1`

	var dto regionDTO

	err := r.db.QueryRowContext(ctx, q, id).Scan(&dto.ID, &dto.CountryID, &dto.Name)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, localErrors.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	region := dto.ToModel()
	return &region, nil

}

func (r *RegionRepo) GetByCountry(ctx context.Context, countryId string) ([]model.Region, error) {
	q := `select id, country_id, name from regions where country_id = $1 order by name`

	rows, err := r.db.QueryContext(ctx, q, countryId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var regions []model.Region

	for rows.Next() {
		var dto regionDTO

		err := rows.Scan(&dto.ID, &dto.CountryID, &dto.Name)
		if err != nil {
			return nil, err
		}

		regions = append(regions, dto.ToModel())
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return regions, nil
}

func (r *RegionRepo) Create(ctx context.Context, region model.Region) (*model.Region, error) {
	q := `insert into regions (country_id, name) values ($1, $2) returning id`

	var id int
	err := r.db.QueryRowContext(ctx, q, region.CountryID, region.Name).Scan(&id)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "fk_unique_country_and_name"):
			return nil, domain.ErrRegionNotUnique
		default:
			return nil, err
		}
	}

	return &model.Region{
		ID:        id,
		CountryID: region.CountryID,
		Name:      region.Name,
	}, nil

}

func (r *RegionRepo) Update(ctx context.Context, id int, region model.Region) error {
	q := `update regions set country_id = $1, name=$2 where id=$3`

	result, err := r.db.ExecContext(ctx, q, region.CountryID, region.Name, id)
	if err != nil {
		return err
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return localErrors.ErrNotFound
	}

	return nil
}

func (r *RegionRepo) Delete(ctx context.Context, id int) error {
	q := `delete from regions where id=$1`

	result, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return err
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return localErrors.ErrNotFound
	}

	return nil
}
