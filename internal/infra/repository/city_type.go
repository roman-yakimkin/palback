package repository

import (
	"context"
	"database/sql"
	"errors"
	
	"palback/internal/domain/model"
	localErrors "palback/internal/pkg/errors"
)

type CityTypeRepo struct {
	db *sql.DB
}

func NewCityTypeRepo(db *sql.DB) *CityTypeRepo {
	return &CityTypeRepo{
		db: db,
	}
}

type cityTypeDTO struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ShortName string `json:"short_name"`
	Weight    int    `json:"weight"`
}

func (dto *cityTypeDTO) ToModel() model.CityType {
	return model.CityType{
		ID:        dto.ID,
		Name:      dto.Name,
		ShortName: dto.ShortName,
		Weight:    dto.Weight,
	}
}

// Get Получить информацию об одном населенном пункте
func (r *CityTypeRepo) Get(ctx context.Context, id int) (*model.CityType, error) {
	q := `select id, name, short_name, weight from city_types where id = $1`

	var dto cityTypeDTO

	err := r.db.QueryRowContext(ctx, q, id).Scan(&dto.ID, &dto.Name, &dto.ShortName, &dto.Weight)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, localErrors.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	cityType := dto.ToModel()
	return &cityType, nil
}

// GetAll Получить информацию обо всех населенных пунктах
func (r *CityTypeRepo) GetAll(ctx context.Context) ([]model.CityType, error) {
	q := `select id, name, short_name, weight from city_types order by weight`

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cityTypes []model.CityType

	for rows.Next() {
		var dto cityTypeDTO

		err := rows.Scan(&dto.ID, &dto.Name, &dto.ShortName, &dto.Weight)
		if err != nil {
			return nil, err
		}

		cityTypes = append(cityTypes, dto.ToModel())
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return cityTypes, nil
}
